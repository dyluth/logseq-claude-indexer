package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dyluth/logseq-claude-indexer/internal/indexer"
)

// WriteReferenceGraph writes the reference graph to a markdown file
func WriteReferenceGraph(graph *indexer.ReferenceGraph, outputDir string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	filePath := filepath.Join(outputDir, "reference-graph.md")

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	// Write header
	fmt.Fprintf(f, "# Logseq Reference Graph\n\n")
	fmt.Fprintf(f, "Generated: %s\n", graph.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(f, "Total Pages: %d\n", len(graph.Nodes))

	// Count total references
	totalRefs := 0
	for _, node := range graph.Nodes {
		totalRefs += len(node.OutboundRefs)
	}
	fmt.Fprintf(f, "Total References: %d\n\n", totalRefs)
	fmt.Fprintf(f, "---\n\n")

	// Write hub pages
	if len(graph.HubPages) > 0 {
		fmt.Fprintf(f, "## Hub Pages (Most Referenced)\n\n")
		for i, pageName := range graph.HubPages {
			node := graph.Nodes[pageName]
			fmt.Fprintf(f, "%d. **[[%s]]** - %d inbound references\n",
				i+1, node.PageName, node.ReferenceCount)
			if node.FilePath != "" {
				fmt.Fprintf(f, "   - File: `%s`\n", node.FilePath)
			} else {
				fmt.Fprintf(f, "   - *Page not yet created*\n")
			}
			fmt.Fprintf(f, "\n")
		}
		fmt.Fprintf(f, "---\n\n")
	}

	// Write detailed page information
	fmt.Fprintf(f, "## Page Details\n\n")

	// Group by pages with most connections first
	type pageEntry struct {
		name       string
		node       *indexer.GraphNode
		totalConns int
	}

	var entries []pageEntry
	for name, node := range graph.Nodes {
		// Only include pages that actually exist (have a file)
		if node.FilePath != "" {
			totalConns := len(node.InboundRefs) + len(node.OutboundRefs)
			entries = append(entries, pageEntry{name, node, totalConns})
		}
	}

	// Sort by total connections descending
	// Simple bubble sort for small datasets
	for i := 0; i < len(entries)-1; i++ {
		for j := 0; j < len(entries)-i-1; j++ {
			if entries[j].totalConns < entries[j+1].totalConns {
				entries[j], entries[j+1] = entries[j+1], entries[j]
			}
		}
	}

	// Write top 20 most connected pages
	limit := 20
	if len(entries) < limit {
		limit = len(entries)
	}

	for i := 0; i < limit; i++ {
		entry := entries[i]
		node := entry.node

		fmt.Fprintf(f, "### [[%s]]\n", node.PageName)
		fmt.Fprintf(f, "- **File**: `%s`\n", node.FilePath)

		if len(node.OutboundRefs) > 0 {
			fmt.Fprintf(f, "- **Outbound References** (%d):\n", len(node.OutboundRefs))
			// Show first 10 references
			displayLimit := 10
			if len(node.OutboundRefs) < displayLimit {
				displayLimit = len(node.OutboundRefs)
			}
			for j := 0; j < displayLimit; j++ {
				fmt.Fprintf(f, "  - [[%s]]\n", node.OutboundRefs[j])
			}
			if len(node.OutboundRefs) > displayLimit {
				fmt.Fprintf(f, "  - *... and %d more*\n", len(node.OutboundRefs)-displayLimit)
			}
		}

		if len(node.InboundRefs) > 0 {
			fmt.Fprintf(f, "- **Inbound References** (%d):\n", len(node.InboundRefs))
			// Show first 10 references
			displayLimit := 10
			if len(node.InboundRefs) < displayLimit {
				displayLimit = len(node.InboundRefs)
			}
			for j := 0; j < displayLimit; j++ {
				fmt.Fprintf(f, "  - [[%s]]\n", node.InboundRefs[j])
			}
			if len(node.InboundRefs) > displayLimit {
				fmt.Fprintf(f, "  - *... and %d more*\n", len(node.InboundRefs)-displayLimit)
			}
		}

		fmt.Fprintf(f, "\n")
	}

	// Note if there are more pages
	if len(entries) > limit {
		fmt.Fprintf(f, "*Showing top %d of %d pages. Pages with fewer connections are omitted.*\n\n", limit, len(entries))
	}

	return nil
}
