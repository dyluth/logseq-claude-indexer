package parser

import (
	"path/filepath"
	"strings"

	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

// ParseReferences extracts all [[page]] references from markdown content
func ParseReferences(content string, filePath string) ([]models.PageReference, error) {
	var refs []models.PageReference
	lines := strings.Split(content, "\n")

	sourcePage := extractPageNameFromPath(filePath)

	for i, line := range lines {
		// Find all page references in this line
		pageRefs := ExtractPageReferences(line)

		for _, targetPage := range pageRefs {
			refs = append(refs, models.PageReference{
				SourceFile: filePath,
				SourcePage: sourcePage,
				TargetPage: targetPage,
				LineNumber: i + 1, // 1-indexed
				Context:    ExtractContext(line, 100),
			})
		}
	}

	return refs, nil
}

// extractPageNameFromPath converts a file path to a page name
// Examples:
//   journals/2025_04_06.md -> "2025_04_06"
//   pages/Hearth Insights.md -> "Hearth Insights"
func extractPageNameFromPath(filePath string) string {
	base := filepath.Base(filePath)
	return strings.TrimSuffix(base, ".md")
}
