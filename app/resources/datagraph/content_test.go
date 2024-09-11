package datagraph

import (
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
	}
}

func TestNewRichText(t *testing.T) {
	t.Run("simple_html", func(t *testing.T) {
		check(t, Content{
			short: `Here's a paragraph. It's pretty neat. Here's the rest of the text. neat photo right? This is quite a long post, the summary...`,
			links: []string{},
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
		})(NewRichText(`<h1>heading</h1>

<p>Here's a paragraph. It's pretty neat.</p>

<a href="https://ao.com/cooking/ovens">here are my favourite ovens</a>
<a href="https://tre.ee/trees/favs">here are my favourite trees</a>
`))
	})

	t.Run("with_uris", func(t *testing.T) {
		mention := utils.Must(xid.FromString("cn2h3gfljatbqvjqctdg"))

		check(t, Content{
			short: `hey @southclaws!`,
			links: []string{},
			sdrs: RefList{
				{Kind: KindProfile, ID: mention},
			},
		})(NewRichText(`<h1>heading</h1><p>hey <a href="sdr:profile/cn2h3gfljatbqvjqctdg">@southclaws</a>!</p>`))
	})
}
