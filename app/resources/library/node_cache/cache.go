package node_cache

import (
	"context"
	"time"

	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/app/resources/mark"
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

// CanonicalKey collapses ID and ID-slug query forms onto the stable node ID.
// Slug-only lookups remain keyed by slug because resolving their ID would
// require reading the node before checking its conditional request cache.
func CanonicalKey(key mark.Queryable) string {
	if id, ok := key.ID().Get(); ok {
		return id.String()
	}

	return key.String()
}

// Invalidate refreshes the stable ID and current slug validators, then removes
// any retired slug aliases. Supplying the previous slug is required when a
// mutation renames or deletes a node.
func (c *Cache) Invalidate(ctx context.Context, id xid.ID, currentSlug string, retiredSlugs ...string) error {
	keys := []string{id.String(), currentSlug}
	seen := make(map[string]struct{}, len(keys)+len(retiredSlugs))
	now := c.clock().UTC()
	value := now.Format(storeTimeFmt)
	values := make(map[string]string, len(keys))

	for _, key := range keys {
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		cacheKey := c.cacheKey(key)
		values[cacheKey] = value
	}

	if err := c.store.SetMany(ctx, values, cacheTTL); err != nil {
		return err
	}

	for _, slug := range retiredSlugs {
		if slug == "" {
			continue
		}
		if _, ok := seen[slug]; ok {
			continue
		}
		seen[slug] = struct{}{}
		if err := c.delete(ctx, slug); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cache) delete(ctx context.Context, key string) error {
	return c.store.Delete(ctx, c.cacheKey(key))
}

func (c *Cache) deleteNode(ctx context.Context, id xid.ID, slug string) error {
	if err := c.delete(ctx, id.String()); err != nil {
		return err
	}

	if slug == "" {
		return nil
	}

	return c.delete(ctx, slug)
}
