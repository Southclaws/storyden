package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type RobotMCPTool struct {
	ent.Schema
}

func (RobotMCPTool) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (RobotMCPTool) Fields() []ent.Field {
	return []ent.Field{
		field.String("server_id").
			GoType(xid.ID{}),

		field.String("tool_id").
			NotEmpty().
			Immutable(),

		field.String("remote_name").
			NotEmpty(),

		field.String("callable_name").
			NotEmpty(),

		field.String("title").
			Optional(),

		field.String("description").
			Optional(),

		field.JSON("input_schema", map[string]any{}).
			Optional(),

		field.JSON("output_schema", map[string]any{}).
			Optional(),

		field.JSON("annotations", map[string]any{}).
			Optional(),

		field.Bool("enabled").
			Default(true),

		field.Time("last_seen_at").
			Default(time.Now),
	}
}

func (RobotMCPTool) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("server", RobotMCPServer.Type).
			Required().
			Ref("tools").
			Field("server_id").
			Unique(),
	}
}

func (RobotMCPTool) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tool_id").Unique(),
		index.Fields("server_id", "callable_name").Unique(),
		index.Fields("server_id", "remote_name").Unique(),
	}
}
