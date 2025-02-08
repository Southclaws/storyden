package library

import (
	"slices"
	"strings"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/internal/ent"
)

type PropertySchemaField struct {
	ID   xid.ID
	Name string
	Type string
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

func (p PropertySchema) FieldNames() []string {
	return dt.Map(p.Fields, func(f *PropertySchemaField) string {
		return f.Name
	})
}

func (p PropertySchema) GetField(name string) (*PropertySchemaField, bool) {
	lookup := lo.KeyBy(p.Fields, func(f *PropertySchemaField) string { return f.Name })
	f, ok := lookup[name]
	return f, ok
}

// Split takes a mutation (a list of properties to update) and splits it into
// two lists, one for properties that need to be added to the schema and current
// schema fields which can be processed as simple property update operations.
// We also need to get the actual field ID for each existing property.
func (p PropertySchema) Split(mutation PropertyMutationList) (newProps PropertyMutationList, existingProps ExistingPropertyMutations) {
	mutationFieldMap := lo.KeyBy(mutation, func(p PropertyMutation) string { return p.Name })
	mutationFieldNames := dt.Map(mutation, func(p PropertyMutation) string { return p.Name })

	// split by existence in the schema. this would be simpler with a DiffBy().
	_, newPropertyNames := lo.Difference(p.FieldNames(), mutationFieldNames)
	existingProperties := lo.Intersect(mutationFieldNames, p.FieldNames())

	existingProps = dt.Map(existingProperties, func(name string) *ExistingPropertyMutation {
		f, ok := p.GetField(name)
		if !ok {
			panic("field not found in schema")
		}
		return &ExistingPropertyMutation{
			PropertySchemaField: *f,
			Value:               mutationFieldMap[name].Value,
		}
	})

	newProps = dt.Map(newPropertyNames, func(name string) PropertyMutation {
		return mutationFieldMap[name]
	})

	return
}

// Property mutations are used to update properties on a node.
type PropertyMutation struct {
	Name  string
	Value string
	Type  opt.Optional[string]
	Sort  opt.Optional[string]
}

type PropertyMutationList []PropertyMutation

type ExistingPropertyMutation struct {
	PropertySchemaField
	Value string
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
	SchemaID xid.ID `db:"schema_id"`
	FieldID  xid.ID `db:"field_id"`
	Name     string `db:"name"`
	Type     string `db:"type"`
	Sort     string `db:"sort"`
	Source   string `db:"source"`
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
