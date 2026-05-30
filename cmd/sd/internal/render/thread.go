package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
)

func ThreadJSON(out io.Writer, thread *openapi.Thread) error {
	return output.JSON(out, thread)
}

func ThreadMarkdown(out io.Writer, thread *openapi.Thread) error {
	if output.IsTerminal(out) {
		view, err := ThreadViewString(out, thread)
		if err != nil {
			return err
		}

		fmt.Fprint(out, view)
		return nil
	}

	markdown, err := ThreadMarkdownString(thread)
	if err != nil {
		return err
	}

	fmt.Fprint(out, markdown)
	return nil
}

func ThreadYAML(out io.Writer, thread *openapi.Thread) error {
	payload := yamlThread{
		Title:       string(thread.Title),
		Slug:        string(thread.Slug),
		Visibility:  string(thread.Visibility),
		Author:      AuthorName(thread.Author),
		Description: stringPtrValue(thread.Description),
		Tags:        tagNames(thread.Tags),
		Replies:     thread.ReplyStatus.Replies,
		Assets:      yamlAssetNames(thread.Assets),
		Body:        string(thread.Body),
	}
	if thread.Category != nil {
		payload.Category = thread.Category.Name
	}
	if thread.Link != nil {
		payload.URL = string(thread.Link.Url)
	}

	return output.YAML(out, payload)
}

func ThreadViewString(out io.Writer, thread *openapi.Thread) (string, error) {
	styles := nodeViewStyles()
	sections := []string{styles.Title.Render(string(thread.Title))}

	if thread.Description != nil && strings.TrimSpace(string(*thread.Description)) != "" {
		sections = append(sections, styles.Description.Render(strings.TrimSpace(string(*thread.Description))))
	}

	details := []nodeField{
		{"Slug", string(thread.Slug)},
		{"ID", string(thread.Id)},
		{"Visibility", string(thread.Visibility)},
		{"Author", AuthorName(thread.Author)},
		{"Created", FormatTime(thread.CreatedAt)},
		{"Updated", FormatTime(thread.UpdatedAt)},
		{"Replies", fmt.Sprintf("%d", thread.ReplyStatus.Replies)},
	}
	if thread.Category != nil {
		details = append(details, nodeField{"Category", thread.Category.Name})
	}
	if thread.Link != nil {
		details = append(details, nodeField{"URL", string(thread.Link.Url)})
	}
	if thread.Pinned != 0 {
		details = append(details, nodeField{"Pinned", fmt.Sprintf("%d", thread.Pinned)})
	}
	sections = append(sections, renderSection("Details", renderFieldList(details, styles), styles))

	if len(thread.Tags) > 0 {
		lines := make([]string, 0, len(thread.Tags))
		for _, tag := range thread.Tags {
			line := styles.Code.Render(tag.Name)
			if tag.Colour != "" {
				line += " " + styles.Code.Render(string(tag.Colour))
			}
			lines = append(lines, line)
		}
		sections = append(sections, renderSection("Tags", renderBulletList(lines, styles), styles))
	}

	if len(thread.Assets) > 0 {
		lines := []string{}
		for _, asset := range thread.Assets {
			lines = append(lines, renderAssetLine(asset, "", styles)...)
		}
		sections = append(sections, renderSection("Assets", renderBulletList(lines, styles), styles))
	}

	body := renderPostBody(string(thread.Body), out)
	if body != "" {
		sections = append(sections, renderSection("Body", body, styles))
	}

	if len(thread.Replies.Replies) > 0 {
		lines := make([]string, 0, len(thread.Replies.Replies))
		for _, reply := range thread.Replies.Replies {
			line := styles.Label.Render(AuthorName(reply.Author))
			if !reply.CreatedAt.IsZero() {
				line += " " + styles.Muted.Render(FormatTime(reply.CreatedAt))
			}
			if reply.Description != nil && strings.TrimSpace(string(*reply.Description)) != "" {
				line += " - " + styles.Value.Render(strings.TrimSpace(string(*reply.Description)))
			}
			lines = append(lines, line)
		}
		sections = append(sections, renderSection("Replies", renderBulletList(lines, styles), styles))
	}

	return strings.Join(sections, "\n\n") + "\n", nil
}

func ThreadMarkdownString(thread *openapi.Thread) (string, error) {
	var b strings.Builder

	fmt.Fprintf(&b, "# %s\n\n", thread.Title)
	if thread.Description != nil && strings.TrimSpace(string(*thread.Description)) != "" {
		fmt.Fprintf(&b, "> %s\n\n", strings.TrimSpace(string(*thread.Description)))
	}

	fmt.Fprintln(&b, "## Details")
	fmt.Fprintln(&b)
	writeMarkdownField(&b, "Slug", string(thread.Slug))
	writeMarkdownField(&b, "ID", string(thread.Id))
	writeMarkdownField(&b, "Visibility", string(thread.Visibility))
	writeMarkdownField(&b, "Author", AuthorName(thread.Author))
	writeMarkdownField(&b, "Created", FormatTime(thread.CreatedAt))
	writeMarkdownField(&b, "Updated", FormatTime(thread.UpdatedAt))
	writeMarkdownField(&b, "Replies", fmt.Sprintf("%d", thread.ReplyStatus.Replies))
	if thread.Category != nil {
		writeMarkdownField(&b, "Category", thread.Category.Name)
	}
	if thread.Link != nil {
		writeMarkdownField(&b, "URL", string(thread.Link.Url))
	}
	if thread.Pinned != 0 {
		writeMarkdownField(&b, "Pinned", fmt.Sprintf("%d", thread.Pinned))
	}
	fmt.Fprintln(&b)

	if len(thread.Tags) > 0 {
		fmt.Fprint(&b, "## Tags\n\n")
		for _, tag := range thread.Tags {
			fmt.Fprintf(&b, "- `%s`", tag.Name)
			if tag.Colour != "" {
				fmt.Fprintf(&b, " `%s`", tag.Colour)
			}
			fmt.Fprintln(&b)
		}
		fmt.Fprintln(&b)
	}

	if len(thread.Assets) > 0 {
		fmt.Fprint(&b, "## Assets\n\n")
		for _, asset := range thread.Assets {
			writeAssetItem(&b, asset, "")
		}
		fmt.Fprintln(&b)
	}

	bodyMarkdown := postBodyMarkdown(string(thread.Body))
	if bodyMarkdown != "" {
		fmt.Fprint(&b, "## Body\n\n")
		fmt.Fprintln(&b, bodyMarkdown)
		fmt.Fprintln(&b)
	}

	if len(thread.Replies.Replies) > 0 {
		fmt.Fprint(&b, "## Replies\n\n")
		for _, reply := range thread.Replies.Replies {
			fmt.Fprintf(&b, "- **%s**", AuthorName(reply.Author))
			if !reply.CreatedAt.IsZero() {
				fmt.Fprintf(&b, " `%s`", FormatTime(reply.CreatedAt))
			}
			if reply.Description != nil && strings.TrimSpace(string(*reply.Description)) != "" {
				fmt.Fprintf(&b, ": %s", markdownInlineValue(string(*reply.Description)))
			}
			fmt.Fprintln(&b)
		}
		fmt.Fprintln(&b)
	}

	return b.String(), nil
}

func renderPostBody(body string, out io.Writer) string {
	bodyMarkdown := postBodyMarkdown(body)
	if bodyMarkdown == "" {
		return ""
	}

	width := output.TerminalWidth(out, 80)
	if width > 10 {
		width -= 2
	}

	return strings.TrimSpace(output.Markdown(bodyMarkdown, out, width))
}

func postBodyMarkdown(body string) string {
	if strings.TrimSpace(body) == "" {
		return ""
	}

	bodyMarkdown, err := HTMLToMarkdown(body)
	if err != nil {
		bodyMarkdown = SanitizePlainText(body)
	}

	return strings.TrimSpace(bodyMarkdown)
}

func stringPtrValue(value *openapi.PostDescription) string {
	if value == nil {
		return ""
	}

	return string(*value)
}

func yamlAssetNames(assets openapi.AssetList) []string {
	names := make([]string, 0, len(assets))
	for _, asset := range assets {
		names = append(names, asset.Filename)
	}

	return names
}

type yamlThread struct {
	Title       string   `yaml:"title"`
	Slug        string   `yaml:"slug"`
	Visibility  string   `yaml:"visibility"`
	Author      string   `yaml:"author"`
	Description string   `yaml:"description,omitempty"`
	Category    string   `yaml:"category,omitempty"`
	URL         string   `yaml:"url,omitempty"`
	Tags        []string `yaml:"tags,omitempty"`
	Replies     int      `yaml:"replies"`
	Assets      []string `yaml:"assets,omitempty"`
	Body        string   `yaml:"body,omitempty"`
}
