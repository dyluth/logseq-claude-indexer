package indexer

import (
	"sort"
	"time"

	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

// ReferenceGraph represents the network of page references
type ReferenceGraph struct {
	GeneratedAt time.Time
	Nodes       map[string]*GraphNode // Page name -> Node
	HubPages    []string              // Most referenced pages (sorted by ref count)
}

// GraphNode represents a page in the reference graph
type GraphNode struct {
	PageName       string
	FilePath       string
	OutboundRefs   []string // Pages this page references
	InboundRefs    []string // Pages that reference this page
	ReferenceCount int      // Total inbound references (for ranking)
}

// BuildReferenceGraph creates a ReferenceGraph from page references and files
func BuildReferenceGraph(refs []models.PageReference, files []models.File) *ReferenceGraph {
	graph := &ReferenceGraph{
		GeneratedAt: time.Now(),
		Nodes:       make(map[string]*GraphNode),
	}

	// Create nodes for all files
	for _, file := range files {
		pageName := extractPageNameFromPath(file.Path)
		graph.Nodes[pageName] = &GraphNode{
			PageName:     pageName,
			FilePath:     file.Path,
			OutboundRefs: []string{},
			InboundRefs:  []string{},
		}
	}

	// Add references
	for _, ref := range refs {
		// Add outbound reference
		if node, exists := graph.Nodes[ref.SourcePage]; exists {
			// Avoid duplicates
			if !contains(node.OutboundRefs, ref.TargetPage) {
				node.OutboundRefs = append(node.OutboundRefs, ref.TargetPage)
			}
		}

		// Add inbound reference (even if target page doesn't exist yet)
		// This handles references to pages that haven't been created
		if _, exists := graph.Nodes[ref.TargetPage]; !exists {
			graph.Nodes[ref.TargetPage] = &GraphNode{
				PageName:     ref.TargetPage,
				FilePath:     "", // No file yet
				OutboundRefs: []string{},
				InboundRefs:  []string{},
			}
		}

		targetNode := graph.Nodes[ref.TargetPage]
		if !contains(targetNode.InboundRefs, ref.SourcePage) {
			targetNode.InboundRefs = append(targetNode.InboundRefs, ref.SourcePage)
			targetNode.ReferenceCount++
		}
	}

	// Identify hub pages
	graph.HubPages = findHubPages(graph.Nodes, 10) // Top 10

	return graph
}

// findHubPages returns the top N most referenced pages
func findHubPages(nodes map[string]*GraphNode, topN int) []string {
	// Create sorted list of nodes by reference count
	type nodeCount struct {
		pageName string
		count    int
	}

	var counts []nodeCount
	for pageName, node := range nodes {
		if node.ReferenceCount > 0 { // Only include pages with references
			counts = append(counts, nodeCount{pageName, node.ReferenceCount})
		}
	}

	// Sort by count descending
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].count > counts[j].count
	})

	// Take top N
	limit := topN
	if len(counts) < limit {
		limit = len(counts)
	}

	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		result[i] = counts[i].pageName
	}

	return result
}

// extractPageNameFromPath converts a file path to page name
func extractPageNameFromPath(filePath string) string {
	// Use the same logic as parser
	// This would be better imported from parser, but to avoid circular deps
	// we duplicate it here (or we could move it to a shared utils package)
	// For now, we'll use the basename logic

	// Find the last slash
	lastSlash := -1
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '/' || filePath[i] == '\\' {
			lastSlash = i
			break
		}
	}

	basename := filePath[lastSlash+1:]

	// Remove .md extension
	if len(basename) > 3 && basename[len(basename)-3:] == ".md" {
		return basename[:len(basename)-3]
	}

	return basename
}

// contains checks if a string slice contains a value
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// GetOrphanPages returns pages with no inbound or outbound references
func (rg *ReferenceGraph) GetOrphanPages() []string {
	var orphans []string

	for pageName, node := range rg.Nodes {
		if len(node.InboundRefs) == 0 && len(node.OutboundRefs) == 0 && node.FilePath != "" {
			orphans = append(orphans, pageName)
		}
	}

	sort.Strings(orphans)
	return orphans
}

// GetUnreferencedPages returns pages with no inbound references (but may have outbound)
func (rg *ReferenceGraph) GetUnreferencedPages() []string {
	var unreferenced []string

	for pageName, node := range rg.Nodes {
		if len(node.InboundRefs) == 0 && node.FilePath != "" {
			unreferenced = append(unreferenced, pageName)
		}
	}

	sort.Strings(unreferenced)
	return unreferenced
}
