package item_search

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/jmoiron/sqlx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/item"
)

type Search interface {
	Search(ctx context.Context, opts ...Option) ([]*datagraph.Item, error)
}

type query struct {
	qs string
}

type Option func(*query)

func WithNameContains(s string) Option {
	return func(q *query) {
		q.qs = s
	}
}

type service struct {
	db  *ent.Client
	raw *sqlx.DB
}

func New(db *ent.Client, raw *sqlx.DB) Search {
	return &service{
		db:  db,
		raw: raw,
	}
}

func (s *service) Search(ctx context.Context, opts ...Option) ([]*datagraph.Item, error) {
	q := &query{}

	for _, fn := range opts {
		fn(q)
	}

	query := s.db.Item.Query().Where(
		item.NameContainsFold(q.qs),
		// TODO: more query/filter params
	).WithOwner()

	r, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	items, err := dt.MapErr(r, datagraph.ItemFromModel)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return items, nil
}
