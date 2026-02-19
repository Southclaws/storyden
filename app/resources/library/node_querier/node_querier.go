package node_querier

import (
	"context"
	"math"
	"slices"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/link"
	"github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/ent/tag"
)

type Querier struct {
	db          *ent.Client
	raw         *sqlx.DB
	aq          *account_querier.Querier
	roleQuerier *role_repo.Repository
}

func New(db *ent.Client, raw *sqlx.DB, aq *account_querier.Querier, roleQuerier *role_repo.Repository) *Querier {
	return &Querier{
		db:          db,
		raw:         raw,
		aq:          aq,
		roleQuerier: roleQuerier,
	}
}

type options struct {
	sortChildrenBy       *ChildSortRule
	searchChildrenBy     opt.Optional[string]
	filterChildrenByTags opt.Optional[[]tag_ref.Name]
	visibilityRules      bool
	requestingAccount    *account.AccountID
}

type Option func(*options)

// WithVisibilityRulesApplied ensures ownership and visibility rules are applied
// if not set the default behaviour is no rules applied, all nodes are returned.
func WithVisibilityRulesApplied(accountID *account.AccountID) Option {
	return func(o *options) {
		o.visibilityRules = true
		o.requestingAccount = accountID
	}
}

func WithSortChildrenBy(field ChildSortRule) Option {
	return func(o *options) {
		o.sortChildrenBy = &field
	}
}

func WithSearchChildren(q string) Option {
	return func(o *options) {
		o.searchChildrenBy = opt.New(q)
	}
}

func WithFilterChildrenByTags(tags ...tag_ref.Name) Option {
	return func(o *options) {
		o.filterChildrenByTags = opt.New(tags)
	}
}

const nodePropertiesQuery = `with
  sibling_properties as (
    select
      ps.id         schema_id,
      min(psf.id)   field_id,
      min(psf.name) name,
      min(psf.type) type,
      min(psf.sort) sort,
      'sibling' as source
    from
      nodes n
      left join nodes sn on sn.parent_node_id = n.parent_node_id
      inner join property_schemas ps on ps.id = sn.property_schema_id
      or ps.id = n.property_schema_id
      inner join property_schema_fields psf on psf.schema_id = ps.id
    where
      n.id = $1
    group by ps.id, psf.id
  ),
  child_properties as (
    select
      ps.id         schema_id,
      min(psf.id)   field_id,
      min(psf.name) name,
      min(psf.type) type,
      min(psf.sort) sort,
      'child' as source
    from
      nodes n
      inner join nodes cn on cn.parent_node_id = n.id
      inner join property_schemas ps on ps.id = cn.property_schema_id
      inner join property_schema_fields psf on psf.schema_id = ps.id
    where
      n.id = $1
    group by ps.id, psf.id
  )
select
  *
from
  sibling_properties
union all
select
  *
from
  child_properties
order by source desc, sort asc
`

func (q *Querier) Get(ctx context.Context, qk library.QueryKey, opts ...Option) (*library.Node, error) {
	query := q.db.Node.Query()

	query.Where(qk.Predicate())

	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	requestingAccount, err := q.getRequestingAccount(ctx, o)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	applyVisibilityRulesPredicate := func(nq *ent.NodeQuery) {
		if !o.visibilityRules {
			return
		}

		// Apply visibility rules:
		// - published nodes are visible to everyone
		// - non-published nodes are not visible to anyone except the owner
		if acc, ok := requestingAccount.Get(); ok {

			canViewInReview := acc.Roles.Permissions().HasAny(rbac.PermissionAdministrator, rbac.PermissionManageLibrary)

			if canViewInReview {
				nq.Where(node.Or(
					node.AccountID(xid.ID(*o.requestingAccount)),
					node.VisibilityIn(node.VisibilityPublished, node.VisibilityReview),
				))
			} else {
				nq.Where(node.Or(
					node.AccountID(xid.ID(*o.requestingAccount)),
					node.VisibilityEQ(node.VisibilityPublished),
				))
			}
		} else {
			nq.Where(node.VisibilityEQ(node.VisibilityPublished))
		}
	}

	query.
		WithOwner().
		WithPrimaryImage(func(aq *ent.AssetQuery) {
			aq.WithParent()
		}).
		WithAssets().
		WithLink(func(lq *ent.LinkQuery) {
			lq.
				WithAssets().
				WithFaviconImage().
				WithPrimaryImage().
				Order(link.ByCreatedAt(sql.OrderDesc()))
		}).
		WithParent(func(cq *ent.NodeQuery) {
			cq.
				WithAssets().
				WithOwner()
		}).
		WithTags().
		WithProperties().
		WithPropertySchema(func(psq *ent.PropertySchemaQuery) {
			psq.WithFields()
		})

	applyVisibilityRulesPredicate(query)

	query.WithNodes(func(cq *ent.NodeQuery) {
		applyVisibilityRulesPredicate(cq)

		cq.
			WithAssets().
			WithOwner().
			WithProperties().
			Order(node.BySort(sql.OrderAsc()))
	})

	col, err := query.Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	propSchema := library.PropertySchemaQueryRows{}
	err = q.raw.SelectContext(ctx, &propSchema, nodePropertiesQuery, col.ID.String())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if o.sortChildrenBy != nil && len(col.Edges.Nodes) > 0 {
		children := dt.Map(col.Edges.Nodes, func(n *ent.Node) string { return n.ID.String() })
		sortmap, err := q.sortedByPropertyValue(ctx, children, *o.sortChildrenBy)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		// For any children that lack values, insert a max int value
		// in order to force these nodes to sort to the end always.
		if len(sortmap) != len(col.Edges.Nodes) {
			for _, node := range col.Edges.Nodes {
				if _, ok := sortmap[node.ID]; !ok {
					sortmap[node.ID] = math.MaxInt
				}
			}
		}

		slices.SortFunc(col.Edges.Nodes, func(a, b *ent.Node) int {
			return sortmap[a.ID] - sortmap[b.ID]
		})
	}

	roleHydrator, err := q.roleQuerier.BuildMultiHydrator(ctx, roleHydrationTargetsFromNode(col))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := library.MapNode(true, propSchema.Map(), roleHydrator.Hydrate)(col)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r, nil
}

func (q *Querier) ListChildren(ctx context.Context, qk library.QueryKey, pp pagination.Parameters, opts ...Option) (*pagination.Result[*library.Node], error) {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	// We need to resolve the parent ID first, as the child query is too complex
	// to use the qk predicate on.
	parentID, err := q.db.Node.Query().Where(qk.Predicate()).OnlyID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	query := q.db.Node.Query().Where(node.ParentNodeID(parentID))

	// Load all relevant edges
	query.
		WithOwner().
		WithPrimaryImage(func(aq *ent.AssetQuery) {
			aq.WithParent()
		}).
		WithAssets().
		WithLink(func(lq *ent.LinkQuery) {
			lq.WithAssets().
				WithFaviconImage().
				WithPrimaryImage().
				Order(link.ByCreatedAt(sql.OrderDesc()))
		}).
		WithTags().
		WithProperties().
		WithPropertySchema(func(psq *ent.PropertySchemaQuery) {
			psq.WithFields()
		})

	// Apply filters
	o.filterChildrenByTags.Call(func(tags []tag_ref.Name) {
		tagNames := dt.Map(tags, func(t tag_ref.Name) string { return t.String() })
		query.Where(node.HasTagsWith(tag.NameIn(tagNames...)))
	})

	o.searchChildrenBy.Call(func(q string) {
		query.Where(node.NameContainsFold(q))
	})

	// Apply visibility rules
	requestingAccount, err := q.getRequestingAccount(ctx, o)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	applyVisibilityRulesPredicate := func(nq *ent.NodeQuery) {
		if !o.visibilityRules {
			return
		}

		// Apply visibility rules:
		// - published nodes are visible to everyone
		// - non-published nodes are not visible to anyone except the owner
		if acc, ok := requestingAccount.Get(); ok {

			canViewInReview := acc.Roles.Permissions().HasAny(rbac.PermissionAdministrator, rbac.PermissionManageLibrary)

			if canViewInReview {
				nq.Where(node.Or(
					node.AccountID(xid.ID(*o.requestingAccount)),
					node.VisibilityIn(node.VisibilityPublished, node.VisibilityReview),
				))
			} else {
				nq.Where(node.Or(
					node.AccountID(xid.ID(*o.requestingAccount)),
					node.VisibilityEQ(node.VisibilityPublished),
				))
			}
		} else {
			nq.Where(node.VisibilityEQ(node.VisibilityPublished))
		}
	}

	applyVisibilityRulesPredicate(query)

	// Apply child-sort rules if present, otherwise fall back to sort by the
	// lexorank sort key field in ascending order.
	if o.sortChildrenBy != nil {
		order := o.sortChildrenBy.OrderClause()
		if o.sortChildrenBy.Fixed {
			switch o.sortChildrenBy.Field {
			case "name":
				query.Order(node.ByName(order))
			case "description":
				query.Order(node.ByDescription(order))
			case "link":
				query.Order(node.ByLinkField("url", order))
			}
		} else {
			// This is vastly simpler due to the post-query sorting with no
			// pagination. If we did perform pagination for this API (and we
			// might have to in future) then this would require performing the
			// actual property query first with pagination parameters and then
			// using the output of that query to pull a fixed set of nodes.
			// I do not envy my future self or contributor who will do that.
			// For now, fall back to sorting by the lexorank sort key field.
			query.Order(node.BySort(order))
		}
	} else {
		query.Order(node.BySort(sql.OrderAsc()))
	}

	nodes, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	propSchema := library.PropertySchemaQueryRows{}
	err = q.raw.SelectContext(ctx, &propSchema, nodePropertiesQuery, parentID.String())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// For non-fixed field sorting predicate, apply the property value sorting.
	if o.sortChildrenBy != nil && !o.sortChildrenBy.Fixed && len(nodes) > 0 {
		// override with the func param, todo: refactor? idk
		o.sortChildrenBy.Page = pp

		children := dt.Map(nodes, func(n *ent.Node) string { return n.ID.String() })
		sortmap, err := q.sortedByPropertyValue(ctx, children, *o.sortChildrenBy)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		slices.SortFunc(nodes, func(a, b *ent.Node) int {
			return sortmap[a.ID] - sortmap[b.ID]
		})
	}

	roleHydrator, err := q.roleQuerier.BuildMultiHydrator(ctx, roleHydrationTargetsFromNodes(nodes))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rs, err := dt.MapErr(nodes, library.MapNode(false, propSchema.Map(), roleHydrator.Hydrate))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r := pagination.NewPageResult(pp, len(rs), rs)

	return &r, nil
}

// Probe does not pull edges, only the node itself, it's fast for quick checks.
// TODO: Provide a more slimmed-down invariant of Node struct for this purpose.
func (q *Querier) Probe(ctx context.Context, id library.NodeID) (*library.Node, error) {
	query := q.db.Node.
		Query().
		Where(node.ID(xid.ID(id))).
		WithOwner()

	col, err := query.Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roleHydrator, err := q.roleQuerier.BuildMultiHydrator(ctx, roleHydrationTargetsFromNode(col))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := library.MapNode(true, nil, roleHydrator.Hydrate)(col)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r, nil
}

// ProbeMany fetches multiple nodes without pulling edges, fast for batch checks.
func (q *Querier) ProbeMany(ctx context.Context, ids ...library.NodeID) ([]*library.Node, error) {
	if len(ids) == 0 {
		return []*library.Node{}, nil
	}

	xids := dt.Map(ids, func(id library.NodeID) xid.ID {
		return xid.ID(id)
	})

	query := q.db.Node.
		Query().
		Where(node.IDIn(xids...)).
		WithOwner()

	nodes, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roleHydrator, err := q.roleQuerier.BuildMultiHydrator(ctx, roleHydrationTargetsFromNodes(nodes))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := dt.MapErr(nodes, library.MapNode(false, nil, roleHydrator.Hydrate))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return result, nil
}

func (q *Querier) getRequestingAccount(ctx context.Context, o *options) (opt.Optional[account.AccountWithEdges], error) {
	if !o.visibilityRules {
		return nil, nil
	}
	if o.requestingAccount == nil {
		return nil, nil
	}

	acc, err := q.aq.GetByID(ctx, *o.requestingAccount)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return opt.New(*acc), nil
}

func roleHydrationTargetsFromNode(n *ent.Node) []*ent.Account {
	return roleHydrationTargetsFromNodes([]*ent.Node{n})
}

func roleHydrationTargetsFromNodes(nodes []*ent.Node) []*ent.Account {
	targets := map[xid.ID]*ent.Account{}
	seenNodes := map[xid.ID]struct{}{}
	stack := append([]*ent.Node{}, nodes...)

	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]
		if current == nil {
			continue
		}

		if _, seen := seenNodes[current.ID]; seen {
			continue
		}
		seenNodes[current.ID] = struct{}{}

		if owner := current.Edges.Owner; owner != nil {
			targets[owner.ID] = owner
		}

		if parent := current.Edges.Parent; parent != nil {
			stack = append(stack, parent)
		}

		for _, child := range current.Edges.Nodes {
			stack = append(stack, child)
		}
	}

	out := make([]*ent.Account, 0, len(targets))
	for _, account := range targets {
		out = append(out, account)
	}

	return out
}
