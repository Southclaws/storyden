package post_querier

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
)

type Querier struct {
	db *ent.Client
}

func New(db *ent.Client) *Querier {
	return &Querier{db: db}
}

func (q *Querier) Probe(ctx context.Context, id post.ID) (*post.PostRef, error) {
	p, err := q.db.Post.
		Query().
		Where(ent_post.IDEQ(xid.ID(id))).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return post.MapRef(p), nil
}
