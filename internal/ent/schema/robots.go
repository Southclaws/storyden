package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type Robot struct {
	ent.Schema
}

func (Robot) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (Robot) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique().
			NotEmpty(),

		field.String("description").
			Optional().
			Comment("Human-readable description of the robot's purpose"),

		field.String("playbook").
			Comment("Primary drive/directive (system prompt) for the agent"),

		field.Strings("tools").
			Optional().
			Comment("A list of tool names that the robot can use"),

		field.JSON("metadata", map[string]any{}).
			Optional().
			Comment("Arbitrary metadata used by clients to store domain specific information"),

		field.String("author_id").GoType(xid.ID{}),
	}
}

func (Robot) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("author", Account.Type).
			Field("author_id").
			Ref("robots").
			Unique().
			Required(),

		edge.To("messages", RobotSessionMessage.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
