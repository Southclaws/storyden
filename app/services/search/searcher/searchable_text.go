package searcher

import (
	"strings"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

func SearchableText(item datagraph.Item) string {
	parts := []string{
		strings.TrimSpace(item.GetContent().Plaintext()),
		strings.TrimSpace(item.GetName()),
		strings.TrimSpace(item.GetDesc()),
	}

	if v, ok := item.(datagraph.WithTagNames); ok {
		parts = append(parts, v.GetTags()...)
	}

	var out []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		out = append(out, part)
	}

	return strings.Join(out, "\n")
}
