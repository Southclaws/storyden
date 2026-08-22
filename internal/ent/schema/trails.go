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

type Trail struct {
	ent.Schema
}

func (Trail) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (Trail) Fields() []ent.Field {
	return []ent.Field{
		field.String("account_id").GoType(xid.ID{}),
		field.String("name").NotEmpty(),
		field.String("description").Default(""),
		field.Enum("status").
			Values("active", "paused", "finished", "archived"),
		field.String("trigger_type"),
		field.JSON("trigger_config", json.RawMessage{}),
		field.Time("next_occurrence_at").Optional().Nillable(),
		field.Time("last_occurrence_at").Optional().Nillable(),
	}
}

func (Trail) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status", "next_occurrence_at"),
		index.Fields("account_id", "updated_at"),
	}
}

func (Trail) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("creator", Account.Type).
			Ref("created_trails").
			Field("account_id").
			Unique().
			Required(),
		edge.To("actions", TrailAction.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("runs", TrailRun.Type).
			Annotations(entsql.OnDelete(entsql.Restrict)),
	}
}
