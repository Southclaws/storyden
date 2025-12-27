package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type Report struct {
	ent.Schema
}

func (Report) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (Report) Fields() []ent.Field {
	return []ent.Field{
		field.String("target_id").
			GoType(xid.ID{}).
			Comment("The ID of the resource being reported. This is not a foreign key as reports can refer to a variety of sources, discriminated by the 'target_kind' field."),

		field.String("target_kind").
			Comment("The datagraph kind of resource being reported."),

		field.String("reported_by_id").
			GoType(xid.ID{}).
			Optional().
			Nillable(),

		field.String("handled_by_id").
			GoType(xid.ID{}).
			Optional().
			Nillable(),

		field.Text("comment").
			Optional().
			Nillable(),

		field.String("reason").
			Optional().
			Nillable(),

		field.String("status").
			Default("submitted"),
	}
}

func (Report) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("reported_by", Account.Type).
			Field("reported_by_id").
			Ref("reports").
			Unique(),

		edge.From("handled_by", Account.Type).
			Field("handled_by_id").
			Ref("handled_reports").
			Unique(),
	}
}
