package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type Notification struct {
	ent.Schema
}

func (Notification) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, DeletedAt{}}
}

func (Notification) Fields() []ent.Field {
	return []ent.Field{
		field.String("event_type"),

		field.String("datagraph_kind").
			Optional().
			Nillable(),

		field.String("datagraph_id").
			GoType(xid.ID{}).
			Optional().
			Nillable().
			Comment("The ID of the resource that this notification relates to. This is not a foreign key as notifications can refer to a variety of sources, discriminated by the 'datagraph_kind' field."),

		field.Bool("read"),

		field.String("owner_account_id").
			GoType(xid.ID{}),

		field.String("source_account_id").
			GoType(xid.ID{}).
			Optional().
			Nillable(),
	}
}

func (Notification) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", Account.Type).
			Field("owner_account_id").
			Ref("notifications").
			Required().
			Unique(),

		edge.From("source", Account.Type).
			Field("source_account_id").
			Ref("triggered_notifications").
			Unique(),
	}
}
