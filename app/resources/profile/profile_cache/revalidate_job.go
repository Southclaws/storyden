package profile_cache

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func (c *Cache) subscribe(ctx context.Context, bus *pubsub.Bus) error {
	if bus == nil {
		return nil
	}

	if _, err := pubsub.Subscribe(ctx, bus, "profile_cache.touch_updated", func(ctx context.Context, evt *rpc.EventAccountUpdated) error {
		return c.Invalidate(ctx, xid.ID(evt.ID))
	}); err != nil {
		return err
	}

	return nil
}
