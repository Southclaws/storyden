package profile_cache

import (
	"context"
	"time"

	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

const (
	cachePrefix  = "profile:last-modified:"
	cacheTTL     = time.Hour * 6
	storeTimeFmt = time.RFC3339Nano
)

type Cache struct {
	store cache.Store
	clock func() time.Time
}

func New(
	ctx context.Context,
	lc fx.Lifecycle,
	store cache.Store,
	bus *pubsub.Bus,
) *Cache {
	c := &Cache{
		store: store,
		clock: time.Now,
	}

	register := func(hook fx.Hook) { lc.Append(hook) }
	register(fx.StartHook(func(hctx context.Context) error {
		return c.subscribe(ctx, bus)
	}))

	return c
}

func (c *Cache) Check(ctx context.Context, cq cachecontrol.Query, id xid.ID) (*cachecontrol.ETag, bool) {
	etag, notModified := cq.Check(func() *time.Time {
		return c.lastModified(ctx, id)
	})

	return etag, notModified
}

func (c *Cache) Store(ctx context.Context, id xid.ID, ts time.Time) error {
	return c.storeTimestamp(ctx, id, ts)
}

func (c *Cache) Invalidate(ctx context.Context, id xid.ID) error {
	now := c.clock().UTC()

	err := c.storeTimestamp(ctx, id, now)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to invalidate cache for profile"))
	}

	return nil
}

func (c *Cache) cacheKey(id xid.ID) string {
	return cachePrefix + id.String()
}

func (c *Cache) lastModified(ctx context.Context, id xid.ID) *time.Time {
	if ts, ok := c.cached(ctx, id); ok {
		return ts
	}

	return nil
}

func (c *Cache) cached(ctx context.Context, id xid.ID) (*time.Time, bool) {
	val, err := c.store.Get(ctx, c.cacheKey(id))
	if err != nil {
		return nil, false
	}

	ts, err := time.Parse(storeTimeFmt, val)
	if err != nil {
		_ = c.store.Delete(ctx, c.cacheKey(id))
		return nil, false
	}

	return &ts, true
}

func (c *Cache) storeTimestamp(ctx context.Context, id xid.ID, ts time.Time) error {
	return c.store.Set(ctx, c.cacheKey(id), ts.UTC().Format(storeTimeFmt), cacheTTL)
}
