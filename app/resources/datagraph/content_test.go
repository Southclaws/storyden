package datagraph

import (
	"encoding/json"
	"testing"

	"github.com/Southclaws/storyden/internal/utils"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func check(t *testing.T, want Content) func(got Content, err error) {
	return func(got Content, err error) {
		require.NoError(t, err)
		assert.Equal(t, want.short, got.short)
		assert.Equal(t, want.links, got.links)
		assert.Equal(t, want.media, got.media)
		assert.Equal(t, want.sdrs, got.sdrs)
	}
}

// mustContent parses raw HTML via NewRichText.
func mustContent(t *testing.T, raw string) Content {
	t.Helper()
	c, err := NewRichText(raw)
	require.NoError(t, err)
	return c
}

func mustContentWithBlocks(t *testing.T, raw string) Content {
	t.Helper()
	c, err := NewRichTextWithBlocks(raw)
	require.NoError(t, err)
	return c.Content
}

func mustContentWithNewBlocks(t *testing.T, raw string) Content {
	t.Helper()
	c, err := NewRichText(raw)
	require.NoError(t, err)
	stable, err := NewRichTextWithNewBlocks(c)
	require.NoError(t, err)
	return stable.Content
}

func TestNewRichText(t *testing.T) {
	t.Run("simple_html", func(t *testing.T) {
		check(t, Content{
			short: `Here's a paragraph. It's pretty neat. Here's the rest of the text. neat photo right? This is quite a long post, the summary...`,
			links: []string{},
			media: []string{"http://image.com"},
		})(NewRichText(`<h1>heading</h1>

<p>Here's a paragraph. It's pretty neat.</p>

<p>Here's the rest of the text.</p>

<img src="http://image.com" />

<p>neat photo right?</p>

<p>This is quite a long post, the summary, should just be the first 128 characters rounded down to the nearest space.</p>`))
	})

	t.Run("pull_links", func(t *testing.T) {
		check(t, Content{
			short: `Here's a paragraph. It's pretty neat. here are my favourite ovens here are my favourite trees`,
			links: []string{"https://ao.com/cooking/ovens", "https://tre.ee/trees/favs"},
			media: []string{},
		})(NewRichText(`<h1>heading</h1>

<p>Here's a paragraph. It's pretty neat.</p>

<a href="https://ao.com/cooking/ovens">here are my favourite ovens</a>
<a href="https://tre.ee/trees/favs">here are my favourite trees</a>
`))
	})

	t.Run("pull_images", func(t *testing.T) {
		check(t, Content{
			short: `Here are some cool photos.`,
			links: []string{},
			media: []string{
				"https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2Fcarters-halt.jpg&w=3840&q=75",
				"https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2F30.jpg&w=3840&q=75",
				"https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2Fboxes.jpg&w=2048&q=75",
			},
		})(NewRichText(`<h1>heading</h1>

<p>Here are some cool photos.</p>

<img src="https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2Fcarters-halt.jpg&w=3840&q=75" />
<img src="https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2F30.jpg&w=3840&q=75" />
<img src="https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2Fboxes.jpg&w=2048&q=75" />
`))
	})

	t.Run("with_uris", func(t *testing.T) {
		mention := utils.Must(xid.FromString("cn2h3gfljatbqvjqctdg"))

		check(t, Content{
			short: `hey @southclaws!`,
			links: []string{},
			media: []string{},
			sdrs: RefList{
				{Kind: KindProfile, ID: mention},
			},
		})(NewRichText(`<h1>heading</h1><p>hey <a href="sdr:profile/cn2h3gfljatbqvjqctdg">@southclaws</a>!</p>`))
	})

	t.Run("json", func(t *testing.T) {
		r := require.New(t)
		a := assert.New(t)

		original, err := NewRichText(`<body><p>a</p></body>`)
		r.NoError(err)
		r.NotEmpty(original)

		encoded, err := json.Marshal(original)
		r.NoError(err)
		r.NotEmpty(encoded)

		// Not asserting the exact JSON string: block IDs are non-deterministic.

		var parsed Content
		err = json.Unmarshal(encoded, &parsed)
		r.NoError(err)
		r.NotEmpty(parsed)

		a.Equal(original, parsed)
	})
}

func TestNewRichTextDoesNotAssignBlockIDs(t *testing.T) {
	c := mustContent(t, `<p>Hello</p><p>World</p>`)

	require.Len(t, c.Blocks(), 2)
	assert.Empty(t, c.Blocks()[0].ID)
	assert.Empty(t, c.Blocks()[1].ID)
	assert.NotContains(t, c.HTML(), ` id=`)
}

func TestNewRichTextStripsIDAttributes(t *testing.T) {
	c := mustContent(t, `<p id="external">Hello</p><p id="sdb_cv3j6fld0aaaab44kk50">World</p>`)

	assert.Equal(t, `<body><p>Hello</p><p>World</p></body>`, c.HTML())
}

func TestNewRichTextWithBlocksPreservesOnlyStorydenIDs(t *testing.T) {
	c := mustContentWithBlocks(t, `<p id="external">Hello</p><p id="sdb_cv3j6fld0aaaab44kk50">World</p>`)

	assert.Equal(t, `<body><p>Hello</p><p id="sdb_cv3j6fld0aaaab44kk50">World</p></body>`, c.HTML())
}

func TestContentAccessors(t *testing.T) {
	c := mustContent(t, `<body><p>Hello <strong>world</strong>.</p><p><a href="https://example.com/docs">Docs</a></p><img src="https://example.com/image.jpg"></body>`)

	assert.Contains(t, c.HTML(), "Hello")
	assert.NotNil(t, c.HTMLTree())
	assert.Equal(t, "Hello world.Docs", c.Plaintext())
	assert.Equal(t, "Hello world.Docs", c.Short())
	assert.Equal(t, []string{"https://example.com/docs"}, c.Links())
	assert.Equal(t, []string{"https://example.com/image.jpg"}, c.Media())
	assert.Empty(t, c.References())
	assert.False(t, c.IsEmpty())
}

func TestContentImageOnlyIsNotEmpty(t *testing.T) {
	c := mustContent(t, `<img src="https://example.com/image.jpg" alt="Example image">`)

	assert.Empty(t, c.Plaintext())
	assert.Empty(t, c.Short())
	assert.Equal(t, []string{"https://example.com/image.jpg"}, c.Media())
	assert.False(t, c.IsEmpty())
}

func TestContentEmptyValue(t *testing.T) {
	var c Content

	assert.Equal(t, EmptyState, c.HTML())
	assert.Nil(t, c.HTMLTree())
	assert.Empty(t, c.Short())
	assert.Empty(t, c.Plaintext())
	assert.Empty(t, c.Links())
	assert.Empty(t, c.Media())
	assert.Empty(t, c.References())
	assert.True(t, c.IsEmpty())
}

func TestContentEmptyParagraphIsEmpty(t *testing.T) {
	c := mustContent(t, `<p></p>`)

	assert.Empty(t, c.Plaintext())
	assert.Empty(t, c.Short())
	assert.Empty(t, c.Links())
	assert.Empty(t, c.Media())
	assert.Empty(t, c.References())
	assert.True(t, c.IsEmpty())
}
