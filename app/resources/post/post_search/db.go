package post_search

import (
	"context"
	"fmt"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/internal/ent"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/react"
)

type database struct {
	db          *ent.Client
	roleQuerier *role_repo.Repository
}

func New(db *ent.Client, roleQuerier *role_repo.Repository) Repository {
	return &database{db: db, roleQuerier: roleQuerier}
}

func roleHydrationTargets(posts []*ent.Post) []*ent.Account {
	targets := make([]*ent.Account, 0, len(posts))
	for _, p := range posts {
		if p == nil {
			continue
		}

		if p.Edges.Author != nil {
			targets = append(targets, p.Edges.Author)
		}

		for _, react := range p.Edges.Reacts {
			if react != nil && react.Edges.Account != nil {
				targets = append(targets, react.Edges.Account)
			}
		}
	}

	return targets
}

func (d *database) Search(ctx context.Context, params pagination.Parameters, filters ...Filter) (*pagination.Result[*post.Post], error) {
	predicate := ent_post.And(
		ent_post.VisibilityEQ(ent_post.VisibilityPublished),
		ent_post.DeletedAtIsNil(),
	)

	q := d.db.Post.
		Query().
		Where(predicate).
		WithAuthor().
		WithReacts(func(rq *ent.ReactQuery) {
			rq.WithAccount().Order(react.ByCreatedAt())
		}).
		WithTags().
		WithRoot().
		Order(ent.Asc(ent_post.FieldCreatedAt)).
		Limit(params.Limit()).
		Offset(params.Offset())

	countQuery := d.db.Post.Query().Where(predicate)

	for _, fn := range filters {
		fn(q)
		fn(countQuery)
	}

	total, err := countQuery.Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := q.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roleHydrator, err := d.roleQuerier.BuildMultiHydrator(ctx, roleHydrationTargets(r))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	posts, err := dt.MapErr(r, func(in *ent.Post) (*post.Post, error) {
		return post.Map(in, roleHydrator.Hydrate)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := pagination.NewPageResult(params, total, posts)

	return &result, nil
}

func (d *database) GetMany(ctx context.Context, ids ...post.ID) ([]*post.Post, error) {
	if len(ids) == 0 {
		return []*post.Post{}, nil
	}

	rawids := dt.Map(ids, func(in post.ID) xid.ID {
		return xid.ID(in)
	})

	q := d.db.Post.
		Query().
		Where(
			ent_post.IDIn(rawids...),
		).
		WithAuthor().
		WithReacts(func(rq *ent.ReactQuery) {
			rq.WithAccount().Order(react.ByCreatedAt())
		}).
		WithTags().
		WithRoot().
		Order(ent.Asc(ent_post.FieldCreatedAt))

	result, err := q.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roleHydrator, err := d.roleQuerier.BuildMultiHydrator(ctx, roleHydrationTargets(result))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	posts, err := dt.MapErr(result, func(in *ent.Post) (*post.Post, error) {
		return post.Map(in, roleHydrator.Hydrate)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return posts, nil
}

type Location struct {
	Slug     string
	Kind     datagraph.Kind
	Index    opt.Optional[int]
	Page     opt.Optional[int]
	Position opt.Optional[int]
}

func (d *database) Locate(ctx context.Context, id post.ID) (*Location, error) {
	// fetch either the post or its root post.
	p, err := d.db.Post.Query().
		Where(
			ent_post.ID(xid.ID(id)),
		).
		WithRoot(func(pq *ent.PostQuery) {
			pq.Select(ent_post.FieldSlug)
		}).
		Select(ent_post.FieldSlug, ent_post.FieldRootPostID).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	// if we hit the id not the root id, it's a thread.
	if p.Edges.Root == nil {
		return &Location{
			Slug: p.Slug,
			Kind: datagraph.KindThread,
		}, nil
	}

	count, err := d.db.Post.
		Query().
		Where(
			ent_post.RootPostIDEQ(*p.RootPostID),
			ent_post.DeletedAtIsNil(),
			ent_post.IDLTE(xid.ID(id)),
		).
		Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	index := count - 1
	page := (index / reply.RepliesPerPage) + 1
	position := index % reply.RepliesPerPage

	return &Location{
		Slug:     fmt.Sprintf("%s#%s", p.Edges.Root.Slug, id.String()),
		Kind:     datagraph.KindReply,
		Index:    opt.New(index),
		Page:     opt.New(page),
		Position: opt.New(position),
	}, nil
}
