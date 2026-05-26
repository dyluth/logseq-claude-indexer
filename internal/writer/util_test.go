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

func TestStripGeneratedLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			// The Generated line is removed; surrounding blank lines stay.
			// "# Title" + \n + "" + \n + "" + \n + "Content"  (3 newlines total)
			name:     "removes bare Generated line",
			input:    "# Title\n\nGenerated: 2024-01-01T00:00:00Z\n\nContent",
			expected: "# Title\n\n\nContent",
		},
		{
			name:     "removes bold Generated line",
			input:    "# Dashboard\n\n**Generated**: 2024-01-01T00:00:00Z\n\nContent",
			expected: "# Dashboard\n\n\nContent",
		},
		{
			name:     "leaves unrelated lines untouched",
			input:    "# Title\n\nSome content\n\nMore content",
			expected: "# Title\n\nSome content\n\nMore content",
		},
		{
			name:     "handles leading whitespace on Generated line",
			input:    "# Title\n\n  Generated: 2024-01-01T00:00:00Z\n\nContent",
			expected: "# Title\n\n\nContent",
		},
		{
			name:     "does not strip a line that merely contains Generated mid-sentence",
			input:    "# Title\n\nThis was Generated: by hand\n\nContent",
			expected: "# Title\n\nThis was Generated: by hand\n\nContent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripGeneratedLine(tt.input)
			if got != tt.expected {
				t.Errorf("stripGeneratedLine():\ngot:  %q\nwant: %q", got, tt.expected)
			}
		})
	}
}

func TestWriteIfChanged_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "new.md")

	content := "# Title\n\nGenerated: 2024-01-01T00:00:00Z\n\nContent\n"
	if err := writeIfChanged(path, content); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}
	if string(got) != content {
		t.Errorf("file content mismatch:\ngot:  %q\nwant: %q", string(got), content)
	}
}

func TestWriteIfChanged_SkipsWhenOnlyTimestampDiffers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "index.md")

	// Write initial version
	v1 := "# Title\n\nGenerated: 2024-01-01T00:00:00Z\n\nReal content\n"
	if err := os.WriteFile(path, []byte(v1), 0644); err != nil {
		t.Fatal(err)
	}

	// Record mtime before the call
	info1, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	// Call writeIfChanged with only the timestamp updated
	v2 := "# Title\n\nGenerated: 2025-06-15T12:00:00Z\n\nReal content\n"
	if err := writeIfChanged(path, v2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info2, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	if !info2.ModTime().Equal(info1.ModTime()) {
		t.Error("file was written when only the timestamp changed; mtime should be unchanged")
	}

	// Verify the file still contains the original content (v1)
	got, _ := os.ReadFile(path)
	if string(got) != v1 {
		t.Errorf("file should still contain v1 content, got: %q", string(got))
	}
}

func TestWriteIfChanged_WritesWhenContentDiffers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "index.md")

	v1 := "# Title\n\nGenerated: 2024-01-01T00:00:00Z\n\nOld content\n"
	if err := os.WriteFile(path, []byte(v1), 0644); err != nil {
		t.Fatal(err)
	}

	v2 := "# Title\n\nGenerated: 2025-06-15T12:00:00Z\n\nNew content\n"
	if err := writeIfChanged(path, v2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := os.ReadFile(path)
	if string(got) != v2 {
		t.Errorf("file should have been updated to v2, got: %q", string(got))
	}
}

// TestWriteIfChanged_SkipsBoldGeneratedVariant covers the dashboard.md format
// which uses **Generated**: instead of Generated:.
func TestWriteIfChanged_SkipsBoldGeneratedVariant(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "dashboard.md")

	v1 := "# Dashboard\n\n**Generated**: 2024-01-01T00:00:00Z\n\nStats\n"
	if err := os.WriteFile(path, []byte(v1), 0644); err != nil {
		t.Fatal(err)
	}
	info1, _ := os.Stat(path)

	v2 := "# Dashboard\n\n**Generated**: 2025-06-15T12:00:00Z\n\nStats\n"
	if err := writeIfChanged(path, v2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info2, _ := os.Stat(path)

	if !info2.ModTime().Equal(info1.ModTime()) {
		t.Error("dashboard: file written when only **Generated** timestamp changed")
	}
}

// TestTimelineRecentDeterministic verifies that WriteTimelineRecent produces
// identical output on consecutive calls (only the timestamp should differ,
// which writeIfChanged ignores).
func TestTimelineRecentDeterministic(t *testing.T) {
	// Build a small timeline index with a journal entry from 3 days ago.
	threeAgo := time.Now().UTC().AddDate(0, 0, -3).Truncate(24 * time.Hour)

	tasks := []models.Task{
		{
			Status:      models.StatusTODO,
			Description: "Write tests",
			SourceFile:  threeAgo.Format("journals/2006_01_02.md"),
			LineNumber:  1,
		},
	}
	files := []models.File{
		{
			Path:         threeAgo.Format("journals/2006_01_02.md"),
			Type:         models.FileTypeJournal,
			AbsolutePath: threeAgo.Format("/repo/journals/2006_01_02.md"),
		},
	}

	dir := t.TempDir()

	// First write (file does not exist yet → always writes).
	idx1 := indexer.BuildTimelineIndex(tasks, files)
	if err := WriteTimelineRecent(idx1, dir); err != nil {
		t.Fatalf("first write: %v", err)
	}
	content1, _ := os.ReadFile(filepath.Join(dir, "timeline-recent.md"))

	// Second write (file exists, content should match except timestamp).
	idx2 := indexer.BuildTimelineIndex(tasks, files)
	if err := WriteTimelineRecent(idx2, dir); err != nil {
		t.Fatalf("second write: %v", err)
	}
	content2, _ := os.ReadFile(filepath.Join(dir, "timeline-recent.md"))

	// The file on disk after the second write should still be the first write's
	// content (writeIfChanged skipped the write), OR at minimum the non-timestamp
	// portions should be identical.
	stripped1 := stripGeneratedLine(string(content1))
	stripped2 := stripGeneratedLine(string(content2))
	if stripped1 != stripped2 {
		t.Errorf("timeline-recent.md content differs between runs (ignoring timestamp):\n--- run1 ---\n%s\n--- run2 ---\n%s",
			stripped1, stripped2)
	}
}

// TestMissingPagesSortDeterministic verifies that BuildMissingPagesIndex
// produces a stable ordering when pages share the same reference count.
func TestMissingPagesSortDeterministic(t *testing.T) {
	// Create a graph where multiple missing pages share the same ref count.
	refs := []models.PageReference{
		{SourcePage: "A", TargetPage: "Zara"},
		{SourcePage: "B", TargetPage: "Zara"},
		{SourcePage: "C", TargetPage: "Zara"},
		{SourcePage: "D", TargetPage: "Zara"},
		{SourcePage: "E", TargetPage: "Zara"},
		// "Alpha" and "Beta" both get 5 refs — tie on count
		{SourcePage: "A", TargetPage: "Alpha"},
		{SourcePage: "B", TargetPage: "Alpha"},
		{SourcePage: "C", TargetPage: "Alpha"},
		{SourcePage: "D", TargetPage: "Alpha"},
		{SourcePage: "E", TargetPage: "Alpha"},
		{SourcePage: "A", TargetPage: "Beta"},
		{SourcePage: "B", TargetPage: "Beta"},
		{SourcePage: "C", TargetPage: "Beta"},
		{SourcePage: "D", TargetPage: "Beta"},
		{SourcePage: "E", TargetPage: "Beta"},
	}
	files := []models.File{} // no files → all pages are "missing"

	// Run the build 10 times and check that the order is always the same.
	var firstOrder []string
	for i := 0; i < 10; i++ {
		graph := indexer.BuildReferenceGraph(refs, files)
		idx := indexer.BuildMissingPagesIndex(graph, 5)

		var names []string
		for _, p := range idx.MissingPages {
			names = append(names, p.Name)
		}

		if firstOrder == nil {
			firstOrder = names
			continue
		}
		if strings.Join(names, ",") != strings.Join(firstOrder, ",") {
			t.Errorf("iteration %d: order changed\ngot:  %v\nwant: %v", i, names, firstOrder)
		}
	}
}
