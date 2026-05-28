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

func TestWriteUpcomingTasksIndex(t *testing.T) {
	now := time.Now()
	fiveDays := now.AddDate(0, 0, 5)
	tenDays := now.AddDate(0, 0, 10)

	tasks := []models.Task{
		{
			Status:        models.StatusTODO,
			Priority:      models.PriorityHigh,
			Description:   "review deployment plan",
			SourceFile:    "journals/2026_05_23.md",
			LineNumber:    12,
			ScheduledDate: &fiveDays,
		},
		{
			Status:        models.StatusLATER,
			Priority:      models.PriorityNone,
			Description:   "check analytics dashboard",
			SourceFile:    "pages/Sprint.md",
			LineNumber:    5,
			ScheduledDate: &tenDays,
		},
	}

	idx := indexer.BuildTaskIndex(tasks)

	tmpDir := t.TempDir()
	if err := WriteUpcomingTasksIndex(idx, tmpDir); err != nil {
		t.Fatalf("WriteUpcomingTasksIndex failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, "tasks-upcoming.md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("tasks-upcoming.md was not created")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	s := string(content)

	for _, want := range []string{
		"# Upcoming Tasks",
		"Total Upcoming",
		"review deployment plan",
		"check analytics dashboard",
		"journals/2026_05_23.md:12",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestWriteUpcomingTasksIndex_Empty(t *testing.T) {
	idx := indexer.BuildTaskIndex([]models.Task{})

	tmpDir := t.TempDir()
	if err := WriteUpcomingTasksIndex(idx, tmpDir); err != nil {
		t.Fatalf("WriteUpcomingTasksIndex failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "tasks-upcoming.md"))
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if !strings.Contains(string(content), "No tasks scheduled") {
		t.Error("expected empty-state message")
	}
}
