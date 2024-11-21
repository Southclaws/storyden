package node_querier

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/link"
	"github.com/Southclaws/storyden/internal/ent/node"
)

type Querier struct {
	db *ent.Client
	aq *account_querier.Querier
}

func New(db *ent.Client, aq *account_querier.Querier) *Querier {
	return &Querier{db, aq}
}

type options struct {
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
		WithTags()

	applyVisibilityRulesPredicate(query)

	query.WithNodes(func(cq *ent.NodeQuery) {
		applyVisibilityRulesPredicate(cq)

		cq.
			WithAssets().
			WithOwner(func(aq *ent.AccountQuery) {
				aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
			}).
			Order(node.ByUpdatedAt(sql.OrderDesc()), node.ByCreatedAt(sql.OrderDesc()))
	})

	col, err := query.Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := library.NodeFromModel(col)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r, nil
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

	r, err := library.NodeFromModel(col)
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
