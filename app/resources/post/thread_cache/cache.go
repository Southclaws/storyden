package thread_cache

import (
	"context"
	"errors"
	"time"

	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

const (
	cachePrefix  = "thread:last-modified:"
	cacheTTL     = time.Hour * 6
	storeTimeFmt = time.RFC3339Nano
)

type Cache struct {
	db    *ent.Client
	store cache.Store
	clock func() time.Time
}

func New(
	ctx context.Context,
	lc fx.Lifecycle,

	db *ent.Client,
	store cache.Store,
	bus *pubsub.Bus,
) *Cache {
	c := &Cache{
		db:    db,
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
	notModified := cq.NotModified(func() *time.Time {
		return c.lastModified(ctx, id)
	})

	return notModified
}

func (c *Cache) LastModified(ctx context.Context, id xid.ID) *time.Time {
	return c.lastModified(ctx, id)
}

func (c *Cache) Store(ctx context.Context, id xid.ID, ts time.Time) error {
	return c.storeTimestamp(ctx, id, ts)
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

func (c *Cache) invalidate(ctx context.Context, id xid.ID) error {
	return c.store.Delete(ctx, c.cacheKey(id))
}

// Replies do not have their own thread cache entry so we need to find the root.
func (c *Cache) touchForReply(ctx context.Context, id xid.ID) error {
	threadID, err := c.threadIDForReply(ctx, id)
	if err != nil {
		if errors.Is(err, errThreadNotFound) {
			return nil
		}
		return err
	}
	return c.touch(ctx, threadID)
}

var errThreadNotFound = fault.New("thread not found")

func (c *Cache) threadIDForReply(ctx context.Context, id xid.ID) (xid.ID, error) {
	postRow, err := c.db.Post.Query().
		Select(post.FieldID, post.FieldRootPostID).
		Where(post.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return xid.ID{}, errThreadNotFound
		}
		return xid.ID{}, err
	}

	if postRow.RootPostID == nil {
		return postRow.ID, nil
	}

	return *postRow.RootPostID, nil
}
