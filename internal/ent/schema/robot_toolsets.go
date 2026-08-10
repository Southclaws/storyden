package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

// RobotToolset is a reusable collection of Robot tools with optional prompt
// guidance. Runtime-provided system and plugin Toolsets live in the service
// registry; this schema stores the Toolsets authored by Storyden members.
type RobotToolset struct {
	ent.Schema
}

func (RobotToolset) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (RobotToolset) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),

		field.String("description").
			Optional().
			Comment("Human-readable summary used for Toolset discovery"),

		field.String("instruction").
			Optional().
			Comment("Prompt guidance injected whenever the Toolset is active"),

		field.Strings("tools").
			Optional().
			Comment("Stable native, MCP, or plugin tool identifiers in this Toolset"),

		field.JSON("metadata", map[string]any{}).
			Optional().
			Comment("Arbitrary metadata used by clients"),

		field.String("author_id").GoType(xid.ID{}),
	}
}

func (RobotToolset) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("author", Account.Type).
			Field("author_id").
			Ref("robot_toolsets").
			Unique().
			Required(),
	}
}

func (RobotToolset) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("author_id", "name").Unique(),
	}
}
