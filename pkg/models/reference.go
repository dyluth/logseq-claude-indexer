package models

// PageReference represents a [[page link]] found in Logseq content
type PageReference struct {
	SourceFile string // File path containing the reference
	SourcePage string // Page name (derived from filename, without .md)
	TargetPage string // Referenced page name (content inside [[...]])
	LineNumber int    // Line number where reference appears (1-indexed)
	Context    string // Surrounding text for context (truncated to reasonable length)
}
