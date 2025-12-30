package sortrule

import "strings"

type SortRule struct {
	field     string
	direction string
}

func Parse(raw string) SortRule {
	if raw == "" {
		return SortRule{field: "", direction: "asc"}
	}

	field := raw
	dir := "asc"

	if strings.HasPrefix(raw, "-") {
		dir = "desc"
		field = raw[1:]
	}

	return SortRule{
		field:     field,
		direction: dir,
	}
}

func (s SortRule) Field() string {
	return s.field
}

func (s SortRule) Direction() string {
	return s.direction
}

func (s SortRule) IsDescending() bool {
	return s.direction == "desc"
}

func (s SortRule) IsAscending() bool {
	return s.direction == "asc"
}
