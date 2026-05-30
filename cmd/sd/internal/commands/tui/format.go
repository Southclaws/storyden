package tui

import "strings"

func titleCase(value string) string {
	if value == "" {
		return ""
	}

	return strings.ToUpper(value[:1]) + value[1:]
}
