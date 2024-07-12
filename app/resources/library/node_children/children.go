package node_children

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/library"
)

type Repository interface {
	Move(ctx context.Context, slug library.NodeSlug, parentSlug library.NodeSlug) (*library.Node, error)
}
