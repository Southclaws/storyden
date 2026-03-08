package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type RobotSessionMessage struct {
	ent.Schema
}

func (RobotSessionMessage) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (RobotSessionMessage) Fields() []ent.Field {
	return []ent.Field{
		field.String("session_id").GoType(xid.ID{}),

		field.String("invocation_id").
			Comment("Invocation ID from ADK Event"),

		field.String("robot_id").
			GoType(xid.ID{}).
			Optional().
			Nillable().
			Comment("Robot ID - which agent generated this message. Null means the built-in Storyden agent or a user message"),

		field.String("account_id").
			GoType(xid.ID{}).
			Optional().
			Nillable().
			Comment("Author account ID from ADK Event, optional for system messages"),

		field.JSON("event_data", map[string]any{}).
			Comment("Full ADK Event object stored as JSON"),
	}
}

func (RobotSessionMessage) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("session_id", "created_at"),
	}
}

func (RobotSessionMessage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("session", RobotSession.Type).
			Field("session_id").
			Ref("messages").
			Unique().
			Required(),

		edge.From("robot", Robot.Type).
			Field("robot_id").
			Ref("messages").
			Unique(),

		edge.From("author", Account.Type).
			Field("account_id").
			Ref("robot_messages").
			Unique(),
	}
}
