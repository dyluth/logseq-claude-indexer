package indexer

import (
	"testing"

	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

func TestBuildMissingPagesIndex(t *testing.T) {
	// Create a reference graph with some missing pages
	refs := []models.PageReference{
		{SourcePage: "Page A", TargetPage: "Missing Person"},
		{SourcePage: "Page B", TargetPage: "Missing Person"},
		{SourcePage: "Page C", TargetPage: "Missing Person"},
		{SourcePage: "Page D", TargetPage: "Missing Person"},
		{SourcePage: "Page E", TargetPage: "Missing Person"},
		{SourcePage: "Page A", TargetPage: "Low Ref Page"},
		{SourcePage: "Page B", TargetPage: "Low Ref Page"},
		{SourcePage: "Page A", TargetPage: "Another Missing"},
		{SourcePage: "Page B", TargetPage: "Another Missing"},
		{SourcePage: "Page C", TargetPage: "Another Missing"},
		{SourcePage: "Page D", TargetPage: "Another Missing"},
		{SourcePage: "Page E", TargetPage: "Another Missing"},
		{SourcePage: "Page F", TargetPage: "Another Missing"},
	}

	files := []models.File{
		{Path: "pages/Page A.md", Type: models.FileTypePage},
		{Path: "pages/Page B.md", Type: models.FileTypePage},
		{Path: "pages/Page C.md", Type: models.FileTypePage},
	}

	graph := BuildReferenceGraph(refs, files)

	// Build missing pages index with threshold of 5
	index := BuildMissingPagesIndex(graph, 5)

	if index.Threshold != 5 {
		t.Errorf("Expected threshold 5, got %d", index.Threshold)
	}

	// Should have 2 missing pages (5+ refs each)
	if len(index.MissingPages) != 2 {
		t.Fatalf("Expected 2 missing pages, got %d", len(index.MissingPages))
	}

	// Check sorting (highest count first)
	if index.MissingPages[0].ReferenceCount < index.MissingPages[1].ReferenceCount {
		t.Error("Missing pages not sorted by reference count descending")
	}

	// Should NOT include "Low Ref Page" (only 2 refs)
	for _, page := range index.MissingPages {
		if page.Name == "Low Ref Page" {
			t.Error("Should not include pages below threshold")
		}
	}

	// Check that ReferencedFrom is populated
	if len(index.MissingPages[0].ReferencedFrom) == 0 {
		t.Error("ReferencedFrom should be populated")
	}
}

func TestBuildMissingPagesIndex_NoMissingPages(t *testing.T) {
	// All referenced pages exist
	refs := []models.PageReference{
		{SourcePage: "Page A", TargetPage: "Page B"},
	}

	files := []models.File{
		{Path: "pages/Page A.md", Type: models.FileTypePage},
		{Path: "pages/Page B.md", Type: models.FileTypePage},
	}

	graph := BuildReferenceGraph(refs, files)
	index := BuildMissingPagesIndex(graph, 5)

	if len(index.MissingPages) != 0 {
		t.Errorf("Expected no missing pages, got %d", len(index.MissingPages))
	}
}

func TestBuildMissingPagesIndex_ThresholdFiltering(t *testing.T) {
	// Create references with different counts
	refs := []models.PageReference{
		{SourcePage: "A", TargetPage: "Missing 10"},
		{SourcePage: "B", TargetPage: "Missing 10"},
		{SourcePage: "C", TargetPage: "Missing 10"},
		{SourcePage: "D", TargetPage: "Missing 10"},
		{SourcePage: "E", TargetPage: "Missing 10"},
		{SourcePage: "F", TargetPage: "Missing 10"},
		{SourcePage: "G", TargetPage: "Missing 10"},
		{SourcePage: "H", TargetPage: "Missing 10"},
		{SourcePage: "I", TargetPage: "Missing 10"},
		{SourcePage: "J", TargetPage: "Missing 10"},
		{SourcePage: "A", TargetPage: "Missing 3"},
		{SourcePage: "B", TargetPage: "Missing 3"},
		{SourcePage: "C", TargetPage: "Missing 3"},
	}

	files := []models.File{
		{Path: "pages/A.md", Type: models.FileTypePage},
	}

	graph := BuildReferenceGraph(refs, files)
	index := BuildMissingPagesIndex(graph, 5)

	// Should only include "Missing 10" (10 refs), not "Missing 3" (3 refs)
	if len(index.MissingPages) != 1 {
		t.Fatalf("Expected 1 missing page with 5+ refs, got %d", len(index.MissingPages))
	}

	if index.MissingPages[0].Name != "Missing 10" {
		t.Errorf("Expected 'Missing 10', got '%s'", index.MissingPages[0].Name)
	}

	if index.MissingPages[0].ReferenceCount != 10 {
		t.Errorf("Expected 10 references, got %d", index.MissingPages[0].ReferenceCount)
	}
}

func TestClassifyPageType(t *testing.T) {
	tests := []struct {
		pageName     string
		expectedType string
	}{
		{"John Smith - Engineering Lead", "person"},
		{"Sarah Chen - Product Manager", "person"},
		{"Alice Johnson", "person"},
		{"Nov 15th, 2025", "date"},
		{"Apr 22nd, 2025", "date"},
		{"Sprint 23", "project"},
		{"Project Phoenix", "project"},
		{"API Design", "concept"},
		{"GraphQL", "concept"},
		{"Microservices Architecture", "concept"},
		{"abc", "concept"},
		{"OneWord", "concept"},
	}

	for _, tt := range tests {
		result := classifyPageType(tt.pageName)
		if result != tt.expectedType {
			t.Errorf("classifyPageType(%q) = %q, expected %q", tt.pageName, result, tt.expectedType)
		}
	}
}

func TestLooksLikeName(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"John Smith", true},
		{"Sarah Chen", true},
		{"Alice Mary Johnson", true},
		{"John", false},                  // Single word
		{"john smith", false},            // Not capitalized
		{"John smith", false},            // Second word not capitalized
		{"API Design", true},             // Technically matches pattern (both capitalized)
		{"This Is Too Many Words Name", false}, // Too many words
		{"", false},                      // Empty
	}

	for _, tt := range tests {
		result := looksLikeName(tt.input)
		if result != tt.expected {
			t.Errorf("looksLikeName(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestMissingPagesIndex_ReferencedFromLimit(t *testing.T) {
	// Create a missing page referenced from 15 different pages
	var refs []models.PageReference
	for i := 0; i < 15; i++ {
		refs = append(refs, models.PageReference{
			SourcePage: string(rune('A' + i)),
			TargetPage: "Popular Missing Page",
		})
	}

	files := []models.File{
		{Path: "pages/A.md", Type: models.FileTypePage},
	}

	graph := BuildReferenceGraph(refs, files)
	index := BuildMissingPagesIndex(graph, 5)

	if len(index.MissingPages) != 1 {
		t.Fatalf("Expected 1 missing page, got %d", len(index.MissingPages))
	}

	// ReferencedFrom should be limited to 10
	if len(index.MissingPages[0].ReferencedFrom) > 10 {
		t.Errorf("ReferencedFrom should be limited to 10, got %d", len(index.MissingPages[0].ReferencedFrom))
	}
}

func TestMissingPagesIndex_Sorting(t *testing.T) {
	// Create missing pages with different reference counts
	refs := []models.PageReference{}

	// Page with 10 refs
	for i := 0; i < 10; i++ {
		refs = append(refs, models.PageReference{
			SourcePage: string(rune('A' + i)),
			TargetPage: "Page 10",
		})
	}

	// Page with 7 refs
	for i := 0; i < 7; i++ {
		refs = append(refs, models.PageReference{
			SourcePage: string(rune('K' + i)),
			TargetPage: "Page 7",
		})
	}

	// Page with 5 refs
	for i := 0; i < 5; i++ {
		refs = append(refs, models.PageReference{
			SourcePage: string(rune('S' + i)),
			TargetPage: "Page 5",
		})
	}

	files := []models.File{
		{Path: "pages/A.md", Type: models.FileTypePage},
	}

	graph := BuildReferenceGraph(refs, files)
	index := BuildMissingPagesIndex(graph, 5)

	if len(index.MissingPages) != 3 {
		t.Fatalf("Expected 3 missing pages, got %d", len(index.MissingPages))
	}

	// Check sorting: should be 10, 7, 5
	expectedOrder := []int{10, 7, 5}
	for i, expected := range expectedOrder {
		if index.MissingPages[i].ReferenceCount != expected {
			t.Errorf("Position %d: expected %d refs, got %d", i, expected, index.MissingPages[i].ReferenceCount)
		}
	}
}
