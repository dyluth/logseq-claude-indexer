package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dyluth/logseq-claude-indexer/internal/indexer"
)

// WriteTasksDone writes all completed tasks to tasks-done.md.
// Tasks are sorted by most recent logbook activity (newest first); tasks without
// logbook entries appear at the end sorted by source file.
func WriteTasksDone(idx *indexer.TaskIndex, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	filePath := filepath.Join(outputDir, "tasks-done.md")

	var buf strings.Builder

	fmt.Fprintf(&buf, "# Completed Tasks\n\n")
	fmt.Fprintf(&buf, "Generated: %s\n\n", idx.GeneratedAt.Format(time.RFC3339))

	// Stats for done tasks
	var doneTimeLogged time.Duration
	trackedCount := 0
	for _, t := range idx.Done {
		if len(t.Logbook) > 0 {
			doneTimeLogged += t.TotalDuration()
			trackedCount++
		}
	}

	fmt.Fprintf(&buf, "**Total Completed**: %d tasks\n", len(idx.Done))
	if trackedCount > 0 {
		fmt.Fprintf(&buf, "**Total Time Logged**: %s across %d tracked tasks\n", formatDuration(doneTimeLogged), trackedCount)
	}
	fmt.Fprintf(&buf, "\n---\n\n")

	if len(idx.Done) == 0 {
		fmt.Fprintf(&buf, "*No completed tasks found.*\n")
		return writeIfChanged(filePath, buf.String())
	}

	// Recently completed: DONE tasks with logbook activity in last 30 days
	thirtyDaysAgo := idx.GeneratedAt.AddDate(0, 0, -30)
	var recent, older []indexEntry

	for i, t := range idx.Done {
		last := indexer.GetMostRecentTime(t)
		if !last.IsZero() && last.After(thirtyDaysAgo) {
			recent = append(recent, indexEntry{i, last})
		} else {
			older = append(older, indexEntry{i, last})
		}
	}

	if len(recent) > 0 {
		fmt.Fprintf(&buf, "## Recently Completed (Last 30 Days) (%d)\n\n", len(recent))
		for _, e := range recent {
			writeFullTask(&buf, idx.Done[e.idx])
		}
		fmt.Fprintf(&buf, "---\n\n")
	}

	fmt.Fprintf(&buf, "## All Completed (%d)\n\n", len(idx.Done))
	for _, t := range idx.Done {
		writeFullTask(&buf, t)
	}

	if err := writeIfChanged(filePath, buf.String()); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}

type indexEntry struct {
	idx  int
	time time.Time
}
