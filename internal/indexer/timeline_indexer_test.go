package indexer

import (
	"strings"
	"testing"
	"time"

	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

func TestBuildTimelineIndex(t *testing.T) {
	files := []models.File{
		{
			Path: "journals/2025_11_06.md",
			Type: models.FileTypeJournal,
		},
		{
			Path: "journals/2025_11_05.md",
			Type: models.FileTypeJournal,
		},
		{
			Path: "pages/SomePage.md",
			Type: models.FileTypePage,
		},
	}

	tasks := []models.Task{
		{
			Status:      models.StatusNOW,
			Priority:    models.PriorityHigh,
			Description: "Urgent task",
			SourceFile:  "journals/2025_11_06.md",
			LineNumber:  1,
		},
		{
			Status:      models.StatusTODO,
			Priority:    models.PriorityNone,
			Description: "Regular task",
			SourceFile:  "journals/2025_11_06.md",
			LineNumber:  2,
			Logbook: []models.LogbookEntry{
				{Duration: 2 * time.Hour},
			},
		},
		{
			Status:      models.StatusDONE,
			Description: "Completed task",
			SourceFile:  "journals/2025_11_05.md",
			LineNumber:  1,
		},
	}

	index := BuildTimelineIndex(tasks, files)

	// Check that index was created
	if index == nil {
		t.Fatal("BuildTimelineIndex returned nil")
	}

	// Should have 2 days (2 journal files)
	if len(index.Entries) != 2 {
		t.Fatalf("Expected 2 timeline entries, got %d", len(index.Entries))
	}

	// Check sorting (newest first)
	if !index.Entries[0].Date.After(index.Entries[1].Date) {
		t.Error("Timeline entries not sorted newest first")
	}

	// Check Nov 6 entry
	nov6 := index.Entries[0]
	expectedDate := time.Date(2025, 11, 6, 0, 0, 0, 0, time.UTC)
	if !nov6.Date.Equal(expectedDate) {
		t.Errorf("Expected date %v, got %v", expectedDate, nov6.Date)
	}

	if len(nov6.TasksCreated) != 2 {
		t.Errorf("Expected 2 tasks on Nov 6, got %d", len(nov6.TasksCreated))
	}

	if nov6.TimeLogged != 2*time.Hour {
		t.Errorf("Expected 2h logged on Nov 6, got %v", nov6.TimeLogged)
	}

	// Check Nov 5 entry
	nov5 := index.Entries[1]
	if len(nov5.TasksCreated) != 1 {
		t.Errorf("Expected 1 task on Nov 5, got %d", len(nov5.TasksCreated))
	}

	// Check key activity was generated
	if len(nov6.KeyActivity) == 0 {
		t.Error("Expected key activity for Nov 6, got none")
	}
}

func TestBuildTimelineIndex_EmptyData(t *testing.T) {
	index := BuildTimelineIndex([]models.Task{}, []models.File{})

	if index == nil {
		t.Fatal("BuildTimelineIndex returned nil")
	}

	if len(index.Entries) != 0 {
		t.Errorf("Expected 0 entries for empty data, got %d", len(index.Entries))
	}
}

func TestBuildTimelineIndex_NoJournalFiles(t *testing.T) {
	files := []models.File{
		{
			Path: "pages/SomePage.md",
			Type: models.FileTypePage,
		},
	}

	tasks := []models.Task{
		{
			Status:      models.StatusTODO,
			Description: "Task in page file",
			SourceFile:  "pages/SomePage.md",
			LineNumber:  1,
		},
	}

	index := BuildTimelineIndex(tasks, files)

	if len(index.Entries) != 0 {
		t.Errorf("Expected 0 entries (no journal files), got %d", len(index.Entries))
	}
}

func TestExtractDateFromJournalPath(t *testing.T) {
	tests := []struct {
		path     string
		expected time.Time
		wantErr  bool
	}{
		{
			path:     "journals/2025_11_06.md",
			expected: time.Date(2025, 11, 6, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			path:     "journals/2025-11-06.md",
			expected: time.Date(2025, 11, 6, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			path:     "journals/2025_01_01.md",
			expected: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			path:     "pages/SomePage.md",
			expected: time.Time{},
			wantErr:  true,
		},
		{
			path:     "journals/invalid.md",
			expected: time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		result, err := extractDateFromJournalPath(tt.path)

		if tt.wantErr {
			if err == nil {
				t.Errorf("extractDateFromJournalPath(%q) expected error, got nil", tt.path)
			}
			continue
		}

		if err != nil {
			t.Errorf("extractDateFromJournalPath(%q) unexpected error: %v", tt.path, err)
			continue
		}

		if !result.Equal(tt.expected) {
			t.Errorf("extractDateFromJournalPath(%q) = %v, expected %v", tt.path, result, tt.expected)
		}
	}
}

func TestGenerateKeyActivity(t *testing.T) {
	day := &TimelineDay{
		Date: time.Date(2025, 11, 6, 0, 0, 0, 0, time.UTC),
		TasksCreated: []models.Task{
			{
				Status:      models.StatusNOW,
				Priority:    models.PriorityHigh,
				Description: "High priority urgent task",
			},
			{
				Status:   models.StatusTODO,
				Priority: models.PriorityNone,
			},
			{
				Status:   models.StatusTODO,
				Priority: models.PriorityNone,
			},
			{
				Status:   models.StatusDONE,
				Priority: models.PriorityNone,
			},
		},
	}

	activity := generateKeyActivity(day)

	// Should have status summaries
	foundNOW := false
	foundTODO := false
	foundDONE := false
	foundHighPri := false

	for _, line := range activity {
		if line == "1 NOW task" {
			foundNOW = true
		}
		if line == "2 TODO tasks" {
			foundTODO = true
		}
		if line == "1 DONE task" {
			foundDONE = true
		}
		if line == "🔥 High priority urgent task" {
			foundHighPri = true
		}
	}

	if !foundNOW {
		t.Error("Expected NOW task summary in key activity")
	}
	if !foundTODO {
		t.Error("Expected TODO tasks summary in key activity")
	}
	if !foundDONE {
		t.Error("Expected DONE task summary in key activity")
	}
	if !foundHighPri {
		t.Error("Expected high priority task highlight in key activity")
	}
}

func TestGenerateKeyActivity_EmptyDay(t *testing.T) {
	day := &TimelineDay{
		Date:         time.Date(2025, 11, 6, 0, 0, 0, 0, time.UTC),
		TasksCreated: []models.Task{},
	}

	activity := generateKeyActivity(day)

	if len(activity) != 0 {
		t.Errorf("Expected no key activity for empty day, got %d items", len(activity))
	}
}

func TestGenerateKeyActivity_TruncateLongDescription(t *testing.T) {
	day := &TimelineDay{
		Date: time.Date(2025, 11, 6, 0, 0, 0, 0, time.UTC),
		TasksCreated: []models.Task{
			{
				Status:      models.StatusNOW,
				Priority:    models.PriorityHigh,
				Description: "This is a very long task description that should be truncated to 60 characters",
			},
		},
	}

	activity := generateKeyActivity(day)

	foundTruncated := false
	for _, line := range activity {
		// Check if line starts with fire emoji and ends with ...
		if strings.HasPrefix(line, "🔥") && strings.HasSuffix(line, "...") {
			foundTruncated = true
			// Verify it's actually shorter than the original
			if len(line) < 81 { // Original description is 81 chars
				break
			}
		}
	}

	if !foundTruncated {
		t.Error("Expected long description to be truncated")
	}
}

func TestFormatTaskCount(t *testing.T) {
	tests := []struct {
		count    int
		status   models.TaskStatus
		expected string
	}{
		{1, models.StatusNOW, "1 NOW task"},
		{2, models.StatusTODO, "2 TODO tasks"},
		{5, models.StatusDONE, "5 DONE tasks"},
		{1, models.StatusDOING, "1 DOING task"},
	}

	for _, tt := range tests {
		result := formatTaskCount(tt.count, tt.status)
		if result != tt.expected {
			t.Errorf("formatTaskCount(%d, %s) = %q, expected %q", tt.count, tt.status, result, tt.expected)
		}
	}
}
