package indexer

import (
	"sort"
	"strings"
)

// MissingPage represents a page that is referenced but doesn't exist
type MissingPage struct {
	Name           string
	ReferenceCount int
	PageType       string // person, date, project, concept
	ReferencedFrom []string
}

// MissingPagesIndex contains pages that should be created
type MissingPagesIndex struct {
	MissingPages []MissingPage
	Threshold    int // Minimum references to be included (5)
}

// BuildMissingPagesIndex identifies high-value pages to create from the reference graph
func BuildMissingPagesIndex(graph *ReferenceGraph, threshold int) *MissingPagesIndex {
	index := &MissingPagesIndex{
		Threshold: threshold,
	}

	// Iterate through all referenced pages
	for pageName, node := range graph.Nodes {
		// Skip pages that already exist as files (FilePath is empty if doesn't exist)
		if node.FilePath != "" {
			continue
		}

		// Only include pages with threshold or more references
		if node.ReferenceCount < threshold {
			continue
		}

		// Collect unique pages that reference this missing page
		// InboundRefs is a slice of page names
		referencedFrom := make([]string, 0, len(node.InboundRefs))
		seen := make(map[string]bool)
		for _, sourcePage := range node.InboundRefs {
			if !seen[sourcePage] {
				referencedFrom = append(referencedFrom, sourcePage)
				seen[sourcePage] = true
			}
		}

		// Limit to first 10 referencing pages
		if len(referencedFrom) > 10 {
			referencedFrom = referencedFrom[:10]
		}

		missingPage := MissingPage{
			Name:           pageName,
			ReferenceCount: node.ReferenceCount,
			PageType:       classifyPageType(pageName),
			ReferencedFrom: referencedFrom,
		}

		index.MissingPages = append(index.MissingPages, missingPage)
	}

	// Sort by reference count descending
	sort.Slice(index.MissingPages, func(i, j int) bool {
		return index.MissingPages[i].ReferenceCount > index.MissingPages[j].ReferenceCount
	})

	return index
}

// classifyPageType determines the likely type of a page based on its name
func classifyPageType(pageName string) string {
	// Person: "FirstName LastName - Title" pattern
	if strings.Contains(pageName, " - ") {
		// Check if it's a person pattern: "Name - Role"
		parts := strings.Split(pageName, " - ")
		if len(parts) == 2 && looksLikeName(parts[0]) {
			return "person"
		}
		return "concept"
	}

	// Date: Contains month names or date patterns
	dateKeywords := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec", "th,", "st,", "nd,", "rd,"}
	for _, keyword := range dateKeywords {
		if strings.Contains(pageName, keyword) {
			return "date"
		}
	}

	// Project: Contains project-related keywords
	projectKeywords := []string{"Sprint", "Project", "Phase", "Release", "Milestone"}
	for _, keyword := range projectKeywords {
		if strings.Contains(pageName, keyword) {
			return "project"
		}
	}

	// Technical/concept keywords that override name detection
	techKeywords := []string{"API", "Service", "Architecture", "System", "Design", "Framework", "Database", "Server", "Team", "Stack", "Platform"}
	for _, keyword := range techKeywords {
		if strings.Contains(pageName, keyword) {
			return "concept"
		}
	}

	// Person: Capitalized words pattern (likely a name)
	// Only if exactly 2-3 words (typical name pattern)
	words := strings.Fields(pageName)
	if len(words) >= 2 && len(words) <= 3 && looksLikeName(pageName) {
		return "person"
	}

	// Default
	return "concept"
}

// looksLikeName checks if a string looks like a person's name
func looksLikeName(s string) bool {
	words := strings.Fields(s)
	if len(words) < 2 || len(words) > 4 {
		return false
	}

	// Check if each word starts with capital letter
	for _, word := range words {
		if len(word) == 0 {
			return false
		}
		if word[0] < 'A' || word[0] > 'Z' {
			return false
		}
	}

	return true
}
