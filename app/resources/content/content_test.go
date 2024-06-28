package content

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func check(t *testing.T, want Rich) func(got Rich, err error) {
	return func(got Rich, err error) {
		assert.Equal(t, want.short, got.short)
		assert.Equal(t, want.links, got.links)
	}
}

func TestNewRichText(t *testing.T) {
	// NOTE: Not using table tests here for easy debugging of individual cases.

	t.Run("simple_html", func(t *testing.T) {
		check(t, Rich{
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
		check(t, Rich{
			short: `Here's a paragraph. It's pretty neat.`,
			links: []string{"https://ao.com/cooking/ovens", "https://tre.ee/trees/favs"},
		})(NewRichText(`<h1>heading</h1>

<p>Here's a paragraph. It's pretty neat.</p>

<a href="https://ao.com/cooking/ovens">here are my favourite ovens</a>
<a href="https://tre.ee/trees/favs">here are my favourite trees</a>
`))
	})
}
