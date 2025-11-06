package parser

import (
	"testing"
	"time"

	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

func TestParseTasks_WithPriority(t *testing.T) {
	content := `# Test Page
- NOW [#A] [[Project Phoenix]] - Complete authentication
- LATER [#B] Research alternatives
- TODO [#C] Update documentation
- DOING High priority task with no marker
- DONE [#A] Completed urgent task`

	tasks, err := ParseTasks(content, "test.md")
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}

	if len(tasks) != 5 {
		t.Fatalf("Expected 5 tasks, got %d", len(tasks))
	}

	// Check priority extraction
	if tasks[0].Priority != models.PriorityHigh {
		t.Errorf("Task 0: expected priority A, got %s", tasks[0].Priority)
	}
	if tasks[1].Priority != models.PriorityMedium {
		t.Errorf("Task 1: expected priority B, got %s", tasks[1].Priority)
	}
	if tasks[2].Priority != models.PriorityLow {
		t.Errorf("Task 2: expected priority C, got %s", tasks[2].Priority)
	}
	if tasks[3].Priority != models.PriorityNone {
		t.Errorf("Task 3: expected no priority, got %s", tasks[3].Priority)
	}

	// Check that priority markers are removed from description
	if tasks[0].Description != "[[Project Phoenix]] - Complete authentication" {
		t.Errorf("Task 0: priority marker not removed, got %s", tasks[0].Description)
	}
}

func TestParseTasks_Basic(t *testing.T) {
	content := `# Test Page
- NOW [[Project A]] - Implement feature X
- LATER [[Project B]] - Research options
- TODO Write documentation
- DOING Fix bug #123
- DONE Setup repository`

	tasks, err := ParseTasks(content, "test.md")
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}

	if len(tasks) != 5 {
		t.Fatalf("Expected 5 tasks, got %d", len(tasks))
	}

	// Check NOW task
	if tasks[0].Status != models.StatusNOW {
		t.Errorf("Task 0: expected NOW, got %s", tasks[0].Status)
	}
	if len(tasks[0].PageRefs) != 1 || tasks[0].PageRefs[0] != "Project A" {
		t.Errorf("Task 0: expected page ref 'Project A', got %v", tasks[0].PageRefs)
	}
	// Check priority defaults to None when not specified
	if tasks[0].Priority != models.PriorityNone {
		t.Errorf("Task 0: expected no priority, got %s", tasks[0].Priority)
	}

	// Check LATER task
	if tasks[1].Status != models.StatusLATER {
		t.Errorf("Task 1: expected LATER, got %s", tasks[1].Status)
	}

	// Check TODO task (no page reference)
	if tasks[2].Status != models.StatusTODO {
		t.Errorf("Task 2: expected TODO, got %s", tasks[2].Status)
	}
	if len(tasks[2].PageRefs) != 0 {
		t.Errorf("Task 2: expected no page refs, got %v", tasks[2].PageRefs)
	}

	// Check line numbers
	if tasks[0].LineNumber != 2 {
		t.Errorf("Task 0: expected line 2, got %d", tasks[0].LineNumber)
	}
}

func TestParseTasks_WithLogbook(t *testing.T) {
	content := `- NOW [[Test Project]] - Task with time tracking
  :LOGBOOK:
  CLOCK: [2025-04-06 Sun 10:00:00]--[2025-04-06 Sun 12:00:00] =>  02:00:00
  CLOCK: [2025-04-07 Mon 14:00:00]--[2025-04-07 Mon 15:30:00] =>  01:30:00
  :END:
- LATER Another task`

	tasks, err := ParseTasks(content, "test.md")
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(tasks))
	}

	// Check first task has logbook
	task := tasks[0]
	if len(task.Logbook) != 2 {
		t.Fatalf("Expected 2 logbook entries, got %d", len(task.Logbook))
	}

	// Verify first logbook entry
	entry := task.Logbook[0]
	if entry.Duration != 2*time.Hour {
		t.Errorf("Expected 2h duration, got %v", entry.Duration)
	}

	// Verify second logbook entry
	entry2 := task.Logbook[1]
	if entry2.Duration != 90*time.Minute {
		t.Errorf("Expected 1h30m duration, got %v", entry2.Duration)
	}

	// Check second task has no logbook
	if len(tasks[1].Logbook) != 0 {
		t.Errorf("Task 1 should have no logbook, got %d entries", len(tasks[1].Logbook))
	}
}

func TestParseTasks_LongDuration(t *testing.T) {
	content := `- NOW [[Project]] - Long running task
  :LOGBOOK:
  CLOCK: [2025-04-06 Sun 12:30:42]--[2025-04-17 Thu 09:17:26] =>  260:46:44
  :END:`

	tasks, err := ParseTasks(content, "test.md")
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tasks))
	}

	if len(tasks[0].Logbook) != 1 {
		t.Fatalf("Expected 1 logbook entry, got %d", len(tasks[0].Logbook))
	}

	duration := tasks[0].Logbook[0].Duration
	expectedDuration := 260*time.Hour + 46*time.Minute + 44*time.Second

	if duration != expectedDuration {
		t.Errorf("Expected duration %v, got %v", expectedDuration, duration)
	}
}

func TestParseReferences(t *testing.T) {
	content := `# Test Page

Some text with [[Page A]] and [[Page B]].
Another line with [[Page A]] again.
- A task with [[Project X]]`

	refs, err := ParseReferences(content, "pages/Test Page.md")
	if err != nil {
		t.Fatalf("ParseReferences failed: %v", err)
	}

	if len(refs) != 4 {
		t.Fatalf("Expected 4 references, got %d", len(refs))
	}

	// Check source page extraction
	if refs[0].SourcePage != "Test Page" {
		t.Errorf("Expected source page 'Test Page', got '%s'", refs[0].SourcePage)
	}

	// Check target pages
	expectedTargets := []string{"Page A", "Page B", "Page A", "Project X"}
	for i, ref := range refs {
		if ref.TargetPage != expectedTargets[i] {
			t.Errorf("Ref %d: expected target '%s', got '%s'",
				i, expectedTargets[i], ref.TargetPage)
		}
	}

	// Check line numbers
	if refs[0].LineNumber != 3 {
		t.Errorf("Expected line 3, got %d", refs[0].LineNumber)
	}

	// Check context is populated
	if refs[0].Context == "" {
		t.Error("Expected context to be populated")
	}
}

func TestExtractPageReferences(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			"No references here",
			[]string{},
		},
		{
			"Single [[reference]]",
			[]string{"reference"},
		},
		{
			"Multiple [[Page A]] and [[Page B]]",
			[]string{"Page A", "Page B"},
		},
		{
			"[[Page with spaces and - dashes]]",
			[]string{"Page with spaces and - dashes"},
		},
	}

	for i, tt := range tests {
		result := ExtractPageReferences(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("Test %d: expected %d refs, got %d", i, len(tt.expected), len(result))
			continue
		}

		for j, exp := range tt.expected {
			if result[j] != exp {
				t.Errorf("Test %d, ref %d: expected '%s', got '%s'", i, j, exp, result[j])
			}
		}
	}
}

func TestParseLogbook(t *testing.T) {
	lines := []string{
		"  :LOGBOOK:",
		"  CLOCK: [2025-04-06 Sun 10:00:00]--[2025-04-06 Sun 12:00:00] =>  02:00:00",
		"  CLOCK: [2025-04-07 Mon 09:00:00]--[2025-04-07 Mon 10:00:00] =>  01:00:00",
		"  :END:",
		"- Next task",
	}

	entries, consumed := ParseLogbook(lines, 0)

	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}

	if consumed != 4 {
		t.Errorf("Expected 4 lines consumed, got %d", consumed)
	}

	// Check first entry
	if entries[0].Duration != 2*time.Hour {
		t.Errorf("Entry 0: expected 2h, got %v", entries[0].Duration)
	}

	// Check second entry
	if entries[1].Duration != 1*time.Hour {
		t.Errorf("Entry 1: expected 1h, got %v", entries[1].Duration)
	}
}

func TestParseLogbook_NotLogbook(t *testing.T) {
	lines := []string{
		"- This is not a logbook",
		"Just some text",
	}

	entries, consumed := ParseLogbook(lines, 0)

	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}

	if consumed != 0 {
		t.Errorf("Expected 0 lines consumed, got %d", consumed)
	}
}

func TestExtractContext(t *testing.T) {
	shortLine := "Short line"
	if ctx := ExtractContext(shortLine, 50); ctx != shortLine {
		t.Errorf("Short line should not be truncated")
	}

	longLine := "This is a very long line that should be truncated to the maximum length specified"
	ctx := ExtractContext(longLine, 30)
	if len(ctx) > 33 { // 30 + "..."
		t.Errorf("Expected truncated length ~33, got %d", len(ctx))
	}
	if ctx[len(ctx)-3:] != "..." {
		t.Errorf("Expected ellipsis at end, got: %s", ctx)
	}
}

func TestTaskTotalDuration(t *testing.T) {
	task := models.Task{
		Logbook: []models.LogbookEntry{
			{Duration: 2 * time.Hour},
			{Duration: 30 * time.Minute},
			{Duration: 15 * time.Minute},
		},
	}

	total := task.TotalDuration()
	expected := 2*time.Hour + 45*time.Minute

	if total != expected {
		t.Errorf("Expected total %v, got %v", expected, total)
	}
}
