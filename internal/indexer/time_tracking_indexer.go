package indexer

import (
	"sort"
	"time"

	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

// ProjectTime represents time logged for a specific project
type ProjectTime struct {
	Project      string
	TimeLogged   time.Duration
	TaskCount    int
	AvgTimePerTask time.Duration
}

// WeeklyTime represents time logged in a specific week
type WeeklyTime struct {
	WeekStart    time.Time // Monday of the week
	TimeLogged   time.Duration
	TaskCount    int
}

// TimeStatistics provides aggregate statistics
type TimeStatistics struct {
	TotalTasks         int
	TasksWithTracking  int
	TasksWithoutTracking int
	AdoptionRate       float64 // Percentage of tasks with time tracking
	AvgTimePerTask     time.Duration
	MostProductiveWeek WeeklyTime
}

// TimeTrackingIndex aggregates all time tracking data
type TimeTrackingIndex struct {
	TotalTimeLogged time.Duration
	ByProject       map[string]time.Duration
	ByWeek          map[string]time.Duration // Key: "2025-11-04" (Monday of week)
	ByPriority      map[models.Priority]time.Duration
	ByStatus        map[models.TaskStatus]time.Duration
	TopProjects     []ProjectTime
	WeeklySummary   []WeeklyTime
	Statistics      TimeStatistics
}

// BuildTimeTrackingIndex creates a time tracking index from all tasks
func BuildTimeTrackingIndex(tasks []models.Task) *TimeTrackingIndex {
	index := &TimeTrackingIndex{
		ByProject:  make(map[string]time.Duration),
		ByWeek:     make(map[string]time.Duration),
		ByPriority: make(map[models.Priority]time.Duration),
		ByStatus:   make(map[models.TaskStatus]time.Duration),
	}

	projectTaskCounts := make(map[string]int)
	weekTaskCounts := make(map[string]int)
	tasksWithTracking := 0

	for _, task := range tasks {
		index.Statistics.TotalTasks++

		// Sum total time logged
		totalTaskTime := time.Duration(0)
		for _, entry := range task.Logbook {
			totalTaskTime += entry.Duration
		}

		if totalTaskTime > 0 {
			tasksWithTracking++
			index.TotalTimeLogged += totalTaskTime

			// Aggregate by project (use first page reference)
			project := "No Project"
			if len(task.PageRefs) > 0 {
				project = task.PageRefs[0]
			}
			index.ByProject[project] += totalTaskTime
			projectTaskCounts[project]++

			// Aggregate by week (group by Monday)
			for _, entry := range task.Logbook {
				weekStart := getWeekStart(entry.Start)
				weekKey := weekStart.Format("2006-01-02")
				index.ByWeek[weekKey] += entry.Duration
				weekTaskCounts[weekKey]++
			}

			// Aggregate by priority
			index.ByPriority[task.Priority] += totalTaskTime

			// Aggregate by status
			index.ByStatus[task.Status] += totalTaskTime
		}
	}

	// Calculate statistics
	index.Statistics.TasksWithTracking = tasksWithTracking
	index.Statistics.TasksWithoutTracking = index.Statistics.TotalTasks - tasksWithTracking
	if index.Statistics.TotalTasks > 0 {
		index.Statistics.AdoptionRate = float64(tasksWithTracking) / float64(index.Statistics.TotalTasks) * 100
	}
	if tasksWithTracking > 0 {
		index.Statistics.AvgTimePerTask = index.TotalTimeLogged / time.Duration(tasksWithTracking)
	}

	// Build TopProjects sorted by time logged
	for project, timeLogged := range index.ByProject {
		taskCount := projectTaskCounts[project]
		avgTime := time.Duration(0)
		if taskCount > 0 {
			avgTime = timeLogged / time.Duration(taskCount)
		}
		index.TopProjects = append(index.TopProjects, ProjectTime{
			Project:        project,
			TimeLogged:     timeLogged,
			TaskCount:      taskCount,
			AvgTimePerTask: avgTime,
		})
	}
	sort.Slice(index.TopProjects, func(i, j int) bool {
		return index.TopProjects[i].TimeLogged > index.TopProjects[j].TimeLogged
	})

	// Build WeeklySummary sorted chronologically
	for weekKey, timeLogged := range index.ByWeek {
		weekStart, _ := time.Parse("2006-01-02", weekKey)
		taskCount := weekTaskCounts[weekKey]
		weekly := WeeklyTime{
			WeekStart:  weekStart,
			TimeLogged: timeLogged,
			TaskCount:  taskCount,
		}
		index.WeeklySummary = append(index.WeeklySummary, weekly)

		// Track most productive week
		if weekly.TimeLogged > index.Statistics.MostProductiveWeek.TimeLogged {
			index.Statistics.MostProductiveWeek = weekly
		}
	}
	sort.Slice(index.WeeklySummary, func(i, j int) bool {
		return index.WeeklySummary[i].WeekStart.After(index.WeeklySummary[j].WeekStart)
	})

	return index
}

// getWeekStart returns the Monday of the week for the given time
func getWeekStart(t time.Time) time.Time {
	// Get the weekday (0 = Sunday, 1 = Monday, ...)
	weekday := int(t.Weekday())

	// Calculate days to subtract to get to Monday
	daysToMonday := (weekday + 6) % 7 // Convert Sunday=0 to Sunday=6, then calculate
	if weekday == 0 {
		daysToMonday = 6 // Sunday -> go back 6 days
	} else {
		daysToMonday = weekday - 1 // Go back to Monday
	}

	monday := t.AddDate(0, 0, -daysToMonday)
	// Normalize to start of day
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())
}
