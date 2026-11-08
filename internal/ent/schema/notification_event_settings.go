package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type NotificationEventSettings struct {
	ent.Schema
}

func (NotificationEventSettings) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}}
}

func (NotificationEventSettings) Fields() []ent.Field {
	return []ent.Field{
		field.String("account_id").GoType(xid.ID{}).Immutable(),
		field.String("event").NotEmpty().Immutable(),
		field.Bool("inbox_enabled"),
		field.Bool("push_enabled"),
	}
}

func (NotificationEventSettings) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Ref("notification_event_settings").
			Field("account_id").
			Required().
			Unique().
			Immutable(),
	}
}

func (NotificationEventSettings) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("account_id", "event").Unique(),
		index.Fields("event", "inbox_enabled", "account_id"),
	}
}
