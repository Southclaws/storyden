package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type Property struct {
	ent.Schema
}

func (Property) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (Property) Fields() []ent.Field {
	return []ent.Field{
		field.String("node_id").GoType(xid.ID{}),
		field.String("field_id").GoType(xid.ID{}),
		field.String("value"),
	}
}

func (Property) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("node", Node.Type).
			Field("node_id").
			Ref("properties").
			Required().
			Unique(),

		edge.From("schema", PropertySchemaField.Type).
			Field("field_id").
			Ref("properties").
			Required().
			Unique().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Property) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("field_id", "node_id").Unique(),
		index.Fields("field_id"),
		index.Fields("node_id"),
	}
}
