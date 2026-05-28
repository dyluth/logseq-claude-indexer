package parser

import (
	"regexp"
	"time"
)

var ordinalRe = regexp.MustCompile(`(\d{1,2})(?:st|nd|rd|th)`)

// parseLogseqDate parses a Logseq date page ref like "Jun 3rd, 2026" or "January 15th, 2026".
func parseLogseqDate(s string) (time.Time, bool) {
	cleaned := ordinalRe.ReplaceAllString(s, "$1")
	for _, format := range []string{"Jan 2, 2006", "January 2, 2006"} {
		if t, err := time.Parse(format, cleaned); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
