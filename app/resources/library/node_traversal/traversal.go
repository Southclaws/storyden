package node_traversal

import (
	"context"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/visibility"
)

type Repository interface {
	Root(ctx context.Context, opts ...Filter) ([]*library.Node, error)
	Subtree(ctx context.Context, id opt.Optional[library.NodeID], opts ...Filter) ([]*library.Node, error)
	FilterSubtree(ctx context.Context, id library.NodeID, filter string) ([]*library.Node, error)
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
