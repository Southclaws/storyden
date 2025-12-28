package thread_cache

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/rs/xid"
)

func (c *Cache) subscribe(ctx context.Context, bus *pubsub.Bus) error {
	if bus == nil {
		return nil
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.touch_published", func(ctx context.Context, evt *message.EventThreadPublished) error {
		return c.Invalidate(ctx, xid.ID(evt.ID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.touch_updated", func(ctx context.Context, evt *message.EventThreadUpdated) error {
		return c.Invalidate(ctx, xid.ID(evt.ID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.touch_unpublished", func(ctx context.Context, evt *message.EventThreadUnpublished) error {
		return c.Invalidate(ctx, xid.ID(evt.ID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.drop_deleted", func(ctx context.Context, evt *message.EventThreadDeleted) error {
		return c.delete(ctx, xid.ID(evt.ID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.reply_created", func(ctx context.Context, evt *message.EventThreadReplyCreated) error {
		return c.Invalidate(ctx, xid.ID(evt.ThreadID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.reply_updated", func(ctx context.Context, evt *message.EventThreadReplyUpdated) error {
		return c.Invalidate(ctx, xid.ID(evt.ThreadID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.reply_deleted", func(ctx context.Context, evt *message.EventThreadReplyDeleted) error {
		return c.Invalidate(ctx, xid.ID(evt.ThreadID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.post_reacted", func(ctx context.Context, evt *message.EventPostReacted) error {
		return c.Invalidate(ctx, xid.ID(evt.RootPostID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.post_unreacted", func(ctx context.Context, evt *message.EventPostUnreacted) error {
		return c.Invalidate(ctx, xid.ID(evt.RootPostID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.post_liked", func(ctx context.Context, evt *message.EventPostLiked) error {
		return c.Invalidate(ctx, xid.ID(evt.RootPostID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.post_unliked", func(ctx context.Context, evt *message.EventPostUnliked) error {
		return c.Invalidate(ctx, xid.ID(evt.RootPostID))
	}); err != nil {
		return err
	}

	return nil
}
