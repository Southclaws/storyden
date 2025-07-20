package post_search

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/react"
)

type database struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &database{db}
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

	posts, err := dt.MapErr(r, post.Map)
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

	posts, err := dt.MapErr(result, post.Map)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return posts, nil
}
