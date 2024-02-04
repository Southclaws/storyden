package cluster_children

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

type Repository interface {
	Move(ctx context.Context, slug datagraph.ClusterSlug, parentSlug datagraph.ClusterSlug, opts ...Option) (*datagraph.Cluster, error)
}

type Option func(*options)

func MoveClusters() Option {
	return func(o *options) {
		o.moveClusters = true
	}
}

func MoveItems() Option {
	return func(o *options) {
		o.moveItems = true
	}
}
