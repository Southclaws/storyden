package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type TrailActionRun struct {
	ent.Schema
}

func (TrailActionRun) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (TrailActionRun) Fields() []ent.Field {
	return []ent.Field{
		field.String("run_id").GoType(xid.ID{}).Immutable(),
		field.String("action_id").GoType(xid.ID{}).Immutable(),
		field.String("kind").Immutable(),
		field.JSON("config", json.RawMessage{}).Immutable(),
		field.Enum("status").
			Values("queued", "running", "completed", "blocked", "failed", "cancelled"),
		field.String("lease_token").Optional().Nillable(),
		field.Time("lease_expires_at").Optional().Nillable(),
		field.JSON("output", json.RawMessage{}).Optional(),
		field.JSON("target", json.RawMessage{}).Optional(),
		field.String("error_text").Optional().Nillable(),
		field.Time("started_at").Optional().Nillable(),
		field.Time("finished_at").Optional().Nillable(),
		field.Time("notified_at").Optional().Nillable(),
	}
}

func (TrailActionRun) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("run_id", "action_id").Unique(),
		index.Fields("status", "lease_expires_at", "created_at"),
	}
}

func (TrailActionRun) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("run", TrailRun.Type).
			Ref("action_runs").
			Field("run_id").
			Immutable().
			Unique().
			Required(),
		edge.From("action", TrailAction.Type).
			Ref("runs").
			Field("action_id").
			Immutable().
			Unique().
			Required(),
	}
}
