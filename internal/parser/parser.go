package parser

import (
	"regexp"
	"strings"
)

var (
	// Match [[page reference]] links
	pageRefRegex = regexp.MustCompile(`\[\[([^\]]+)\]\]`)
)

// ExtractPageReferences finds all [[page]] references in a line of text
func ExtractPageReferences(line string) []string {
	matches := pageRefRegex.FindAllStringSubmatch(line, -1)

	var refs []string
	for _, match := range matches {
		if len(match) > 1 {
			refs = append(refs, match[1])
		}
	}
	return refs
}

// ExtractContext returns a substring of the line for context, truncated to maxLen
func ExtractContext(line string, maxLen int) string {
	line = strings.TrimSpace(line)
	if len(line) <= maxLen {
		return line
	}
	return line[:maxLen] + "..."
}
