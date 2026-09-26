package node_cache

import (
	"context"

	"github.com/rs/xid"
)

type Invalidation struct {
	ID           xid.ID
	Slug         string
	PreviousSlug string
}

func (c *Cache) InvalidateMany(ctx context.Context, changes []Invalidation) error {
	if len(changes) == 0 {
		return nil
	}

	values := make(map[string]string, len(changes)*3)
	for _, change := range changes {
		if change.PreviousSlug != "" && change.PreviousSlug != change.Slug {
			// An invalid timestamp makes retired aliases miss without a Delete per slug.
			values[c.cacheKey(change.PreviousSlug)] = ""
		}
	}

	value := c.clock().UTC().Format(storeTimeFmt)
	for _, change := range changes {
		values[c.cacheKey(change.ID.String())] = value
		if change.Slug != "" {
			values[c.cacheKey(change.Slug)] = value
		}
	}

	return c.store.SetMany(ctx, values, cacheTTL)
}
