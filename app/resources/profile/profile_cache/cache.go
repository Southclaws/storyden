package profile_cache

import (
	"context"
	"time"

	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/kv"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/spanner"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

const (
	cachePrefix  = "profile:last-modified:"
	cacheTTL     = time.Hour * 6
	storeTimeFmt = time.RFC3339Nano
)

type Cache struct {
	ins   spanner.Instrumentation
	store cache.Store
	clock func() time.Time
}

func New(
	ctx context.Context,
	lc fx.Lifecycle,
	ins spanner.Builder,
	store cache.Store,
	bus *pubsub.Bus,
) *Cache {
	c := &Cache{
		ins:   ins.Build(),
		store: store,
		clock: time.Now,
	}

	register := func(hook fx.Hook) { lc.Append(hook) }
	register(fx.StartHook(func(hctx context.Context) error {
		return c.subscribe(hctx, bus)
	}))

	return c
}

func (c *Cache) IsNotModified(ctx context.Context, cq cachecontrol.Query, id xid.ID) bool {
	ctx, span := c.ins.Instrument(ctx, kv.String("id", id.String()))
	defer span.End()

	notModified := cq.NotModified(func() *time.Time {
		return c.lastModified(ctx, id)
	})

	return notModified
}

func (c *Cache) LastModified(ctx context.Context, id xid.ID) *time.Time {
	ctx, span := c.ins.Instrument(ctx, kv.String("id", id.String()))
	defer span.End()

	return c.lastModified(ctx, id)
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

func (c *Cache) touch(ctx context.Context, id xid.ID) error {
	return c.storeTimestamp(ctx, id, c.clock().UTC())
}
