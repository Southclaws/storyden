package library

import (
	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/internal/ent"
)

type Property struct {
	PropertySchema
	Value opt.Optional[string]
}

type PropertyTable []Property

func MapProperty(in *ent.Property) Property {
	return Property{
		PropertySchema: MapPropertyValueSchema(in),
		Value:          opt.New(in.Value),
	}
}

type PropertySchema struct {
	Name string
	Type string
}

type PropertySchemas []PropertySchema

func MapPropertySchema(in PropertySchemaQueryRow) PropertySchema {
	return PropertySchema{
		Name: in.Name,
		Type: in.Type,
	}
}

func MapPropertyValueSchema(in *ent.Property) PropertySchema {
	return PropertySchema{
		Name: in.Name,
		Type: in.Type,
	}
}

// PropertySchemaQueryRow is a row from the property schema query which pulls
// all the property schemas for both sibling and child properties of a node.
type PropertySchemaQueryRow struct {
	Name   string `db:"name"`
	Type   string `db:"type"`
	Source string `db:"source"`
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

// SiblingProperties yields the properties that are set on the node and also the
// properties that don't have values by merging in the unused property schemas.
func (r *PropertySchemaTable) SiblingProperties(in []*ent.Property) PropertyTable {
	if r == nil {
		return nil
	}

	propMap := lo.KeyBy(r.siblingSchemas, func(r PropertySchemaQueryRow) string { return r.Name })

	out := make(PropertyTable, len(in))

	// Add all the properties that have values.
	for i, p := range in {
		if _, ok := propMap[p.Name]; ok {
			delete(propMap, p.Name)
		}
		out[i] = MapProperty(p)
	}

	// Add the remaining property schemas that do not have values.
	for _, p := range propMap {
		out = append(out, Property{
			PropertySchema: MapPropertySchema(p),
		})
	}
	return out
}

func (r PropertySchemaTable) ChildSchemas() PropertySchemas {
	return dt.Map(r.childSchemas, func(r PropertySchemaQueryRow) PropertySchema {
		return PropertySchema{
			Name: r.Name,
			Type: r.Type,
		}
	})
}
