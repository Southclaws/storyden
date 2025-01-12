package profile_cache

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/account"
)

type Cache struct {
	db *ent.Client
}

func New(db *ent.Client) *Cache {
	return &Cache{
		db: db,
	}
}

func (c *Cache) IsNotModified(ctx context.Context, cq cachecontrol.Query, id xid.ID) bool {
	r, err := c.db.Account.Query().Select(account.FieldUpdatedAt).Where(account.ID(id)).Only(ctx)
	if err != nil {
		return false
	}

	notModified := cq.NotModified(r.UpdatedAt)

	return notModified
}
