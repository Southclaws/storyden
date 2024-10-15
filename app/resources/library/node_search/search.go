package node_search

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/jmoiron/sqlx"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/node"
)

type Search interface {
	Search(ctx context.Context, opts ...Option) ([]*library.Node, error)
}

type query struct {
	qs         string
	visibility []visibility.Visibility
}

type Option func(*query)

func WithNameContains(s string) Option {
	return func(q *query) {
		q.qs = s
	}
}

func WithVisibility(v []visibility.Visibility) Option {
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

func (s *service) Search(ctx context.Context, opts ...Option) ([]*library.Node, error) {
	q := &query{}

	for _, fn := range opts {
		fn(q)
	}

	query := s.db.Node.Query().
		Where(
			node.NameContainsFold(q.qs),
			// TODO: more query/filter params
		).
		WithOwner(func(aq *ent.AccountQuery) {
			aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
		}).
		WithNodes(func(cq *ent.NodeQuery) {
			cq.WithOwner(func(aq *ent.AccountQuery) {
				aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
			})
		}).
		WithPrimaryImage().
		WithAssets().
		Order(node.ByUpdatedAt(sql.OrderDesc()), node.ByCreatedAt(sql.OrderDesc()))

	r, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	nodes, err := dt.MapErr(r, library.NodeFromModel)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return nodes, nil
}
