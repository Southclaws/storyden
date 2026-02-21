package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
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
		field.Float("sort_key").Annotations(entsql.Default("0.0")),
		field.JSON("metadata", map[string]any{}).Optional(),
	}
}

func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("accounts", Account.Type).
			Through("account_roles", AccountRoles.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
