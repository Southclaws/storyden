package cluster_traversal

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
)

type Repository interface {
	Root(ctx context.Context, opts ...Filter) ([]*datagraph.Cluster, error)
	Subtree(ctx context.Context, id datagraph.ClusterID, opts ...Filter) ([]*datagraph.Cluster, error)
	FilterSubtree(ctx context.Context, id datagraph.ClusterID, filter string) ([]*datagraph.Cluster, error)
}

type filters struct {
	accountSlug *string
	visibility  []post.Visibility
}

type Filter func(*filters)

func WithOwner(v string) Filter {
	return func(f *filters) {
		f.accountSlug = &v
	}
}

func WithVisibility(v []post.Visibility) Filter {
	return func(f *filters) {
		f.visibility = v
	}
}
