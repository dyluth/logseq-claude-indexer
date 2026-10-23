package writer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dyluth/logseq-claude-indexer/internal/indexer"
	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

// WriteTimelineRecent writes the recent timeline (last 7 days detailed)
func WriteTimelineRecent(index *indexer.TimelineIndex, outputDir string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	filePath := filepath.Join(outputDir, "timeline-recent.md")

	var buf strings.Builder

	// Write header
	fmt.Fprintf(&buf, "# Recent Activity Timeline\n\n")
	fmt.Fprintf(&buf, "Generated: %s\n\n", index.GeneratedAt.Format(time.RFC3339))

	// Get last 7 days.  Snap the cutoff to UTC midnight so two runs that
	// are seconds apart always produce the same set of included days.
	// Journal dates are also parsed as UTC midnight (extractDateFromJournalPath),
	// so an exact midnight-to-midnight comparison is what we want.
	now := index.GeneratedAt.UTC()
	sevenDaysAgo := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -7)
	var recentDays []indexer.TimelineDay
	for _, day := range index.Entries {
		if day.Date.After(sevenDaysAgo) || day.Date.Equal(sevenDaysAgo) {
			recentDays = append(recentDays, day)
		}
	}

	if len(recentDays) == 0 {
		fmt.Fprintf(&buf, "*No activity in the last 7 days.*\n")
		return writeIfChanged(filePath, buf.String())
	}

	fmt.Fprintf(&buf, "**Last 7 Days**: %d days with activity\n\n", len(recentDays))
	fmt.Fprintf(&buf, "---\n\n")

	// Write each day in detail
	for _, day := range recentDays {
		writeDayDetail(&buf, day)
	}

	if err := writeIfChanged(filePath, buf.String()); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}

// WriteTimelineFull writes the complete timeline history
func WriteTimelineFull(index *indexer.TimelineIndex, outputDir string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	filePath := filepath.Join(outputDir, "timeline-full.md")

	var buf strings.Builder

	// Write header
	fmt.Fprintf(&buf, "# Complete Activity Timeline\n\n")
	fmt.Fprintf(&buf, "Generated: %s\n\n", index.GeneratedAt.Format(time.RFC3339))

	if len(index.Entries) == 0 {
		fmt.Fprintf(&buf, "*No activity recorded.*\n")
		return writeIfChanged(filePath, buf.String())
	}

	fmt.Fprintf(&buf, "**Total Days**: %d days with activity\n\n", len(index.Entries))
	fmt.Fprintf(&buf, "---\n\n")

	// Write each day with condensed format
	for _, day := range index.Entries {
		writeDayCondensed(&buf, day)
	}

	if err := writeIfChanged(filePath, buf.String()); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}

// writeDayDetail writes a single day with full task details
func writeDayDetail(f io.Writer, day indexer.TimelineDay) {
	// Date header
	fmt.Fprintf(f, "## %s\n\n", day.Date.Format("Monday, January 2, 2006"))
	fmt.Fprintf(f, "**Journal**: `%s`\n\n", day.JournalPath)

	// Key activity summary
	if len(day.KeyActivity) > 0 {
		fmt.Fprintf(f, "**Activity**:\n")
		for _, activity := range day.KeyActivity {
			fmt.Fprintf(f, "- %s\n", activity)
		}
		fmt.Fprintf(f, "\n")
	}

	// Time logged
	if day.TimeLogged > 0 {
		fmt.Fprintf(f, "**Time Logged**: %s\n\n", formatDuration(day.TimeLogged))
	}

	// List tasks
	if len(day.TasksCreated) > 0 {
		fmt.Fprintf(f, "**Tasks** (%d):\n", len(day.TasksCreated))
		for _, task := range day.TasksCreated {
			writeTimelineTask(f, task)
		}
	}

	fmt.Fprintf(f, "---\n\n")
}

// writeDayCondensed writes a single day in condensed format
func writeDayCondensed(f io.Writer, day indexer.TimelineDay) {
	// Date header (shorter format)
	fmt.Fprintf(f, "## %s\n\n", day.Date.Format("2006-01-02 (Mon)"))

	// Key activity only
	if len(day.KeyActivity) > 0 {
		for _, activity := range day.KeyActivity {
			fmt.Fprintf(f, "- %s\n", activity)
		}
	} else {
		fmt.Fprintf(f, "- *No tasks*\n")
	}

	// Time logged (inline)
	if day.TimeLogged > 0 {
		fmt.Fprintf(f, "- ⏱ %s logged\n", formatDuration(day.TimeLogged))
	}

	fmt.Fprintf(f, "\n")
}

// writeTimelineTask writes a task in lean timeline format
func writeTimelineTask(f io.Writer, task models.Task) {
	description := task.Description
	if len(description) > 80 {
		description = description[:77] + "..."
	}

	// Format: - [STATUS] Description [#A] ⏱ 2h
	line := fmt.Sprintf("- **[%s]** %s", task.Status, description)

	if task.Priority != models.PriorityNone {
		line += fmt.Sprintf(" [#%s]", task.Priority)
	}

	if len(task.Logbook) > 0 {
		line += fmt.Sprintf(" ⏱ %s", formatDuration(task.TotalDuration()))
	}

	fmt.Fprintf(f, "%s\n", line)
}
