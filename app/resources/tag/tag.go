package tag

import (
	"github.com/Southclaws/dt"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/internal/ent"
)

type Tag struct {
	tag_ref.Tag

	Posts []*thread.Thread
	Nodes []*library.Node
}

type Tags []*Tag

func Map(in *ent.Tag) (*Tag, error) {
	postsEdge, err := in.Edges.PostsOrErr()
	if err != nil {
		return nil, err
	}

	nodesEdge, err := in.Edges.NodesOrErr()
	if err != nil {
		return nil, err
	}

	tag := tag_ref.Map(in)

	posts, err := dt.MapErr(postsEdge, thread.FromModel(nil, nil))
	if err != nil {
		return nil, err
	}

	nodes, err := dt.MapErr(nodesEdge, library.NodeFromModel)
	if err != nil {
		return nil, err
	}

	return &Tag{
		Tag:   *tag,
		Posts: posts,
		Nodes: nodes,
	}, nil
}
