package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dyluth/logseq-claude-indexer/internal/indexer"
	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

// WriteUpcomingTasksIndex writes tasks scheduled 3–14 days from now to tasks-upcoming.md.
func WriteUpcomingTasksIndex(index *indexer.TaskIndex, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	filePath := filepath.Join(outputDir, "tasks-upcoming.md")

	var buf strings.Builder

	fmt.Fprintf(&buf, "# Upcoming Tasks (Next 3–14 Days)\n\n")
	fmt.Fprintf(&buf, "Generated: %s\n\n", index.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(&buf, "**Total Upcoming**: %d tasks\n\n", len(index.Upcoming))

	if len(index.Upcoming) == 0 {
		fmt.Fprintf(&buf, "*No tasks scheduled in the next 3–14 days.*\n")
		return writeIfChanged(filePath, buf.String())
	}

	fmt.Fprintf(&buf, "---\n\n")

	// Group by date
	type dateGroup struct {
		date  time.Time
		label string
		tasks []models.Task
	}

	var groups []dateGroup
	dateIndex := map[time.Time]int{}

	for _, task := range index.Upcoming {
		day := time.Date(task.ScheduledDate.Year(), task.ScheduledDate.Month(), task.ScheduledDate.Day(), 0, 0, 0, 0, time.UTC)
		if idx, ok := dateIndex[day]; ok {
			groups[idx].tasks = append(groups[idx].tasks, task)
		} else {
			dateIndex[day] = len(groups)
			groups = append(groups, dateGroup{
				date:  day,
				label: task.ScheduledDate.Format("Mon Jan 2, 2006"),
				tasks: []models.Task{task},
			})
		}
	}

	for _, g := range groups {
		fmt.Fprintf(&buf, "## %s (%d)\n\n", g.label, len(g.tasks))
		for _, task := range g.tasks {
			priorityIndicator := ""
			if task.Priority != models.PriorityNone {
				priorityIndicator = fmt.Sprintf(" [#%s]", task.Priority)
			}
			fmt.Fprintf(&buf, "- **%s** (%s)%s `%s:%d`\n",
				task.Description, task.Status, priorityIndicator, task.SourceFile, task.LineNumber)
		}
		fmt.Fprintf(&buf, "\n")
	}

	if err := writeIfChanged(filePath, buf.String()); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}
