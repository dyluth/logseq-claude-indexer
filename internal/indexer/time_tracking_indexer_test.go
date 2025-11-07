package indexer

import (
	"testing"
	"time"

	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

func TestBuildTimeTrackingIndex_EmptyTasks(t *testing.T) {
	index := BuildTimeTrackingIndex([]models.Task{})

	if index.TotalTimeLogged != 0 {
		t.Errorf("Expected 0 total time, got %v", index.TotalTimeLogged)
	}
	if len(index.TopProjects) != 0 {
		t.Errorf("Expected 0 projects, got %d", len(index.TopProjects))
	}
	if index.Statistics.TotalTasks != 0 {
		t.Errorf("Expected 0 total tasks, got %d", index.Statistics.TotalTasks)
	}
}

func TestBuildTimeTrackingIndex_TasksWithoutTracking(t *testing.T) {
	tasks := []models.Task{
		{
			Description: "Task without tracking",
			Status:      models.StatusNOW,
			Priority:    models.PriorityHigh,
			PageRefs:    []string{"Project A"},
		},
		{
			Description: "Another task without tracking",
			Status:      models.StatusTODO,
			Priority:    models.PriorityMedium,
			PageRefs:    []string{"Project B"},
		},
	}

	index := BuildTimeTrackingIndex(tasks)

	if index.TotalTimeLogged != 0 {
		t.Errorf("Expected 0 total time, got %v", index.TotalTimeLogged)
	}
	if index.Statistics.TotalTasks != 2 {
		t.Errorf("Expected 2 total tasks, got %d", index.Statistics.TotalTasks)
	}
	if index.Statistics.TasksWithTracking != 0 {
		t.Errorf("Expected 0 tasks with tracking, got %d", index.Statistics.TasksWithTracking)
	}
	if index.Statistics.AdoptionRate != 0 {
		t.Errorf("Expected 0%% adoption, got %.1f%%", index.Statistics.AdoptionRate)
	}
}

func TestBuildTimeTrackingIndex_SingleTask(t *testing.T) {
	start := time.Date(2025, 11, 6, 9, 0, 0, 0, time.UTC)
	end := time.Date(2025, 11, 6, 11, 30, 0, 0, time.UTC)

	tasks := []models.Task{
		{
			Description: "Test task",
			Status:      models.StatusDONE,
			Priority:    models.PriorityHigh,
			PageRefs:    []string{"Project Alpha"},
			Logbook: []models.LogbookEntry{
				{Start: start, End: end, Duration: 2*time.Hour + 30*time.Minute},
			},
		},
	}

	index := BuildTimeTrackingIndex(tasks)

	expectedDuration := 2*time.Hour + 30*time.Minute

	if index.TotalTimeLogged != expectedDuration {
		t.Errorf("Expected %v total time, got %v", expectedDuration, index.TotalTimeLogged)
	}
	if index.Statistics.TotalTasks != 1 {
		t.Errorf("Expected 1 total task, got %d", index.Statistics.TotalTasks)
	}
	if index.Statistics.TasksWithTracking != 1 {
		t.Errorf("Expected 1 task with tracking, got %d", index.Statistics.TasksWithTracking)
	}
	if index.Statistics.AdoptionRate != 100.0 {
		t.Errorf("Expected 100%% adoption, got %.1f%%", index.Statistics.AdoptionRate)
	}
	if index.ByProject["Project Alpha"] != expectedDuration {
		t.Errorf("Expected %v for Project Alpha, got %v", expectedDuration, index.ByProject["Project Alpha"])
	}
	if index.ByStatus[models.StatusDONE] != expectedDuration {
		t.Errorf("Expected %v for DONE status, got %v", expectedDuration, index.ByStatus[models.StatusDONE])
	}
	if index.ByPriority[models.PriorityHigh] != expectedDuration {
		t.Errorf("Expected %v for Priority High, got %v", expectedDuration, index.ByPriority[models.PriorityHigh])
	}
}

func TestBuildTimeTrackingIndex_MultipleProjects(t *testing.T) {
	tasks := []models.Task{
		{
			Description: "Task 1",
			Status:      models.StatusDONE,
			PageRefs:    []string{"Project A"},
			Logbook: []models.LogbookEntry{
				{Duration: 5 * time.Hour},
			},
		},
		{
			Description: "Task 2",
			Status:      models.StatusDONE,
			PageRefs:    []string{"Project B"},
			Logbook: []models.LogbookEntry{
				{Duration: 3 * time.Hour},
			},
		},
		{
			Description: "Task 3",
			Status:      models.StatusDONE,
			PageRefs:    []string{"Project A"},
			Logbook: []models.LogbookEntry{
				{Duration: 2 * time.Hour},
			},
		},
	}

	index := BuildTimeTrackingIndex(tasks)

	if index.TotalTimeLogged != 10*time.Hour {
		t.Errorf("Expected 10h total time, got %v", index.TotalTimeLogged)
	}

	// Project A should have 7 hours (5 + 2)
	if index.ByProject["Project A"] != 7*time.Hour {
		t.Errorf("Expected 7h for Project A, got %v", index.ByProject["Project A"])
	}

	// Project B should have 3 hours
	if index.ByProject["Project B"] != 3*time.Hour {
		t.Errorf("Expected 3h for Project B, got %v", index.ByProject["Project B"])
	}

	// TopProjects should be sorted by time (Project A first)
	if len(index.TopProjects) != 2 {
		t.Fatalf("Expected 2 projects, got %d", len(index.TopProjects))
	}
	if index.TopProjects[0].Project != "Project A" {
		t.Errorf("Expected Project A first, got %s", index.TopProjects[0].Project)
	}
	if index.TopProjects[0].TimeLogged != 7*time.Hour {
		t.Errorf("Expected 7h for top project, got %v", index.TopProjects[0].TimeLogged)
	}
	if index.TopProjects[0].TaskCount != 2 {
		t.Errorf("Expected 2 tasks for Project A, got %d", index.TopProjects[0].TaskCount)
	}
}

func TestBuildTimeTrackingIndex_WeeklyAggregation(t *testing.T) {
	// Week of Nov 3-9, 2025 (Monday Nov 3)
	// Nov 4 is Tuesday, Nov 6 is Wednesday, Nov 9 is Sunday
	tuesday := time.Date(2025, 11, 4, 10, 0, 0, 0, time.UTC)
	wednesday := time.Date(2025, 11, 6, 14, 0, 0, 0, time.UTC)
	sunday := time.Date(2025, 11, 9, 9, 0, 0, 0, time.UTC)

	// Different week (Oct 27 - Nov 2, Monday Oct 27)
	previousMonday := time.Date(2025, 10, 27, 10, 0, 0, 0, time.UTC)

	// Next week (Nov 10-16, Monday Nov 10)
	nextMonday := time.Date(2025, 11, 10, 11, 0, 0, 0, time.UTC)

	tasks := []models.Task{
		{
			Description: "Task on Tuesday (week of Nov 3)",
			Status:      models.StatusDONE,
			PageRefs:    []string{"Project A"},
			Logbook: []models.LogbookEntry{
				{Start: tuesday, Duration: 4 * time.Hour},
			},
		},
		{
			Description: "Task on Wednesday (week of Nov 3)",
			Status:      models.StatusDONE,
			PageRefs:    []string{"Project A"},
			Logbook: []models.LogbookEntry{
				{Start: wednesday, Duration: 3 * time.Hour},
			},
		},
		{
			Description: "Task on Sunday (week of Nov 3)",
			Status:      models.StatusDONE,
			PageRefs:    []string{"Project B"},
			Logbook: []models.LogbookEntry{
				{Start: sunday, Duration: 2 * time.Hour},
			},
		},
		{
			Description: "Task in previous week (Oct 27)",
			Status:      models.StatusDONE,
			PageRefs:    []string{"Project C"},
			Logbook: []models.LogbookEntry{
				{Start: previousMonday, Duration: 5 * time.Hour},
			},
		},
		{
			Description: "Task in next week (Nov 10)",
			Status:      models.StatusDONE,
			PageRefs:    []string{"Project D"},
			Logbook: []models.LogbookEntry{
				{Start: nextMonday, Duration: 1 * time.Hour},
			},
		},
	}

	index := BuildTimeTrackingIndex(tasks)

	// Should have 3 weeks
	if len(index.WeeklySummary) != 3 {
		t.Errorf("Expected 3 weeks, got %d", len(index.WeeklySummary))
	}

	// Week of Nov 3 should have 9 hours total (4 + 3 + 2)
	week1Key := "2025-11-03" // Monday Nov 3
	if index.ByWeek[week1Key] != 9*time.Hour {
		t.Errorf("Expected 9h for week of Nov 3, got %v", index.ByWeek[week1Key])
	}

	// Week of Oct 27 should have 5 hours
	week2Key := "2025-10-27" // Monday Oct 27
	if index.ByWeek[week2Key] != 5*time.Hour {
		t.Errorf("Expected 5h for week of Oct 27, got %v", index.ByWeek[week2Key])
	}

	// Week of Nov 10 should have 1 hour
	week3Key := "2025-11-10" // Monday Nov 10
	if index.ByWeek[week3Key] != 1*time.Hour {
		t.Errorf("Expected 1h for week of Nov 10, got %v", index.ByWeek[week3Key])
	}

	// Most productive week should be Nov 3 week
	if index.Statistics.MostProductiveWeek.TimeLogged != 9*time.Hour {
		t.Errorf("Expected most productive week to have 9h, got %v",
			index.Statistics.MostProductiveWeek.TimeLogged)
	}
}

func TestBuildTimeTrackingIndex_ByStatusAndPriority(t *testing.T) {
	tasks := []models.Task{
		{
			Description: "High priority done task",
			Status:      models.StatusDONE,
			Priority:    models.PriorityHigh,
			PageRefs:    []string{"Project A"},
			Logbook: []models.LogbookEntry{{Duration: 10 * time.Hour}},
		},
		{
			Description: "Medium priority in progress",
			Status:      models.StatusNOW,
			Priority:    models.PriorityMedium,
			PageRefs:    []string{"Project A"},
			Logbook: []models.LogbookEntry{{Duration: 5 * time.Hour}},
		},
		{
			Description: "Another high priority done",
			Status:      models.StatusDONE,
			Priority:    models.PriorityHigh,
			PageRefs:    []string{"Project B"},
			Logbook: []models.LogbookEntry{{Duration: 3 * time.Hour}},
		},
		{
			Description: "No priority TODO",
			Status:      models.StatusTODO,
			Priority:    models.PriorityNone,
			PageRefs:    []string{"Project C"},
			Logbook: []models.LogbookEntry{{Duration: 2 * time.Hour}},
		},
	}

	index := BuildTimeTrackingIndex(tasks)

	// Check by status
	if index.ByStatus[models.StatusDONE] != 13*time.Hour {
		t.Errorf("Expected 13h for DONE status, got %v", index.ByStatus[models.StatusDONE])
	}
	if index.ByStatus[models.StatusNOW] != 5*time.Hour {
		t.Errorf("Expected 5h for NOW status, got %v", index.ByStatus[models.StatusNOW])
	}
	if index.ByStatus[models.StatusTODO] != 2*time.Hour {
		t.Errorf("Expected 2h for TODO status, got %v", index.ByStatus[models.StatusTODO])
	}

	// Check by priority
	if index.ByPriority[models.PriorityHigh] != 13*time.Hour {
		t.Errorf("Expected 13h for Priority High, got %v", index.ByPriority[models.PriorityHigh])
	}
	if index.ByPriority[models.PriorityMedium] != 5*time.Hour {
		t.Errorf("Expected 5h for Priority Medium, got %v", index.ByPriority[models.PriorityMedium])
	}
	if index.ByPriority[models.PriorityNone] != 2*time.Hour {
		t.Errorf("Expected 2h for no priority, got %v", index.ByPriority[models.PriorityNone])
	}
}

func TestBuildTimeTrackingIndex_MixedTracking(t *testing.T) {
	tasks := []models.Task{
		{
			Description: "With tracking",
			Status:      models.StatusDONE,
			PageRefs:    []string{"Project A"},
			Logbook: []models.LogbookEntry{{Duration: 5 * time.Hour}},
		},
		{
			Description: "Without tracking",
			Status:      models.StatusTODO,
			PageRefs:    []string{"Project A"},
		},
		{
			Description: "With tracking",
			Status:      models.StatusDONE,
			PageRefs:    []string{"Project B"},
			Logbook: []models.LogbookEntry{{Duration: 3 * time.Hour}},
		},
	}

	index := BuildTimeTrackingIndex(tasks)

	if index.Statistics.TotalTasks != 3 {
		t.Errorf("Expected 3 total tasks, got %d", index.Statistics.TotalTasks)
	}
	if index.Statistics.TasksWithTracking != 2 {
		t.Errorf("Expected 2 tasks with tracking, got %d", index.Statistics.TasksWithTracking)
	}
	if index.Statistics.TasksWithoutTracking != 1 {
		t.Errorf("Expected 1 task without tracking, got %d", index.Statistics.TasksWithoutTracking)
	}

	expectedAdoption := (2.0 / 3.0) * 100.0
	if index.Statistics.AdoptionRate < expectedAdoption-0.1 || index.Statistics.AdoptionRate > expectedAdoption+0.1 {
		t.Errorf("Expected %.1f%% adoption, got %.1f%%", expectedAdoption, index.Statistics.AdoptionRate)
	}

	expectedAvg := 4 * time.Hour // 8h / 2 tasks
	if index.Statistics.AvgTimePerTask != expectedAvg {
		t.Errorf("Expected %v avg time, got %v", expectedAvg, index.Statistics.AvgTimePerTask)
	}
}

func TestGetWeekStart(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "Monday",
			input:    time.Date(2025, 11, 3, 15, 30, 0, 0, time.UTC),
			expected: time.Date(2025, 11, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Wednesday",
			input:    time.Date(2025, 11, 5, 10, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 11, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Sunday",
			input:    time.Date(2025, 11, 9, 20, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 11, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Next Monday",
			input:    time.Date(2025, 11, 10, 9, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 11, 10, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getWeekStart(tt.input)
			if !result.Equal(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
