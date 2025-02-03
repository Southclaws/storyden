package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type PropertySchemaField struct {
	ent.Schema
}

func (PropertySchemaField) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}}
}

func (PropertySchemaField) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("type"),
		field.String("schema_id").GoType(xid.ID{}),
	}
}

func (PropertySchemaField) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("schema", PropertySchema.Type).
			Field("schema_id").
			Ref("fields").
			Required().
			Unique(),

		edge.To("properties", Property.Type),
	}
}

func (PropertySchemaField) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("schema_id", "name").Unique(),
		index.Fields("name"),
	}
}
