package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type RobotMCPServer struct {
	ent.Schema
}

func (RobotMCPServer) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (RobotMCPServer) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),

		field.String("slug").
			NotEmpty().
			Immutable(),

		field.String("description").
			Optional(),

		field.String("endpoint_url").
			NotEmpty(),

		field.String("oauth_remote_connection_id").
			GoType(xid.ID{}).
			Optional().
			Nillable(),

		field.Bool("enabled").
			Default(true),

		field.String("bearer_token").
			Optional().
			Sensitive(),

		field.Time("last_refreshed_at").
			Optional().
			Nillable(),

		field.String("last_error").
			Optional().
			Nillable(),

		field.String("added_by").
			GoType(xid.ID{}),
	}
}

func (RobotMCPServer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Required().
			Ref("robot_mcp_servers").
			Field("added_by").
			Unique(),

		edge.From("oauth_remote_connection", OAuthRemoteConnection.Type).
			Ref("robot_mcp_servers").
			Field("oauth_remote_connection_id").
			Unique(),

		edge.To("tools", RobotMCPTool.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (RobotMCPServer) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("slug").Unique(),
	}
}
