package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type OAuthAuthorisationCode struct {
	ent.Schema
}

func (OAuthAuthorisationCode) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (OAuthAuthorisationCode) Fields() []ent.Field {
	return []ent.Field{
		field.String("client_id").GoType(xid.ID{}).Immutable(),
		field.String("account_id").GoType(xid.ID{}).Immutable(),
		field.String("code_hash").NotEmpty().Immutable(),
		field.String("redirect_uri").NotEmpty().Immutable(),
		field.String("scope"),
		field.String("nonce").Optional().Nillable().Immutable(),
		field.String("code_challenge").NotEmpty().Immutable(),
		field.Enum("code_challenge_method").Values("S256").Default("S256").Immutable(),
		field.Time("expires_at").Immutable(),
		field.Time("consumed_at").Optional().Nillable(),
	}
}

func (OAuthAuthorisationCode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("client", OAuthClient.Type).
			Ref("authorisation_codes").
			Field("client_id").
			Required().
			Immutable().
			Unique(),

		edge.From("account", Account.Type).
			Ref("oauth_authorisation_codes").
			Field("account_id").
			Required().
			Immutable().
			Unique(),
	}
}

func (OAuthAuthorisationCode) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code_hash").Unique(),
	}
}
