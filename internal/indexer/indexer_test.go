package indexer

import (
	"testing"
	"time"

	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

func TestBuildTaskIndex(t *testing.T) {
	tasks := []models.Task{
		{
			Status:      models.StatusNOW,
			Description: "Task 1",
			PageRefs:    []string{"Project A"},
			SourceFile:  "test.md",
			LineNumber:  1,
		},
		{
			Status:      models.StatusNOW,
			Description: "Task 2",
			PageRefs:    []string{"Project B"},
			SourceFile:  "test.md",
			LineNumber:  2,
		},
		{
			Status:      models.StatusLATER,
			Description: "Task 3",
			PageRefs:    []string{"Project A"},
			SourceFile:  "test.md",
			LineNumber:  3,
		},
		{
			Status:      models.StatusDONE,
			Description: "Task 4",
			PageRefs:    []string{"Project C"},
			SourceFile:  "test.md",
			LineNumber:  4,
		},
	}

	index := BuildTaskIndex(tasks)

	// Check total
	if index.TotalTasks != 4 {
		t.Errorf("Expected 4 total tasks, got %d", index.TotalTasks)
	}

	// Check by status
	if len(index.ByStatus[models.StatusNOW]) != 2 {
		t.Errorf("Expected 2 NOW tasks, got %d", len(index.ByStatus[models.StatusNOW]))
	}
	if len(index.ByStatus[models.StatusLATER]) != 1 {
		t.Errorf("Expected 1 LATER task, got %d", len(index.ByStatus[models.StatusLATER]))
	}
	if len(index.ByStatus[models.StatusDONE]) != 1 {
		t.Errorf("Expected 1 DONE task, got %d", len(index.ByStatus[models.StatusDONE]))
	}
	if len(index.ByStatus[models.StatusTODO]) != 0 {
		t.Errorf("Expected 0 TODO tasks, got %d", len(index.ByStatus[models.StatusTODO]))
	}

	// Check by project
	if len(index.ByProject["Project A"]) != 2 {
		t.Errorf("Expected 2 tasks for Project A, got %d", len(index.ByProject["Project A"]))
	}
	if len(index.ByProject["Project B"]) != 1 {
		t.Errorf("Expected 1 task for Project B, got %d", len(index.ByProject["Project B"]))
	}
	if len(index.ByProject["Project C"]) != 1 {
		t.Errorf("Expected 1 task for Project C, got %d", len(index.ByProject["Project C"]))
	}
}

func TestBuildTaskIndex_RecentTasks(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	twoMonthsAgo := now.AddDate(0, -2, 0)

	tasks := []models.Task{
		{
			Status:      models.StatusNOW,
			Description: "Recent task",
			SourceFile:  "test.md",
			LineNumber:  1,
			Logbook: []models.LogbookEntry{
				{
					Start:    yesterday,
					End:      now,
					Duration: time.Hour,
				},
			},
		},
		{
			Status:      models.StatusDONE,
			Description: "Old task",
			SourceFile:  "test.md",
			LineNumber:  2,
			Logbook: []models.LogbookEntry{
				{
					Start:    twoMonthsAgo.Add(-time.Hour),
					End:      twoMonthsAgo,
					Duration: time.Hour,
				},
			},
		},
		{
			Status:      models.StatusLATER,
			Description: "Task without logbook",
			SourceFile:  "test.md",
			LineNumber:  3,
		},
	}

	index := BuildTaskIndex(tasks)

	// Only the recent task should appear in Recent
	if len(index.Recent) != 1 {
		t.Errorf("Expected 1 recent task, got %d", len(index.Recent))
	}

	if len(index.Recent) > 0 && index.Recent[0].Description != "Recent task" {
		t.Errorf("Expected 'Recent task' in recent, got '%s'", index.Recent[0].Description)
	}
}

func TestGetProjectSummaries(t *testing.T) {
	tasks := []models.Task{
		{Status: models.StatusNOW, PageRefs: []string{"Project A"}},
		{Status: models.StatusLATER, PageRefs: []string{"Project A"}},
		{Status: models.StatusDONE, PageRefs: []string{"Project A"}},
		{Status: models.StatusNOW, PageRefs: []string{"Project B"}},
		{Status: models.StatusTODO, PageRefs: []string{"Project B"}},
	}

	index := BuildTaskIndex(tasks)
	summaries := index.GetProjectSummaries()

	if len(summaries) != 2 {
		t.Fatalf("Expected 2 project summaries, got %d", len(summaries))
	}

	// Should be sorted by total tasks (Project A has 3, B has 2)
	if summaries[0].ProjectName != "Project A" {
		t.Errorf("Expected Project A first, got %s", summaries[0].ProjectName)
	}

	if summaries[0].TotalTasks != 3 {
		t.Errorf("Expected 3 total tasks for Project A, got %d", summaries[0].TotalTasks)
	}

	// Check status counts
	if summaries[0].ByStatus[models.StatusNOW] != 1 {
		t.Errorf("Expected 1 NOW task for Project A, got %d", summaries[0].ByStatus[models.StatusNOW])
	}
	if summaries[0].ByStatus[models.StatusLATER] != 1 {
		t.Errorf("Expected 1 LATER task for Project A, got %d", summaries[0].ByStatus[models.StatusLATER])
	}
	if summaries[0].ByStatus[models.StatusDONE] != 1 {
		t.Errorf("Expected 1 DONE task for Project A, got %d", summaries[0].ByStatus[models.StatusDONE])
	}
}

func TestBuildReferenceGraph(t *testing.T) {
	files := []models.File{
		{Path: "pages/Page A.md", Type: models.FileTypePage},
		{Path: "pages/Page B.md", Type: models.FileTypePage},
		{Path: "pages/Page C.md", Type: models.FileTypePage},
	}

	refs := []models.PageReference{
		{SourcePage: "Page A", TargetPage: "Page B"},
		{SourcePage: "Page A", TargetPage: "Page C"},
		{SourcePage: "Page B", TargetPage: "Page C"},
		{SourcePage: "Page B", TargetPage: "Page C"}, // Duplicate
	}

	graph := BuildReferenceGraph(refs, files)

	// Check nodes exist
	if len(graph.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(graph.Nodes))
	}

	// Check Page A outbound
	nodeA := graph.Nodes["Page A"]
	if len(nodeA.OutboundRefs) != 2 {
		t.Errorf("Page A: expected 2 outbound refs, got %d", len(nodeA.OutboundRefs))
	}

	// Check Page C inbound (should have 2 unique sources: A and B)
	nodeC := graph.Nodes["Page C"]
	if len(nodeC.InboundRefs) != 2 {
		t.Errorf("Page C: expected 2 inbound refs, got %d", len(nodeC.InboundRefs))
	}
	if nodeC.ReferenceCount != 2 {
		t.Errorf("Page C: expected reference count 2, got %d", nodeC.ReferenceCount)
	}

	// Check hub pages (Page C should be top)
	if len(graph.HubPages) == 0 {
		t.Fatal("Expected hub pages, got none")
	}
	if graph.HubPages[0] != "Page C" {
		t.Errorf("Expected Page C as top hub, got %s", graph.HubPages[0])
	}
}

func TestBuildReferenceGraph_NonExistentTargets(t *testing.T) {
	files := []models.File{
		{Path: "pages/Page A.md", Type: models.FileTypePage},
	}

	// Reference to non-existent page
	refs := []models.PageReference{
		{SourcePage: "Page A", TargetPage: "Non Existent Page"},
	}

	graph := BuildReferenceGraph(refs, files)

	// Should create node for non-existent page
	if _, exists := graph.Nodes["Non Existent Page"]; !exists {
		t.Error("Expected node for non-existent page")
	}

	// Non-existent page should have no file path
	if graph.Nodes["Non Existent Page"].FilePath != "" {
		t.Error("Non-existent page should have empty file path")
	}

	// Page A should reference the non-existent page
	nodeA := graph.Nodes["Page A"]
	if len(nodeA.OutboundRefs) != 1 || nodeA.OutboundRefs[0] != "Non Existent Page" {
		t.Error("Page A should reference Non Existent Page")
	}
}

func TestGetOrphanPages(t *testing.T) {
	files := []models.File{
		{Path: "pages/Connected.md"},
		{Path: "pages/Orphan.md"},
	}

	refs := []models.PageReference{
		{SourcePage: "Connected", TargetPage: "Another"},
	}

	graph := BuildReferenceGraph(refs, files)
	orphans := graph.GetOrphanPages()

	// Orphan page should appear
	found := false
	for _, p := range orphans {
		if p == "Orphan" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find Orphan in orphan pages")
	}

	// Connected should not be an orphan (has outbound ref)
	for _, p := range orphans {
		if p == "Connected" {
			t.Error("Connected should not be an orphan")
		}
	}
}

func TestExtractPageNameFromPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"pages/Test Page.md", "Test Page"},
		{"journals/2025_04_06.md", "2025_04_06"},
		{"pages/subfolder/Nested.md", "Nested"},
		{"Test.md", "Test"},
	}

	for _, tt := range tests {
		result := extractPageNameFromPath(tt.input)
		if result != tt.expected {
			t.Errorf("extractPageNameFromPath(%s) = %s, want %s",
				tt.input, result, tt.expected)
		}
	}
}
