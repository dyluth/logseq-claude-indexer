package writer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dyluth/logseq-claude-indexer/internal/indexer"
)

func TestWriteMissingPages(t *testing.T) {
	// Create test missing pages index
	index := &indexer.MissingPagesIndex{
		Threshold: 5,
		MissingPages: []indexer.MissingPage{
			{
				Name:           "John Smith - Engineer",
				ReferenceCount: 10,
				PageType:       "person",
				ReferencedFrom: []string{"Page A", "Page B", "Page C"},
			},
			{
				Name:           "Project Phoenix",
				ReferenceCount: 8,
				PageType:       "project",
				ReferencedFrom: []string{"Page D", "Page E"},
			},
			{
				Name:           "GraphQL",
				ReferenceCount: 6,
				PageType:       "concept",
				ReferencedFrom: []string{"Page F"},
			},
		},
	}

	// Write to temp directory
	tmpDir := t.TempDir()
	err := WriteMissingPages(index, tmpDir)
	if err != nil {
		t.Fatalf("WriteMissingPages failed: %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(tmpDir, "missing-pages.md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("missing-pages.md was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// Check for expected sections
	expectedStrings := []string{
		"# Missing Pages to Create",
		"Pages with 5+ references",
		"## People",
		"## Projects",
		"## Concepts",
		"[[John Smith - Engineer]]",
		"**References**: 10",
		"[[Project Phoenix]]",
		"**References**: 8",
		"[[GraphQL]]",
		"**References**: 6",
		"**Referenced from**",
		"[[Page A]]",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Expected content to contain %q, but it didn't", expected)
		}
	}
}

func TestWriteMissingPages_NoMissingPages(t *testing.T) {
	index := &indexer.MissingPagesIndex{
		Threshold:    5,
		MissingPages: []indexer.MissingPage{},
	}

	tmpDir := t.TempDir()
	err := WriteMissingPages(index, tmpDir)
	if err != nil {
		t.Fatalf("WriteMissingPages failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "missing-pages.md"))
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Should indicate no missing pages
	if !strings.Contains(string(content), "No missing pages with 5+ references found") {
		t.Error("Should indicate no missing pages")
	}
}

func TestWriteMissingPages_GroupByType(t *testing.T) {
	// Create index with multiple page types
	index := &indexer.MissingPagesIndex{
		Threshold: 5,
		MissingPages: []indexer.MissingPage{
			{
				Name:           "Alice Johnson",
				ReferenceCount: 10,
				PageType:       "person",
			},
			{
				Name:           "Bob Smith",
				ReferenceCount: 8,
				PageType:       "person",
			},
			{
				Name:           "Sprint 24",
				ReferenceCount: 7,
				PageType:       "project",
			},
			{
				Name:           "Nov 15th, 2025",
				ReferenceCount: 6,
				PageType:       "date",
			},
			{
				Name:           "Microservices",
				ReferenceCount: 5,
				PageType:       "concept",
			},
		},
	}

	tmpDir := t.TempDir()
	err := WriteMissingPages(index, tmpDir)
	if err != nil {
		t.Fatalf("WriteMissingPages failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "missing-pages.md"))
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// Should have all type sections
	if !strings.Contains(contentStr, "## People (2)") {
		t.Error("Should have People section with count")
	}

	if !strings.Contains(contentStr, "## Projects (1)") {
		t.Error("Should have Projects section with count")
	}

	if !strings.Contains(contentStr, "## Dates (1)") {
		t.Error("Should have Dates section with count")
	}

	if !strings.Contains(contentStr, "## Concepts (1)") {
		t.Error("Should have Concepts section with count")
	}

	// Check that pages appear in their sections
	if !strings.Contains(contentStr, "[[Alice Johnson]]") {
		t.Error("Should contain Alice Johnson")
	}

	if !strings.Contains(contentStr, "[[Sprint 24]]") {
		t.Error("Should contain Sprint 24")
	}
}

func TestWriteMissingPages_ReferencedFromList(t *testing.T) {
	index := &indexer.MissingPagesIndex{
		Threshold: 5,
		MissingPages: []indexer.MissingPage{
			{
				Name:           "Popular Page",
				ReferenceCount: 5,
				PageType:       "concept",
				ReferencedFrom: []string{"Source 1", "Source 2", "Source 3"},
			},
		},
	}

	tmpDir := t.TempDir()
	err := WriteMissingPages(index, tmpDir)
	if err != nil {
		t.Fatalf("WriteMissingPages failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "missing-pages.md"))
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// Should list all referenced-from pages
	if !strings.Contains(contentStr, "[[Source 1]]") {
		t.Error("Should contain Source 1")
	}

	if !strings.Contains(contentStr, "[[Source 2]]") {
		t.Error("Should contain Source 2")
	}

	if !strings.Contains(contentStr, "[[Source 3]]") {
		t.Error("Should contain Source 3")
	}

	// Should be comma-separated
	if !strings.Contains(contentStr, "[[Source 1]], [[Source 2]], [[Source 3]]") {
		t.Error("Sources should be comma-separated")
	}
}

func TestWriteMissingPage(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	page := indexer.MissingPage{
		Name:           "Test Missing Page",
		ReferenceCount: 15,
		PageType:       "concept",
		ReferencedFrom: []string{"Page A", "Page B"},
	}

	writeMissingPage(tmpFile, page)
	tmpFile.Sync()

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)

	// Check expected format
	if !strings.Contains(contentStr, "### [[Test Missing Page]]") {
		t.Error("Should have page name as header")
	}

	if !strings.Contains(contentStr, "**References**: 15") {
		t.Error("Should show reference count")
	}

	if !strings.Contains(contentStr, "**Referenced from**") {
		t.Error("Should show referenced from label")
	}

	if !strings.Contains(contentStr, "[[Page A]]") {
		t.Error("Should show first source page")
	}

	if !strings.Contains(contentStr, "[[Page B]]") {
		t.Error("Should show second source page")
	}
}
