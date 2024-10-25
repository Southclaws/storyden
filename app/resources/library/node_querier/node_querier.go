package node_querier

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/link"
	"github.com/Southclaws/storyden/internal/ent/node"
)

type Querier struct {
	db *ent.Client
}

func New(db *ent.Client) *Querier {
	return &Querier{db}
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

	if o.visibilityRules {
		if o.requestingAccount == nil {
			query.Where(node.VisibilityEQ(node.VisibilityPublished))
		} else {
			query.Where(node.Or(
				node.AccountID(xid.ID(*o.requestingAccount)),
				node.VisibilityEQ(node.VisibilityPublished),
			))
		}
	}

	query.WithNodes(func(cq *ent.NodeQuery) {
		if o.visibilityRules {
			// Apply visibility rules:
			// - published nodes are visible to everyone
			// - non-published nodes are not visible to anyone except the owner
			if o.requestingAccount == nil {
				cq.Where(node.VisibilityEQ(node.VisibilityPublished))
			} else {
				cq.Where(node.Or(
					node.AccountID(xid.ID(*o.requestingAccount)),
					node.VisibilityEQ(node.VisibilityPublished),
				))
			}
		}

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
