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

type TrailRun struct {
	ent.Schema
}

func (TrailRun) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (TrailRun) Fields() []ent.Field {
	return []ent.Field{
		field.String("trail_id").GoType(xid.ID{}),
		field.String("initiated_by_id").GoType(xid.ID{}).Optional().Nillable(),
		field.Enum("kind").Values("scheduled", "manual"),
		field.JSON("trigger_payload", json.RawMessage{}),
		field.Time("scheduled_for").Optional().Nillable(),
		field.Enum("status").
			Values("queued", "running", "completed", "attention_required", "cancelled", "skipped"),
		field.Time("finished_at").Optional().Nillable(),
	}
}

func (TrailRun) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("trail_id", "created_at"),
		index.Fields("trail_id", "scheduled_for").Unique(),
		index.Fields("status", "created_at"),
	}
}

func (TrailRun) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("trail", Trail.Type).
			Ref("runs").
			Field("trail_id").
			Unique().
			Required(),
		edge.From("initiator", Account.Type).
			Ref("initiated_trail_runs").
			Field("initiated_by_id").
			Unique(),
		edge.To("action_runs", TrailActionRun.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
