package node_cache

import (
	"context"
	"time"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

const (
	cachePrefix  = "node:last-modified:"
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

func (c *Cache) Check(ctx context.Context, cq cachecontrol.Query, key string) (*cachecontrol.ETag, bool) {
	etag, notModified := cq.Check(func() *time.Time {
		return c.lastModified(ctx, key)
	})

	return etag, notModified
}

func (c *Cache) Store(ctx context.Context, key string, ts time.Time) error {
	return c.storeTimestamp(ctx, key, ts)
}

func (c *Cache) cacheKey(key string) string {
	return cachePrefix + key
}

func (c *Cache) lastModified(ctx context.Context, key string) *time.Time {
	if ts, ok := c.cached(ctx, key); ok {
		return ts
	}

	return nil
}

func (c *Cache) cached(ctx context.Context, key string) (*time.Time, bool) {
	val, err := c.store.Get(ctx, c.cacheKey(key))
	if err != nil {
		return nil, false
	}

	ts, err := time.Parse(storeTimeFmt, val)
	if err != nil {
		_ = c.store.Delete(ctx, c.cacheKey(key))
		return nil, false
	}

	return &ts, true
}

func (c *Cache) storeTimestamp(ctx context.Context, key string, ts time.Time) error {
	return c.store.Set(ctx, c.cacheKey(key), ts.UTC().Format(storeTimeFmt), cacheTTL)
}

func (c *Cache) Invalidate(ctx context.Context, key string) error {
	now := c.clock().UTC()

	return c.storeTimestamp(ctx, key, now)
}

func (c *Cache) delete(ctx context.Context, key string) error {
	return c.store.Delete(ctx, c.cacheKey(key))
}
