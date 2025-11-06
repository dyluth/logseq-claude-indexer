package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dyluth/logseq-claude-indexer/internal/indexer"
)

// WriteMissingPages writes the missing pages index to a markdown file
func WriteMissingPages(index *indexer.MissingPagesIndex, outputDir string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	filePath := filepath.Join(outputDir, "missing-pages.md")

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	// Write header
	fmt.Fprintf(f, "# Missing Pages to Create\n\n")
	fmt.Fprintf(f, "Generated: %s\n\n", time.Now().Format(time.RFC3339))

	if len(index.MissingPages) == 0 {
		fmt.Fprintf(f, "*No missing pages with %d+ references found.*\n", index.Threshold)
		return nil
	}

	fmt.Fprintf(f, "**Pages with %d+ references that don't exist yet**: %d\n\n",
		index.Threshold, len(index.MissingPages))

	fmt.Fprintf(f, "---\n\n")

	// Group by page type
	byType := make(map[string][]indexer.MissingPage)
	for _, page := range index.MissingPages {
		byType[page.PageType] = append(byType[page.PageType], page)
	}

	// Write pages by type in priority order
	typeOrder := []string{"person", "project", "date", "concept"}
	typeLabels := map[string]string{
		"person":  "People",
		"project": "Projects",
		"date":    "Dates",
		"concept": "Concepts",
	}

	for _, pageType := range typeOrder {
		pages, exists := byType[pageType]
		if !exists || len(pages) == 0 {
			continue
		}

		fmt.Fprintf(f, "## %s (%d)\n\n", typeLabels[pageType], len(pages))

		for _, page := range pages {
			writeMissingPage(f, page)
		}

		fmt.Fprintf(f, "---\n\n")
	}

	return nil
}

// writeMissingPage writes a single missing page entry
func writeMissingPage(f *os.File, page indexer.MissingPage) {
	fmt.Fprintf(f, "### [[%s]]\n", page.Name)
	fmt.Fprintf(f, "- **References**: %d\n", page.ReferenceCount)

	// Show first few pages that reference this
	if len(page.ReferencedFrom) > 0 {
		fmt.Fprintf(f, "- **Referenced from**: ")
		for i, sourcePage := range page.ReferencedFrom {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "[[%s]]", sourcePage)
		}
		fmt.Fprintf(f, "\n")
	}

	fmt.Fprintf(f, "\n")
}
