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

func TestWriteTaskIndex(t *testing.T) {
	// Create test index
	tasks := []models.Task{
		{
			Status:      models.StatusNOW,
			Priority:    models.PriorityHigh,
			Description: "Test task 1",
			PageRefs:    []string{"Project A"},
			SourceFile:  "test.md",
			LineNumber:  1,
			Logbook: []models.LogbookEntry{
				{
					Start:    time.Now().Add(-2 * time.Hour),
					End:      time.Now(),
					Duration: 2 * time.Hour,
				},
			},
		},
		{
			Status:      models.StatusLATER,
			Priority:    models.PriorityNone,
			Description: "Test task 2",
			SourceFile:  "test.md",
			LineNumber:  2,
		},
	}

	index := indexer.BuildTaskIndex(tasks)

	// Write to temp directory
	tmpDir := t.TempDir()
	err := WriteTaskIndex(index, tmpDir)
	if err != nil {
		t.Fatalf("WriteTaskIndex failed: %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(tmpDir, "tasks-by-status.md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("tasks-by-status.md was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// Check for expected sections
	expectedStrings := []string{
		"# Tasks by Status",
		"## Statistics",
		"Total Tasks",
		"Completion Rate",
		"Time Tracking",
		"By Priority",
		"By Status",
		"## NOW (1)",
		"## LATER (1)",
		"Test task 1",
		"Test task 2",
		"test.md:1",
		"test.md:2",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Expected content to contain %q, but it didn't", expected)
		}
	}
}

func TestWriteReferenceGraph(t *testing.T) {
	// Create test graph
	files := []models.File{
		{Path: "pages/Page A.md", Type: models.FileTypePage},
		{Path: "pages/Page B.md", Type: models.FileTypePage},
	}

	refs := []models.PageReference{
		{SourcePage: "Page A", TargetPage: "Page B"},
	}

	graph := indexer.BuildReferenceGraph(refs, files)

	// Write to temp directory
	tmpDir := t.TempDir()
	err := WriteReferenceGraph(graph, tmpDir)
	if err != nil {
		t.Fatalf("WriteReferenceGraph failed: %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(tmpDir, "reference-graph.md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("reference-graph.md was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// Check for expected sections
	expectedStrings := []string{
		"# Logseq Reference Graph",
		"Total Pages: 2",
		"Total References: 1",
		"## Hub Pages (Most Referenced)",
		"[[Page B]]",
		"1 inbound references",
		"## Page Details",
		"**Outbound References**",
		"**Inbound References**",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Expected content to contain %q, but it didn't", expected)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{30 * time.Second, "30s"},
		{2 * time.Minute, "2m"},
		{2*time.Minute + 30*time.Second, "2m 30s"},
		{1 * time.Hour, "1h"},
		{2 * time.Hour, "2h"},
		{2*time.Hour + 30*time.Minute, "2h 30m"},
		{100 * time.Hour, "100h"},
		{100*time.Hour + 15*time.Minute, "100h 15m"},
	}

	for _, tt := range tests {
		result := formatDuration(tt.duration)
		if result != tt.expected {
			t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, result, tt.expected)
		}
	}
}

func TestWriteTaskIndex_EmptyIndex(t *testing.T) {
	// Create empty index
	index := indexer.BuildTaskIndex([]models.Task{})

	// Write to temp directory
	tmpDir := t.TempDir()
	err := WriteTaskIndex(index, tmpDir)
	if err != nil {
		t.Fatalf("WriteTaskIndex failed: %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(tmpDir, "tasks-by-status.md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("tasks-by-status.md was not created")
	}

	// Read content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// Should still have header with Total Tasks: 0
	if !strings.Contains(contentStr, "Total Tasks") {
		t.Error("Expected 'Total Tasks' in empty index")
	}

	// Should have statistics section
	if !strings.Contains(contentStr, "## Statistics") {
		t.Error("Expected statistics section in empty index")
	}
}

func TestWriteReferenceGraph_NoFiles(t *testing.T) {
	// Create empty graph
	graph := indexer.BuildReferenceGraph([]models.PageReference{}, []models.File{})

	// Write to temp directory
	tmpDir := t.TempDir()
	err := WriteReferenceGraph(graph, tmpDir)
	if err != nil {
		t.Fatalf("WriteReferenceGraph failed: %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(tmpDir, "reference-graph.md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("reference-graph.md was not created")
	}

	// Read content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Should still have header
	if !strings.Contains(string(content), "# Logseq Reference Graph") {
		t.Error("Expected header in empty graph")
	}
}

func TestWriteTaskIndex_CreateDirectory(t *testing.T) {
	// Use nested directory that doesn't exist
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "nested", "output")

	index := indexer.BuildTaskIndex([]models.Task{})

	err := WriteTaskIndex(index, outputDir)
	if err != nil {
		t.Fatalf("WriteTaskIndex should create directory: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Error("Output directory was not created")
	}
}

func TestWritePriorityIndex(t *testing.T) {
	// Create test tasks with various priorities
	tasks := []models.Task{
		{
			Status:      models.StatusNOW,
			Priority:    models.PriorityHigh,
			Description: "Urgent task with high priority",
			PageRefs:    []string{"Project Phoenix"},
			SourceFile:  "journals/2025-11-06.md",
			LineNumber:  10,
			Logbook: []models.LogbookEntry{
				{Duration: 2 * time.Hour},
			},
		},
		{
			Status:      models.StatusTODO,
			Priority:    models.PriorityHigh,
			Description: "Another high priority task",
			SourceFile:  "pages/tasks.md",
			LineNumber:  20,
		},
		{
			Status:      models.StatusTODO,
			Priority:    models.PriorityMedium,
			Description: "Medium priority task",
			SourceFile:  "pages/tasks.md",
			LineNumber:  30,
		},
		{
			Status:      models.StatusDONE,
			Priority:    models.PriorityNone,
			Description: "Completed task without priority",
			SourceFile:  "pages/tasks.md",
			LineNumber:  40,
		},
	}

	index := indexer.BuildTaskIndex(tasks)

	// Write to temp directory
	tmpDir := t.TempDir()
	err := WritePriorityIndex(index, tmpDir)
	if err != nil {
		t.Fatalf("WritePriorityIndex failed: %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(tmpDir, "tasks-by-priority.md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("tasks-by-priority.md was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// Check for expected sections
	expectedStrings := []string{
		"# High Priority Tasks [#A]",
		"Total High Priority",
		"Urgent task with high priority",
		"Another high priority task",
		"journals/2025-11-06.md:10",
		"pages/tasks.md:20",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Expected content to contain %q, but it didn't", expected)
		}
	}

	// Should NOT contain medium or no-priority tasks
	unexpectedStrings := []string{
		"Medium priority task",
		"Completed task without priority",
	}

	for _, unexpected := range unexpectedStrings {
		if strings.Contains(contentStr, unexpected) {
			t.Errorf("Expected content NOT to contain %q, but it did", unexpected)
		}
	}
}

func TestWritePriorityIndex_NoHighPriority(t *testing.T) {
	// Create test tasks with no high priority
	tasks := []models.Task{
		{
			Status:      models.StatusTODO,
			Priority:    models.PriorityMedium,
			Description: "Medium priority task",
			SourceFile:  "test.md",
			LineNumber:  1,
		},
	}

	index := indexer.BuildTaskIndex(tasks)

	// Write to temp directory
	tmpDir := t.TempDir()
	err := WritePriorityIndex(index, tmpDir)
	if err != nil {
		t.Fatalf("WritePriorityIndex failed: %v", err)
	}

	// Read content
	filePath := filepath.Join(tmpDir, "tasks-by-priority.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// Should indicate no high priority tasks
	if !strings.Contains(contentStr, "No high priority tasks found") {
		t.Error("Expected message about no high priority tasks")
	}
}
