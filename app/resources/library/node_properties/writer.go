package node_properties

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/ent/propertyschemafield"
)

type SchemaWriter struct {
	db *ent.Client
}

func New(db *ent.Client) (*SchemaWriter, *Writer) {
	return &SchemaWriter{
			db: db,
		}, &Writer{
			db: db,
		}
}

type SchemaFieldMutation struct {
	ID   opt.Optional[xid.ID]
	Name string
	Type string
	Sort string
}

type FieldSchemaMutations []*SchemaFieldMutation

func (w SchemaWriter) CreateForNode(ctx context.Context, nodeID library.NodeID, schemas FieldSchemaMutations) (*library.PropertySchema, error) {
	node, err := w.db.Node.Get(ctx, xid.ID(nodeID))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	schemaID, err := w.doSchemaUpdates(ctx, node.Edges.PropertySchema, schemas, node)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &library.PropertySchema{
		ID: *schemaID,
	}, nil
}

func (w *SchemaWriter) UpdateChildren(ctx context.Context, qk library.QueryKey, schemas FieldSchemaMutations) (*library.PropertySchema, error) {
	parent, err := w.db.Node.Query().Where(qk.Predicate()).WithNodes(func(nq *ent.NodeQuery) {
		nq.WithPropertySchema(func(psq *ent.PropertySchemaQuery) {
			psq.WithFields()
		})
	}).Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	children := parent.Edges.Nodes
	if len(children) == 0 {
		// no children to update, no-op.
		return &library.PropertySchema{}, nil
	}

	return w.updateNodes(ctx, schemas, children...)
}

func (w *SchemaWriter) UpdateSiblings(ctx context.Context, qk library.QueryKey, schemas FieldSchemaMutations) (*library.PropertySchema, error) {
	current, err := w.db.Node.Query().Where(
		node.Or(qk.Predicate()),
	).Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	children, err := w.db.Node.Query().Where(
		node.HasParentWith(node.ID(current.ParentNodeID)),
	).All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return w.updateNodes(ctx, schemas, children...)
}

func (w *SchemaWriter) updateNodes(ctx context.Context, schemas FieldSchemaMutations, nodes ...*ent.Node) (*library.PropertySchema, error) {
	if len(nodes) == 0 {
		return &library.PropertySchema{}, nil
	}

	grouping := lo.GroupBy(nodes, func(n *ent.Node) string {
		return n.PropertySchemaID.String()
	})

	if len(grouping) > 1 {
		// TODO: Self heal by picking the most common schema and re-assigning.
		panic("schema mismatch")
	}

	currentSchema := nodes[0].Edges.PropertySchema

	schemaID, err := w.doSchemaUpdates(ctx, currentSchema, schemas, nodes...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// Mutations finished, query the final result for returning.

	return w.Get(ctx, *schemaID)
}

func (w *SchemaWriter) Get(ctx context.Context, schemaID xid.ID) (*library.PropertySchema, error) {
	schemaFields, err := w.db.PropertySchemaField.Query().
		Where(propertyschemafield.SchemaID(schemaID)).
		Order(ent.Asc(propertyschemafield.FieldSort)).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	updatedSchemas := dt.Map(schemaFields, func(f *ent.PropertySchemaField) *library.PropertySchemaField {
		return &library.PropertySchemaField{
			ID:   f.ID,
			Name: f.Name,
			Type: f.Type,
			Sort: f.Sort,
		}
	})

	return &library.PropertySchema{
		ID:     schemaID,
		Fields: updatedSchemas,
	}, nil
}

func (w *SchemaWriter) AddFields(ctx context.Context, schemaID xid.ID, schemas FieldSchemaMutations) (*library.PropertySchema, error) {
	fields := []*ent.PropertySchemaFieldCreate{}
	for _, s := range schemas {
		fields = append(fields, w.db.PropertySchemaField.Create().SetName(s.Name).SetSort(s.Sort).SetType(s.Type).SetSchemaID(schemaID))
	}

	err := w.db.PropertySchemaField.CreateBulk(fields...).Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return w.Get(ctx, schemaID)
}

func (w *SchemaWriter) doSchemaUpdates(ctx context.Context, currentSchema *ent.PropertySchema, schemas FieldSchemaMutations, children ...*ent.Node) (*xid.ID, error) {
	creates := FieldSchemaMutations{}
	updates := FieldSchemaMutations{}
	deletes := map[xid.ID]*ent.PropertySchemaField{}

	if currentSchema != nil {
		deletes = lo.KeyBy(currentSchema.Edges.Fields, func(f *ent.PropertySchemaField) xid.ID { return f.ID })
	}

	for _, s := range schemas {
		id, ok := s.ID.Get()
		delete(deletes, id)
		if !ok {
			creates = append(creates, s)
			continue
		}

		updates = append(updates, s)

	}

	tx, err := w.db.Tx(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	defer func() {
		err = tx.Rollback()
	}()

	// Create schema if it doesn't exist
	if currentSchema == nil {
		currentSchema, err = tx.PropertySchema.Create().Save(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		childIDs := dt.Map(children, func(n *ent.Node) xid.ID { return n.ID })

		// assign schema to all child nodes
		err = tx.Node.Update().Where(node.IDIn(childIDs...)).SetPropertySchema(currentSchema).Exec(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	// Delete fields

	if len(deletes) > 0 {
		deleteIDs := dt.Map(lo.Values(deletes), func(f *ent.PropertySchemaField) xid.ID { return f.ID })
		_, err = tx.PropertySchemaField.Delete().Where(propertyschemafield.IDIn(deleteIDs...)).Exec(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	// Update fields
	if len(updates) > 0 {
		for _, s := range updates {
			// we know this is non-zero already.
			id := s.ID.OrZero()

			err = tx.PropertySchemaField.UpdateOneID(id).SetName(s.Name).SetSort(s.Sort).SetType(s.Type).Exec(ctx)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
		}
	}

	// Create fields
	if len(creates) > 0 {
		for _, s := range creates {
			err = tx.PropertySchemaField.Create().SetName(s.Name).SetSort(s.Sort).SetType(s.Type).SetSchemaID(currentSchema.ID).Exec(ctx)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &currentSchema.ID, nil
}
