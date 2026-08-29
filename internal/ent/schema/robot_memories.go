package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type RobotMemory struct {
	ent.Schema
}

func (RobotMemory) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (RobotMemory) Fields() []ent.Field {
	return []ent.Field{
		field.String("robot_ref").NotEmpty(),
		field.String("parent_id").GoType(xid.ID{}).Optional().Nillable(),
		field.String("content").NotEmpty(),

		// Basic RDF "facts" model.
		field.String("subject").Optional().Nillable(),
		field.String("predicate").Optional().Nillable(),
		field.String("object").Optional().Nillable(),

		field.Enum("state").Values("active", "superseded", "archived").Default("active"),
		field.Time("last_accessed_at").Default(time.Now),
		field.Uint64("access_count").Default(0),
	}
}

func (RobotMemory) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("robot_ref", "parent_id", "state"),
		index.Fields("robot_ref", "state", "updated_at"),
		index.Fields("subject"),
		index.Fields("predicate"),
		index.Fields("object"),
	}
}

func (RobotMemory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", RobotMemory.Type).
			From("parent").
			Unique().
			Field("parent_id"),
	}
}
