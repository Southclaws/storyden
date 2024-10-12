package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/ent"
)

func Test_fromEnt(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)

	in := []*ent.Setting{
		{ID: "Title", Value: "Storyden"},
		{ID: "Description", Value: "A forum for the modern age."},
		{ID: "Content", Value: "<body><h1>Welcome to Storyden</h1></body>"},
		{ID: "AccentColour", Value: "27482225"},
	}

	out, err := fromEnt(in)
	r.NoError(err)
	r.NotNil(out)

	a.Equal("Storyden", out.Title.value)
	a.Equal("A forum for the modern age.", out.Description.value)
	a.Equal("<body><h1>Welcome to Storyden</h1></body>", out.Content.value.HTML())
	a.Equal("27482225", out.AccentColour.value)
}
