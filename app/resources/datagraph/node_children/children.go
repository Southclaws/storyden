package node_children

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

type Repository interface {
	Move(ctx context.Context, slug datagraph.NodeSlug, parentSlug datagraph.NodeSlug, opts ...Option) (*datagraph.Node, error)
}

type Option func(*options)

func MoveNodes() Option {
	return func(o *options) {
		o.moveNodes = true
	}
}
