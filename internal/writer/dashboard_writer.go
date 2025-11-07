package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dyluth/logseq-claude-indexer/internal/indexer"
	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

// WriteDashboard generates a dashboard.md file with aggregated overview
func WriteDashboard(
	taskIndex *indexer.TaskIndex,
	graphIndex *indexer.ReferenceGraph,
	timelineIndex *indexer.TimelineIndex,
	missingPagesIndex *indexer.MissingPagesIndex,
	timeTrackingIndex *indexer.TimeTrackingIndex,
	outputDir string,
) error {
	outputPath := filepath.Join(outputDir, "dashboard.md")
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create dashboard file: %w", err)
	}
	defer f.Close()

	// Header
	fmt.Fprintf(f, "# Knowledge Dashboard\n\n")
	fmt.Fprintf(f, "**Generated**: %s\n\n", time.Now().UTC().Format(time.RFC3339))

	// Quick Stats
	fmt.Fprintf(f, "## 📊 Quick Stats\n\n")
	fmt.Fprintf(f, "- **Total Tasks**: %d\n", taskIndex.TotalTasks)

	// Calculate completed tasks from status breakdown
	doneCount := 0
	if count, exists := taskIndex.Statistics.StatusBreakdown[models.StatusDONE]; exists {
		doneCount = count
	}
	fmt.Fprintf(f, "- **Completion Rate**: %.1f%% (%d DONE)\n",
		taskIndex.Statistics.CompletionRate, doneCount)

	if timeTrackingIndex.Statistics.TasksWithTracking > 0 {
		fmt.Fprintf(f, "- **Time Tracking**: %.1f%% adoption, %s logged\n",
			timeTrackingIndex.Statistics.AdoptionRate,
			formatDuration(timeTrackingIndex.TotalTimeLogged))
	}

	// Calculate total references
	totalRefs := 0
	for _, node := range graphIndex.Nodes {
		totalRefs += len(node.InboundRefs)
	}
	fmt.Fprintf(f, "- **Knowledge Graph**: %d pages, %d references\n",
		len(graphIndex.Nodes), totalRefs)
	fmt.Fprintf(f, "\n")

	// Current Priorities (High priority NOW and TODO tasks)
	highPriorityTasks := []models.Task{}
	if tasks, exists := taskIndex.ByPriority[models.PriorityHigh]; exists {
		for _, task := range tasks {
			if task.Status == models.StatusNOW || task.Status == models.StatusTODO {
				highPriorityTasks = append(highPriorityTasks, task)
			}
		}
	}

	if len(highPriorityTasks) > 0 {
		fmt.Fprintf(f, "## 🎯 Current Priorities [#A]\n\n")
		count := 0
		for _, task := range highPriorityTasks {
			if count >= 5 { // Show max 5
				break
			}
			desc := task.Description
			if len(desc) > 80 {
				desc = desc[:77] + "..."
			}
			fmt.Fprintf(f, "- **[%s]** %s `%s:%d`\n",
				task.Status, desc, task.SourceFile, task.LineNumber)
			count++
		}
		highPriorityCount := 0
		if count, exists := taskIndex.Statistics.PriorityBreakdown[models.PriorityHigh]; exists {
			highPriorityCount = count
		}
		if highPriorityCount > 5 {
			fmt.Fprintf(f, "\n*+%d more high priority tasks*\n", highPriorityCount-5)
		}
		fmt.Fprintf(f, "\n")
	}

	// Recent Activity (last 3 days)
	if len(timelineIndex.Entries) > 0 {
		fmt.Fprintf(f, "## 📅 Recent Activity\n\n")
		limit := 3
		if len(timelineIndex.Entries) < limit {
			limit = len(timelineIndex.Entries)
		}
		for i := 0; i < limit; i++ {
			day := timelineIndex.Entries[i]
			fmt.Fprintf(f, "### %s\n", day.Date.Format("Monday, Jan 2"))

			statusCounts := make(map[string]int)
			for _, task := range day.TasksCreated {
				statusCounts[string(task.Status)]++
			}

			// Show task counts by status
			if len(statusCounts) > 0 {
				for _, status := range []string{"NOW", "TODO", "DONE", "DOING", "LATER"} {
					if count := statusCounts[status]; count > 0 {
						fmt.Fprintf(f, "- %d %s task%s\n", count, status, pluralize(count))
					}
				}
			}

			if day.TimeLogged > 0 {
				fmt.Fprintf(f, "- ⏱ %s logged\n", formatDuration(day.TimeLogged))
			}

			// Show key activity bullets
			if len(day.KeyActivity) > 0 {
				fmt.Fprintf(f, "\n")
				for _, activity := range day.KeyActivity {
					fmt.Fprintf(f, "%s\n", activity)
				}
			}
			fmt.Fprintf(f, "\n")
		}
	}

	// Top Projects
	if len(timeTrackingIndex.TopProjects) > 0 || len(taskIndex.ByProject) > 0 {
		fmt.Fprintf(f, "## 📁 Top Projects\n\n")

		// Merge data from both indexes
		projectData := make(map[string]struct {
			TaskCount  int
			TimeLogged time.Duration
		})

		// Count active tasks by project
		for project, tasks := range taskIndex.ByProject {
			activeTasks := 0
			for _, task := range tasks {
				if task.Status != models.StatusDONE {
					activeTasks++
				}
			}
			data := projectData[project]
			data.TaskCount = activeTasks
			projectData[project] = data
		}

		// Add time tracking data
		for _, proj := range timeTrackingIndex.TopProjects {
			data := projectData[proj.Project]
			data.TimeLogged = proj.TimeLogged
			projectData[proj.Project] = data
		}

		// Show top 5 projects (by time logged)
		count := 0
		for _, proj := range timeTrackingIndex.TopProjects {
			if count >= 5 {
				break
			}
			if proj.Project == "No Project" {
				continue
			}
			data := projectData[proj.Project]
			if data.TimeLogged > 0 {
				fmt.Fprintf(f, "- **%s**: %s", proj.Project, formatDuration(data.TimeLogged))
			} else {
				fmt.Fprintf(f, "- **%s**", proj.Project)
			}
			if data.TaskCount > 0 {
				fmt.Fprintf(f, " (%d active task%s)", data.TaskCount, pluralize(data.TaskCount))
			}
			fmt.Fprintf(f, "\n")
			count++
		}
		fmt.Fprintf(f, "\n")
	}

	// Top Missing Pages
	if len(missingPagesIndex.MissingPages) > 0 {
		fmt.Fprintf(f, "## 📝 Pages to Create\n\n")
		fmt.Fprintf(f, "*Pages with 5+ references that don't exist yet*\n\n")

		limit := 5
		if len(missingPagesIndex.MissingPages) < limit {
			limit = len(missingPagesIndex.MissingPages)
		}

		for i := 0; i < limit; i++ {
			page := missingPagesIndex.MissingPages[i]
			fmt.Fprintf(f, "- **%s** (%d refs, %s)\n",
				page.Name, page.ReferenceCount, page.PageType)
		}

		if len(missingPagesIndex.MissingPages) > 5 {
			fmt.Fprintf(f, "\n*+%d more suggested pages*\n", len(missingPagesIndex.MissingPages)-5)
		}
		fmt.Fprintf(f, "\n")
	}

	// Quick Links
	fmt.Fprintf(f, "## 🔗 Detailed Reports\n\n")
	fmt.Fprintf(f, "- [Tasks by Status](./tasks-by-status.md) - All tasks organized by workflow stage\n")
	fmt.Fprintf(f, "- [Tasks by Priority](./tasks-by-priority.md) - High priority tasks requiring attention\n")
	fmt.Fprintf(f, "- [Timeline (Recent)](./timeline-recent.md) - Activity from last 7 days\n")
	fmt.Fprintf(f, "- [Timeline (Full)](./timeline-full.md) - Complete activity history\n")
	fmt.Fprintf(f, "- [Missing Pages](./missing-pages.md) - Suggested pages to create\n")
	fmt.Fprintf(f, "- [Time Tracking](./time-tracking.md) - Time allocation analytics\n")
	fmt.Fprintf(f, "- [Reference Graph](./reference-graph.md) - Page connections and relationships\n")
	fmt.Fprintf(f, "\n")

	return nil
}

// pluralize adds "s" if count != 1
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
