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
	}
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

	t.Run("pull_images_relative", func(t *testing.T) {
		check(t, Content{
			short: `Here are some cool photos.`,
			links: []string{},
			media: []string{
				"https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2Fcarters-halt.jpg&w=3840&q=75",
				"https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2F30.jpg&w=3840&q=75",
				"https://barney.is/_next/image?url=%2Fphotography%2Fcity-of-london%2Fboxes.jpg&w=2048&q=75",
			},
		})(NewRichTextWithOptions(`<h1>heading</h1>

<p>Here are some cool photos.</p>

<img src="/_next/image?url=%2Fphotography%2Fcity-of-london%2Fcarters-halt.jpg&w=3840&q=75" />
<img src="/_next/image?url=%2Fphotography%2Fcity-of-london%2F30.jpg&w=3840&q=75" />
<img src="/_next/image?url=%2Fphotography%2Fcity-of-london%2Fboxes.jpg&w=2048&q=75" />
`, WithBaseURL("https://barney.is")))
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

		a.Equal(`"\u003cbody\u003e\u003cp\u003ea\u003c/p\u003e\u003c/body\u003e"`, string(encoded))

		var parsed Content
		err = json.Unmarshal(encoded, &parsed)
		r.NoError(err)
		r.NotEmpty(parsed)

		a.Equal(original, parsed)
	})
}
