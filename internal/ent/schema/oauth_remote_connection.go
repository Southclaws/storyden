package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type OAuthRemoteConnection struct {
	ent.Schema
}

func (OAuthRemoteConnection) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (OAuthRemoteConnection) Fields() []ent.Field {
	return []ent.Field{
		field.String("resource_url").NotEmpty(),
		field.String("resource").Optional(),
		field.String("resource_name").Optional(),
		field.JSON("protected_resource_metadata", map[string]any{}).Optional(),
		field.String("authorization_server").Optional(),
		field.JSON("authorization_server_metadata", map[string]any{}).Optional(),
		field.Enum("mode").Values("cimd", "dcr", "manual").Default("manual"),
		field.Enum("status").Values("pending", "connected", "error").Default("pending"),
		field.String("client_id").Optional(),
		field.String("client_secret").Optional().Sensitive(),
		field.String("authorization_endpoint").Optional(),
		field.String("token_endpoint").Optional(),
		field.String("registration_endpoint").Optional(),
		field.String("token_endpoint_auth_method").Optional(),
		field.JSON("redirect_uris", []string{}).Optional(),
		field.String("redirect_uri").Optional(),
		field.String("scope").Optional(),
		field.String("access_token").Optional().Sensitive(),
		field.String("refresh_token").Optional().Sensitive(),
		field.String("token_type").Optional(),
		field.Time("token_expiry").Optional().Nillable(),
		field.Time("token_refresh_started_at").Optional().Nillable(),
		field.String("last_error").Optional().Nillable(),
		field.String("added_by").GoType(xid.ID{}),
	}
}

func (OAuthRemoteConnection) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Ref("oauth_remote_connections").
			Field("added_by").
			Required().
			Unique(),

		edge.To("authorisation_flows", OAuthRemoteAuthorisationFlow.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("robot_mcp_servers", RobotMCPServer.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (OAuthRemoteConnection) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("resource_url", "authorization_server", "added_by").Unique(),
	}
}
