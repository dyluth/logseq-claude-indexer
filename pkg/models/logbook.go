package models

import "time"

// LogbookEntry represents a single time tracking entry from Logseq's CLOCK format
// Example: CLOCK: [2025-04-06 Sun 10:00:00]--[2025-04-06 Sun 12:00:00] =>  02:00:00
type LogbookEntry struct {
	Start    time.Time     // Start time of the clock entry
	End      time.Time     // End time of the clock entry
	Duration time.Duration // Calculated duration
}
