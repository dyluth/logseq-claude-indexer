package indexer

import (
	"testing"
	"time"

	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

func TestBuildTaskIndex_Upcoming(t *testing.T) {
	now := time.Now()

	makeDate := func(daysFromNow int) *time.Time {
		d := now.AddDate(0, 0, daysFromNow)
		return &d
	}

	tasks := []models.Task{
		// 5 days out — should be upcoming
		{Status: models.StatusTODO, Description: "due in 5 days", SourceFile: "t.md", LineNumber: 1, ScheduledDate: makeDate(5)},
		// 14 days out — edge, should be upcoming
		{Status: models.StatusLATER, Description: "due in 14 days", SourceFile: "t.md", LineNumber: 2, ScheduledDate: makeDate(14)},
		// 3 days out — edge, should be upcoming
		{Status: models.StatusNOW, Description: "due in 3 days", SourceFile: "t.md", LineNumber: 3, ScheduledDate: makeDate(3)},
		// 2 days out — too soon, not upcoming
		{Status: models.StatusTODO, Description: "due in 2 days", SourceFile: "t.md", LineNumber: 4, ScheduledDate: makeDate(2)},
		// 15 days out — too far, not upcoming
		{Status: models.StatusTODO, Description: "due in 15 days", SourceFile: "t.md", LineNumber: 5, ScheduledDate: makeDate(15)},
		// 5 days out but DONE — not upcoming
		{Status: models.StatusDONE, Description: "done but dated", SourceFile: "t.md", LineNumber: 6, ScheduledDate: makeDate(5)},
		// No date — not upcoming
		{Status: models.StatusTODO, Description: "no date", SourceFile: "t.md", LineNumber: 7},
	}

	index := BuildTaskIndex(tasks, time.Now())

	if len(index.Upcoming) != 3 {
		t.Errorf("expected 3 upcoming tasks, got %d", len(index.Upcoming))
		for _, u := range index.Upcoming {
			t.Logf("  upcoming: %s (date: %v)", u.Description, u.ScheduledDate)
		}
	}

	// Should be sorted ascending by date (3 days, 5 days, 14 days)
	if len(index.Upcoming) == 3 {
		if index.Upcoming[0].Description != "due in 3 days" {
			t.Errorf("expected first upcoming to be 'due in 3 days', got %q", index.Upcoming[0].Description)
		}
		if index.Upcoming[2].Description != "due in 14 days" {
			t.Errorf("expected last upcoming to be 'due in 14 days', got %q", index.Upcoming[2].Description)
		}
	}
}
