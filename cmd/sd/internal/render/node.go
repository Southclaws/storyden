package render

import (
	"fmt"
	"html"
	"io"
	"regexp"
	"strings"
	"unicode"

	"charm.land/lipgloss/v2"
	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/microcosm-cc/bluemonday"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
)

func NodeJSON(out io.Writer, node *openapi.NodeWithChildren) error {
	return output.JSON(out, node)
}

func NodeMarkdown(out io.Writer, node *openapi.NodeWithChildren) error {
	if output.IsTerminal(out) {
		view, err := NodeViewString(out, node)
		if err != nil {
			return err
		}

		fmt.Fprint(out, view)
		return nil
	}

	markdown, err := NodeMarkdownString(node)
	if err != nil {
		return err
	}

	fmt.Fprint(out, markdown)
	return nil
}

func NodeYAML(out io.Writer, node *openapi.NodeWithChildren) error {
	payload := yamlNode{
		Name:        string(node.Name),
		Slug:        string(node.Slug),
		Visibility:  string(node.Visibility),
		Owner:       AuthorName(node.Owner),
		Description: string(node.Description),
		Tags:        tagNames(node.Tags),
		Properties:  yamlProperties(node.Properties),
	}
	if node.Parent != nil {
		payload.Parent = string(node.Parent.Slug)
	}
	if node.Content != nil {
		payload.Content = string(*node.Content)
	}

	return output.YAML(out, payload)
}

func NodeViewString(out io.Writer, node *openapi.NodeWithChildren) (string, error) {
	styles := nodeViewStyles()
	sections := []string{styles.Title.Render(string(node.Name))}

	if node.Description != "" {
		sections = append(sections, styles.Description.Render(strings.TrimSpace(string(node.Description))))
	}

	details := []nodeField{
		{"Slug", string(node.Slug)},
		{"ID", string(node.Id)},
		{"Visibility", string(node.Visibility)},
		{"Owner", AuthorName(node.Owner)},
		{"Created", node.CreatedAt.Local().Format("2006-01-02 15:04")},
		{"Updated", node.UpdatedAt.Local().Format("2006-01-02 15:04")},
		{"Children", fmt.Sprintf("%d", len(node.Children))},
	}
	if node.Parent != nil {
		details = append(details, nodeField{"Parent", string(node.Parent.Slug)})
	}
	if node.HideChildTree {
		details = append(details, nodeField{"Child tree", "hidden"})
	}
	sections = append(sections, renderSection("Details", renderFieldList(details, styles), styles))

	if len(node.Tags) > 0 {
		lines := make([]string, 0, len(node.Tags))
		for _, tag := range node.Tags {
			line := styles.Code.Render(tag.Name)
			if tag.Colour != "" {
				line += " " + styles.Code.Render(string(tag.Colour))
			}
			lines = append(lines, line)
		}
		sections = append(sections, renderSection("Tags", renderBulletList(lines, styles), styles))
	}

	if len(node.Properties) > 0 {
		lines := make([]string, 0, len(node.Properties))
		for _, prop := range node.Properties {
			lines = append(lines, fmt.Sprintf(
				"%s %s: %s",
				styles.Label.Render(markdownInlineValue(string(prop.Name))),
				styles.Code.Render(markdownInlineValue(string(prop.Type))),
				styles.Value.Render(markdownInlineValue(displayValue(string(prop.Value)))),
			))
		}
		sections = append(sections, renderSection("Properties", renderBulletList(lines, styles), styles))
	}

	if node.PrimaryImage != nil || len(node.Assets) > 0 {
		lines := []string{}
		if node.PrimaryImage != nil {
			lines = append(lines, renderAssetLine(*node.PrimaryImage, "primary image", styles)...)
		}
		for _, asset := range node.Assets {
			if node.PrimaryImage != nil && asset.Id == node.PrimaryImage.Id {
				continue
			}
			lines = append(lines, renderAssetLine(asset, "", styles)...)
		}
		sections = append(sections, renderSection("Assets", renderBulletList(lines, styles), styles))
	}

	if len(node.Children) > 0 {
		lines := make([]string, 0, len(node.Children))
		for _, child := range node.Children {
			line := styles.Label.Render(string(child.Name)) + " " + styles.Code.Render(string(child.Slug))
			if child.Visibility != "" {
				line += " " + styles.Code.Render(string(child.Visibility))
			}
			if child.Description != "" {
				line += " - " + styles.Value.Render(strings.TrimSpace(string(child.Description)))
			}
			lines = append(lines, line)
		}
		sections = append(sections, renderSection("Children", renderBulletList(lines, styles), styles))
	}

	if node.Content != nil {
		contentMarkdown, err := HTMLToMarkdown(string(*node.Content))
		if err != nil {
			contentMarkdown = SanitizePlainText(string(*node.Content))
		}

		contentMarkdown = strings.TrimSpace(contentMarkdown)
		if contentMarkdown != "" {
			width := output.TerminalWidth(out, 80)
			if width > 10 {
				width -= 2
			}
			sections = append(sections, renderSection("Content", strings.TrimSpace(output.Markdown(contentMarkdown, out, width)), styles))
		}
	}

	return strings.Join(sections, "\n\n") + "\n", nil
}

type nodeField struct {
	Label string
	Value string
}

type nodeStyles struct {
	Title       lipgloss.Style
	Description lipgloss.Style
	Section     lipgloss.Style
	Label       lipgloss.Style
	Value       lipgloss.Style
	Muted       lipgloss.Style
	Code        lipgloss.Style
	Bullet      lipgloss.Style
}

func nodeViewStyles() nodeStyles {
	return nodeStyles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("228")).
			Background(lipgloss.Color("63")).
			Padding(0, 1),
		Description: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("245")).
			PaddingLeft(1),
		Section: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")),
		Label: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("252")),
		Value: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),
		Muted: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")),
		Code: lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")).
			Background(lipgloss.Color("236")).
			Padding(0, 1),
		Bullet: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")),
	}
}

func renderSection(title string, body string, styles nodeStyles) string {
	if strings.TrimSpace(body) == "" {
		return ""
	}

	return styles.Section.Render(title) + "\n\n" + body
}

func renderFieldList(fields []nodeField, styles nodeStyles) string {
	width := 0
	for _, field := range fields {
		width = max(width, lipgloss.Width(field.Label))
	}

	lines := make([]string, 0, len(fields))
	for _, field := range fields {
		label := styles.Label.Width(width).Render(field.Label)
		lines = append(lines, fmt.Sprintf("%s %s  %s", styles.Bullet.Render("•"), label, styles.Value.Render(markdownInlineValue(field.Value))))
	}

	return strings.Join(lines, "\n")
}

func renderBulletList(lines []string, styles nodeStyles) string {
	rendered := make([]string, 0, len(lines))
	for _, line := range lines {
		rendered = append(rendered, fmt.Sprintf("%s %s", styles.Bullet.Render("•"), line))
	}

	return strings.Join(rendered, "\n")
}

func renderAssetLine(asset openapi.Asset, note string, styles nodeStyles) []string {
	name := asset.Filename
	if note != "" {
		name += " (" + note + ")"
	}

	line := styles.Label.Render(name)
	if asset.MimeType != "" {
		line += " " + styles.Code.Render(asset.MimeType)
	}
	if dimensions := assetDimensions(asset); dimensions != "" {
		line += " " + styles.Code.Render(dimensions)
	}

	lines := []string{line}
	if asset.Path != "" {
		lines = append(lines, styles.Muted.Render("Path:")+" "+styles.Code.Render(asset.Path))
	}

	return lines
}

func NodeMarkdownString(node *openapi.NodeWithChildren) (string, error) {
	var b strings.Builder

	fmt.Fprintf(&b, "# %s\n\n", node.Name)

	if node.Description != "" {
		fmt.Fprintf(&b, "> %s\n\n", strings.TrimSpace(string(node.Description)))
	}

	fmt.Fprintln(&b, "## Details")
	fmt.Fprintln(&b)
	writeMarkdownField(&b, "Slug", string(node.Slug))
	writeMarkdownField(&b, "ID", string(node.Id))
	writeMarkdownField(&b, "Visibility", string(node.Visibility))
	writeMarkdownField(&b, "Owner", AuthorName(node.Owner))
	writeMarkdownField(&b, "Created", node.CreatedAt.Local().Format("2006-01-02 15:04"))
	writeMarkdownField(&b, "Updated", node.UpdatedAt.Local().Format("2006-01-02 15:04"))
	writeMarkdownField(&b, "Children", fmt.Sprintf("%d", len(node.Children)))
	if node.Parent != nil {
		writeMarkdownField(&b, "Parent", string(node.Parent.Slug))
	}
	if node.HideChildTree {
		writeMarkdownField(&b, "Child tree", "hidden")
	}
	fmt.Fprintln(&b)

	if len(node.Tags) > 0 {
		fmt.Fprint(&b, "## Tags\n\n")
		for _, tag := range node.Tags {
			fmt.Fprintf(&b, "- `%s`", tag.Name)
			if tag.Colour != "" {
				fmt.Fprintf(&b, " `%s`", tag.Colour)
			}
			fmt.Fprintln(&b)
		}
		fmt.Fprintln(&b)
	}

	if len(node.Properties) > 0 {
		fmt.Fprint(&b, "## Properties\n\n")
		for _, prop := range node.Properties {
			writeMarkdownProperty(&b, prop)
		}
		fmt.Fprintln(&b)
	}

	if node.PrimaryImage != nil || len(node.Assets) > 0 {
		fmt.Fprint(&b, "## Assets\n\n")
		if node.PrimaryImage != nil {
			writeAssetItem(&b, *node.PrimaryImage, "primary image")
		}
		for _, asset := range node.Assets {
			if node.PrimaryImage != nil && asset.Id == node.PrimaryImage.Id {
				continue
			}
			writeAssetItem(&b, asset, "")
		}
		fmt.Fprintln(&b)
	}

	if len(node.Children) > 0 {
		fmt.Fprint(&b, "## Children\n\n")
		for _, child := range node.Children {
			fmt.Fprintf(&b, "- **%s** `%s`", child.Name, child.Slug)
			if child.Visibility != "" {
				fmt.Fprintf(&b, " `%s`", child.Visibility)
			}
			if child.Description != "" {
				fmt.Fprintf(&b, " - %s", strings.TrimSpace(string(child.Description)))
			}
			fmt.Fprintln(&b)
		}
		fmt.Fprintln(&b)
	}

	if node.Content != nil {
		contentStr := string(*node.Content)
		if contentStr != "" {
			contentMarkdown, err := HTMLToMarkdown(contentStr)
			if err != nil {
				contentMarkdown = SanitizePlainText(contentStr)
			}

			contentMarkdown = strings.TrimSpace(contentMarkdown)
			if contentMarkdown != "" {
				fmt.Fprint(&b, "## Content\n\n")
				fmt.Fprintln(&b, contentMarkdown)
				fmt.Fprintln(&b)
			}
		}
	}

	return b.String(), nil
}

func writeMarkdownField(b *strings.Builder, label string, value string) {
	fmt.Fprintf(b, "- **%s:** %s\n", label, markdownInlineValue(value))
}

func writeMarkdownProperty(b *strings.Builder, prop openapi.Property) {
	fmt.Fprintf(
		b,
		"- **%s** `%s`: %s\n",
		markdownInlineValue(string(prop.Name)),
		markdownInlineValue(string(prop.Type)),
		markdownInlineValue(displayValue(string(prop.Value))),
	)
}

func markdownInlineValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}

	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")
	value = strings.Join(strings.Fields(value), " ")
	value = strings.ReplaceAll(value, `\`, `\\`)

	return value
}

func displayValue(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}

	return value
}

func writeAssetItem(b *strings.Builder, asset openapi.Asset, note string) {
	filename := asset.Filename
	if note != "" {
		filename += " (" + note + ")"
	}

	fmt.Fprintf(b, "- **%s**", filename)
	if asset.MimeType != "" {
		fmt.Fprintf(b, " `%s`", asset.MimeType)
	}
	if dimensions := assetDimensions(asset); dimensions != "" {
		fmt.Fprintf(b, " `%s`", dimensions)
	}
	if asset.Path != "" {
		fmt.Fprintf(b, "\n  - Path: `%s`", asset.Path)
	}
	fmt.Fprintln(b)
}

func assetDimensions(asset openapi.Asset) string {
	if asset.Width <= 0 || asset.Height <= 0 {
		return ""
	}

	return fmt.Sprintf("%.0fx%.0f", asset.Width, asset.Height)
}

func HTMLToMarkdown(htmlContent string) (string, error) {
	safeHTML := bluemonday.UGCPolicy().Sanitize(htmlContent)
	markdown, err := htmltomarkdown.ConvertString(safeHTML)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(markdown), nil
}

func AuthorName(author openapi.ProfileReference) string {
	if author.Name != "" {
		return author.Name
	}
	if author.Handle != "" {
		return "@" + author.Handle
	}

	return author.Id
}

func tagNames(tags openapi.TagReferenceList) []string {
	names := make([]string, len(tags))
	for i, tag := range tags {
		names[i] = tag.Name
	}

	return names
}

type yamlNode struct {
	Name        string         `yaml:"name"`
	Slug        string         `yaml:"slug"`
	Visibility  string         `yaml:"visibility"`
	Owner       string         `yaml:"owner"`
	Description string         `yaml:"description,omitempty"`
	Tags        []string       `yaml:"tags,omitempty"`
	Parent      string         `yaml:"parent,omitempty"`
	Properties  []yamlProperty `yaml:"properties,omitempty"`
	Content     string         `yaml:"content,omitempty"`
}

type yamlProperty struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

func yamlProperties(properties []openapi.Property) []yamlProperty {
	if len(properties) == 0 {
		return nil
	}

	out := make([]yamlProperty, len(properties))
	for i, prop := range properties {
		out[i] = yamlProperty{
			Name:  string(prop.Name),
			Type:  string(prop.Type),
			Value: string(prop.Value),
		}
	}

	return out
}

var ansiEscapePattern = regexp.MustCompile(`\x1b(?:\[[0-?]*[ -/]*[@-~]|\][^\a]*(?:\a|\x1b\\)|[@-Z\\-_])`)

func SanitizePlainText(content string) string {
	policy := bluemonday.StrictPolicy()
	plainText := policy.Sanitize(content)
	plainText = html.UnescapeString(plainText)
	plainText = ansiEscapePattern.ReplaceAllString(plainText, "")

	return strings.Map(func(r rune) rune {
		switch r {
		case '\n', '\r', '\t':
			return r
		}
		if unicode.IsControl(r) {
			return -1
		}

		return r
	}, plainText)
}
