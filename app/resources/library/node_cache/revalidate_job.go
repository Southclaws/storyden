package node_cache

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func (c *Cache) subscribe(ctx context.Context, bus *pubsub.Bus) error {
	if _, err := pubsub.Subscribe(ctx, bus, "node_cache.invalidate_deleted", func(ctx context.Context, evt *rpc.EventNodeDeleted) error {
		return c.deleteNode(ctx, xid.ID(evt.ID), evt.Slug)
	}); err != nil {
		return err
	}

	return nil
}
