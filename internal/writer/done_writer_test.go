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

func TestWriteTasksDone(t *testing.T) {
	recent := time.Now().Add(-7 * 24 * time.Hour)
	old := time.Now().Add(-60 * 24 * time.Hour)

	tasks := []models.Task{
		{
			Status:      models.StatusDONE,
			Priority:    models.PriorityHigh,
			Description: "Deploy new API version",
			PageRefs:    []string{"Project Phoenix"},
			SourceFile:  "journals/2025_04_10.md",
			LineNumber:  5,
			Logbook: []models.LogbookEntry{
				{Start: recent.Add(-2 * time.Hour), End: recent, Duration: 2 * time.Hour},
			},
		},
		{
			Status:      models.StatusDONE,
			Priority:    models.PriorityNone,
			Description: "Write unit tests for parser",
			SourceFile:  "journals/2025_02_01.md",
			LineNumber:  3,
			Logbook: []models.LogbookEntry{
				{Start: old.Add(-time.Hour), End: old, Duration: time.Hour},
			},
		},
		{
			Status:      models.StatusDONE,
			Priority:    models.PriorityNone,
			Description: "No logbook task",
			SourceFile:  "pages/backlog.md",
			LineNumber:  10,
		},
		// Active task — should NOT appear in tasks-done.md
		{
			Status:      models.StatusNOW,
			Priority:    models.PriorityNone,
			Description: "Active task",
			SourceFile:  "journals/2026_05_28.md",
			LineNumber:  1,
		},
	}

	idx := indexer.BuildTaskIndex(tasks)
	tmpDir := t.TempDir()

	if err := WriteTasksDone(idx, tmpDir); err != nil {
		t.Fatalf("WriteTasksDone failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, "tasks-done.md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("tasks-done.md was not created")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	s := string(content)

	for _, want := range []string{
		"# Completed Tasks",
		"Total Completed",
		"Total Time Logged",
		"## Recently Completed",
		"## All Completed",
		"Deploy new API version",
		"Write unit tests for parser",
		"No logbook task",
		"journals/2025_04_10.md:5",
		"Project Phoenix",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("expected %q in output", want)
		}
	}

	// Active task must not appear
	if strings.Contains(s, "Active task") {
		t.Error("active (non-DONE) task should not appear in tasks-done.md")
	}
}

func TestWriteTasksDone_Empty(t *testing.T) {
	idx := indexer.BuildTaskIndex([]models.Task{})
	tmpDir := t.TempDir()

	if err := WriteTasksDone(idx, tmpDir); err != nil {
		t.Fatalf("WriteTasksDone failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "tasks-done.md"))
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if !strings.Contains(string(content), "No completed tasks found") {
		t.Error("expected empty-state message")
	}
}

func TestWriteTasksDone_NoneRecent(t *testing.T) {
	old := time.Now().Add(-60 * 24 * time.Hour)

	tasks := []models.Task{
		{
			Status:      models.StatusDONE,
			Description: "Old completed task",
			SourceFile:  "journals/old.md",
			LineNumber:  1,
			Logbook: []models.LogbookEntry{
				{Start: old.Add(-time.Hour), End: old, Duration: time.Hour},
			},
		},
	}

	idx := indexer.BuildTaskIndex(tasks)
	tmpDir := t.TempDir()

	if err := WriteTasksDone(idx, tmpDir); err != nil {
		t.Fatalf("WriteTasksDone failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "tasks-done.md"))
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	s := string(content)

	// No recently completed section when nothing is recent
	if strings.Contains(s, "## Recently Completed") {
		t.Error("should not have Recently Completed section when nothing is recent")
	}

	if !strings.Contains(s, "## All Completed") {
		t.Error("expected All Completed section")
	}
}
