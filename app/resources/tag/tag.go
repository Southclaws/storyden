package tag

import (
	"github.com/Southclaws/dt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/internal/ent"
)

type Tag struct {
	tag_ref.Tag

	Items []datagraph.Item
}

type Tags []*Tag

func Map(in *ent.Tag, roleHydratorFn func(accID xid.ID) (held.Roles, error)) (*Tag, error) {
	postsEdge, err := in.Edges.PostsOrErr()
	if err != nil {
		return nil, err
	}

	nodesEdge, err := in.Edges.NodesOrErr()
	if err != nil {
		return nil, err
	}

	tag := tag_ref.Map(nil)(in)

	tag.ItemCount = len(postsEdge) + len(nodesEdge)

	posts, err := dt.MapErr(postsEdge, func(in *ent.Post) (*thread.Thread, error) {
		return thread.Map(in, roleHydratorFn)
	})
	if err != nil {
		return nil, err
	}

	nodes, err := dt.MapErr(nodesEdge, library.MapNode(true, nil, roleHydratorFn))
	if err != nil {
		return nil, err
	}

	items := make([]datagraph.Item, 0, len(posts)+len(nodes))
	for _, post := range posts {
		items = append(items, post)
	}
	for _, node := range nodes {
		items = append(items, node)
	}

	return &Tag{
		Tag:   *tag,
		Items: items,
	}, nil
}
