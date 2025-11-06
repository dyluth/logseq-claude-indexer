package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

// Scanner walks a Logseq repository and finds all markdown files
type Scanner struct {
	repoPath string
}

// New creates a new Scanner for the given repository path
func New(repoPath string) *Scanner {
	return &Scanner{repoPath: repoPath}
}

// Scan walks the repository and returns all markdown files from pages/ and journals/
func (s *Scanner) Scan() ([]models.File, error) {
	var files []models.File

	// Scan journals/
	journalFiles, err := s.scanDirectory("journals", models.FileTypeJournal)
	if err != nil {
		// If journals directory doesn't exist, that's okay - just skip it
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("scanning journals: %w", err)
		}
	} else {
		files = append(files, journalFiles...)
	}

	// Scan pages/
	pageFiles, err := s.scanDirectory("pages", models.FileTypePage)
	if err != nil {
		// If pages directory doesn't exist, that's okay - just skip it
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("scanning pages: %w", err)
		}
	} else {
		files = append(files, pageFiles...)
	}

	return files, nil
}

// scanDirectory walks a specific directory and finds all .md files
func (s *Scanner) scanDirectory(dir string, fileType models.FileType) ([]models.File, error) {
	dirPath := filepath.Join(s.repoPath, dir)

	// Check if directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, err
	}

	var files []models.File

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Log error but continue walking
			return nil
		}

		// Skip directories
		if d.IsDir() {
			// Skip hidden directories and known exclusions
			name := d.Name()
			if strings.HasPrefix(name, ".") || name == ".recycle" || name == "bak" {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip non-markdown files
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		// Skip hidden files
		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		// Get file info for ModTime
		info, err := d.Info()
		if err != nil {
			// Skip files we can't stat
			return nil
		}

		// Calculate relative path from repo root
		relPath, err := filepath.Rel(s.repoPath, path)
		if err != nil {
			// This shouldn't happen, but skip if it does
			return nil
		}

		files = append(files, models.File{
			Path:         relPath,
			AbsolutePath: path,
			Type:         fileType,
			ModTime:      info.ModTime(),
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking directory %s: %w", dir, err)
	}

	return files, nil
}
