package reply_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/asset"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
)

type Querier struct {
	db          *ent.Client
	roleQuerier *role_repo.Repository
}

func New(db *ent.Client, roleQuerier *role_repo.Repository) *Querier {
	return &Querier{db: db, roleQuerier: roleQuerier}
}

func roleHydrationTargets(posts ...*ent.Post) []*ent.Account {
	targets := make([]*ent.Account, 0)
	for _, p := range posts {
		if p == nil {
			continue
		}

		if p.Edges.Author != nil {
			targets = append(targets, p.Edges.Author)
		}

		if p.Edges.Root != nil && p.Edges.Root.Edges.Author != nil {
			targets = append(targets, p.Edges.Root.Edges.Author)
		}

		if p.Edges.ReplyTo != nil && p.Edges.ReplyTo.Edges.Author != nil {
			targets = append(targets, p.Edges.ReplyTo.Edges.Author)
		}
	}

	return targets
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

	roleHydrator, err := d.roleQuerier.BuildMultiHydrator(ctx, roleHydrationTargets(p))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return reply.Map(p, roleHydrator.Hydrate)
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
	roleHydrator, err := d.roleQuerier.BuildMultiHydrator(ctx, roleHydrationTargets(posts...))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	for _, p := range posts {
		r, err := reply.Map(p, roleHydrator.Hydrate)
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
