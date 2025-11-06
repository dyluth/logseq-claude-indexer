package parser

import (
	"regexp"
	"strings"

	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

// Task status markers in Logseq
var taskStatuses = []models.TaskStatus{
	models.StatusNOW,
	models.StatusLATER,
	models.StatusTODO,
	models.StatusDOING,
	models.StatusDONE,
}

// ParseTasks extracts all tasks from markdown content
func ParseTasks(content string, filePath string) ([]models.Task, error) {
	var tasks []models.Task
	lines := strings.Split(content, "\n")

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Check if line is a bullet point (task candidate)
		if !isTaskLine(line) {
			continue
		}

		// Check for task status marker
		status, found := extractTaskStatus(line)
		if !found {
			continue
		}

		// Extract priority marker (if present)
		priority := extractPriority(line)

		// Extract task description (everything after status and priority markers)
		description := extractTaskDescription(line, status, priority)

		// Extract page references
		pageRefs := ExtractPageReferences(line)

		// Create task
		task := models.Task{
			Status:      status,
			Priority:    priority,
			Description: description,
			PageRefs:    pageRefs,
			SourceFile:  filePath,
			LineNumber:  i + 1, // 1-indexed
		}

		// Check if next line starts a logbook
		if i+1 < len(lines) && strings.Contains(lines[i+1], ":LOGBOOK:") {
			logbook, consumed := ParseLogbook(lines, i+1)
			task.Logbook = logbook
			i += consumed // Skip past the logbook lines
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// isTaskLine checks if a line looks like a task (starts with bullet)
func isTaskLine(line string) bool {
	trimmed := strings.TrimSpace(line)

	// Check for common bullet markers
	return strings.HasPrefix(trimmed, "- ") ||
		strings.HasPrefix(trimmed, "* ") ||
		strings.HasPrefix(trimmed, "+ ")
}

// extractTaskStatus finds the task status marker in a line
func extractTaskStatus(line string) (models.TaskStatus, bool) {
	for _, status := range taskStatuses {
		// Look for "- NOW " or "- LATER " etc.
		marker := string(status) + " "
		if strings.Contains(line, marker) {
			return status, true
		}
	}
	return "", false
}

// extractPriority extracts the priority marker from a task line
// Matches [#A], [#B], [#C] patterns
func extractPriority(line string) models.Priority {
	re := regexp.MustCompile(`\[#([ABC])\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		return models.Priority(matches[1])
	}
	return models.PriorityNone
}

// extractTaskDescription extracts the text after the status and priority markers
func extractTaskDescription(line string, status models.TaskStatus, priority models.Priority) string {
	// Find the status marker
	marker := string(status) + " "
	idx := strings.Index(line, marker)
	if idx == -1 {
		return strings.TrimSpace(line)
	}

	// Get everything after the status marker
	description := line[idx+len(marker):]

	// Remove priority marker if present (e.g., "[#A] ")
	if priority != models.PriorityNone {
		priorityMarker := "[#" + string(priority) + "] "
		description = strings.Replace(description, priorityMarker, "", 1)
	}

	return strings.TrimSpace(description)
}
