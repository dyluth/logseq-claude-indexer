package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

	var buf strings.Builder

	// Write header
	fmt.Fprintf(&buf, "# Logseq Reference Graph\n\n")
	fmt.Fprintf(&buf, "Generated: %s\n", graph.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(&buf, "Total Pages: %d\n", len(graph.Nodes))

	// Count total references
	totalRefs := 0
	for _, node := range graph.Nodes {
		totalRefs += len(node.OutboundRefs)
	}
	fmt.Fprintf(&buf, "Total References: %d\n\n", totalRefs)
	fmt.Fprintf(&buf, "---\n\n")

	// Write hub pages
	if len(graph.HubPages) > 0 {
		fmt.Fprintf(&buf, "## Hub Pages (Most Referenced)\n\n")
		for i, pageName := range graph.HubPages {
			node := graph.Nodes[pageName]
			fmt.Fprintf(&buf, "%d. **[[%s]]** - %d inbound references\n",
				i+1, node.PageName, node.ReferenceCount)
			if node.FilePath != "" {
				fmt.Fprintf(&buf, "   - File: `%s`\n", node.FilePath)
			} else {
				fmt.Fprintf(&buf, "   - *Page not yet created*\n")
			}
			fmt.Fprintf(&buf, "\n")
		}
		fmt.Fprintf(&buf, "---\n\n")
	}

	// Write detailed page information
	fmt.Fprintf(&buf, "## Page Details\n\n")

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

	// Sort by total connections descending; use page name as tiebreaker so
	// the result is stable across runs regardless of map-iteration order.
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].totalConns != entries[j].totalConns {
			return entries[i].totalConns > entries[j].totalConns
		}
		return entries[i].name < entries[j].name
	})

	// Write top 20 most connected pages
	limit := 20
	if len(entries) < limit {
		limit = len(entries)
	}

	for i := 0; i < limit; i++ {
		entry := entries[i]
		node := entry.node

		fmt.Fprintf(&buf, "### [[%s]]\n", node.PageName)
		fmt.Fprintf(&buf, "- **File**: `%s`\n", node.FilePath)

		if len(node.OutboundRefs) > 0 {
			fmt.Fprintf(&buf, "- **Outbound References** (%d):\n", len(node.OutboundRefs))
			// Show first 10 references
			displayLimit := 10
			if len(node.OutboundRefs) < displayLimit {
				displayLimit = len(node.OutboundRefs)
			}
			for j := 0; j < displayLimit; j++ {
				fmt.Fprintf(&buf, "  - [[%s]]\n", node.OutboundRefs[j])
			}
			if len(node.OutboundRefs) > displayLimit {
				fmt.Fprintf(&buf, "  - *... and %d more*\n", len(node.OutboundRefs)-displayLimit)
			}
		}

		if len(node.InboundRefs) > 0 {
			fmt.Fprintf(&buf, "- **Inbound References** (%d):\n", len(node.InboundRefs))
			// Show first 10 references
			displayLimit := 10
			if len(node.InboundRefs) < displayLimit {
				displayLimit = len(node.InboundRefs)
			}
			for j := 0; j < displayLimit; j++ {
				fmt.Fprintf(&buf, "  - [[%s]]\n", node.InboundRefs[j])
			}
			if len(node.InboundRefs) > displayLimit {
				fmt.Fprintf(&buf, "  - *... and %d more*\n", len(node.InboundRefs)-displayLimit)
			}
		}

		fmt.Fprintf(&buf, "\n")
	}

	// Note if there are more pages
	if len(entries) > limit {
		fmt.Fprintf(&buf, "*Showing top %d of %d pages. Pages with fewer connections are omitted.*\n\n", limit, len(entries))
	}

	if err := writeIfChanged(filePath, buf.String()); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}
