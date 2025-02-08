package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
)

type PropertySchema struct {
	ent.Schema
}

func (PropertySchema) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}}
}

func (PropertySchema) Fields() []ent.Field {
	return []ent.Field{}
}

func (PropertySchema) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("node", Node.Type),
		edge.To("fields", PropertySchemaField.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
