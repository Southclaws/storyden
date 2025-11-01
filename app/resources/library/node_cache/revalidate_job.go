package node_cache

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func (c *Cache) subscribe(ctx context.Context, bus *pubsub.Bus) error {
	if _, err := pubsub.Subscribe(ctx, bus, "node_cache.touch_created", func(ctx context.Context, evt *message.EventNodeCreated) error {
		return c.touch(ctx, eventKey(evt.Slug, xid.ID(evt.ID)))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "node_cache.touch_updated", func(ctx context.Context, evt *message.EventNodeUpdated) error {
		return c.touch(ctx, eventKey(evt.Slug, xid.ID(evt.ID)))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "node_cache.touch_published", func(ctx context.Context, evt *message.EventNodePublished) error {
		return c.touch(ctx, eventKey(evt.Slug, xid.ID(evt.ID)))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "node_cache.touch_submitted", func(ctx context.Context, evt *message.EventNodeSubmittedForReview) error {
		return c.touch(ctx, eventKey(evt.Slug, xid.ID(evt.ID)))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "node_cache.touch_unpublished", func(ctx context.Context, evt *message.EventNodeUnpublished) error {
		return c.touch(ctx, eventKey(evt.Slug, xid.ID(evt.ID)))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "node_cache.invalidate_deleted", func(ctx context.Context, evt *message.EventNodeDeleted) error {
		return c.invalidate(ctx, eventKey(evt.Slug, xid.ID(evt.ID)))
	}); err != nil {
		return err
	}

	return nil
}

func eventKey(markValue string, id xid.ID) string {
	if markValue != "" {
		return markValue
	}

	return mark.NewQueryKeyID(id).String()
}
