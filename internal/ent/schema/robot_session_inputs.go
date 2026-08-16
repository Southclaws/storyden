package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

// RobotSessionInput is durable work waiting to be consumed by a Robot turn.
// Multiple compatible inputs may be claimed by one turn.
type RobotSessionInput struct {
	ent.Schema
}

func (RobotSessionInput) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (RobotSessionInput) Fields() []ent.Field {
	return []ent.Field{
		field.String("session_id").GoType(xid.ID{}),
		field.String("account_id").GoType(xid.ID{}),
		field.String("turn_id").GoType(xid.ID{}).Optional().Nillable(),
		field.Uint64("sequence").Comment("Monotonic submission order within the session."),
		field.String("source_kind"),
		field.String("batch_key").Comment("Inputs with the same key may be claimed by one turn."),
		field.JSON("input_data", json.RawMessage{}),
		field.Time("not_before").
			Optional().
			Nillable().
			Comment("Earliest time this input may be claimed by a turn."),
		field.Enum("status").Values("queued", "claimed", "cancelled").Default("queued"),
	}
}

func (RobotSessionInput) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("session_id", "status", "not_before", "sequence"),
		index.Fields("session_id", "turn_id", "sequence"),
	}
}

func (RobotSessionInput) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("session", RobotSession.Type).
			Field("session_id").
			Ref("inputs").
			Unique().
			Required(),
		edge.From("account", Account.Type).
			Field("account_id").
			Ref("robot_session_inputs").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("turn", RobotSessionTurn.Type).
			Field("turn_id").
			Ref("inputs").
			Unique(),
	}
}
