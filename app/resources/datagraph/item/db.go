package item

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
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/cluster"
	"github.com/Southclaws/storyden/internal/ent/item"
	"github.com/Southclaws/storyden/internal/ent/link"
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
	slug string,
	desc string,
	opts ...Option,
) (*datagraph.Item, error) {
	create := d.db.Item.Create()
	mutate := create.Mutation()

	mutate.SetOwnerID(xid.ID(owner))
	mutate.SetName(name)
	mutate.SetSlug(slug)
	mutate.SetDescription(desc)

	for _, fn := range opts {
		fn(mutate)
	}

	item, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err,
				fctx.With(ctx),
				ftag.With(ftag.AlreadyExists),
				fmsg.WithDesc("already exists", "The item URL slug must be unique and the specified slug is already in use."),
			)
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return d.Get(ctx, datagraph.ItemSlug(item.Slug))
}

func (d *database) List(ctx context.Context, filters ...Filter) ([]*datagraph.Item, error) {
	q := d.db.Item.
		Query().
		WithOwner().
		WithAssets().
		WithLinks(func(lq *ent.LinkQuery) {
			lq.WithAssets().Order(link.ByCreatedAt(sql.OrderDesc()))
		}).
		Order(item.ByUpdatedAt(sql.OrderDesc()), item.ByCreatedAt(sql.OrderDesc()))

	for _, fn := range filters {
		fn(q)
	}

	cols, err := q.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	all, err := dt.MapErr(cols, datagraph.ItemFromModel)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return all, nil
}

func (d *database) Get(ctx context.Context, slug datagraph.ItemSlug) (*datagraph.Item, error) {
	item, err := d.db.Item.
		Query().
		Where(item.Slug(string(slug))).
		WithOwner().
		WithAssets().
		WithLinks(func(lq *ent.LinkQuery) {
			lq.WithAssets().Order(link.ByCreatedAt(sql.OrderDesc()))
		}).
		WithClusters(func(cq *ent.ClusterQuery) {
			cq.
				WithAssets().
				WithOwner().Order(cluster.ByUpdatedAt(sql.OrderDesc()), cluster.ByCreatedAt(sql.OrderDesc()))
		}).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := datagraph.ItemFromModel(item)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r, nil
}

func (d *database) GetByID(ctx context.Context, id datagraph.ItemID) (*datagraph.Item, error) {
	item, err := d.db.Item.
		Query().
		Where(item.ID(xid.ID(id))).
		WithOwner().
		WithAssets().
		WithLinks(func(lq *ent.LinkQuery) {
			lq.WithAssets().Order(link.ByCreatedAt(sql.OrderDesc()))
		}).
		WithClusters(func(cq *ent.ClusterQuery) {
			cq.
				WithAssets().
				WithOwner().Order(cluster.ByUpdatedAt(sql.OrderDesc()), cluster.ByCreatedAt(sql.OrderDesc()))
		}).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := datagraph.ItemFromModel(item)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r, nil
}

func (d *database) Update(ctx context.Context, id datagraph.ItemID, opts ...Option) (*datagraph.Item, error) {
	create := d.db.Item.UpdateOneID(xid.ID(id))
	mutate := create.Mutation()

	for _, fn := range opts {
		fn(mutate)
	}

	c, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return d.Get(ctx, datagraph.ItemSlug(c.Slug))
}

func (d *database) Delete(ctx context.Context, slug datagraph.ItemSlug) (*datagraph.Item, error) {
	_, err := d.db.Item.Delete().Where(item.Slug(string(slug))).Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return d.Get(ctx, slug)
}
