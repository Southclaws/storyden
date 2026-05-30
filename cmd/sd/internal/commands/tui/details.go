package tui

import (
	"fmt"
	"strings"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/render"
)

func nodeDetails(node openapi.NodeWithChildren) string {
	var lines []string
	lines = append(lines,
		"Name: "+node.Name,
		"Slug: "+node.Slug,
		"ID: "+node.Id,
		"Visibility: "+string(node.Visibility),
		"Owner: "+render.AuthorName(node.Owner),
		"Updated: "+render.FormatTime(node.UpdatedAt),
		fmt.Sprintf("Children: %d", len(node.Children)),
		fmt.Sprintf("Assets: %d", len(node.Assets)),
	)
	if node.PrimaryImage != nil {
		lines = append(lines, "Primary image: "+node.PrimaryImage.Id)
	}
	if node.Description != "" {
		lines = append(lines, "", strings.TrimSpace(string(node.Description)))
	}
	if node.Content != nil {
		content, err := render.NodeMarkdownString(&node)
		if err == nil {
			lines = append(lines, "", render.ClampLines(content, 12))
		}
	}

	return strings.Join(lines, "\n")
}

func threadDetails(thread openapi.ThreadReference) string {
	lines := []string{
		"Title: " + thread.Title,
		"Slug: " + thread.Slug,
		"ID: " + thread.Id,
		"Visibility: " + string(thread.Visibility),
		"Author: " + render.AuthorName(thread.Author),
		"Updated: " + render.FormatTime(thread.UpdatedAt),
		fmt.Sprintf("Replies: %d", thread.ReplyStatus.Replies),
	}
	if thread.Category != nil {
		lines = append(lines, "Category: "+thread.Category.Name)
	}
	if thread.Description != nil && *thread.Description != "" {
		lines = append(lines, "", strings.TrimSpace(string(*thread.Description)))
	}
	body, err := render.HTMLToMarkdown(string(thread.Body))
	if err == nil && body != "" {
		lines = append(lines, "", body)
	}

	return strings.Join(lines, "\n")
}
