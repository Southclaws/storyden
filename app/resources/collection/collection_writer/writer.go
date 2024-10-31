package collection_writer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/collection/collection_querier"
	"github.com/Southclaws/storyden/internal/ent"
)

type Writer struct {
	db      *ent.Client
	querier *collection_querier.Querier
}

func New(db *ent.Client, querier *collection_querier.Querier) *Writer {
	return &Writer{
		db:      db,
		querier: querier,
	}
}

type Option func(*ent.CollectionMutation)

func WithID(id collection.CollectionID) Option {
	return func(c *ent.CollectionMutation) {
		c.SetID(xid.ID(id))
	}
}

func WithName(v string) Option {
	return func(c *ent.CollectionMutation) {
		c.SetName(v)
	}
}

func WithDescription(v string) Option {
	return func(c *ent.CollectionMutation) {
		c.SetDescription(v)
	}
}

func (w *Writer) Create(ctx context.Context, owner account.AccountID, name string, opts ...Option) (*collection.CollectionWithItems, error) {
	create := w.db.Collection.Create()
	mutate := create.Mutation()

	mutate.SetOwnerID(xid.ID(owner))
	mutate.SetName(name)

	for _, fn := range opts {
		fn(mutate)
	}

	col, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return w.querier.Get(ctx, collection.CollectionID(col.ID))
}

func (w *Writer) Update(ctx context.Context, id collection.CollectionID, opts ...Option) (*collection.CollectionWithItems, error) {
	create := w.db.Collection.UpdateOneID(xid.ID(id))
	mutate := create.Mutation()

	for _, fn := range opts {
		fn(mutate)
	}

	_, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return w.querier.Get(ctx, id)
}

func (w *Writer) Delete(ctx context.Context, id collection.CollectionID) error {
	err := w.db.Collection.DeleteOneID(xid.ID(id)).Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
