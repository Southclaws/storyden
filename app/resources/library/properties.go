package library

import (
	"slices"
	"strings"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/internal/ent"
)

type PropertySchemaField struct {
	ID   xid.ID
	Name string
	Type PropertyType
	Sort string
}

type PropertySchemaFields []*PropertySchemaField

type Property struct {
	Field PropertySchemaField
	Value opt.Optional[string]
}

type Properties []*Property

type PropertyTable struct {
	Schema     PropertySchema
	Properties Properties
}

type PropertySchema struct {
	ID     xid.ID
	Fields PropertySchemaFields
}

func (p PropertySchema) FieldIDs() []xid.ID {
	return dt.Map(p.Fields, func(f *PropertySchemaField) xid.ID { return f.ID })
}

func (p PropertySchema) GetField(id xid.ID) (*PropertySchemaField, bool) {
	lookup := lo.KeyBy(p.Fields, func(f *PropertySchemaField) xid.ID { return f.ID })
	f, ok := lookup[id]
	return f, ok
}

type PropertySchemaMutation struct {
	NewProps      PropertyMutationList
	ExistingProps ExistingPropertyMutations
	RemovedProps  ExistingPropertyMutations
}

// Split takes a mutation (a list of properties to update) and splits it into
// two lists, one for properties that need to be added to the schema and current
// schema fields which can be processed as simple property update operations.
// We also need to get the actual field ID for each existing property.
func (p PropertySchema) Split(mutation PropertyMutationList) (*PropertySchemaMutation, error) {
	fids := lo.FilterMap(mutation, func(p *PropertyMutation, _ int) (xid.ID, bool) { return p.ID.Get() })

	existingProperties, newProps := lo.FilterReject(mutation, func(m *PropertyMutation, _ int) bool {
		return m.ID.Ok()
	})

	// split by existence in the schema. this would be simpler with a DiffBy().
	removedIDs, _ := lo.Difference(p.FieldIDs(), fids)

	existingProps, err := dt.MapErr(existingProperties, p.getSchemaMutationFromPropertyMutation)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	removedProps, err := dt.MapErr(removedIDs, func(fid xid.ID) (*ExistingPropertyMutation, error) {
		f, ok := p.GetField(fid)
		if !ok {
			return nil, fault.Wrap(fault.Newf("field ID '%v' not found in schema", fid), ftag.With(ftag.InvalidArgument))
		}
		return &ExistingPropertyMutation{
			PropertySchemaField: *f,
		}, nil
	})
	if err != nil {
		return nil, fault.Wrap(err)
	}

	// Check for missing IDs, where the request contains a named property that
	// exists but the ID is not set. The above code will interpret this as the
	// field being removed and added again. This is not what we want so error.

	if len(removedProps) > 0 && len(newProps) > 0 {
		removedNames := dt.Map(removedProps, func(p *ExistingPropertyMutation) string { return p.Name })
		newNames := dt.Map(newProps, func(p *PropertyMutation) string { return p.Name })

		// If the intersection of the two lists is not empty, then we have a
		// conflict where a property is being removed and added again in the same
		// mutation. This is not allowed.
		intersect := lo.Intersect(removedNames, newNames)
		if len(intersect) > 0 {
			return nil, fault.Wrap(
				fault.Newf("cannot remove and add the same property in the same mutation: %v did you miss a field ID?", intersect),
				ftag.With(ftag.InvalidArgument),
			)
		}
	}

	return &PropertySchemaMutation{
		NewProps:      newProps,
		ExistingProps: existingProps,
		RemovedProps:  removedProps,
	}, nil
}

func (p PropertySchema) getSchemaMutationFromPropertyMutation(pm *PropertyMutation) (*ExistingPropertyMutation, error) {
	f, ok := p.GetField(pm.ID.OrZero())
	if !ok {
		return nil, fault.Wrap(fault.Newf("field '%v' not found in schema", pm.ID), ftag.With(ftag.InvalidArgument))
	}

	// During a property mutation, the request may also change the schema by
	// changing the name, type or sort properties. Mark as changed if so.
	isChanged := false
	// NOTE: A small hack here checking if the mutation name is empty. This will
	// permit property assignments by just the fid but as a result creates some
	// undefined behaviour when explicitly setting the field name to "". This
	// whole code path needs a refactor.
	if f.Name != pm.Name && pm.Name != "" {
		isChanged = true
		f.Name = pm.Name
	}
	if t, ok := pm.Type.Get(); ok && t != f.Type {
		isChanged = true
		f.Type = t
	}
	if s, ok := pm.Sort.Get(); ok && s != f.Sort {
		isChanged = true
		f.Sort = s
	}

	return &ExistingPropertyMutation{
		PropertySchemaField: *f,
		IsSchemaChanged:     isChanged,
		Value:               pm.Value,
	}, nil
}

// Property mutations are used to update properties on a node.
type PropertyMutation struct {
	// ID is optional, when set the mutation is modifying an existing field and
	// when not set, the mutation assumes it's a new field.
	ID    opt.Optional[xid.ID]
	Name  string
	Value string
	Type  opt.Optional[PropertyType]
	Sort  opt.Optional[string]
}

type PropertyMutationList []*PropertyMutation

type ExistingPropertyMutation struct {
	PropertySchemaField
	IsSchemaChanged bool
	Value           string
}

type ExistingPropertyMutations []*ExistingPropertyMutation

func MapPropertyFieldSchema(in PropertySchemaQueryRow) PropertySchemaField {
	return PropertySchemaField{
		ID:   in.FieldID,
		Name: in.Name,
		Type: in.Type,
		Sort: in.Sort,
	}
}

// PropertySchemaQueryRow is a row from the property schema query which pulls
// all the property schemas for both sibling and child properties of a node.
type PropertySchemaQueryRow struct {
	SchemaID xid.ID       `db:"schema_id"`
	FieldID  xid.ID       `db:"field_id"`
	Name     string       `db:"name"`
	Type     PropertyType `db:"type"`
	Sort     string       `db:"sort"`
	Source   string       `db:"source"`
}

type PropertySchemaQueryRows []PropertySchemaQueryRow

type PropertySchemaTable struct {
	siblingSchemas PropertySchemaQueryRows
	childSchemas   PropertySchemaQueryRows
}

// Map harmonises the splits the raw rows into sibling and child schemas.
func (r PropertySchemaQueryRows) Map() *PropertySchemaTable {
	siblings, children := lo.FilterReject(r, func(r PropertySchemaQueryRow, _ int) bool {
		return r.Source == "sibling"
	})

	return &PropertySchemaTable{
		siblingSchemas: siblings,
		childSchemas:   children,
	}
}

// BuildPropertyTable yields the properties that are set for the node and also
// properties that don't have values by merging in the unused property schemas.
func (r *PropertySchemaTable) BuildPropertyTable(in []*ent.Property, isRoot bool) *PropertyTable {
	if r == nil {
		return nil
	}

	// When mapping a node with children, we fetch the entire list of schemas
	// from the perspective of the root fetched node. So when mapping properties
	// we need to switch the source of schemas depending on the mapping context.
	schemas := r.siblingSchemas
	if !isRoot {
		schemas = r.childSchemas
	}

	if len(schemas) == 0 {
		return nil
	}

	// Properties are name-unique within a schema (name + schema_id as an index)
	propMap := lo.KeyBy(schemas, func(r PropertySchemaQueryRow) xid.ID { return r.FieldID })

	// Assumption: all schemas for all children are identical, select the first
	// field to retrieve the schema's ID.
	schemaID := schemas[0].SchemaID

	fields := []*Property{}
	schema := PropertySchema{
		ID: schemaID,
	}

	// Add all the properties that have values.
	for _, p := range in {
		if s, ok := propMap[p.FieldID]; ok {
			delete(propMap, p.FieldID)
			fieldSchema := MapPropertyFieldSchema(s)
			fields = append(fields, &Property{
				Field: fieldSchema,
				Value: opt.New(p.Value),
			})
			schema.Fields = append(schema.Fields, &fieldSchema)
		}

		// If a property was not in the schema, ignore it. The member might move
		// a node back to a parent that had a different schema so we retain data
	}

	// Add the remaining property schemas that do not have values.
	for _, p := range propMap {
		fieldSchema := MapPropertyFieldSchema(p)
		schema.Fields = append(schema.Fields, &fieldSchema)
		fields = append(fields, &Property{
			Field: fieldSchema,
		})
	}

	slices.SortFunc(fields, func(i, j *Property) int {
		return strings.Compare(i.Field.Sort, j.Field.Sort)
	})

	slices.SortFunc(schema.Fields, func(i, j *PropertySchemaField) int {
		return strings.Compare(i.Sort, j.Sort)
	})

	return &PropertyTable{
		Schema:     schema,
		Properties: fields,
	}
}

func (r PropertySchemaTable) ChildSchemas() *PropertySchema {
	if len(r.childSchemas) == 0 {
		return nil
	}

	// Assumption: all schemas for all children are identical, select the first
	// field to retrieve the schema's ID.
	schemaID := r.childSchemas[0].SchemaID

	fields := dt.Map(r.childSchemas, func(s PropertySchemaQueryRow) *PropertySchemaField {
		return &PropertySchemaField{
			ID:   s.FieldID,
			Name: s.Name,
			Type: s.Type,
			Sort: s.Sort,
		}
	})

	slices.SortFunc(fields, func(i, j *PropertySchemaField) int {
		return strings.Compare(i.Sort, j.Sort)
	})

	return &PropertySchema{
		ID:     schemaID,
		Fields: fields,
	}
}
