package get

import (
	"bytes"
	"regexp"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/render"
)

func TestSanitizePlainText(t *testing.T) {
	a := assert.New(t)

	got := render.SanitizePlainText("<b>Hello</b>\x1b[31m<script>alert(1)</script>\x00")

	a.Equal("Hello", got)
}

func TestNewDefaultsToMarkdown(t *testing.T) {
	command := New(nil)

	format := (*cobra.Command)(command).Flags().Lookup("format")
	require.NotNil(t, format)
	assert.Equal(t, formatMarkdown, format.DefValue)
}

func TestNodeMarkdownIncludesNodeContextAndConvertedContent(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	node := testNode(t)

	markdown, err := render.NodeMarkdownString(node)
	r.NoError(err)

	a.Contains(markdown, "# Documentation Hub")
	a.Contains(markdown, "- **Slug:** documentation-hub")
	a.NotContains(markdown, "Layout")
	a.Contains(markdown, "## Properties")
	a.Contains(markdown, "- **Status** `text`: active")
	a.Contains(markdown, "## Assets")
	a.Contains(markdown, "- **hero.png** `image/png` `1280x720`")
	a.Contains(markdown, "- Path: `/api/assets/hero.png`")
	a.Contains(markdown, "## Children")
	a.Contains(markdown, "**Child Page** `child-page`")
	a.Contains(markdown, "## Content")
	a.Contains(markdown, "Welcome")
	a.Contains(markdown, "**HTML**")
}

func TestNodeViewRendersMetadataWithoutMarkdownTables(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	view, err := render.NodeViewString(&bytes.Buffer{}, testNode(t))
	r.NoError(err)

	ansiEscapePattern := regexp.MustCompile(`\x1b(?:\[[0-?]*[ -/]*[@-~]|\][^\a]*(?:\a|\x1b\\)|[@-Z\\-_])`)
	plain := ansiEscapePattern.ReplaceAllString(view, "")
	a.Contains(plain, "Documentation Hub")
	a.Contains(plain, "Details")
	a.Contains(plain, "Slug")
	a.Contains(plain, "documentation-hub")
	a.Contains(plain, "Properties")
	a.Contains(plain, "Status")
	a.NotContains(plain, "| --- |")
	a.NotContains(plain, "| Field | Value |")
}

func testNode(t *testing.T) *openapi.NodeWithChildren {
	t.Helper()

	content := openapi.PostContent("<body><h1>Welcome</h1><p>This is <strong>HTML</strong> content.</p></body>")

	return &openapi.NodeWithChildren{
		Assets: openapi.AssetList{{
			Filename: "hero.png",
			Id:       "asset_1",
			MimeType: "image/png",
			Path:     "/api/assets/hero.png",
			Width:    1280,
			Height:   720,
		}},
		Children: []openapi.NodeWithChildren{{
			Name:        "Child Page",
			Slug:        "child-page",
			Visibility:  openapi.VisibilityPublished,
			Description: "A child page",
		}},
		Content:     &content,
		CreatedAt:   mustParseTime(t, "2026-05-22T15:31:53+07:00"),
		Description: "Documentation landing page",
		Id:          "node_1",
		Meta: openapi.Metadata{
			"layout": map[string]interface{}{
				"blocks": []interface{}{
					map[string]interface{}{"type": "cover"},
					map[string]interface{}{"type": "content"},
				},
			},
		},
		Name:       "Documentation Hub",
		Owner:      openapi.ProfileReference{Handle: "southclaws"},
		Properties: openapi.PropertyList{{Name: "Status", Type: openapi.PropertyTypeText, Value: "active"}},
		Slug:       "documentation-hub",
		Tags:       openapi.TagReferenceList{{Name: "docs", Colour: "#abcdef"}},
		UpdatedAt:  mustParseTime(t, "2026-05-29T10:56:07+07:00"),
		Visibility: openapi.VisibilityPublished,
	}
}

func TestRenderYAMLUsesStructuredEncoder(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	content := openapi.PostContent("body: with # characters\nand newlines")
	node := &openapi.NodeWithChildren{
		Name:        "name: with # characters",
		Slug:        "node-slug",
		Visibility:  openapi.VisibilityPublished,
		Owner:       openapi.ProfileReference{Name: "owner: with # characters"},
		Description: " description: with # characters\nand newlines ",
		Tags:        openapi.TagReferenceList{{Name: "tag: with # characters"}},
		Properties: openapi.PropertyList{{
			Name:  "property: with # characters",
			Type:  openapi.PropertyTypeText,
			Value: "value: with # characters",
		}},
		Content: &content,
	}
	var out bytes.Buffer

	r.NoError(render.NodeYAML(&out, node))

	var decoded map[string]any
	r.NoError(yaml.Unmarshal(out.Bytes(), &decoded))
	a.Equal("name: with # characters", decoded["name"])
	a.Equal("owner: with # characters", decoded["owner"])
	a.Equal("body: with # characters\nand newlines", decoded["content"])
}

func mustParseTime(t *testing.T, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse(time.RFC3339, value)
	require.NoError(t, err)

	return parsed
}
