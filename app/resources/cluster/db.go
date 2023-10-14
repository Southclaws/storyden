package cluster

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
	"github.com/Southclaws/storyden/internal/ent/cluster"
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
) (*datagraph.Cluster, error) {
	create := d.db.Cluster.Create()
	mutate := create.Mutation()

	mutate.SetOwnerID(xid.ID(owner))
	mutate.SetName(name)
	mutate.SetSlug(slug)
	mutate.SetDescription(desc)

	for _, fn := range opts {
		fn(mutate)
	}

	col, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err,
				fctx.With(ctx),
				ftag.With(ftag.AlreadyExists),
				fmsg.WithDesc("already exists", "The cluster URL slug must be unique and the specified slug is already in use."),
			)
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return d.Get(ctx, datagraph.ClusterSlug(col.Slug))
}

func (d *database) List(ctx context.Context, filters ...Filter) ([]*datagraph.Cluster, error) {
	q := d.db.Cluster.
		Query().
		WithOwner()

	for _, fn := range filters {
		fn(q)
	}

	cols, err := q.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	all, err := dt.MapErr(cols, datagraph.ClusterFromModel)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return all, nil
}

func (d *database) Get(ctx context.Context, slug datagraph.ClusterSlug) (*datagraph.Cluster, error) {
	col, err := d.db.Cluster.
		Query().
		Where(cluster.Slug(string(slug))).
		WithOwner().
		WithItems().
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := datagraph.ClusterFromModel(col)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r, nil
}

func (d *database) Update(ctx context.Context, id datagraph.ClusterID, opts ...Option) (*datagraph.Cluster, error) {
	create := d.db.Cluster.UpdateOneID(xid.ID(id))
	mutate := create.Mutation()

	for _, fn := range opts {
		fn(mutate)
	}

	c, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return d.Get(ctx, datagraph.ClusterSlug(c.Slug))
}

func (d *database) Archive(ctx context.Context, slug datagraph.ClusterSlug) (*datagraph.Cluster, error) {
	update := d.db.Cluster.Update().Where(cluster.Slug(string(slug)))
	update.SetDeletedAt(time.Now())

	_, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return d.Get(ctx, slug)
}
