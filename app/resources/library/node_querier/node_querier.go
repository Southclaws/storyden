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
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/link"
	"github.com/Southclaws/storyden/internal/ent/node"
)

type Querier struct {
	db  *ent.Client
	raw *sqlx.DB
	aq  *account_querier.Querier
}

func New(db *ent.Client, raw *sqlx.DB, aq *account_querier.Querier) *Querier {
	return &Querier{db, raw, aq}
}

type options struct {
	sortChildrenBy    *ChildSortRule
	visibilityRules   bool
	requestingAccount *account.AccountID
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
		WithOwner(func(aq *ent.AccountQuery) {
			aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
		}).
		WithPrimaryImage(func(aq *ent.AssetQuery) {
			aq.WithParent()
		}).
		WithAssets().
		WithLink(func(lq *ent.LinkQuery) {
			lq.WithAssets().Order(link.ByCreatedAt(sql.OrderDesc()))
		}).
		WithParent(func(cq *ent.NodeQuery) {
			cq.
				WithAssets().
				WithOwner(func(aq *ent.AccountQuery) {
					aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
				})
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
			WithOwner(func(aq *ent.AccountQuery) {
				aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
			}).
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

	r, err := library.MapNode(true, propSchema.Map())(col)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r, nil
}

func (q *Querier) ListChildren(ctx context.Context, qk library.QueryKey, pp pagination.Parameters, opts ...Option) (*pagination.Result[*library.Node], error) {
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

	applyVisibilityRulesPredicate(query)

	query.WithNodes(func(cq *ent.NodeQuery) {
		applyVisibilityRulesPredicate(cq)

		cq.
			WithOwner(func(aq *ent.AccountQuery) {
				aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
			}).
			WithPrimaryImage(func(aq *ent.AssetQuery) {
				aq.WithParent()
			}).
			WithAssets().
			WithLink(func(lq *ent.LinkQuery) {
				lq.WithAssets().Order(link.ByCreatedAt(sql.OrderDesc()))
			}).
			WithTags().
			WithProperties().
			WithPropertySchema(func(psq *ent.PropertySchemaQuery) {
				psq.WithFields()
			}).
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

	nodes := col.Edges.Nodes

	if o.sortChildrenBy != nil && len(nodes) > 0 {
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

	rs, err := dt.MapErr(nodes, library.MapNode(false, propSchema.Map()))
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
		WithOwner(func(aq *ent.AccountQuery) {
			aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
		})

	col, err := query.Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := library.MapNode(true, nil)(col)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r, nil
}

func (q *Querier) getRequestingAccount(ctx context.Context, o *options) (opt.Optional[account.Account], error) {
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
