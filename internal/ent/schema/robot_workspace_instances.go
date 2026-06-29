package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type RobotWorkspaceInstance struct {
	ent.Schema
}

func (RobotWorkspaceInstance) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (RobotWorkspaceInstance) Fields() []ent.Field {
	return []ent.Field{
		field.String("workspace_id").
			GoType(xid.ID{}),

		field.String("created_by").
			GoType(xid.ID{}),

		field.Enum("provider").
			Values("local", "sprites").
			Default("local").
			Comment("Workspace provider type"),

		field.JSON("provider_state", map[string]any{}).
			Optional().
			Comment("Provider-specific live instance state"),

		field.JSON("metadata", map[string]any{}).
			Optional().
			Comment("Arbitrary metadata used by clients to store domain specific information"),
	}
}

func (RobotWorkspaceInstance) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("workspace", RobotWorkspace.Type).
			Field("workspace_id").
			Ref("instances").
			Unique().
			Required(),

		edge.From("creator", Account.Type).
			Field("created_by").
			Ref("robot_workspace_instances").
			Unique().
			Required(),
	}
}
