package thread_cache

import (
	"context"
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/kv"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/spanner"
)

type Cache struct {
	ins spanner.Instrumentation
	db  *ent.Client
}

func New(ins spanner.Builder, db *ent.Client) *Cache {
	return &Cache{
		ins: ins.Build(),
		db:  db,
	}
}

func (c *Cache) IsNotModified(ctx context.Context, cq cachecontrol.Query, id xid.ID) bool {
	ctx, span := c.ins.Instrument(ctx, kv.String("id", id.String()))
	defer span.End()

	notModified := cq.NotModified(func() *time.Time {
		r, err := c.db.Post.Query().Select(post.FieldUpdatedAt).Where(post.ID(id)).Only(ctx)
		if err != nil {
			return nil
		}
		return &r.UpdatedAt
	})

	return notModified
}
