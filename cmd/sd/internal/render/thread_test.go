package render

import (
	"bytes"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func TestThreadMarkdownIncludesThreadContextAndConvertedBody(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	markdown, err := ThreadMarkdownString(testThread(t))
	r.NoError(err)

	a.Contains(markdown, "# Launch Notes")
	a.Contains(markdown, "- **Slug:** thread_1-launch-notes")
	a.Contains(markdown, "- **Category:** Announcements")
	a.Contains(markdown, "## Tags")
	a.Contains(markdown, "`release`")
	a.Contains(markdown, "## Assets")
	a.Contains(markdown, "- **diagram.png** `image/png` `800x400`")
	a.Contains(markdown, "## Body")
	a.Contains(markdown, "**Storyden CLI**")
	a.Contains(markdown, "## Replies")
	a.Contains(markdown, "**@odin**")
}

func TestThreadViewRendersMetadataWithoutMarkdownTables(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	view, err := ThreadViewString(&bytes.Buffer{}, testThread(t))
	r.NoError(err)

	ansiEscapePattern := regexp.MustCompile(`\x1b(?:\[[0-?]*[ -/]*[@-~]|\][^\a]*(?:\a|\x1b\\)|[@-Z\\-_])`)
	plain := ansiEscapePattern.ReplaceAllString(view, "")
	a.Contains(plain, "Launch Notes")
	a.Contains(plain, "Details")
	a.Contains(plain, "thread_1-launch-notes")
	a.Contains(plain, "Body")
	a.Contains(plain, "Storyden CLI")
	a.NotContains(plain, "| --- |")
	a.NotContains(plain, "| Field | Value |")
}

func testThread(t *testing.T) *openapi.Thread {
	t.Helper()

	description := openapi.PostDescription("Release announcement")
	body := openapi.PostContent("<body><p><strong>Storyden CLI</strong> is ready.</p></body>")
	replyDescription := openapi.PostDescription("Looks good")

	return &openapi.Thread{
		Assets: openapi.AssetList{{
			Filename: "diagram.png",
			Id:       "asset_1",
			MimeType: "image/png",
			Path:     "/api/assets/diagram.png",
			Width:    800,
			Height:   400,
		}},
		Author:      openapi.ProfileReference{Handle: "southclaws"},
		Body:        body,
		Category:    &openapi.CategoryReference{Name: "Announcements", Slug: "announcements"},
		CreatedAt:   mustParseThreadTime(t, "2026-05-22T15:31:53+07:00"),
		Description: &description,
		Id:          "thread_1",
		Replies: openapi.PaginatedReplyList{
			Replies: openapi.ReplyList{{
				Author:      openapi.ProfileReference{Handle: "odin"},
				CreatedAt:   mustParseThreadTime(t, "2026-05-22T16:31:53+07:00"),
				Description: &replyDescription,
			}},
		},
		ReplyStatus: openapi.ReplyStatus{Replies: 1},
		Slug:        "thread_1-launch-notes",
		Tags:        openapi.TagReferenceList{{Name: "release", Colour: "#abcdef"}},
		Title:       "Launch Notes",
		UpdatedAt:   mustParseThreadTime(t, "2026-05-23T15:31:53+07:00"),
		Visibility:  openapi.Published,
	}
}

func mustParseThreadTime(t *testing.T, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse(time.RFC3339, value)
	require.NoError(t, err)

	return parsed
}
