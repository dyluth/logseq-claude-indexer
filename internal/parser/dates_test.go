package parser

import (
	"testing"
	"time"
)

func TestParseLogseqDate(t *testing.T) {
	tests := []struct {
		input   string
		wantOK  bool
		wantDay int
		wantMon time.Month
		wantYr  int
	}{
		{"Jun 3rd, 2026", true, 3, time.June, 2026},
		{"Jan 1st, 2026", true, 1, time.January, 2026},
		{"March 15th, 2026", true, 15, time.March, 2026},
		{"December 31st, 2025", true, 31, time.December, 2025},
		{"Feb 22nd, 2027", true, 22, time.February, 2027},
		{"Project Name", false, 0, 0, 0},
		{"Some Task", false, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, ok := parseLogseqDate(tt.input)
			if ok != tt.wantOK {
				t.Fatalf("parseLogseqDate(%q) ok = %v, want %v", tt.input, ok, tt.wantOK)
			}
			if !tt.wantOK {
				return
			}
			if got.Day() != tt.wantDay || got.Month() != tt.wantMon || got.Year() != tt.wantYr {
				t.Errorf("parseLogseqDate(%q) = %v, want %d %v %d",
					tt.input, got, tt.wantDay, tt.wantMon, tt.wantYr)
			}
		})
	}
}

func TestParseTasksExtractsScheduledDate(t *testing.T) {
	content := "- TODO [[Jun 3rd, 2026]] review deployment plan\n- NOW [[Project Alpha]] unrelated task\n"
	tasks, err := ParseTasks(content, "test.md")
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	// First task has a date ref
	if tasks[0].ScheduledDate == nil {
		t.Error("expected ScheduledDate to be set for task with date ref")
	} else {
		if tasks[0].ScheduledDate.Month() != time.June || tasks[0].ScheduledDate.Day() != 3 {
			t.Errorf("unexpected ScheduledDate: %v", tasks[0].ScheduledDate)
		}
	}

	// Second task has no date ref
	if tasks[1].ScheduledDate != nil {
		t.Errorf("expected ScheduledDate to be nil for task without date ref, got %v", tasks[1].ScheduledDate)
	}
}
