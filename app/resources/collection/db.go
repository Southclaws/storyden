package collection

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/collection"
)

type database struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &database{db}
}

func (d *database) Create(ctx context.Context, owner account.AccountID, name string, desc string, opts ...Option) (*Collection, error) {
	create := d.db.Collection.Create()
	mutate := create.Mutation()

	mutate.SetOwnerID(xid.ID(owner))
	mutate.SetName(name)
	mutate.SetDescription(desc)

	for _, fn := range opts {
		fn(mutate)
	}

	col, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return d.Get(ctx, CollectionID(col.ID))
}

func (d *database) List(ctx context.Context, filters ...Filter) ([]*Collection, error) {
	q := d.db.Collection.
		Query().
		WithOwner().
		WithPosts(func(pq *ent.PostQuery) {
			pq.WithAuthor()
			pq.WithCategory()
			pq.WithTags()
		}).
		WithNodes(func(nq *ent.NodeQuery) {
			nq.WithOwner()
			nq.WithAssets()
			nq.WithTags()
		})

	for _, fn := range filters {
		fn(q)
	}

	cols, err := q.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	all, err := dt.MapErr(cols, FromModel)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return all, nil
}

func (d *database) Get(ctx context.Context, id CollectionID) (*Collection, error) {
	col, err := d.db.Collection.
		Query().
		Where(collection.ID(xid.ID(id))).
		WithOwner().
		WithPosts(func(pq *ent.PostQuery) {
			pq.WithAuthor()
			pq.WithCategory()
			pq.WithTags()
			pq.WithRoot()
		}).
		WithNodes(func(nq *ent.NodeQuery) {
			nq.WithOwner()
			nq.WithAssets()
			nq.WithTags()
		}).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return FromModel(col)
}

func (d *database) Update(ctx context.Context, id CollectionID, opts ...Option) (*Collection, error) {
	create := d.db.Collection.UpdateOneID(xid.ID(id))
	mutate := create.Mutation()

	for _, fn := range opts {
		fn(mutate)
	}

	_, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return d.Get(ctx, id)
}

func (d *database) Delete(ctx context.Context, id CollectionID) error {
	err := d.db.Collection.DeleteOneID(xid.ID(id)).Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
