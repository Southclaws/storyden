package thread_cache

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/rs/xid"
)

func (c *Cache) subscribe(ctx context.Context, bus *pubsub.Bus) error {
	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.drop_deleted", func(ctx context.Context, evt *message.EventThreadDeleted) error {
		return c.delete(ctx, xid.ID(evt.ID))
	}); err != nil {
		return err
	}

	return nil
}
