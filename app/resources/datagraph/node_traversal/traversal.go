package node_traversal

import (
	"context"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/visibility"
)

type Repository interface {
	Root(ctx context.Context, opts ...Filter) ([]*datagraph.Node, error)
	Subtree(ctx context.Context, id opt.Optional[datagraph.NodeID], opts ...Filter) ([]*datagraph.Node, error)
	FilterSubtree(ctx context.Context, id datagraph.NodeID, filter string) ([]*datagraph.Node, error)
}

type filters struct {
	accountSlug *string
	visibility  []visibility.Visibility
	depth       *uint
}

type Filter func(*filters)

func WithOwner(v string) Filter {
	return func(f *filters) {
		f.accountSlug = &v
	}
}

func WithVisibility(v ...visibility.Visibility) Filter {
	return func(f *filters) {
		f.visibility = v
	}
}

func WithDepth(v uint) Filter {
	return func(f *filters) {
		f.depth = &v
	}
}
