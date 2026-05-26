package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type OAuthClient struct {
	ent.Schema
}

func (OAuthClient) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (OAuthClient) Fields() []ent.Field {
	return []ent.Field{
		field.String("account_id").GoType(xid.ID{}).Optional().Nillable().Immutable(),
		field.String("client_id").Unique().NotEmpty(),
		field.String("client_secret_hash").Optional().Nillable(),
		field.String("name").NotEmpty(),
		field.Enum("type").Values("public", "confidential").Default("public"),
		field.Enum("scope_policy").Values("explicit", "inherit").Default("explicit"),
		field.JSON("redirect_uris", []string{}),
		field.JSON("allowed_scopes", []string{}),
		field.JSON("allowed_grants", []string{}),
	}
}

func (OAuthClient) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Ref("oauth_clients").
			Field("account_id").
			Immutable().
			Unique(),

		edge.To("authorisation_codes", OAuthAuthorisationCode.Type),
		edge.To("authorisation_requests", OAuthAuthorisationRequest.Type),
		edge.To("device_authorisations", OAuthDeviceAuthorisation.Type),
		edge.To("refresh_tokens", OAuthRefreshToken.Type),
	}
}

func (OAuthClient) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("client_id").Unique(),
	}
}
