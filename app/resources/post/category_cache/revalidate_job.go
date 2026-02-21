package category_cache

import (
	"context"

	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func (c *Cache) subscribe(ctx context.Context, bus *pubsub.Bus) error {
	if bus == nil {
		return nil
	}

	if _, err := pubsub.Subscribe(ctx, bus, "category_cache.touch_updated", func(ctx context.Context, evt *rpc.EventCategoryUpdated) error {
		return c.Invalidate(ctx, evt.Slug)
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "category_cache.drop_deleted", func(ctx context.Context, evt *rpc.EventCategoryDeleted) error {
		return c.delete(ctx, evt.Slug)
	}); err != nil {
		return err
	}

	return nil
}
