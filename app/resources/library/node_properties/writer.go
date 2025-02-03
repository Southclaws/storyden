package node_properties

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/propertyschemafield"
	"github.com/rs/xid"
	"github.com/samber/lo"
)

type SchemaWriter struct {
	db *ent.Client
}

func New(db *ent.Client) *SchemaWriter {
	return &SchemaWriter{
		db: db,
	}
}

type SchemaMutation struct {
	ID   opt.Optional[xid.ID]
	Name string
	Type string
	Sort string
}

type SchemaMutations []*SchemaMutation

func (w *SchemaWriter) UpdateChildren(ctx context.Context, qk library.QueryKey, schemas SchemaMutations) (library.PropertySchemas, error) {
	parent, err := w.db.Node.Query().Where(qk.Predicate()).WithNodes(func(nq *ent.NodeQuery) {
		nq.WithPropertySchemas(func(psq *ent.PropertySchemaQuery) {
			psq.WithFields()
		})
	}).Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	children := parent.Edges.Nodes
	if len(children) == 0 {
		// no children to update, no-op.w
		return library.PropertySchemas{}, nil
	}

	grouping := lo.GroupBy(children, func(n *ent.Node) string {
		return n.PropertySchemaID.String()
	})

	if len(grouping) > 1 {
		// TODO: Self heal by picking the most common schema and re-assigning.
		panic("schema mismatch")
	}

	currentSchema := children[0].Edges.PropertySchemas

	creates := SchemaMutations{}
	updates := SchemaMutations{}
	deletes := lo.KeyBy(currentSchema.Edges.Fields, func(f *ent.PropertySchemaField) xid.ID { return f.ID })
	for _, s := range schemas {
		id, ok := s.ID.Get()
		if !ok {
			creates = append(creates, s)
			continue
		}

		updates = append(updates, s)

		delete(deletes, id)
	}

	deleteIDs := dt.Map(lo.Values(deletes), func(f *ent.PropertySchemaField) xid.ID { return f.ID })
	_, err = w.db.PropertySchemaField.Delete().Where(propertyschemafield.IDIn(deleteIDs...)).Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tx, err := w.db.Tx(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	defer func() {
		err = tx.Rollback()
	}()

	// Update fields
	for _, s := range updates {
		// we know this is non-zero already.
		id := s.ID.OrZero()

		err = tx.PropertySchemaField.UpdateOneID(id).SetName(s.Name).SetSort(s.Sort).SetType(s.Type).Exec(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	// Create fields
	for _, s := range creates {
		err = tx.PropertySchemaField.Create().SetName(s.Name).SetSort(s.Sort).SetType(s.Type).Exec(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// Mutations finished, query the final result for returning.

	schemaFields, err := w.db.PropertySchemaField.Query().Where(propertyschemafield.SchemaID(currentSchema.ID)).All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	updatedSchemas := dt.Map(schemaFields, func(f *ent.PropertySchemaField) *library.PropertySchema {
		return &library.PropertySchema{
			Name: f.Name,
			Type: f.Type,
			Sort: f.Sort,
		}
	})

	return updatedSchemas, nil
}
