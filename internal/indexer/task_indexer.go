package indexer

import (
	"sort"
	"time"

	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

// TaskIndex organizes tasks by status, priority, and project for easy querying
type TaskIndex struct {
	GeneratedAt time.Time
	TotalTasks  int
	ByStatus    map[models.TaskStatus][]models.Task
	ByPriority  map[models.Priority][]models.Task // Grouped by priority level
	ByProject   map[string][]models.Task           // Keyed by first page reference
	Recent      []models.Task                      // Last 30 days
	Statistics  TaskStatistics                     // Summary statistics
}

// TaskStatistics provides summary statistics for tasks
type TaskStatistics struct {
	CompletionRate    float64                        // Percentage of DONE tasks
	PriorityBreakdown map[models.Priority]int        // Count by priority
	StatusBreakdown   map[models.TaskStatus]int      // Count by status
	WithTimeTracking  int                            // Count of tasks with LOGBOOK
	TotalTimeLogged   time.Duration                  // Sum of all LOGBOOK time
	TrackingAdoption  float64                        // Percentage with time tracking
}

// BuildTaskIndex creates a TaskIndex from a list of tasks
func BuildTaskIndex(tasks []models.Task) *TaskIndex {
	index := &TaskIndex{
		GeneratedAt: time.Now(),
		TotalTasks:  len(tasks),
		ByStatus:    make(map[models.TaskStatus][]models.Task),
		ByPriority:  make(map[models.Priority][]models.Task),
		ByProject:   make(map[string][]models.Task),
		Statistics: TaskStatistics{
			PriorityBreakdown: make(map[models.Priority]int),
			StatusBreakdown:   make(map[models.TaskStatus]int),
		},
	}

	// Initialize status maps with empty slices
	for _, status := range []models.TaskStatus{
		models.StatusNOW,
		models.StatusLATER,
		models.StatusTODO,
		models.StatusDOING,
		models.StatusDONE,
	} {
		index.ByStatus[status] = []models.Task{}
	}

	// Initialize priority maps with empty slices
	for _, priority := range []models.Priority{
		models.PriorityHigh,
		models.PriorityMedium,
		models.PriorityLow,
		models.PriorityNone,
	} {
		index.ByPriority[priority] = []models.Task{}
	}

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	for _, task := range tasks {
		// Group by status
		index.ByStatus[task.Status] = append(index.ByStatus[task.Status], task)
		index.Statistics.StatusBreakdown[task.Status]++

		// Group by priority
		index.ByPriority[task.Priority] = append(index.ByPriority[task.Priority], task)
		index.Statistics.PriorityBreakdown[task.Priority]++

		// Group by project (first page reference)
		if len(task.PageRefs) > 0 {
			project := task.PageRefs[0]
			index.ByProject[project] = append(index.ByProject[project], task)
		}

		// Add to recent if has recent logbook entries
		if hasRecentActivity(task, thirtyDaysAgo) {
			index.Recent = append(index.Recent, task)
		}

		// Track time logging statistics
		if len(task.Logbook) > 0 {
			index.Statistics.WithTimeTracking++
			index.Statistics.TotalTimeLogged += task.TotalDuration()
		}
	}

	// Sort recent tasks by most recent activity
	sort.Slice(index.Recent, func(i, j int) bool {
		return getMostRecentTime(index.Recent[i]).After(getMostRecentTime(index.Recent[j]))
	})

	// Calculate statistics
	if len(tasks) > 0 {
		index.Statistics.CompletionRate = float64(index.Statistics.StatusBreakdown[models.StatusDONE]) / float64(len(tasks)) * 100
		index.Statistics.TrackingAdoption = float64(index.Statistics.WithTimeTracking) / float64(len(tasks)) * 100
	}

	return index
}

// hasRecentActivity checks if a task has logbook entries within the time window
func hasRecentActivity(task models.Task, since time.Time) bool {
	for _, entry := range task.Logbook {
		if entry.Start.After(since) || entry.End.After(since) {
			return true
		}
	}
	return false
}

// getMostRecentTime returns the most recent timestamp from a task's logbook
func getMostRecentTime(task models.Task) time.Time {
	if len(task.Logbook) == 0 {
		return time.Time{} // Zero time for tasks without logbook
	}

	mostRecent := task.Logbook[0].End
	for _, entry := range task.Logbook[1:] {
		if entry.End.After(mostRecent) {
			mostRecent = entry.End
		}
	}
	return mostRecent
}

// GetProjectSummary returns a summary of tasks for a project
type ProjectSummary struct {
	ProjectName string
	TotalTasks  int
	ByStatus    map[models.TaskStatus]int
}

// GetProjectSummaries generates summaries for all projects in the index
func (ti *TaskIndex) GetProjectSummaries() []ProjectSummary {
	var summaries []ProjectSummary

	for projectName, tasks := range ti.ByProject {
		summary := ProjectSummary{
			ProjectName: projectName,
			TotalTasks:  len(tasks),
			ByStatus:    make(map[models.TaskStatus]int),
		}

		for _, task := range tasks {
			summary.ByStatus[task.Status]++
		}

		summaries = append(summaries, summary)
	}

	// Sort by total tasks descending
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].TotalTasks > summaries[j].TotalTasks
	})

	return summaries
}
