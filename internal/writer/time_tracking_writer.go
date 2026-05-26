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

// WriteTimeTracking generates a time-tracking.md file with analytics
func WriteTimeTracking(index *indexer.TimeTrackingIndex, outputDir string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	outputPath := filepath.Join(outputDir, "time-tracking.md")

	var buf strings.Builder

	// Header
	fmt.Fprintf(&buf, "# Time Tracking Analytics\n\n")
	fmt.Fprintf(&buf, "Generated: %s\n\n", time.Now().UTC().Format(time.RFC3339))

	// Overall Statistics
	fmt.Fprintf(&buf, "## Summary\n\n")
	fmt.Fprintf(&buf, "- **Total Time Logged**: %s\n", formatDuration(index.TotalTimeLogged))
	fmt.Fprintf(&buf, "- **Tasks Tracked**: %d / %d (%.1f%% adoption)\n",
		index.Statistics.TasksWithTracking,
		index.Statistics.TotalTasks,
		index.Statistics.AdoptionRate)
	if index.Statistics.AvgTimePerTask > 0 {
		fmt.Fprintf(&buf, "- **Avg Time/Task**: %s\n", formatDuration(index.Statistics.AvgTimePerTask))
	}
	if index.Statistics.MostProductiveWeek.TimeLogged > 0 {
		fmt.Fprintf(&buf, "- **Most Productive Week**: %s (%s)\n",
			index.Statistics.MostProductiveWeek.WeekStart.Format("2006-01-02"),
			formatDuration(index.Statistics.MostProductiveWeek.TimeLogged))
	}
	fmt.Fprintf(&buf, "\n---\n\n")

	// Top Projects
	if len(index.TopProjects) > 0 {
		fmt.Fprintf(&buf, "## Top Projects\n\n")
		// Show top 10 projects
		limit := 10
		if len(index.TopProjects) < limit {
			limit = len(index.TopProjects)
		}
		for i := 0; i < limit; i++ {
			proj := index.TopProjects[i]
			fmt.Fprintf(&buf, "### %s\n", proj.Project)
			fmt.Fprintf(&buf, "- **Time**: %s (%d tasks, avg %s/task)\n",
				formatDuration(proj.TimeLogged),
				proj.TaskCount,
				formatDuration(proj.AvgTimePerTask))
		}
		fmt.Fprintf(&buf, "\n---\n\n")
	}

	// Weekly Breakdown (last 8 weeks only for token efficiency)
	if len(index.WeeklySummary) > 0 {
		fmt.Fprintf(&buf, "## Weekly Breakdown\n\n")
		limit := 8
		if len(index.WeeklySummary) < limit {
			limit = len(index.WeeklySummary)
		}
		for i := 0; i < limit; i++ {
			week := index.WeeklySummary[i]
			fmt.Fprintf(&buf, "- **Week of %s**: %s (%d tasks)\n",
				week.WeekStart.Format("2006-01-02"),
				formatDuration(week.TimeLogged),
				week.TaskCount)
		}
		if len(index.WeeklySummary) > limit {
			fmt.Fprintf(&buf, "\n*Showing last %d weeks of %d total*\n", limit, len(index.WeeklySummary))
		}
		fmt.Fprintf(&buf, "\n---\n\n")
	}

	// By Priority
	if len(index.ByPriority) > 0 {
		fmt.Fprintf(&buf, "## By Priority\n\n")
		priorities := []models.Priority{models.PriorityHigh, models.PriorityMedium, models.PriorityLow, models.PriorityNone}
		for _, priority := range priorities {
			if duration, exists := index.ByPriority[priority]; exists && duration > 0 {
				label := string(priority)
				if label == "" {
					label = "None"
				}
				fmt.Fprintf(&buf, "- **[#%s]**: %s\n", label, formatDuration(duration))
			}
		}
		fmt.Fprintf(&buf, "\n---\n\n")
	}

	// By Status
	if len(index.ByStatus) > 0 {
		fmt.Fprintf(&buf, "## By Status\n\n")
		statuses := []models.TaskStatus{
			models.StatusDONE,
			models.StatusNOW,
			models.StatusDOING,
			models.StatusTODO,
			models.StatusLATER,
		}
		for _, status := range statuses {
			if duration, exists := index.ByStatus[status]; exists && duration > 0 {
				fmt.Fprintf(&buf, "- **%s**: %s\n", status, formatDuration(duration))
			}
		}
		fmt.Fprintf(&buf, "\n")
	}

	if err := writeIfChanged(outputPath, buf.String()); err != nil {
		return fmt.Errorf("writing time tracking file: %w", err)
	}
	return nil
}
