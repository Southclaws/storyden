package reply_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/asset"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
)

type Querier struct {
	db *ent.Client
}

func New(db *ent.Client) *Querier {
	return &Querier{db: db}
}

func (d *Querier) Get(ctx context.Context, id post.ID) (*reply.Reply, error) {
	p, err := d.db.Post.
		Query().
		Where(ent_post.IDEQ(xid.ID(id))).
		WithAuthor().
		WithRoot(func(pq *ent.PostQuery) {
			pq.WithAuthor()
		}).
		WithAssets(func(aq *ent.AssetQuery) {
			aq.Order(asset.ByUpdatedAt(), asset.ByCreatedAt())
		}).
		WithReplyTo(func(pq *ent.PostQuery) {
			pq.WithAuthor()
			pq.Where(ent_post.DeletedAtIsNil())
		}).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return reply.Map(p)
}

func (d *Querier) GetMany(ctx context.Context, ids ...post.ID) ([]*reply.Reply, error) {
	if len(ids) == 0 {
		return []*reply.Reply{}, nil
	}

	xids := dt.Map(ids, func(id post.ID) xid.ID { return xid.ID(id) })

	posts, err := d.db.Post.
		Query().
		Where(ent_post.IDIn(xids...)).
		WithAuthor().
		WithRoot(func(pq *ent.PostQuery) {
			pq.WithAuthor()
		}).
		WithAssets(func(aq *ent.AssetQuery) {
			aq.Order(asset.ByUpdatedAt(), asset.ByCreatedAt())
		}).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	replies := make([]*reply.Reply, 0, len(posts))
	for _, p := range posts {
		r, err := reply.Map(p)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		replies = append(replies, r)
	}

	return replies, nil
}

func (d *Querier) Probe(ctx context.Context, id post.ID) (*reply.ReplyRef, error) {
	p, err := d.db.Post.
		Query().
		Where(ent_post.IDEQ(xid.ID(id))).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return reply.MapRef(p), nil
}
