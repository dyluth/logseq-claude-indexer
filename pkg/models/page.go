package models

import "time"

// FileType distinguishes between journal entries and regular pages
type FileType int

const (
	FileTypeJournal FileType = iota
	FileTypePage
)

func (ft FileType) String() string {
	switch ft {
	case FileTypeJournal:
		return "journal"
	case FileTypePage:
		return "page"
	default:
		return "unknown"
	}
}

// File represents a Logseq markdown file discovered by the scanner
type File struct {
	Path         string    // Relative path from repo root (e.g., "pages/Project.md")
	AbsolutePath string    // Absolute filesystem path
	Type         FileType  // Journal or Page
	ModTime      time.Time // Last modified timestamp
}
