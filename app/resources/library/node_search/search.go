package node_search

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/jmoiron/sqlx"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/node"
)

type Search interface {
	Search(ctx context.Context, params pagination.Parameters, opts ...Option) (*pagination.Result[*library.Node], error)
}

type query struct {
	nameContains    string
	contentContains string
	visibility      []visibility.Visibility
}

type Option func(*query)

func WithNameContains(s string) Option {
	return func(q *query) {
		q.nameContains = s
	}
}

func WithContentContains(s string) Option {
	return func(q *query) {
		q.contentContains = s
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

func (s *service) Search(ctx context.Context, params pagination.Parameters, opts ...Option) (*pagination.Result[*library.Node], error) {
	q := &query{}

	for _, fn := range opts {
		fn(q)
	}

	predicate := node.And(
		node.Or(
			node.NameContainsFold(q.nameContains),
			node.ContentContainsFold(q.contentContains),
		),
		node.VisibilityEQ(node.VisibilityPublished),
		node.DeletedAtIsNil(),
	)

	total, err := s.db.Node.Query().Where(predicate).Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	query := s.db.Node.Query().
		Where(predicate).
		WithOwner().
		WithNodes(func(cq *ent.NodeQuery) {
			cq.WithOwner()
		}).
		WithPrimaryImage().
		Order(node.ByUpdatedAt(sql.OrderDesc()), node.ByCreatedAt(sql.OrderDesc())).
		Limit(params.Limit()).
		Offset(params.Offset())

	r, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	nodes, err := dt.MapErr(r, library.MapNode(true, nil))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := pagination.NewPageResult(params, total, nodes)

	return &result, nil
}
