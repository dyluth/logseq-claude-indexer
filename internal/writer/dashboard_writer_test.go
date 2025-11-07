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

func TestWriteDashboard_Minimal(t *testing.T) {
	tmpDir := t.TempDir()

	taskIndex := &indexer.TaskIndex{
		TotalTasks: 10,
		ByStatus:   make(map[models.TaskStatus][]models.Task),
		ByPriority: make(map[models.Priority][]models.Task),
		ByProject:  make(map[string][]models.Task),
		Statistics: indexer.TaskStatistics{
			StatusBreakdown: map[models.TaskStatus]int{
				models.StatusDONE: 5,
			},
			CompletionRate:    50.0,
			PriorityBreakdown: make(map[models.Priority]int),
		},
	}

	graphIndex := &indexer.ReferenceGraph{
		Nodes: make(map[string]*indexer.GraphNode),
	}

	timelineIndex := &indexer.TimelineIndex{
		Entries: []indexer.TimelineDay{},
	}

	missingPagesIndex := &indexer.MissingPagesIndex{
		MissingPages: []indexer.MissingPage{},
	}

	timeTrackingIndex := &indexer.TimeTrackingIndex{
		TopProjects: []indexer.ProjectTime{},
		ByProject:   make(map[string]time.Duration),
		ByWeek:      make(map[string]time.Duration),
		ByPriority:  make(map[models.Priority]time.Duration),
		ByStatus:    make(map[models.TaskStatus]time.Duration),
		Statistics: indexer.TimeStatistics{
			TotalTasks:        10,
			TasksWithTracking: 0,
			AdoptionRate:      0,
		},
	}

	err := WriteDashboard(taskIndex, graphIndex, timelineIndex, missingPagesIndex, timeTrackingIndex, tmpDir)
	if err != nil {
		t.Fatalf("WriteDashboard failed: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "dashboard.md")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read dashboard file: %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "# Knowledge Dashboard") {
		t.Error("Expected dashboard title")
	}
	if !strings.Contains(output, "Total Tasks**: 10") {
		t.Error("Expected total tasks count")
	}
	if !strings.Contains(output, "50.0%") {
		t.Error("Expected completion rate")
	}
}

func TestWriteDashboard_WithAllData(t *testing.T) {
	tmpDir := t.TempDir()

	// Create task index with high priority tasks
	taskIndex := &indexer.TaskIndex{
		TotalTasks: 50,
		ByStatus:   make(map[models.TaskStatus][]models.Task),
		ByPriority: map[models.Priority][]models.Task{
			models.PriorityHigh: {
				{
					Status:      models.StatusNOW,
					Priority:    models.PriorityHigh,
					Description: "Fix critical bug in authentication",
					SourceFile:  "journals/2025_11_06.md",
					LineNumber:  10,
				},
				{
					Status:      models.StatusNOW,
					Priority:    models.PriorityHigh,
					Description: "Deploy hotfix to production",
					SourceFile:  "journals/2025_11_06.md",
					LineNumber:  15,
				},
			},
		},
		ByProject: map[string][]models.Task{
			"Project Alpha": {
				{Status: models.StatusNOW},
				{Status: models.StatusNOW},
				{Status: models.StatusNOW},
				{Status: models.StatusNOW},
				{Status: models.StatusNOW},
			},
			"Project Beta": {
				{Status: models.StatusNOW},
				{Status: models.StatusNOW},
				{Status: models.StatusNOW},
			},
		},
		Statistics: indexer.TaskStatistics{
			StatusBreakdown: map[models.TaskStatus]int{
				models.StatusDONE: 30,
			},
			PriorityBreakdown: map[models.Priority]int{
				models.PriorityHigh: 2,
			},
			CompletionRate: 60.0,
		},
	}

	graphIndex := &indexer.ReferenceGraph{
		Nodes: map[string]*indexer.GraphNode{
			"Page 1": {InboundRefs: []string{"Page 2", "Page 3"}},
			"Page 2": {InboundRefs: []string{"Page 3"}},
			"Page 3": {InboundRefs: []string{}},
		},
	}

	timelineIndex := &indexer.TimelineIndex{
		Entries: []indexer.TimelineDay{
			{
				Date: time.Date(2025, 11, 6, 0, 0, 0, 0, time.UTC),
				TasksCreated: []models.Task{
					{Status: models.StatusNOW},
					{Status: models.StatusNOW},
					{Status: models.StatusDONE},
				},
				TimeLogged: 5 * time.Hour,
				KeyActivity: []string{
					"- **[NOW]** Fix critical bug",
					"🔥 Deploy to production",
				},
			},
			{
				Date: time.Date(2025, 11, 5, 0, 0, 0, 0, time.UTC),
				TasksCreated: []models.Task{
					{Status: models.StatusTODO},
				},
				TimeLogged: 2 * time.Hour,
			},
		},
	}

	missingPagesIndex := &indexer.MissingPagesIndex{
		MissingPages: []indexer.MissingPage{
			{Name: "Important Feature", ReferenceCount: 10, PageType: "project"},
			{Name: "John Doe - Engineer", ReferenceCount: 8, PageType: "person"},
			{Name: "Q4 Planning", ReferenceCount: 7, PageType: "concept"},
		},
	}

	timeTrackingIndex := &indexer.TimeTrackingIndex{
		TotalTimeLogged: 100 * time.Hour,
		TopProjects: []indexer.ProjectTime{
			{Project: "Project Alpha", TimeLogged: 50 * time.Hour, TaskCount: 5},
			{Project: "Project Beta", TimeLogged: 30 * time.Hour, TaskCount: 3},
			{Project: "No Project", TimeLogged: 20 * time.Hour, TaskCount: 2},
		},
		ByProject:  make(map[string]time.Duration),
		ByWeek:     make(map[string]time.Duration),
		ByPriority: make(map[models.Priority]time.Duration),
		ByStatus:   make(map[models.TaskStatus]time.Duration),
		Statistics: indexer.TimeStatistics{
			TotalTasks:        50,
			TasksWithTracking: 25,
			AdoptionRate:      50.0,
		},
	}

	err := WriteDashboard(taskIndex, graphIndex, timelineIndex, missingPagesIndex, timeTrackingIndex, tmpDir)
	if err != nil {
		t.Fatalf("WriteDashboard failed: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "dashboard.md")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read dashboard file: %v", err)
	}

	output := string(content)

	// Check Quick Stats
	if !strings.Contains(output, "Total Tasks**: 50") {
		t.Error("Expected 50 total tasks")
	}
	if !strings.Contains(output, "60.0%") {
		t.Error("Expected 60% completion rate")
	}
	if !strings.Contains(output, "50.0% adoption") {
		t.Error("Expected 50% time tracking adoption")
	}
	if !strings.Contains(output, "3 pages") {
		t.Error("Expected 3 pages in knowledge graph")
	}

	// Check Current Priorities
	if !strings.Contains(output, "Current Priorities") {
		t.Error("Expected Current Priorities section")
	}
	if !strings.Contains(output, "Fix critical bug") {
		t.Error("Expected high priority task")
	}

	// Check Recent Activity
	if !strings.Contains(output, "Recent Activity") {
		t.Error("Expected Recent Activity section")
	}
	if !strings.Contains(output, "2 NOW") {
		t.Error("Expected NOW task count")
	}
	if !strings.Contains(output, "5h") {
		t.Error("Expected time logged")
	}
	if !strings.Contains(output, "Deploy to production") {
		t.Error("Expected key activity")
	}

	// Check Top Projects
	if !strings.Contains(output, "Top Projects") {
		t.Error("Expected Top Projects section")
	}
	if !strings.Contains(output, "Project Alpha") {
		t.Error("Expected Project Alpha")
	}
	if !strings.Contains(output, "50h") {
		t.Error("Expected project time")
	}

	// Check Missing Pages
	if !strings.Contains(output, "Pages to Create") {
		t.Error("Expected Pages to Create section")
	}
	if !strings.Contains(output, "Important Feature") {
		t.Error("Expected missing page")
	}
	if !strings.Contains(output, "10 refs") {
		t.Error("Expected reference count")
	}

	// Check Quick Links
	if !strings.Contains(output, "Detailed Reports") {
		t.Error("Expected Detailed Reports section")
	}
	if !strings.Contains(output, "tasks-by-status.md") {
		t.Error("Expected link to tasks by status")
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		count    int
		expected string
	}{
		{0, "s"},
		{1, ""},
		{2, "s"},
		{100, "s"},
	}

	for _, tt := range tests {
		result := pluralize(tt.count)
		if result != tt.expected {
			t.Errorf("pluralize(%d) = %q, want %q", tt.count, result, tt.expected)
		}
	}
}
