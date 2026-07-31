package cligenconv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
)

func TestConvertNodeWithChildren(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	now := time.Now().UTC().Truncate(time.Second)
	src := openapi.NodeWithChildren{
		Id:            "node-1",
		Name:          "Example",
		Slug:          "example",
		Description:   "desc",
		HideChildTree: false,
		Visibility:    openapi.VisibilityPublished,
		Owner: openapi.ProfileReference{
			Id:     "acc-1",
			Handle: "barney",
			Name:   "Barney",
		},
		CreatedAt: now,
		Children: []openapi.NodeWithChildren{
			{
				Id:   "child-1",
				Name: "Child",
				Slug: "child",
			},
		},
	}

	out, err := Convert[cligen.NodeWithChildren](src)
	r.NoError(err)

	a.Equal("node-1", out.ID)
	a.Equal("Example", out.Name)
	a.Equal("example", out.Slug)
	a.Equal(cligen.VisibilityPublished, out.Visibility)
	a.Equal("barney", out.Owner.Handle)
	r.NotNil(out.CreatedAt)
	a.True(now.Equal(*out.CreatedAt))
	r.Len(out.Children, 1)
	a.Equal("child-1", out.Children[0].ID)
}
