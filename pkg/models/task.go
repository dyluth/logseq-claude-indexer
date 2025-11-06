package models

import "time"

// TaskStatus represents the state of a task in Logseq
type TaskStatus string

const (
	StatusNOW   TaskStatus = "NOW"
	StatusLATER TaskStatus = "LATER"
	StatusTODO  TaskStatus = "TODO"
	StatusDOING TaskStatus = "DOING"
	StatusDONE  TaskStatus = "DONE"
)

// Priority represents the urgency level of a task in Logseq
// Extracted from markers like [#A], [#B], [#C]
type Priority string

const (
	PriorityHigh   Priority = "A"
	PriorityMedium Priority = "B"
	PriorityLow    Priority = "C"
	PriorityNone   Priority = ""
)

// Task represents a task extracted from Logseq markdown
// Example: - NOW [#A] [[Project Name]] - Task description
type Task struct {
	Status      TaskStatus      // The task's status marker
	Priority    Priority        // The task's priority level ([#A], [#B], [#C])
	Description string          // Full task text (without status/priority markers)
	PageRefs    []string        // [[Page Name]] references found in the task
	SourceFile  string          // Relative path to file containing this task
	LineNumber  int             // Line number where task appears (1-indexed)
	Logbook     []LogbookEntry  // Time tracking entries (if :LOGBOOK: present)
}

// TotalDuration calculates the sum of all logbook entry durations
func (t *Task) TotalDuration() time.Duration {
	var total time.Duration
	for _, entry := range t.Logbook {
		total += entry.Duration
	}
	return total
}
