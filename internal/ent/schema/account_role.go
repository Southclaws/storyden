package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type AccountRoles struct {
	ent.Schema
}

func (AccountRoles) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (AccountRoles) Fields() []ent.Field {
	return []ent.Field{
		field.String("account_id").GoType(xid.ID{}),
		field.String("role_id").GoType(xid.ID{}),
		field.Bool("badge").Optional().Nillable(),
	}
}

func (AccountRoles) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("account", Account.Type).
			Field("account_id").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("role", Role.Type).
			Field("role_id").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (AccountRoles) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("account_id", "role_id").
			Unique().
			StorageKey("unique_account_role"),

		index.Fields("account_id", "badge").
			Unique(),
	}
}
