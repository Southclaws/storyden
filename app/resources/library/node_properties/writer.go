package node_properties

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_cache"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/ent/propertyschemafield"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type SchemaWriter struct {
	db    *ent.Client
	cache *node_cache.Cache
	bus   *pubsub.Bus
}

func New(
	db *ent.Client,
	cache *node_cache.Cache,
	bus *pubsub.Bus,
) (*SchemaWriter, *Writer) {
	return &SchemaWriter{
			db:    db,
			cache: cache,
			bus:   bus,
		}, &Writer{
			db: db,
		}
}

type SchemaFieldMutation struct {
	ID   opt.Optional[xid.ID]
	Name string
	Type library.PropertyType
	Sort string
}

type FieldSchemaMutations []*SchemaFieldMutation

func (w SchemaWriter) CreateForNode(ctx context.Context, nodeID library.NodeID, schemas FieldSchemaMutations) (*library.PropertySchema, error) {
	node, err := w.db.Node.Get(ctx, xid.ID(nodeID))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := w.cache.Invalidate(ctx, node.Slug); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	schemaID, err := w.doSchemaUpdates(ctx, node.Edges.PropertySchema, schemas, node)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	w.bus.Publish(ctx, rpc.EventNodeUpdated{
		ID:   library.NodeID(node.ID),
		Slug: node.Slug,
	})

	return w.Get(ctx, *schemaID)
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

	events := []rpc.EventNodeUpdated{
		{
			ID:   library.NodeID(parent.ID),
			Slug: parent.Slug,
		},
	}

	if err := w.cache.Invalidate(ctx, parent.Slug); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	for _, node := range children {
		if err := w.cache.Invalidate(ctx, node.Slug); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		events = append(events, rpc.EventNodeUpdated{
			ID:   library.NodeID(node.ID),
			Slug: node.Slug,
		})
	}

	schema, err := w.updateNodes(ctx, schemas, children...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	w.bus.PublishMany(ctx, events)

	return schema, nil
}

func (w *SchemaWriter) UpdateSiblings(ctx context.Context, qk library.QueryKey, schemas FieldSchemaMutations) (*library.PropertySchema, error) {
	current, err := w.db.Node.Query().
		Where(
			node.Or(qk.Predicate()),
		).
		WithPropertySchema().
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	siblings, err := w.db.Node.Query().
		Where(
			node.HasParentWith(node.ID(current.ParentNodeID)),
		).
		WithPropertySchema().
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	events := []rpc.EventNodeUpdated{
		{
			ID:   library.NodeID(current.ID),
			Slug: current.Slug,
		},
	}

	if err := w.cache.Invalidate(ctx, current.Slug); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	for _, node := range siblings {
		if err := w.cache.Invalidate(ctx, node.Slug); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		events = append(events, rpc.EventNodeUpdated{
			ID:   library.NodeID(node.ID),
			Slug: node.Slug,
		})
	}

	schema, err := w.updateNodes(ctx, schemas, siblings...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	w.bus.PublishMany(ctx, events)

	return schema, nil
}

func (w *SchemaWriter) updateNodes(ctx context.Context, schemas FieldSchemaMutations, nodes ...*ent.Node) (*library.PropertySchema, error) {
	if len(nodes) == 0 {
		// NOTE: This is wrong... but awkward to fix. If the caller provides
		// arguments that result in zero nodes to update, it does not result
		// in a new schema being returned (we can't query the schema here as
		// there's no node to query from, and the current node in the caller
		// may not have a schema). So right now, callers just need to be careful
		// to not pass in zero nodes and only re-assign the schema if present.
		return &library.PropertySchema{}, nil
	}

	currentSchema, err := w.ensureSiblingSchemaConsistency(ctx, nodes)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	schemaID, err := w.doSchemaUpdates(ctx, currentSchema, schemas, nodes...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return w.Get(ctx, *schemaID)
}

func (w *SchemaWriter) ensureSiblingSchemaConsistency(ctx context.Context, nodes []*ent.Node) (*ent.PropertySchema, error) {
	var targetSchema *ent.PropertySchema
	targetSchemaCount := 0

	grouping := lo.GroupBy(nodes, func(n *ent.Node) *xid.ID {
		return n.PropertySchemaID
	})

	for schemaID, nodes := range grouping {
		if schemaID == nil {
			continue
		}

		nodesWithSchema := len(nodes)
		if nodesWithSchema > targetSchemaCount {
			targetSchema = nodes[0].Edges.PropertySchema
			targetSchemaCount = nodesWithSchema
		}
	}

	if targetSchema == nil {
		return nil, nil
	}

	// gather all nodes which do NOT have targetSchema as their schema.
	nodesWithWrongSchema := []*ent.Node{}
	for _, nodes := range grouping {
		if nodes[0].PropertySchemaID == nil || *nodes[0].PropertySchemaID != targetSchema.ID {
			nodesWithWrongSchema = append(nodesWithWrongSchema, nodes...)
		}
	}

	if len(nodesWithWrongSchema) > 0 {
		tx, err := w.db.Tx(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		defer tx.Rollback()

		nodeIDs := dt.Map(nodesWithWrongSchema, func(n *ent.Node) xid.ID { return n.ID })

		err = tx.Node.Update().
			Where(node.IDIn(nodeIDs...)).
			SetPropertySchema(targetSchema).
			Exec(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if err := tx.Commit(); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	return targetSchema, nil
}

func (w *SchemaWriter) Get(ctx context.Context, schemaID xid.ID) (*library.PropertySchema, error) {
	schemaFields, err := w.db.PropertySchemaField.Query().
		Where(propertyschemafield.SchemaID(schemaID)).
		Order(ent.Asc(propertyschemafield.FieldSort)).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	updatedSchemas, err := dt.MapErr(schemaFields, func(f *ent.PropertySchemaField) (*library.PropertySchemaField, error) {
		t, err := library.NewPropertyType(f.Type)
		if err != nil {
			return nil, err
		}

		return &library.PropertySchemaField{
			ID:   f.ID,
			Name: f.Name,
			Type: t,
			Sort: f.Sort,
		}, nil
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &library.PropertySchema{
		ID:     schemaID,
		Fields: updatedSchemas,
	}, nil
}

func (w *SchemaWriter) AddFields(ctx context.Context, schemaID xid.ID, schemas FieldSchemaMutations) (*library.PropertySchema, error) {
	fields := []*ent.PropertySchemaFieldCreate{}
	for _, s := range schemas {
		fields = append(fields, w.db.PropertySchemaField.Create().SetName(s.Name).SetSort(s.Sort).SetType(s.Type.String()).SetSchemaID(schemaID))
	}

	err := w.db.PropertySchemaField.CreateBulk(fields...).Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return w.Get(ctx, schemaID)
}

func (w *SchemaWriter) RemoveFields(ctx context.Context, schemaID xid.ID, schemas FieldSchemaMutations) (*library.PropertySchema, error) {
	tx, err := w.db.Tx(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	defer tx.Rollback()

	for _, s := range schemas {
		_, err = tx.PropertySchemaField.Delete().
			Where(
				propertyschemafield.SchemaID(schemaID),
				propertyschemafield.Name(s.Name),
			).
			Exec(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	if err := tx.Commit(); err != nil {
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

			err = tx.PropertySchemaField.
				UpdateOneID(id).
				SetName(s.Name).
				SetSort(s.Sort).
				SetType(s.Type.String()).
				Exec(ctx)
			if err != nil {
				if ent.IsConstraintError(err) {
					err = fault.Wrap(err, ftag.With(ftag.AlreadyExists), fmsg.WithDesc("constraint error",
						"A property with this name already exists."))
				}
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
		}
	}

	// Create fields
	if len(creates) > 0 {
		for _, s := range creates {
			err = tx.PropertySchemaField.Create().
				SetName(s.Name).
				SetSort(s.Sort).
				SetType(s.Type.String()).
				SetSchemaID(currentSchema.ID).
				Exec(ctx)
			if err != nil {
				if ent.IsConstraintError(err) {
					err = fault.Wrap(err, ftag.With(ftag.AlreadyExists), fmsg.WithDesc("constraint error",
						"A property with this name already exists."))
				}
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &currentSchema.ID, nil
}
