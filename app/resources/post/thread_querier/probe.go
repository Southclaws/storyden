package thread_querier

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
)

func (q *Querier) Probe(ctx context.Context, id post.ID) (*thread.ThreadRef, error) {
	p, err := q.db.Post.
		Query().
		Where(ent_post.IDEQ(xid.ID(id))).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return thread.MapRef(p), nil
}
