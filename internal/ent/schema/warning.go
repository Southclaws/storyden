package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type Warning struct {
	ent.Schema
}

func (Warning) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (Warning) Fields() []ent.Field {
	return []ent.Field{
		field.String("account_id").
			GoType(xid.ID{}),
		field.String("author_id").
			GoType(xid.ID{}).
			Optional().
			Nillable(),
		field.String("reason").
			NotEmpty().
			MaxLen(2000),
	}
}

func (Warning) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Field("account_id").
			Ref("warnings").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("author", Account.Type).
			Field("author_id").
			Ref("authored_warnings").
			Unique().
			Annotations(entsql.OnDelete(entsql.SetNull)),
	}
}
