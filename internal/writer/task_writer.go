package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dyluth/logseq-claude-indexer/internal/indexer"
	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

// WriteTaskIndex writes the task index to tasks-by-status.md with statistics
func WriteTaskIndex(index *indexer.TaskIndex, outputDir string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	filePath := filepath.Join(outputDir, "tasks-by-status.md")

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	// Write header
	fmt.Fprintf(f, "# Tasks by Status\n\n")
	fmt.Fprintf(f, "Generated: %s\n\n", index.GeneratedAt.Format(time.RFC3339))

	// Write statistics section
	writeStatistics(f, index)

	// Write tasks by status in priority order
	statuses := []struct {
		status models.TaskStatus
		label  string
	}{
		{models.StatusNOW, "NOW"},
		{models.StatusDOING, "DOING"},
		{models.StatusTODO, "TODO"},
		{models.StatusLATER, "LATER"},
		{models.StatusDONE, "DONE"},
	}

	for _, s := range statuses {
		tasks := index.ByStatus[s.status]
		if len(tasks) > 0 {
			fmt.Fprintf(f, "## %s (%d)\n\n", s.label, len(tasks))
			for _, task := range tasks {
				writeLeanTask(f, task)
			}
			fmt.Fprintf(f, "---\n\n")
		}
	}

	return nil
}

// WritePriorityIndex writes high priority tasks to tasks-by-priority.md
func WritePriorityIndex(index *indexer.TaskIndex, outputDir string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	filePath := filepath.Join(outputDir, "tasks-by-priority.md")

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	// Write header
	fmt.Fprintf(f, "# High Priority Tasks [#A]\n\n")
	fmt.Fprintf(f, "Generated: %s\n\n", index.GeneratedAt.Format(time.RFC3339))

	// Get high priority tasks
	highPriorityTasks := index.ByPriority[models.PriorityHigh]
	fmt.Fprintf(f, "**Total High Priority**: %d tasks\n\n", len(highPriorityTasks))

	if len(highPriorityTasks) == 0 {
		fmt.Fprintf(f, "*No high priority tasks found.*\n")
		return nil
	}

	fmt.Fprintf(f, "---\n\n")

	// Group by status
	statuses := []struct {
		status models.TaskStatus
		label  string
	}{
		{models.StatusNOW, "NOW"},
		{models.StatusDOING, "DOING"},
		{models.StatusTODO, "TODO"},
		{models.StatusLATER, "LATER"},
		{models.StatusDONE, "DONE"},
	}

	for _, s := range statuses {
		// Filter high priority tasks by status
		var statusTasks []models.Task
		for _, task := range highPriorityTasks {
			if task.Status == s.status {
				statusTasks = append(statusTasks, task)
			}
		}

		if len(statusTasks) > 0 {
			fmt.Fprintf(f, "## %s (%d)\n\n", s.label, len(statusTasks))
			for _, task := range statusTasks {
				writeFullTask(f, task)
			}
			fmt.Fprintf(f, "---\n\n")
		}
	}

	return nil
}

// writeStatistics writes the statistics section
func writeStatistics(f *os.File, index *indexer.TaskIndex) {
	stats := index.Statistics

	fmt.Fprintf(f, "## Statistics\n\n")
	fmt.Fprintf(f, "- **Total Tasks**: %d\n", index.TotalTasks)
	fmt.Fprintf(f, "- **Completion Rate**: %.1f%% (%d DONE)\n",
		stats.CompletionRate, stats.StatusBreakdown[models.StatusDONE])
	fmt.Fprintf(f, "- **Time Tracking**: %d tasks (%.1f%% adoption)\n",
		stats.WithTimeTracking, stats.TrackingAdoption)
	fmt.Fprintf(f, "- **Total Time Logged**: %s\n", formatDuration(stats.TotalTimeLogged))

	// Priority breakdown
	fmt.Fprintf(f, "\n**By Priority**:\n")
	priorityLabels := []struct {
		priority models.Priority
		label    string
	}{
		{models.PriorityHigh, "High [#A]"},
		{models.PriorityMedium, "Medium [#B]"},
		{models.PriorityLow, "Low [#C]"},
		{models.PriorityNone, "None"},
	}
	for _, p := range priorityLabels {
		count := stats.PriorityBreakdown[p.priority]
		if count > 0 {
			fmt.Fprintf(f, "- %s: %d\n", p.label, count)
		}
	}

	// Status breakdown
	fmt.Fprintf(f, "\n**By Status**:\n")
	for _, s := range []models.TaskStatus{
		models.StatusNOW,
		models.StatusDOING,
		models.StatusTODO,
		models.StatusLATER,
		models.StatusDONE,
	} {
		count := stats.StatusBreakdown[s]
		if count > 0 {
			fmt.Fprintf(f, "- %s: %d\n", s, count)
		}
	}

	fmt.Fprintf(f, "\n---\n\n")
}

// writeLeanTask writes a task with truncated description (token-optimized)
func writeLeanTask(f *os.File, task models.Task) {
	description := task.Description
	if len(description) > 100 {
		description = description[:97] + "..."
	}

	// Single line format with priority indicator
	priorityIndicator := ""
	if task.Priority != models.PriorityNone {
		priorityIndicator = fmt.Sprintf(" [#%s]", task.Priority)
	}

	timeInfo := ""
	if len(task.Logbook) > 0 {
		timeInfo = fmt.Sprintf(" ⏱ %s", formatDuration(task.TotalDuration()))
	}

	fmt.Fprintf(f, "- **%s**%s%s `%s:%d`\n",
		description, priorityIndicator, timeInfo, task.SourceFile, task.LineNumber)
}

// writeFullTask writes a task with full details (for high priority tasks)
func writeFullTask(f *os.File, task models.Task) {
	fmt.Fprintf(f, "### %s\n", task.Description)
	fmt.Fprintf(f, "- **File**: `%s:%d`\n", task.SourceFile, task.LineNumber)

	// Write page references if present
	if len(task.PageRefs) > 0 {
		fmt.Fprintf(f, "- **References**: ")
		for i, ref := range task.PageRefs {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "[[%s]]", ref)
		}
		fmt.Fprintf(f, "\n")
	}

	// Write time tracking info if present
	if len(task.Logbook) > 0 {
		totalDuration := task.TotalDuration()
		fmt.Fprintf(f, "- **Time Logged**: %s (%d entries)\n",
			formatDuration(totalDuration), len(task.Logbook))

		// Show most recent entry
		mostRecent := task.Logbook[len(task.Logbook)-1]
		fmt.Fprintf(f, "- **Last Activity**: %s\n", mostRecent.End.Format("2006-01-02 15:04"))
	}

	fmt.Fprintf(f, "\n")
}

// formatDuration converts a duration to human-readable format
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		if minutes > 0 {
			return fmt.Sprintf("%dh %dm", hours, minutes)
		}
		return fmt.Sprintf("%dh", hours)
	}

	if minutes > 0 {
		if seconds > 0 {
			return fmt.Sprintf("%dm %ds", minutes, seconds)
		}
		return fmt.Sprintf("%dm", minutes)
	}

	return fmt.Sprintf("%ds", seconds)
}
