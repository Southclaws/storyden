package category_cache

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func (c *Cache) subscribe(ctx context.Context, bus *pubsub.Bus) error {
	if bus == nil {
		return nil
	}

	if _, err := pubsub.Subscribe(ctx, bus, "category_cache.touch_updated", func(ctx context.Context, evt *message.EventCategoryUpdated) error {
		return c.touch(ctx, evt.Slug)
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "category_cache.drop_deleted", func(ctx context.Context, evt *message.EventCategoryDeleted) error {
		return c.invalidate(ctx, evt.Slug)
	}); err != nil {
		return err
	}

	return nil
}
