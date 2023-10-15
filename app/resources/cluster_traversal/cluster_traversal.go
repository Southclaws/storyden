package cluster_traversal

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

type Repository interface {
	Root(ctx context.Context) ([]*datagraph.Cluster, error)
	Subtree(ctx context.Context, id datagraph.ClusterID) ([]*datagraph.Cluster, error)
	FilterSubtree(ctx context.Context, id datagraph.ClusterID, filter string) ([]*datagraph.Cluster, error)
}
