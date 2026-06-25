package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type RobotWorkspace struct {
	ent.Schema
}

func (RobotWorkspace) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (RobotWorkspace) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Comment("Human-readable name of the workspace template"),

		field.String("description").
			Optional().
			Comment("Human-readable description of the workspace template"),

		field.Enum("provider").
			Values("local").
			Default("local").
			Comment("Workspace provider type"),

		field.JSON("config", map[string]any{}).
			Optional().
			Comment("Provider-specific workspace template configuration"),

		field.JSON("metadata", map[string]any{}).
			Optional().
			Comment("Arbitrary metadata used by clients to store domain specific information"),

		field.String("created_by").
			GoType(xid.ID{}),
	}
}

func (RobotWorkspace) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("creator", Account.Type).
			Field("created_by").
			Ref("robot_workspaces").
			Unique().
			Required(),

		edge.To("robots", Robot.Type).
			Annotations(entsql.OnDelete(entsql.SetNull)),

		edge.To("instances", RobotWorkspaceInstance.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
