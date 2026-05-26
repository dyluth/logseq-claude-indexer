package writer

import (
	"os"
	"path/filepath"
	"testing"
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
