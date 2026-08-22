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

type TrailAction struct {
	ent.Schema
}

func (TrailAction) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (TrailAction) Fields() []ent.Field {
	return []ent.Field{
		field.String("trail_id").GoType(xid.ID{}),
		field.String("kind"),
		field.Int("position").Default(0),
		field.JSON("config", json.RawMessage{}),
		field.Time("archived_at").Optional().Nillable(),
	}
}

func (TrailAction) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("trail_id", "position"),
		index.Fields("trail_id", "archived_at"),
	}
}

func (TrailAction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("trail", Trail.Type).
			Ref("actions").
			Field("trail_id").
			Unique().
			Required(),
		edge.To("runs", TrailActionRun.Type).
			Annotations(entsql.OnDelete(entsql.Restrict)),
	}
}
