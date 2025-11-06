package indexer

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

// TimelineDay represents activity on a specific date
type TimelineDay struct {
	Date         time.Time
	JournalPath  string
	TasksCreated []models.Task
	TimeLogged   time.Duration
	KeyActivity  []string // Summary bullets
}

// TimelineIndex organizes activity chronologically by date
type TimelineIndex struct {
	GeneratedAt time.Time
	Entries     []TimelineDay // Sorted newest first
}

// BuildTimelineIndex creates a timeline from tasks and files
func BuildTimelineIndex(tasks []models.Task, files []models.File) *TimelineIndex {
	index := &TimelineIndex{
		GeneratedAt: time.Now(),
	}

	// Group tasks by date (from journal file path or logbook)
	dayMap := make(map[string]*TimelineDay)

	// First, create entries for all journal files
	for _, file := range files {
		if file.Type == models.FileTypeJournal {
			date, err := extractDateFromJournalPath(file.Path)
			if err != nil {
				continue
			}

			dateKey := date.Format("2006-01-02")
			if _, exists := dayMap[dateKey]; !exists {
				dayMap[dateKey] = &TimelineDay{
					Date:        date,
					JournalPath: file.Path,
				}
			}
		}
	}

	// Add tasks to their respective days
	for _, task := range tasks {
		// Extract date from source file path
		date, err := extractDateFromJournalPath(task.SourceFile)
		if err != nil {
			// Skip tasks not in journal files
			continue
		}

		dateKey := date.Format("2006-01-02")
		if day, exists := dayMap[dateKey]; exists {
			day.TasksCreated = append(day.TasksCreated, task)
			day.TimeLogged += task.TotalDuration()
		}
	}

	// Generate key activity summaries
	for _, day := range dayMap {
		day.KeyActivity = generateKeyActivity(day)
	}

	// Convert map to sorted slice (newest first)
	for _, day := range dayMap {
		index.Entries = append(index.Entries, *day)
	}

	sort.Slice(index.Entries, func(i, j int) bool {
		return index.Entries[i].Date.After(index.Entries[j].Date)
	})

	return index
}

// extractDateFromJournalPath extracts date from journal file path
// journals/2025_11_06.md -> Nov 6, 2025
// journals/2025-11-06.md -> Nov 6, 2025
func extractDateFromJournalPath(path string) (time.Time, error) {
	filename := filepath.Base(path)
	filename = strings.TrimSuffix(filename, ".md")

	// Try underscore format: 2025_11_06
	if strings.Contains(filename, "_") {
		t, err := time.Parse("2006_01_02", filename)
		if err == nil {
			return t, nil
		}
	}

	// Try dash format: 2025-11-06
	if strings.Contains(filename, "-") {
		t, err := time.Parse("2006-01-02", filename)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, &time.ParseError{
		Layout:     "2006_01_02 or 2006-01-02",
		Value:      filename,
		LayoutElem: "journal filename",
	}
}

// generateKeyActivity creates summary bullets for a day
func generateKeyActivity(day *TimelineDay) []string {
	var activity []string

	if len(day.TasksCreated) == 0 {
		return activity
	}

	// Count by status
	statusCounts := make(map[models.TaskStatus]int)
	var highPriorityTasks []models.Task

	for _, task := range day.TasksCreated {
		statusCounts[task.Status]++
		if task.Priority == models.PriorityHigh {
			highPriorityTasks = append(highPriorityTasks, task)
		}
	}

	// Add status summary
	if statusCounts[models.StatusNOW] > 0 {
		activity = append(activity, formatTaskCount(statusCounts[models.StatusNOW], "NOW"))
	}
	if statusCounts[models.StatusDOING] > 0 {
		activity = append(activity, formatTaskCount(statusCounts[models.StatusDOING], "DOING"))
	}
	if statusCounts[models.StatusTODO] > 0 {
		activity = append(activity, formatTaskCount(statusCounts[models.StatusTODO], "TODO"))
	}
	if statusCounts[models.StatusLATER] > 0 {
		activity = append(activity, formatTaskCount(statusCounts[models.StatusLATER], "LATER"))
	}
	if statusCounts[models.StatusDONE] > 0 {
		activity = append(activity, formatTaskCount(statusCounts[models.StatusDONE], "DONE"))
	}

	// Highlight high priority tasks (up to 2)
	for i, task := range highPriorityTasks {
		if i >= 2 {
			break
		}
		desc := task.Description
		if len(desc) > 60 {
			desc = desc[:57] + "..."
		}
		activity = append(activity, "🔥 "+desc)
	}

	return activity
}

// formatTaskCount formats a count with proper pluralization
func formatTaskCount(count int, status models.TaskStatus) string {
	if count == 1 {
		return "1 " + string(status) + " task"
	}
	return fmt.Sprintf("%d %s tasks", count, status)
}
