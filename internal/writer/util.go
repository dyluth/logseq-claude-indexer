package writer

import (
	"os"
	"strings"
)

// writeIfChanged writes newContent to path only when the content has
// meaningfully changed.  The "Generated:" / "**Generated**:" timestamp line is
// excluded from the comparison so that a timestamp-only difference does not
// trigger a write (and therefore does not dirty the git working tree).
//
// If the file does not yet exist the content is always written.
func writeIfChanged(path, newContent string) error {
	if existing, err := os.ReadFile(path); err == nil {
		if stripGeneratedLine(string(existing)) == stripGeneratedLine(newContent) {
			return nil // nothing of substance changed; leave file untouched
		}
	}
	return os.WriteFile(path, []byte(newContent), 0644)
}

// stripGeneratedLine removes every line whose non-whitespace content starts
// with "Generated:" or "**Generated**:" so that two versions of a file can be
// compared without being thrown off by the timestamp.
func stripGeneratedLine(content string) string {
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Generated:") || strings.HasPrefix(trimmed, "**Generated**:") {
			continue
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}
