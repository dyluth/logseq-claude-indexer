package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

func TestScanner_Scan(t *testing.T) {
	// Create temp directory with test structure
	tmpDir := t.TempDir()

	// Create pages and journals directories
	pagesDir := filepath.Join(tmpDir, "pages")
	journalsDir := filepath.Join(tmpDir, "journals")
	os.MkdirAll(pagesDir, 0755)
	os.MkdirAll(journalsDir, 0755)

	// Create test files
	testFiles := map[string]string{
		"pages/Test Page.md":        "# Test Page",
		"pages/Another Page.md":     "# Another",
		"journals/2025_04_06.md":    "# Journal",
		"journals/2025_04_07.md":    "# Journal 2",
		"pages/.hidden.md":          "# Hidden (should be skipped)",
		"pages/not-markdown.txt":    "Not markdown (should be skipped)",
		"pages/.recycle/deleted.md": "# Deleted (should be skipped)",
	}

	for relPath, content := range testFiles {
		fullPath := filepath.Join(tmpDir, relPath)
		os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", relPath, err)
		}
	}

	// Run scanner
	scanner := New(tmpDir)
	files, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}

	// Verify results
	if len(files) != 4 {
		t.Errorf("Expected 4 files, got %d", len(files))
		for _, f := range files {
			t.Logf("Found: %s", f.Path)
		}
	}

	// Count by type
	pageCount := 0
	journalCount := 0
	for _, f := range files {
		if f.Type == models.FileTypePage {
			pageCount++
		} else if f.Type == models.FileTypeJournal {
			journalCount++
		}

		// Verify paths are relative
		if filepath.IsAbs(f.Path) {
			t.Errorf("Expected relative path, got absolute: %s", f.Path)
		}

		// Verify absolute paths are absolute
		if !filepath.IsAbs(f.AbsolutePath) {
			t.Errorf("Expected absolute path, got relative: %s", f.AbsolutePath)
		}
	}

	if pageCount != 2 {
		t.Errorf("Expected 2 page files, got %d", pageCount)
	}
	if journalCount != 2 {
		t.Errorf("Expected 2 journal files, got %d", journalCount)
	}
}

func TestScanner_EmptyRepository(t *testing.T) {
	// Create empty temp directory (no pages or journals dirs)
	tmpDir := t.TempDir()

	scanner := New(tmpDir)
	files, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() should not fail on empty repo: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files in empty repo, got %d", len(files))
	}
}

func TestScanner_OnlyPages(t *testing.T) {
	tmpDir := t.TempDir()

	// Create only pages directory
	pagesDir := filepath.Join(tmpDir, "pages")
	os.MkdirAll(pagesDir, 0755)
	os.WriteFile(filepath.Join(pagesDir, "Test.md"), []byte("# Test"), 0644)

	scanner := New(tmpDir)
	files, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}

	if len(files) > 0 && files[0].Type != models.FileTypePage {
		t.Errorf("Expected FileTypePage, got %v", files[0].Type)
	}
}

func TestScanner_OnlyJournals(t *testing.T) {
	tmpDir := t.TempDir()

	// Create only journals directory
	journalsDir := filepath.Join(tmpDir, "journals")
	os.MkdirAll(journalsDir, 0755)
	os.WriteFile(filepath.Join(journalsDir, "2025_04_06.md"), []byte("# Journal"), 0644)

	scanner := New(tmpDir)
	files, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}

	if len(files) > 0 && files[0].Type != models.FileTypeJournal {
		t.Errorf("Expected FileTypeJournal, got %v", files[0].Type)
	}
}

func TestScanner_NestedDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested structure
	nestedDir := filepath.Join(tmpDir, "pages", "subfolder")
	os.MkdirAll(nestedDir, 0755)
	os.WriteFile(filepath.Join(nestedDir, "Nested.md"), []byte("# Nested"), 0644)

	scanner := New(tmpDir)
	files, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file (nested), got %d", len(files))
	}

	if len(files) > 0 {
		expectedPath := filepath.Join("pages", "subfolder", "Nested.md")
		if files[0].Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, files[0].Path)
		}
	}
}
