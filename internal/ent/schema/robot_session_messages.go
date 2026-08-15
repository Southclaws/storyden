package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
	adksession "google.golang.org/adk/v2/session"
)

type RobotSessionEvent adksession.Event

func NewRobotSessionEvent(event adksession.Event) RobotSessionEvent {
	return RobotSessionEvent(event)
}

func (event RobotSessionEvent) ADK() adksession.Event {
	return adksession.Event(event)
}

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

		field.String("branch").
			Optional().
			Comment("ADK agent branch path used to attribute delegated work"),

		field.String("isolation_scope").
			Optional().
			Comment("ADK exact-match history scope for private sub-agent events"),

		field.String("robot_id").
			GoType(xid.ID{}).
			Optional().
			Nillable().
			Comment("Database Robot ID that generated this message. Null for built-in Robot or human-authored messages."),

		field.String("builtin_robot").
			Optional().
			Nillable().
			Comment("Built-in Robot ID that generated this message. Set only when robot_id and account_id are null."),

		field.String("account_id").
			GoType(xid.ID{}).
			Optional().
			Nillable().
			Comment("Author account ID from ADK Event, optional for system messages"),

		field.JSON("event_data", RobotSessionEvent{}).
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
