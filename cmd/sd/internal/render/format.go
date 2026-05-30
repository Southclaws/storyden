package render

import (
	"strings"
	"time"
)

func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Local().Format("2006-01-02 15:04")
}

func ClampLines(text string, limit int) string {
	if limit <= 0 {
		return ""
	}

	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) <= limit {
		return strings.Join(lines, "\n")
	}

	return strings.Join(append(lines[:limit], "..."), "\n")
}
