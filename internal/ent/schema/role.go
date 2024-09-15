package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Role struct {
	ent.Schema
}

func (Role) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
		field.String("colour").Default("hsl(157, 65%, 44%)"),
		field.Strings("permissions"),
	}
}

func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("accounts", Account.Type),
	}
}
