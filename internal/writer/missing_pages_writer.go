package writer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

	var buf strings.Builder

	// Write header
	fmt.Fprintf(&buf, "# Missing Pages to Create\n\n")
	fmt.Fprintf(&buf, "Generated: %s\n\n", time.Now().Format(time.RFC3339))

	if len(index.MissingPages) == 0 {
		fmt.Fprintf(&buf, "*No missing pages with %d+ references found.*\n", index.Threshold)
		return writeIfChanged(filePath, buf.String())
	}

	fmt.Fprintf(&buf, "**Pages with %d+ references that don't exist yet**: %d\n\n",
		index.Threshold, len(index.MissingPages))

	fmt.Fprintf(&buf, "---\n\n")

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

		fmt.Fprintf(&buf, "## %s (%d)\n\n", typeLabels[pageType], len(pages))

		for _, page := range pages {
			writeMissingPage(&buf, page)
		}

		fmt.Fprintf(&buf, "---\n\n")
	}

	if err := writeIfChanged(filePath, buf.String()); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}

// writeMissingPage writes a single missing page entry
func writeMissingPage(f io.Writer, page indexer.MissingPage) {
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
