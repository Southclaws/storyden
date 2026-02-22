package question

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role/role_querier"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/question"
	"github.com/rs/xid"
)

type Repository struct {
	db          *ent.Client
	roleQuerier *role_querier.Querier
}

func New(db *ent.Client, roleQuerier *role_querier.Querier) *Repository {
	return &Repository{db: db, roleQuerier: roleQuerier}
}

func (r *Repository) Store(ctx context.Context,
	query string,
	result datagraph.Content,
	accountID opt.Optional[account.AccountID],
	parentID opt.Optional[xid.ID],
) (*Question, error) {
	create := r.db.Question.Create()
	mutate := create.Mutation()

	slug := mark.Slugify(query)

	mutate.SetSlug(slug)
	mutate.SetQuery(query)
	mutate.SetResult(result.HTML())

	accountID.Call(func(id account.AccountID) {
		mutate.SetAccountID(xid.ID(id))
	})

	parentID.Call(func(id xid.ID) {
		mutate.SetParentID(xid.ID(id))
	})

	create.OnConflictColumns("slug").UpdateNewValues()

	res, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	q, err := r.db.Question.Query().
		Where(question.ID(res.ID)).
		WithAuthor().
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if q.Edges.Author != nil {
		if err := r.roleQuerier.HydrateRoleEdges(ctx, q.Edges.Author); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	return Map(q)
}

func (r *Repository) Get(ctx context.Context, id xid.ID) (*Question, error) {
	q, err := r.db.Question.Query().
		Where(question.ID(id)).
		WithAuthor().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	if q.Edges.Author != nil {
		if err := r.roleQuerier.HydrateRoleEdges(ctx, q.Edges.Author); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	return Map(q)
}

func (r *Repository) GetByQuerySlug(ctx context.Context, query string) (*Question, error) {
	slug := mark.Slugify(query)

	q, err := r.db.Question.Query().
		Where(question.Slug(slug)).
		WithAuthor().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	if q.Edges.Author != nil {
		if err := r.roleQuerier.HydrateRoleEdges(ctx, q.Edges.Author); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	return Map(q)
}
