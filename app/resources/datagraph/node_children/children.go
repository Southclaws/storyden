package node_children

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

type Repository interface {
	Move(ctx context.Context, slug datagraph.NodeSlug, parentSlug datagraph.NodeSlug) (*datagraph.Node, error)
}
