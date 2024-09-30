package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type EventParticipant struct {
	ent.Schema
}

func (EventParticipant) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (EventParticipant) Fields() []ent.Field {
	return []ent.Field{
		field.String("role"),
		field.String("status"),
		field.String("account_id").GoType(xid.ID{}),
		field.String("event_id").GoType(xid.ID{}),
	}
}

func (EventParticipant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Ref("events").
			Field("account_id").
			Unique().
			Required(),

		edge.From("event", Event.Type).
			Ref("participants").
			Field("event_id").
			Unique().
			Required(),
	}
}

func (EventParticipant) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("account_id", "event_id").
			Unique().
			StorageKey("unique_event_participant"),
	}
}
