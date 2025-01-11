package thread_cache

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/post"
)

type Cache struct {
	db *ent.Client
}

func New(db *ent.Client) *Cache {
	return &Cache{
		db: db,
	}
}

func (c *Cache) IsNotModified(ctx context.Context, cq opt.Optional[cachecontrol.Query], id xid.ID) (bool, error) {
	query, ok := cq.Get()
	if !ok {
		return false, nil
	}

	r, err := c.db.Post.Query().Select(post.FieldUpdatedAt).Where(post.ID(id)).Only(ctx)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}

	notModified := query.NotModified(r.UpdatedAt)

	return notModified, nil
}
