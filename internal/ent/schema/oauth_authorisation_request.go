package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type OAuthAuthorisationRequest struct {
	ent.Schema
}

func (OAuthAuthorisationRequest) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (OAuthAuthorisationRequest) Fields() []ent.Field {
	return []ent.Field{
		field.String("client_id").GoType(xid.ID{}).Immutable(),
		field.String("account_id").GoType(xid.ID{}).Immutable(),
		field.String("request_id_hash").NotEmpty().Immutable(),
		field.String("redirect_uri").NotEmpty().Immutable(),
		field.String("scope"),
		field.String("state").Optional().Nillable().Immutable(),
		field.String("nonce").Optional().Nillable().Immutable(),
		field.String("code_challenge").NotEmpty().Immutable(),
		field.Enum("code_challenge_method").Values("S256").Default("S256").Immutable(),
		field.Time("expires_at").Immutable(),
		field.Time("approved_at").Optional().Nillable(),
		field.Time("denied_at").Optional().Nillable(),
	}
}

func (OAuthAuthorisationRequest) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("client", OAuthClient.Type).
			Ref("authorisation_requests").
			Field("client_id").
			Required().
			Immutable().
			Unique(),

		edge.From("account", Account.Type).
			Ref("oauth_authorisation_requests").
			Field("account_id").
			Required().
			Immutable().
			Unique(),
	}
}

func (OAuthAuthorisationRequest) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("request_id_hash").Unique(),
	}
}
