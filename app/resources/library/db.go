package library

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/link"
	"github.com/Southclaws/storyden/internal/ent/node"
)

type database struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &database{db}
}

func (d *database) Create(
	ctx context.Context,
	owner account.AccountID,
	name string,
	slug mark.Slug,
	opts ...Option,
) (*Node, error) {
	create := d.db.Node.Create()
	mutate := create.Mutation()

	mutate.SetOwnerID(xid.ID(owner))
	mutate.SetName(name)
	mutate.SetSlug(slug.String())

	for _, fn := range opts {
		fn(mutate)
	}

	col, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err,
				fctx.With(ctx),
				ftag.With(ftag.AlreadyExists),
				fmsg.WithDesc("already exists", "The node URL slug must be unique and the specified slug is already in use."),
			)
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return d.Get(ctx, QueryKey{mark.NewQueryKeyID(col.ID)})
}

func (d *database) List(ctx context.Context, filters ...Filter) ([]*Node, error) {
	q := d.db.Node.
		Query().
		WithOwner(func(aq *ent.AccountQuery) {
			aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
		}).
		WithAssets().
		WithLink(func(lq *ent.LinkQuery) {
			lq.WithAssets().Order(link.ByCreatedAt(sql.OrderDesc()))
		}).
		Order(node.ByUpdatedAt(sql.OrderDesc()), node.ByCreatedAt(sql.OrderDesc()))

	for _, fn := range filters {
		fn(q)
	}

	cols, err := q.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	all, err := dt.MapErr(cols, NodeFromModel)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return all, nil
}

func (d *database) Get(ctx context.Context, qk QueryKey) (*Node, error) {
	query := d.db.Node.Query()

	query.Where(qk.Predicate())

	query.
		WithOwner(func(aq *ent.AccountQuery) {
			aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
		}).
		WithAssets().
		WithLink(func(lq *ent.LinkQuery) {
			lq.WithAssets().Order(link.ByCreatedAt(sql.OrderDesc()))
		}).
		WithNodes(func(cq *ent.NodeQuery) {
			cq.
				WithAssets().
				WithOwner(func(aq *ent.AccountQuery) {
					aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
				}).
				Order(node.ByUpdatedAt(sql.OrderDesc()), node.ByCreatedAt(sql.OrderDesc()))
		}).
		WithParent(func(cq *ent.NodeQuery) {
			cq.
				WithAssets().
				WithOwner(func(aq *ent.AccountQuery) {
					aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
				})
		})

	col, err := query.Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := NodeFromModel(col)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r, nil
}

func (d *database) GetByID(ctx context.Context, id NodeID) (*Node, error) {
	col, err := d.db.Node.
		Query().
		Where(node.ID(xid.ID(id))).
		WithOwner(func(aq *ent.AccountQuery) {
			aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
		}).
		WithAssets().
		WithLink(func(lq *ent.LinkQuery) {
			lq.WithAssets().Order(link.ByCreatedAt(sql.OrderDesc()))
		}).
		WithNodes(func(cq *ent.NodeQuery) {
			cq.
				WithAssets().
				WithOwner(func(aq *ent.AccountQuery) {
					aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
				}).
				Order(node.ByUpdatedAt(sql.OrderDesc()), node.ByCreatedAt(sql.OrderDesc()))
		}).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := NodeFromModel(col)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r, nil
}

func (d *database) Update(ctx context.Context, qk QueryKey, opts ...Option) (*Node, error) {
	// NOTE: Should be a probe not a full read. Query is necessary because of
	// the Mark being used (id or slug) for updates. Cannot use UpdateOneID.
	pre, err := d.Get(ctx, qk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	update := d.db.Node.Update()
	update.Where(qk.Predicate())

	mutate := update.Mutation()
	for _, fn := range opts {
		fn(mutate)
	}

	err = update.Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	qk = QueryKey{mark.NewQueryKeyID(pre.Mark.ID())}

	return d.Get(ctx, qk)
}

func (d *database) Delete(ctx context.Context, qk QueryKey) error {
	delete := d.db.Node.Delete()

	delete.Where(qk.Predicate())

	_, err := delete.Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
