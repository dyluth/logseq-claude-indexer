package parser

import (
	"regexp"
	"strings"
	"time"

	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

// Logseq timestamp format: [2025-04-06 Sun 12:30:42]
const logseqTimeFormat = "2006-01-02 Mon 15:04:05"

var (
	// Match CLOCK: [timestamp]--[timestamp] => duration
	clockLineRegex = regexp.MustCompile(`CLOCK:\s*\[([^\]]+)\]--\[([^\]]+)\]\s*=>\s*(.+)`)
)

// ParseLogbook extracts time tracking entries from a :LOGBOOK: block
// Returns the logbook entries and the number of lines consumed
func ParseLogbook(lines []string, startIdx int) ([]models.LogbookEntry, int) {
	var entries []models.LogbookEntry
	linesConsumed := 0

	// Check if we're starting at a :LOGBOOK: line
	if startIdx >= len(lines) || !strings.Contains(lines[startIdx], ":LOGBOOK:") {
		return entries, 0
	}

	linesConsumed++ // Count the :LOGBOOK: line

	// Parse CLOCK entries until we hit :END:
	for i := startIdx + 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		linesConsumed++

		// Check for end of logbook
		if strings.Contains(line, ":END:") {
			break
		}

		// Try to parse as CLOCK entry
		if entry, ok := parseClockLine(line); ok {
			entries = append(entries, entry)
		}
	}

	return entries, linesConsumed
}

// parseClockLine parses a single CLOCK line
// Example: CLOCK: [2025-04-06 Sun 10:00:00]--[2025-04-06 Sun 12:00:00] =>  02:00:00
func parseClockLine(line string) (models.LogbookEntry, bool) {
	matches := clockLineRegex.FindStringSubmatch(line)
	if len(matches) != 4 {
		return models.LogbookEntry{}, false
	}

	startStr := matches[1]
	endStr := matches[2]
	durationStr := matches[3]

	// Parse start time
	start, err := time.Parse(logseqTimeFormat, startStr)
	if err != nil {
		return models.LogbookEntry{}, false
	}

	// Parse end time
	end, err := time.Parse(logseqTimeFormat, endStr)
	if err != nil {
		return models.LogbookEntry{}, false
	}

	// Parse duration (format: HH:MM:SS or DDD:HH:MM:SS for >24h)
	duration, err := parseLogseqDuration(durationStr)
	if err != nil {
		// If duration parsing fails, calculate from timestamps
		duration = end.Sub(start)
	}

	return models.LogbookEntry{
		Start:    start,
		End:      end,
		Duration: duration,
	}, true
}

// parseLogseqDuration parses Logseq's duration format
// Formats: "02:30:15" (2h30m15s) or "260:46:44" (260h46m44s)
func parseLogseqDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	parts := strings.Split(s, ":")

	if len(parts) != 3 {
		return 0, nil
	}

	var hours, minutes, seconds int
	_, err := parseIntegers(parts[0], &hours, parts[1], &minutes, parts[2], &seconds)
	if err != nil {
		return 0, err
	}

	return time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second, nil
}

// parseIntegers is a helper to parse multiple integer strings
func parseIntegers(pairs ...interface{}) ([]int, error) {
	var results []int
	for i := 0; i < len(pairs); i += 2 {
		str := pairs[i].(string)
		ptr := pairs[i+1].(*int)

		var val int
		_, err := parseIntString(str, &val)
		if err != nil {
			return nil, err
		}
		*ptr = val
		results = append(results, val)
	}
	return results, nil
}

// parseIntString parses a string to int, handling leading spaces
func parseIntString(s string, dest *int) (int, error) {
	s = strings.TrimSpace(s)
	n, err := parseSimpleInt(s)
	if err != nil {
		return 0, err
	}
	*dest = n
	return n, nil
}

// parseSimpleInt is a simple integer parser without using fmt.Sscanf
func parseSimpleInt(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}

	result := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, nil
		}
		result = result*10 + int(c-'0')
	}
	return result, nil
}
