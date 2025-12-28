package category_cache

import (
	"context"
	"time"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

const (
	cachePrefix  = "category:last-modified:"
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

func (c *Cache) Check(ctx context.Context, cq cachecontrol.Query, slug string) (*cachecontrol.ETag, bool) {
	etag, notModified := cq.Check(func() *time.Time {
		return c.lastModified(ctx, slug)
	})

	return etag, notModified
}

func (c *Cache) Store(ctx context.Context, slug string, ts time.Time) error {
	return c.storeTimestamp(ctx, slug, ts)
}

func (c *Cache) cacheKey(slug string) string {
	return cachePrefix + slug
}

func (c *Cache) lastModified(ctx context.Context, slug string) *time.Time {
	if ts, ok := c.cached(ctx, slug); ok {
		return ts
	}

	return nil
}

func (c *Cache) cached(ctx context.Context, slug string) (*time.Time, bool) {
	val, err := c.store.Get(ctx, c.cacheKey(slug))
	if err != nil {
		return nil, false
	}

	ts, err := time.Parse(storeTimeFmt, val)
	if err != nil {
		_ = c.store.Delete(ctx, c.cacheKey(slug))
		return nil, false
	}

	return &ts, true
}

func (c *Cache) storeTimestamp(ctx context.Context, slug string, ts time.Time) error {
	return c.store.Set(ctx, c.cacheKey(slug), ts.UTC().Format(storeTimeFmt), cacheTTL)
}

func (c *Cache) Invalidate(ctx context.Context, slug string) error {
	now := c.clock().UTC()

	return c.storeTimestamp(ctx, slug, now)
}

func (c *Cache) delete(ctx context.Context, slug string) error {
	return c.store.Delete(ctx, c.cacheKey(slug))
}
