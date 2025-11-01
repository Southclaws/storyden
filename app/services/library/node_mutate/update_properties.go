package node_mutate

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_properties"
)

func (s *Manager) applyPropertyMutations(ctx context.Context, n *library.Node, properties library.PropertyMutationList) (*library.PropertyTable, error) {
	schema, hasSchema := n.Properties.Get()

	migration, err := schema.Schema.Split(properties)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !hasSchema {
		mutations, err := dt.MapErr(migration.NewProps, mapNewPropertyMutation)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if len(mutations) > 0 {
			newSchema, err := s.schemaWriter.CreateForNode(ctx, library.NodeID(n.Mark.ID()), mutations)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			schema.Schema = *newSchema
		}
	} else {
		schemaUpdates := []*node_properties.SchemaFieldMutation{}

		schemaUpdates = lo.FilterMap(migration.ExistingProps, func(pm *library.ExistingPropertyMutation, _ int) (*node_properties.SchemaFieldMutation, bool) {
			if !pm.IsSchemaChanged {
				return nil, false
			}

			return &node_properties.SchemaFieldMutation{
				ID:   opt.New(pm.ID),
				Name: pm.Name,
				Type: pm.Type,
				Sort: pm.Sort,
			}, true
		})

		if len(schemaUpdates) > 0 {
			newSchema, err := s.schemaWriter.UpdateSiblings(ctx, library.NewQueryKey(n.Mark), schemaUpdates)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			schema.Schema = *newSchema
		}
	}

	for _, newProp := range migration.NewProps {
		newSchemaProp, found := lo.Find(schema.Schema.Fields, func(f *library.PropertySchemaField) bool {
			return f.Name == newProp.Name
		})
		if !found {
			continue
		}
		for i, mutProp := range properties {
			if newProp.Name == mutProp.Name {
				properties[i].ID = opt.New(newSchemaProp.ID)
			}
		}
	}

	// TODO: Remove all this code below and move other migrations into the
	// above UpdateSiblings call. Currently that call only does existing
	// field updates but it should perform all migrations. This means we
	// would no longer need to call .Split() twice and mutate properties.

	// re-validate the schema properties mutation plan.
	migration, err = schema.Schema.Split(properties)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if len(migration.NewProps) > 0 {
		newSchemaFields, err := dt.MapErr(migration.NewProps, mapNewPropertyMutation)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if len(newSchemaFields) > 0 {
			newSchema, err := s.schemaWriter.AddFields(ctx, schema.Schema.ID, newSchemaFields)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			schema.Schema = *newSchema

			for _, newProp := range migration.NewProps {
				newSchemaProp, found := lo.Find(schema.Schema.Fields, func(f *library.PropertySchemaField) bool {
					return f.Name == newProp.Name
				})
				if !found {
					continue
				}

				migration.ExistingProps = append(migration.ExistingProps, &library.ExistingPropertyMutation{
					PropertySchemaField: *newSchemaProp,
					Value:               newProp.Value,
				})
			}
		}
	}

	if len(migration.RemovedProps) > 0 {
		removedSchemaFields, err := dt.MapErr(migration.RemovedProps, mapExistingPropertyMutation)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		newSchema, err := s.schemaWriter.RemoveFields(ctx, schema.Schema.ID, removedSchemaFields)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		schema.Schema = *newSchema
	}

	// Assumption: all schema changes are done by this point. Update no
	// longer needs to actually check the schema, just write the data.
	updated, err := s.propWriter.Update(ctx, library.NodeID(n.GetID()), schema.Schema, migration.ExistingProps)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return updated, nil
}

func mapNewPropertyMutation(pm *library.PropertyMutation) (*node_properties.SchemaFieldMutation, error) {
	ft, ok := pm.Type.Get()
	if !ok {
		return nil, fault.Wrap(fault.New("no type on new field"), ftag.With(ftag.InvalidArgument), fmsg.WithDesc("missing type", "You must provide a field type when adding a new property."))
	}
	return &node_properties.SchemaFieldMutation{
		Name: pm.Name,
		Type: ft,
		Sort: pm.Sort.OrZero(),
	}, nil
}

func mapExistingPropertyMutation(pm *library.ExistingPropertyMutation) (*node_properties.SchemaFieldMutation, error) {
	return &node_properties.SchemaFieldMutation{
		ID:   opt.New(pm.ID),
		Name: pm.Name,
		Type: pm.Type,
		Sort: pm.Sort,
	}, nil
}
