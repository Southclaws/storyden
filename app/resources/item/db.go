package item

import (
	"context"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/item"
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
		WithOwner()

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
		WithClusters().
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

func (d *database) Archive(ctx context.Context, slug datagraph.ItemSlug) (*datagraph.Item, error) {
	update := d.db.Item.Update().Where(item.Slug(string(slug)))
	update.SetDeletedAt(time.Now())

	_, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return d.Get(ctx, slug)
}
