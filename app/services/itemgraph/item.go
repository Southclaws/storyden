package itemgraph

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

type ItemManager interface {
	// Link adds an item to a cluster.
	Link(ctx context.Context, item datagraph.ItemSlug, cluster datagraph.ClusterSlug) (*datagraph.Item, error)

	// Sever removes an item from a cluster if it was a member. If it was not a
	// member, then the return value is (nil, false, nil).
	Sever(ctx context.Context, item datagraph.ItemSlug, cluster datagraph.ClusterSlug) (*datagraph.Item, bool, error)
}
