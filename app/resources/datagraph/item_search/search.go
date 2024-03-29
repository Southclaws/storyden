package item_search

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/jmoiron/sqlx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/item"
)

type Search interface {
	Search(ctx context.Context, opts ...Option) ([]*datagraph.Item, error)
}

type query struct {
	qs         string
	visibility []post.Visibility
}

type Option func(*query)

func WithNameContains(s string) Option {
	return func(q *query) {
		q.qs = s
	}
}

func WithVisibility(v []post.Visibility) Option {
	return func(q *query) {
		q.visibility = v
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

	query := s.db.Item.Query().
		Where(
			item.NameContainsFold(q.qs),
			// TODO: more query/filter params
		).
		WithOwner().
		WithClusters(func(cq *ent.ClusterQuery) {
			cq.WithOwner()
		}).
		WithAssets().
		Order(item.ByUpdatedAt(sql.OrderDesc()), item.ByCreatedAt(sql.OrderDesc()))

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
