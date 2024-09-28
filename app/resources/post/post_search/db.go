package post_search

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

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

func (d *database) Search(ctx context.Context, filters ...Filter) ([]*post.Post, error) {
	if len(filters) == 0 {
		return []*post.Post{}, nil
	}

	q := d.db.Post.
		Query().
		WithAuthor(func(aq *ent.AccountQuery) {
			aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
		}).
		WithReacts(func(rq *ent.ReactQuery) {
			rq.WithAccount(func(aq *ent.AccountQuery) {
				aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
			}).Order(react.ByCreatedAt())
		}).
		WithTags().
		WithRoot().
		Order(ent.Asc(ent_post.FieldCreatedAt))

	for _, fn := range filters {
		fn(q)
	}

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
		WithAuthor(func(aq *ent.AccountQuery) {
			aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
		}).
		WithReacts(func(rq *ent.ReactQuery) {
			rq.WithAccount(func(aq *ent.AccountQuery) {
				aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
			}).Order(react.ByCreatedAt())
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
