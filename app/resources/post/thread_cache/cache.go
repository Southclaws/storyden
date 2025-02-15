package thread_cache

import (
	"context"
	"time"

	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/dt"
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
		// TODO: Be more clever about this. This query runs for every request
		// with a If-Modified-Since header, and in the worst case it will not
		// result in a cache hit resulting in a full query to the database as
		// well as this query (which is 2 + n queries where n is the post repo
		// query count, which is surprisingly high...) a good fix for this could
		// be to store the last_replied_at timestamp on the post itself and also
		// store this value in the cache so that conditional requests are fast.
		r, err := c.db.Debug().Post.
			Query().
			// Select(post.FieldUpdatedAt).
			WithPosts(func(pq *ent.PostQuery) {
				pq.Where(post.DeletedAtIsNil())
			}).
			Where(post.ID(id)).
			Only(ctx)
		if err != nil {
			return nil
		}
		dates := append(dt.Map(r.Edges.Posts, func(r *ent.Post) time.Time { return r.CreatedAt }), r.UpdatedAt)

		latest := lo.MaxBy(dates, func(a time.Time, b time.Time) bool { return a.After(b) })

		return &latest
	})

	return notModified
}
