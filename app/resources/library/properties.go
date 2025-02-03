package library

import (
	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/internal/ent"
)

type PropertySchema struct {
	Name string
	Type string
	Sort string
}

type Property struct {
	PropertySchema
	Value opt.Optional[string]
}

type PropertyTable []*Property

type PropertySchemas []*PropertySchema

func MapPropertySchema(in PropertySchemaQueryRow) PropertySchema {
	return PropertySchema{
		Name: in.Name,
		Type: in.Type,
		Sort: in.Sort,
	}
}

// PropertySchemaQueryRow is a row from the property schema query which pulls
// all the property schemas for both sibling and child properties of a node.
type PropertySchemaQueryRow struct {
	FieldID string `db:"id"`
	Name    string `db:"name"`
	Type    string `db:"type"`
	Sort    string `db:"sort"`
	Source  string `db:"source"`
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
func (r *PropertySchemaTable) BuildPropertyTable(in []*ent.Property, isRoot bool) PropertyTable {
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
	propMap := lo.KeyBy(schemas, func(r PropertySchemaQueryRow) string { return r.FieldID })

	out := PropertyTable{}

	// Add all the properties that have values.
	for _, p := range in {
		if s, ok := propMap[p.FieldID.String()]; ok {
			delete(propMap, p.FieldID.String())
			out = append(out, &Property{
				PropertySchema: MapPropertySchema(s),
				Value:          opt.New(p.Value),
			})
		}

		// If a property was not in the schema, ignore it. The member might move
		// a node back to a parent that had a different schema so we retain data
	}

	// Add the remaining property schemas that do not have values.
	for _, p := range propMap {
		out = append(out, &Property{
			PropertySchema: MapPropertySchema(p),
		})
	}
	return out
}

func (r PropertySchemaTable) ChildSchemas() PropertySchemas {
	return dt.Map(r.childSchemas, func(r PropertySchemaQueryRow) *PropertySchema {
		return &PropertySchema{
			Name: r.Name,
			Type: r.Type,
			Sort: r.Sort,
		}
	})
}
