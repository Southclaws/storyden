package thread_cache

import (
	"context"

	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
	"github.com/rs/xid"
)

func (c *Cache) subscribe(ctx context.Context, bus *pubsub.Bus) error {
	if bus == nil {
		return nil
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.touch_published", func(ctx context.Context, evt *rpc.EventThreadPublished) error {
		return c.Invalidate(ctx, xid.ID(evt.ID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.touch_updated", func(ctx context.Context, evt *rpc.EventThreadUpdated) error {
		return c.Invalidate(ctx, xid.ID(evt.ID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.touch_unpublished", func(ctx context.Context, evt *rpc.EventThreadUnpublished) error {
		return c.Invalidate(ctx, xid.ID(evt.ID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.drop_deleted", func(ctx context.Context, evt *rpc.EventThreadDeleted) error {
		return c.delete(ctx, xid.ID(evt.ID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.reply_created", func(ctx context.Context, evt *rpc.EventThreadReplyCreated) error {
		return c.Invalidate(ctx, xid.ID(evt.ThreadID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.reply_updated", func(ctx context.Context, evt *rpc.EventThreadReplyUpdated) error {
		return c.Invalidate(ctx, xid.ID(evt.ThreadID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.reply_deleted", func(ctx context.Context, evt *rpc.EventThreadReplyDeleted) error {
		return c.Invalidate(ctx, xid.ID(evt.ThreadID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.post_reacted", func(ctx context.Context, evt *rpc.EventPostReacted) error {
		return c.Invalidate(ctx, xid.ID(evt.RootPostID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.post_unreacted", func(ctx context.Context, evt *rpc.EventPostUnreacted) error {
		return c.Invalidate(ctx, xid.ID(evt.RootPostID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.post_liked", func(ctx context.Context, evt *rpc.EventPostLiked) error {
		return c.Invalidate(ctx, xid.ID(evt.RootPostID))
	}); err != nil {
		return err
	}

	if _, err := pubsub.Subscribe(ctx, bus, "thread_cache.post_unliked", func(ctx context.Context, evt *rpc.EventPostUnliked) error {
		return c.Invalidate(ctx, xid.ID(evt.RootPostID))
	}); err != nil {
		return err
	}

	return nil
}
