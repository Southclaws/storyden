package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

// RobotSessionTurn is one request to invoke the Robot runtime. A turn
// exists while queued, before a worker acquires the session execution lease.
type RobotSessionTurn struct {
	ent.Schema
}

func (RobotSessionTurn) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (RobotSessionTurn) Fields() []ent.Field {
	return []ent.Field{
		field.String("session_id").GoType(xid.ID{}),
		field.String("initiated_by_account_id").
			GoType(xid.ID{}).
			Optional().
			Nillable(),
		field.String("source_kind").
			Comment("Extensible source such as user_message, tool_result, timer, task, or plugin."),
		field.String("robot_ref"),
		field.JSON("input_data", json.RawMessage{}).
			Optional(),
		field.Enum("status").
			Values("queued", "running", "completed", "blocked", "failed", "cancelled").
			Default("queued"),
		field.String("continuation_of_turn_id").
			GoType(xid.ID{}).
			Optional().
			Nillable(),
		field.Time("started_at").
			Optional().
			Nillable(),
		field.Time("finished_at").
			Optional().
			Nillable(),
		field.String("error_text").
			Optional().
			Nillable(),
	}
}

func (RobotSessionTurn) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("session_id", "status", "created_at"),
		index.Fields("initiated_by_account_id", "created_at"),
	}
}

func (RobotSessionTurn) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("session", RobotSession.Type).
			Field("session_id").
			Ref("turns").
			Unique().
			Required(),
		edge.From("initiator", Account.Type).
			Field("initiated_by_account_id").
			Ref("initiated_robot_turns").
			Unique(),
	}
}
