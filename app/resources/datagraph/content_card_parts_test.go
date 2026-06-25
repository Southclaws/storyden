package datagraph

import (
	"testing"

	"github.com/Southclaws/storyden/internal/utils"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitCardParts(t *testing.T) {
	nodeID := utils.Must(xid.FromString("crk0gvqfunp7891n7ah0"))

	t.Run("no_sdr_links", func(t *testing.T) {
		r := require.New(t)

		c, err := NewRichText(`<body>
<p>Para 1.</p>
<p>Para 2.</p>
<p>Para 3.</p>
</body>`)
		r.NoError(err)

		parts := c.SplitCardParts()
		r.Len(parts, 1)
		assert.Contains(t, parts[0].HTML(), "Para 1")
		assert.Contains(t, parts[0].HTML(), "Para 3")
		assert.Empty(t, parts[0].References())
	})

	t.Run("sdr_in_middle", func(t *testing.T) {
		r := require.New(t)
		a := assert.New(t)

		c, err := NewRichText(`<body>
<p>Para 1.</p>
<p>Para 2.</p>
<p><a href="sdr:node/crk0gvqfunp7891n7ah0">A node</a></p>
<p>Para 4.</p>
<p>Para 5.</p>
<p>Para 6.</p>
</body>`)
		r.NoError(err)

		parts := c.SplitCardParts()
		r.Len(parts, 3)

		a.Contains(parts[0].HTML(), "Para 1")
		a.Contains(parts[0].HTML(), "Para 2")
		a.Empty(parts[0].References())

		a.Len(parts[1].References(), 1)
		a.Equal(KindNode, parts[1].References()[0].Kind)
		a.Equal(nodeID, parts[1].References()[0].ID)

		a.Contains(parts[2].HTML(), "Para 4")
		a.Contains(parts[2].HTML(), "Para 6")
		a.Empty(parts[2].References())
	})

	t.Run("sdr_at_start", func(t *testing.T) {
		r := require.New(t)
		a := assert.New(t)

		c, err := NewRichText(`<body>
<p><a href="sdr:node/crk0gvqfunp7891n7ah0">A node</a></p>
<p>Para 2.</p>
<p>Para 3.</p>
</body>`)
		r.NoError(err)

		parts := c.SplitCardParts()
		r.Len(parts, 2)

		a.Len(parts[0].References(), 1)
		a.Equal(KindNode, parts[0].References()[0].Kind)
		a.Equal(nodeID, parts[0].References()[0].ID)

		a.Contains(parts[1].HTML(), "Para 2")
		a.Contains(parts[1].HTML(), "Para 3")
	})

	t.Run("sdr_at_end", func(t *testing.T) {
		r := require.New(t)
		a := assert.New(t)

		c, err := NewRichText(`<body>
<p>Para 1.</p>
<p>Para 2.</p>
<p><a href="sdr:node/crk0gvqfunp7891n7ah0">A node</a></p>
</body>`)
		r.NoError(err)

		parts := c.SplitCardParts()
		r.Len(parts, 2)

		a.Contains(parts[0].HTML(), "Para 1")
		a.Contains(parts[0].HTML(), "Para 2")
		a.Empty(parts[0].References())

		a.Len(parts[1].References(), 1)
		a.Equal(KindNode, parts[1].References()[0].Kind)
		a.Equal(nodeID, parts[1].References()[0].ID)
	})

	t.Run("multiple_sdr_links", func(t *testing.T) {
		r := require.New(t)
		a := assert.New(t)

		threadID := utils.Must(xid.FromString("cn2h3gfljatbqvjqctdg"))

		c, err := NewRichText(`<body>
<p>Intro paragraph.</p>
<p><a href="sdr:node/crk0gvqfunp7891n7ah0">Node link</a></p>
<p>Middle paragraph.</p>
<p><a href="sdr:profile/cn2h3gfljatbqvjqctdg">Profile link</a></p>
<p>Outro paragraph.</p>
</body>`)
		r.NoError(err)

		parts := c.SplitCardParts()
		r.Len(parts, 5)

		a.Contains(parts[0].HTML(), "Intro")

		a.Len(parts[1].References(), 1)
		a.Equal(KindNode, parts[1].References()[0].Kind)
		a.Equal(nodeID, parts[1].References()[0].ID)

		a.Contains(parts[2].HTML(), "Middle")

		a.Len(parts[3].References(), 1)
		a.Equal(KindProfile, parts[3].References()[0].Kind)
		a.Equal(threadID, parts[3].References()[0].ID)

		a.Contains(parts[4].HTML(), "Outro")
	})

	t.Run("sdr_with_surrounding_text_in_para_is_not_split", func(t *testing.T) {
		r := require.New(t)

		c, err := NewRichText(`<body>
<p>See <a href="sdr:node/crk0gvqfunp7891n7ah0">this node</a> for details.</p>
<p>Another paragraph.</p>
</body>`)
		r.NoError(err)

		parts := c.SplitCardParts()
		r.Len(parts, 1)
		assert.Contains(t, parts[0].HTML(), "Another paragraph")
	})

	t.Run("empty_content", func(t *testing.T) {
		c := Content{}
		assert.Nil(t, c.SplitCardParts())
	})
}
