package node_cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
)

const (
	revisionKey     = "node:revision"
	validatorPrefix = "node:validator:"
	cacheTTL        = time.Hour * 6
)

// Cache ties per-lookup validators to one library revision so changes to a
// node also invalidate representations derived from it, including descendants.
type Cache struct {
	logger *slog.Logger
	store  cache.Store
}

func New(logger *slog.Logger, store cache.Store) *Cache {
	return &Cache{logger: logger, store: store}
}

func (c *Cache) Check(ctx context.Context, query cachecontrol.Query, key string) (*cachecontrol.ETag, bool) {
	stored, err := c.store.Get(ctx, c.validatorKey(key))

	revision, ok := c.revision(ctx)
	if !ok {
		return nil, false
	}

	etag := validator(revision, key)
	if err != nil || stored != etag.Value {
		return etag, false
	}

	return etag, query.MatchesETag(etag.String())
}

func (c *Cache) Prepare(ctx context.Context, key string) (*cachecontrol.ETag, error) {
	revision, ok := c.revision(ctx)
	if !ok {
		revision = xid.New().String()
		if err := c.store.Set(ctx, revisionKey, revision, cacheTTL); err != nil {
			return nil, err
		}
	}

	return validator(revision, key), nil
}

func (c *Cache) Store(ctx context.Context, key string, etag *cachecontrol.ETag) error {
	return c.store.Set(ctx, c.validatorKey(key), etag.Value, cacheTTL)
}

func (c *Cache) Invalidate(ctx context.Context) error {
	return c.store.Set(ctx, revisionKey, xid.New().String(), cacheTTL)
}

func (c *Cache) InvalidateAfterWrite(ctx context.Context) {
	if err := c.Invalidate(ctx); err != nil {
		c.logger.ErrorContext(ctx, "failed to invalidate node cache after write", slog.String("error", err.Error()))
	}
}

func (c *Cache) revision(ctx context.Context) (string, bool) {
	revision, err := c.store.Get(ctx, revisionKey)
	if err != nil {
		return "", false
	}

	return revision, true
}

func (c *Cache) validatorKey(key string) string {
	return validatorPrefix + key
}

func validator(revision, key string) *cachecontrol.ETag {
	hash := sha256.Sum256([]byte(key))

	return cachecontrol.NewETagValue("r-" + revision + "-" + hex.EncodeToString(hash[:]))
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
