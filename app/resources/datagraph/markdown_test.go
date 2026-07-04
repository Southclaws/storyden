package datagraph

import (
	"testing"

	"github.com/Southclaws/storyden/internal/utils"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRichTextFromMarkdown(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		fmd, err := NewRichTextFromMarkdown(`To start with data science, it is essential to begin with a practical and step-by-step approach. First, explore platforms like Kaggle, which provide datasets that are accessible for beginners. While Kaggle may appear daunting at first, you can choose beginner-friendly tutorials and datasets that interest you. Begin by downloading and inspecting these datasets to get familiar with their structure and content. Concurrently, work on crafting questions from the data to guide your exploration—this helps in developing a problem-solving mindset.

Consistency in practicing with data, asking for advice, and seeking support, such as shared links or files, are also key steps. Keep in mind that experience and understanding grow steadily through practice rather than seeking perfection right away.

References:
- sdr:thread/cto7n8ifunp55p1bujv0: Emphasized the importance of staying practical and using beginner tutorials and platforms like Kaggle.
- sdr:thread/cto7nm2funp55p1bujvg: Provided advice on starting with data, forming questions, and the value of consistent practice.
`)

		check(t, Content{
			short: `To start with data science, it is essential to begin with a practical and step-by-step approach. First, explore platforms like...`,
			links: []string{},
			media: []string{},
		})(fmd, err)

		assert.Equal(t, `<body><p>To start with data science, it is essential to begin with a practical and step-by-step approach. First, explore platforms like Kaggle, which provide datasets that are accessible for beginners. While Kaggle may appear daunting at first, you can choose beginner-friendly tutorials and datasets that interest you. Begin by downloading and inspecting these datasets to get familiar with their structure and content. Concurrently, work on crafting questions from the data to guide your exploration—this helps in developing a problem-solving mindset.</p>

<p>Consistency in practicing with data, asking for advice, and seeking support, such as shared links or files, are also key steps. Keep in mind that experience and understanding grow steadily through practice rather than seeking perfection right away.</p>

<p>References:</p>

<ul>
<li>sdr:thread/cto7n8ifunp55p1bujv0: Emphasized the importance of staying practical and using beginner tutorials and platforms like Kaggle.</li>
<li>sdr:thread/cto7nm2funp55p1bujvg: Provided advice on starting with data, forming questions, and the value of consistent practice.</li>
</ul>
</body>`, fmd.HTML())
	})

	t.Run("html_anchor_links", func(t *testing.T) {
		check(t, Content{
			short: `Check out this great resource and also this one.`,
			links: []string{"https://example.com/resource", "https://another.com/page"},
			media: []string{},
		})(NewRichTextFromMarkdown(`Check out <a href="https://example.com/resource">this great resource</a> and also <a href="https://another.com/page">this one</a>.`))
	})

	t.Run("html_anchor_sdr_links", func(t *testing.T) {
		r := require.New(t)
		a := assert.New(t)

		nodeID := utils.Must(xid.FromString("crk0gvqfunp7891n7ah0"))

		got, err := NewRichTextFromMarkdown(`See <a href="sdr:node/crk0gvqfunp7891n7ah0">this node</a> for details.`)
		r.NoError(err)

		refs := got.References()
		r.Len(refs, 1)
		a.Equal(KindNode, refs[0].Kind)
		a.Equal(nodeID, refs[0].ID)
	})

	t.Run("markdown_links", func(t *testing.T) {
		check(t, Content{
			short: `Visit the homepage and also the docs.`,
			links: []string{"https://storyden.org", "https://docs.storyden.org"},
			media: []string{},
		})(NewRichTextFromMarkdown(`Visit the [homepage](https://storyden.org) and also the [docs](https://docs.storyden.org).`))
	})

	t.Run("markdown_sdr_links", func(t *testing.T) {
		r := require.New(t)
		a := assert.New(t)

		nodeID := utils.Must(xid.FromString("crk0gvqfunp7891n7ah0"))

		got, err := NewRichTextFromMarkdown(`See [Test](sdr:node/crk0gvqfunp7891n7ah0) for more details.`)
		r.NoError(err)

		refs := got.References()
		r.Len(refs, 1)
		a.Equal(KindNode, refs[0].Kind)
		a.Equal(nodeID, refs[0].ID)
	})
}
