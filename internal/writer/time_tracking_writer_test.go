package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dyluth/logseq-claude-indexer/internal/indexer"
	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

func TestWriteTimeTracking_EmptyIndex(t *testing.T) {
	tmpDir := t.TempDir()
	index := &indexer.TimeTrackingIndex{
		ByProject:  make(map[string]time.Duration),
		ByWeek:     make(map[string]time.Duration),
		ByPriority: make(map[models.Priority]time.Duration),
		ByStatus:   make(map[models.TaskStatus]time.Duration),
	}

	err := WriteTimeTracking(index, tmpDir)
	if err != nil {
		t.Fatalf("WriteTimeTracking failed: %v", err)
	}

	// Verify file was created
	outputPath := filepath.Join(tmpDir, "time-tracking.md")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "# Time Tracking Analytics") {
		t.Error("Expected title in output")
	}
	if !strings.Contains(output, "0 / 0") {
		t.Error("Expected 0 / 0 tasks tracked")
	}
}

func TestWriteTimeTracking_WithData(t *testing.T) {
	tmpDir := t.TempDir()

	index := &indexer.TimeTrackingIndex{
		TotalTimeLogged: 20 * time.Hour,
		TopProjects: []indexer.ProjectTime{
			{
				Project:        "Project Alpha",
				TimeLogged:     12 * time.Hour,
				TaskCount:      3,
				AvgTimePerTask: 4 * time.Hour,
			},
			{
				Project:        "Project Beta",
				TimeLogged:     8 * time.Hour,
				TaskCount:      2,
				AvgTimePerTask: 4 * time.Hour,
			},
		},
		WeeklySummary: []indexer.WeeklyTime{
			{
				WeekStart:  time.Date(2025, 11, 3, 0, 0, 0, 0, time.UTC),
				TimeLogged: 15 * time.Hour,
				TaskCount:  4,
			},
			{
				WeekStart:  time.Date(2025, 10, 27, 0, 0, 0, 0, time.UTC),
				TimeLogged: 5 * time.Hour,
				TaskCount:  1,
			},
		},
		ByPriority: map[models.Priority]time.Duration{
			models.PriorityHigh:   12 * time.Hour,
			models.PriorityMedium: 8 * time.Hour,
		},
		ByStatus: map[models.TaskStatus]time.Duration{
			models.StatusDONE: 18 * time.Hour,
			models.StatusNOW:  2 * time.Hour,
		},
		ByProject: make(map[string]time.Duration),
		ByWeek:    make(map[string]time.Duration),
		Statistics: indexer.TimeStatistics{
			TotalTasks:        10,
			TasksWithTracking: 5,
			AdoptionRate:      50.0,
			AvgTimePerTask:    4 * time.Hour,
			MostProductiveWeek: indexer.WeeklyTime{
				WeekStart:  time.Date(2025, 11, 3, 0, 0, 0, 0, time.UTC),
				TimeLogged: 15 * time.Hour,
			},
		},
	}

	err := WriteTimeTracking(index, tmpDir)
	if err != nil {
		t.Fatalf("WriteTimeTracking failed: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "time-tracking.md")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	output := string(content)

	// Check summary section
	if !strings.Contains(output, "20h") {
		t.Error("Expected total time of 20h")
	}
	if !strings.Contains(output, "5 / 10") {
		t.Error("Expected 5 / 10 tasks tracked")
	}
	if !strings.Contains(output, "50.0%") {
		t.Error("Expected 50% adoption rate")
	}
	if !strings.Contains(output, "4h") {
		t.Error("Expected average time per task")
	}

	// Check top projects
	if !strings.Contains(output, "Project Alpha") {
		t.Error("Expected Project Alpha in output")
	}
	if !strings.Contains(output, "12h") {
		t.Error("Expected 12h for Project Alpha")
	}
	if !strings.Contains(output, "Project Beta") {
		t.Error("Expected Project Beta in output")
	}

	// Check weekly breakdown
	if !strings.Contains(output, "2025-11-03") {
		t.Error("Expected week of 2025-11-03")
	}
	if !strings.Contains(output, "15h") {
		t.Error("Expected 15h for most productive week")
	}

	// Check by priority
	if !strings.Contains(output, "[#A]") {
		t.Error("Expected Priority A")
	}
	if !strings.Contains(output, "[#B]") {
		t.Error("Expected Priority B")
	}

	// Check by status
	if !strings.Contains(output, "DONE") {
		t.Error("Expected DONE status")
	}
	if !strings.Contains(output, "NOW") {
		t.Error("Expected NOW status")
	}
}

func TestWriteTimeTracking_ManyProjects(t *testing.T) {
	tmpDir := t.TempDir()

	// Create 15 projects, but only top 10 should be shown
	projects := []indexer.ProjectTime{}
	for i := 15; i > 0; i-- {
		projects = append(projects, indexer.ProjectTime{
			Project:        fmt.Sprintf("Project %d", i),
			TimeLogged:     time.Duration(i) * time.Hour,
			TaskCount:      1,
			AvgTimePerTask: time.Duration(i) * time.Hour,
		})
	}

	index := &indexer.TimeTrackingIndex{
		TopProjects: projects,
		ByProject:   make(map[string]time.Duration),
		ByWeek:      make(map[string]time.Duration),
		ByPriority:  make(map[models.Priority]time.Duration),
		ByStatus:    make(map[models.TaskStatus]time.Duration),
	}

	err := WriteTimeTracking(index, tmpDir)
	if err != nil {
		t.Fatalf("WriteTimeTracking failed: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "time-tracking.md")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	output := string(content)

	// Should show Project 15 through Project 6 (top 10)
	if !strings.Contains(output, "Project 15") {
		t.Error("Expected Project 15 (highest)")
	}
	if !strings.Contains(output, "Project 6") {
		t.Error("Expected Project 6 (10th highest)")
	}
	// Should NOT show Project 5 or lower
	if strings.Contains(output, "Project 5") {
		t.Error("Should not show Project 5 (below top 10)")
	}
}

func TestWriteTimeTracking_ManyWeeks(t *testing.T) {
	tmpDir := t.TempDir()

	// Create 12 weeks, but only last 8 should be shown
	weeks := []indexer.WeeklyTime{}
	baseDate := time.Date(2025, 11, 3, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 12; i++ {
		weeks = append(weeks, indexer.WeeklyTime{
			WeekStart:  baseDate.AddDate(0, 0, -7*i), // Go backwards in time
			TimeLogged: time.Duration(i+1) * time.Hour,
			TaskCount:  1,
		})
	}

	index := &indexer.TimeTrackingIndex{
		WeeklySummary: weeks,
		ByProject:     make(map[string]time.Duration),
		ByWeek:        make(map[string]time.Duration),
		ByPriority:    make(map[models.Priority]time.Duration),
		ByStatus:      make(map[models.TaskStatus]time.Duration),
	}

	err := WriteTimeTracking(index, tmpDir)
	if err != nil {
		t.Fatalf("WriteTimeTracking failed: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "time-tracking.md")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	output := string(content)

	// Should show last 8 weeks with a note about total
	if !strings.Contains(output, "Showing last 8 weeks of 12 total") {
		t.Error("Expected note about showing 8 of 12 weeks")
	}
}
