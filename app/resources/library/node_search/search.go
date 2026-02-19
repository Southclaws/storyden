package node_search

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account "github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/node"
	ent_tag "github.com/Southclaws/storyden/internal/ent/tag"
)

type Search interface {
	Search(ctx context.Context, params pagination.Parameters, opts ...Option) (*pagination.Result[*library.Node], error)
}

type query struct {
	nameContains    string
	contentContains string
	visibility      []visibility.Visibility
	authors         []account.AccountID
	tags            []tag_ref.Name
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

func WithAuthors(ids ...account.AccountID) Option {
	return func(q *query) {
		q.authors = ids
	}
}

func WithTags(names ...tag_ref.Name) Option {
	return func(q *query) {
		q.tags = names
	}
}

type service struct {
	db          *ent.Client
	raw         *sqlx.DB
	roleQuerier *role_repo.Repository
}

func New(db *ent.Client, raw *sqlx.DB, roleQuerier *role_repo.Repository) Search {
	return &service{
		db:          db,
		raw:         raw,
		roleQuerier: roleQuerier,
	}
}

func (s *service) Search(ctx context.Context, params pagination.Parameters, opts ...Option) (*pagination.Result[*library.Node], error) {
	q := &query{}

	for _, fn := range opts {
		fn(q)
	}

	baseQuery := s.db.Node.Query().Where(
		node.Or(
			node.NameContainsFold(q.nameContains),
			node.ContentContainsFold(q.contentContains),
		),
		node.VisibilityEQ(node.VisibilityPublished),
		node.DeletedAtIsNil(),
	)

	if len(q.authors) > 0 {
		authorIDs := dt.Map(q.authors, func(id account.AccountID) xid.ID {
			return xid.ID(id)
		})
		baseQuery = baseQuery.Where(node.HasOwnerWith(ent_account.IDIn(authorIDs...)))
	}

	if len(q.tags) > 0 {
		for _, tag := range q.tags {
			baseQuery = baseQuery.Where(node.HasTagsWith(ent_tag.NameEQ(tag.String())))
		}
	}

	total, err := baseQuery.Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	query := baseQuery.
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

	roleLookup, err := s.roleQuerier.BuildMultiHydrator(ctx, roleHydrationTargets(r))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	nodes, err := dt.MapErr(r, library.MapNode(true, nil, roleLookup.Hydrate))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := pagination.NewPageResult(params, total, nodes)

	return &result, nil
}

func roleHydrationTargets(nodes []*ent.Node) []*ent.Account {
	targets := map[xid.ID]*ent.Account{}
	visited := map[xid.ID]struct{}{}
	queue := append([]*ent.Node{}, nodes...)

	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]

		if n == nil {
			continue
		}
		if _, ok := visited[n.ID]; ok {
			continue
		}
		visited[n.ID] = struct{}{}

		if owner := n.Edges.Owner; owner != nil {
			targets[owner.ID] = owner
		}

		queue = append(queue, n.Edges.Nodes...)
	}

	out := make([]*ent.Account, 0, len(targets))
	for _, acc := range targets {
		out = append(out, acc)
	}

	return out
}
