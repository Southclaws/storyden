package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type RobotSession struct {
	ent.Schema
}

func (RobotSession) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (RobotSession) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			DefaultFunc(func() string {
				return "Untitled (" + time.Now().Format("January 2, 2006 at 3:04 PM") + ")"
			}).
			Comment("The name of the session (generated from the message history shortly after the session begins.)"),

		field.String("account_id").
			GoType(xid.ID{}).
			Comment("The account that created this session."),

		field.JSON("state", map[string]any{}).
			Optional().
			Comment("Session state from ADK"),

		field.Enum("execution_status").
			Values("idle", "running", "blocked").
			Default("idle").
			Comment("Serialises Robot turns for this shared session."),

		field.String("active_turn_id").
			GoType(xid.ID{}).
			Optional().
			Nillable().
			Comment("Turn currently holding or blocking this session."),

		field.String("lease_token").
			Optional().
			Nillable().
			Comment("Anonymous capability held by the active execution."),

		field.Uint64("lease_generation").
			Default(0).
			Comment("Monotonic fencing generation for session executions."),

		field.Uint64("next_event_sequence").
			Default(0).
			Comment("Allocates ordered events within this session."),

		field.Time("lease_expires_at").
			Optional().
			Nillable().
			Comment("Expiry used to recover a session after an interrupted execution."),
	}
}

func (RobotSession) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("creator", Account.Type).
			Ref("created_robot_sessions").
			Field("account_id").
			Unique().
			Required(),

		edge.To("views", RobotSessionView.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("messages", RobotSessionMessage.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("inputs", RobotSessionInput.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("turns", RobotSessionTurn.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
