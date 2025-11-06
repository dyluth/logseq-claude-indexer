package writer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dyluth/logseq-claude-indexer/internal/indexer"
	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

func TestWriteTimelineRecent(t *testing.T) {
	// Create test timeline with recent dates
	now := time.Now()
	index := &indexer.TimelineIndex{
		GeneratedAt: now,
		Entries: []indexer.TimelineDay{
			{
				Date:        now.AddDate(0, 0, -1), // Yesterday
				JournalPath: "journals/yesterday.md",
				TasksCreated: []models.Task{
					{
						Status:      models.StatusNOW,
						Priority:    models.PriorityHigh,
						Description: "Urgent task from yesterday",
					},
				},
				TimeLogged: 2 * time.Hour,
				KeyActivity: []string{
					"1 NOW task",
					"🔥 Urgent task from yesterday",
				},
			},
			{
				Date:         now.AddDate(0, 0, -3), // 3 days ago
				JournalPath:  "journals/three-days-ago.md",
				TasksCreated: []models.Task{},
				TimeLogged:   0,
				KeyActivity:  []string{},
			},
			{
				Date:        now.AddDate(0, 0, -10), // 10 days ago (should not appear)
				JournalPath: "journals/ten-days-ago.md",
				TasksCreated: []models.Task{
					{
						Status:      models.StatusTODO,
						Description: "Old task",
					},
				},
			},
		},
	}

	// Write to temp directory
	tmpDir := t.TempDir()
	err := WriteTimelineRecent(index, tmpDir)
	if err != nil {
		t.Fatalf("WriteTimelineRecent failed: %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(tmpDir, "timeline-recent.md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("timeline-recent.md was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// Check for expected content
	if !strings.Contains(contentStr, "# Recent Activity Timeline") {
		t.Error("Missing header")
	}

	if !strings.Contains(contentStr, "Last 7 Days") {
		t.Error("Missing last 7 days summary")
	}

	if !strings.Contains(contentStr, "Urgent task from yesterday") {
		t.Error("Missing task from yesterday")
	}

	if !strings.Contains(contentStr, "journals/yesterday.md") {
		t.Error("Missing journal path")
	}

	if !strings.Contains(contentStr, "Time Logged") {
		t.Error("Missing time logged")
	}

	// Should NOT contain 10-day-old task
	if strings.Contains(contentStr, "Old task") {
		t.Error("Should not include tasks older than 7 days")
	}
}

func TestWriteTimelineRecent_NoRecentActivity(t *testing.T) {
	// Create timeline with only old dates
	index := &indexer.TimelineIndex{
		GeneratedAt: time.Now(),
		Entries: []indexer.TimelineDay{
			{
				Date:        time.Now().AddDate(0, 0, -10),
				JournalPath: "journals/old.md",
			},
		},
	}

	tmpDir := t.TempDir()
	err := WriteTimelineRecent(index, tmpDir)
	if err != nil {
		t.Fatalf("WriteTimelineRecent failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "timeline-recent.md"))
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Should indicate no recent activity
	if !strings.Contains(string(content), "No activity in the last 7 days") {
		t.Error("Should indicate no recent activity")
	}
}

func TestWriteTimelineFull(t *testing.T) {
	// Create test timeline
	index := &indexer.TimelineIndex{
		GeneratedAt: time.Now(),
		Entries: []indexer.TimelineDay{
			{
				Date:        time.Date(2025, 11, 6, 0, 0, 0, 0, time.UTC),
				JournalPath: "journals/2025_11_06.md",
				TasksCreated: []models.Task{
					{
						Status:      models.StatusNOW,
						Description: "Task 1",
					},
				},
				TimeLogged: 1 * time.Hour,
				KeyActivity: []string{
					"1 NOW task",
				},
			},
			{
				Date:        time.Date(2025, 11, 5, 0, 0, 0, 0, time.UTC),
				JournalPath: "journals/2025_11_05.md",
				TasksCreated: []models.Task{
					{
						Status:      models.StatusDONE,
						Description: "Task 2",
					},
				},
				KeyActivity: []string{
					"1 DONE task",
				},
			},
		},
	}

	// Write to temp directory
	tmpDir := t.TempDir()
	err := WriteTimelineFull(index, tmpDir)
	if err != nil {
		t.Fatalf("WriteTimelineFull failed: %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(tmpDir, "timeline-full.md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("timeline-full.md was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// Check for expected content
	if !strings.Contains(contentStr, "# Complete Activity Timeline") {
		t.Error("Missing header")
	}

	if !strings.Contains(contentStr, "Total Days") {
		t.Error("Missing total days summary")
	}

	// Should contain both dates
	if !strings.Contains(contentStr, "2025-11-06") {
		t.Error("Missing Nov 6 entry")
	}

	if !strings.Contains(contentStr, "2025-11-05") {
		t.Error("Missing Nov 5 entry")
	}

	// Should contain key activity
	if !strings.Contains(contentStr, "1 NOW task") {
		t.Error("Missing NOW task activity")
	}

	if !strings.Contains(contentStr, "1 DONE task") {
		t.Error("Missing DONE task activity")
	}
}

func TestWriteTimelineFull_EmptyTimeline(t *testing.T) {
	index := &indexer.TimelineIndex{
		GeneratedAt: time.Now(),
		Entries:     []indexer.TimelineDay{},
	}

	tmpDir := t.TempDir()
	err := WriteTimelineFull(index, tmpDir)
	if err != nil {
		t.Fatalf("WriteTimelineFull failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "timeline-full.md"))
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Should indicate no activity
	if !strings.Contains(string(content), "No activity recorded") {
		t.Error("Should indicate no activity")
	}
}

func TestWriteTimelineTask_Truncation(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Create a task with a very long description
	longDesc := strings.Repeat("a", 100)
	task := models.Task{
		Status:      models.StatusTODO,
		Priority:    models.PriorityNone,
		Description: longDesc,
	}

	writeTimelineTask(tmpFile, task)
	tmpFile.Sync()

	// Read content
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Should be truncated to 80 chars
	if !strings.Contains(string(content), "...") {
		t.Error("Long description should be truncated with ...")
	}
}

func TestWriteTimelineTask_WithPriorityAndTime(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	task := models.Task{
		Status:      models.StatusNOW,
		Priority:    models.PriorityHigh,
		Description: "Important task",
		Logbook: []models.LogbookEntry{
			{Duration: 3 * time.Hour},
		},
	}

	writeTimelineTask(tmpFile, task)
	tmpFile.Sync()

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)

	// Should show status
	if !strings.Contains(contentStr, "[NOW]") {
		t.Error("Should show status")
	}

	// Should show priority
	if !strings.Contains(contentStr, "[#A]") {
		t.Error("Should show priority")
	}

	// Should show time
	if !strings.Contains(contentStr, "⏱") {
		t.Error("Should show time tracking indicator")
	}

	if !strings.Contains(contentStr, "3h") {
		t.Error("Should show duration")
	}
}

func TestWriteDayDetail(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	day := indexer.TimelineDay{
		Date:        time.Date(2025, 11, 6, 0, 0, 0, 0, time.UTC),
		JournalPath: "journals/2025_11_06.md",
		TasksCreated: []models.Task{
			{
				Status:      models.StatusNOW,
				Description: "Test task",
			},
		},
		TimeLogged: 2 * time.Hour,
		KeyActivity: []string{
			"1 NOW task",
		},
	}

	writeDayDetail(tmpFile, day)
	tmpFile.Sync()

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)

	// Check all expected elements
	expectedStrings := []string{
		"November 6, 2025", // Date (without checking day of week)
		"Journal",
		"journals/2025_11_06.md",
		"Activity",
		"1 NOW task",
		"Time Logged",
		"2h",
		"Tasks",
		"Test task",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Expected content to contain %q, but it didn't", expected)
		}
	}
}

func TestWriteDayCondensed(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	day := indexer.TimelineDay{
		Date:        time.Date(2025, 11, 6, 0, 0, 0, 0, time.UTC),
		JournalPath: "journals/2025_11_06.md",
		TasksCreated: []models.Task{
			{Status: models.StatusNOW},
		},
		TimeLogged: 1 * time.Hour,
		KeyActivity: []string{
			"1 NOW task",
		},
	}

	writeDayCondensed(tmpFile, day)
	tmpFile.Sync()

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)

	// Check condensed format (date format, without checking exact day of week)
	if !strings.Contains(contentStr, "2025-11-06") {
		t.Error("Should use condensed date format with year-month-day")
	}

	if !strings.Contains(contentStr, "1 NOW task") {
		t.Error("Should show key activity")
	}

	if !strings.Contains(contentStr, "⏱ 1h logged") {
		t.Error("Should show inline time logged")
	}
}

func TestWriteDayCondensed_NoActivity(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	day := indexer.TimelineDay{
		Date:         time.Date(2025, 11, 6, 0, 0, 0, 0, time.UTC),
		JournalPath:  "journals/2025_11_06.md",
		TasksCreated: []models.Task{},
		KeyActivity:  []string{},
	}

	writeDayCondensed(tmpFile, day)
	tmpFile.Sync()

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Should indicate no tasks
	if !strings.Contains(string(content), "No tasks") {
		t.Error("Should indicate no tasks")
	}
}
