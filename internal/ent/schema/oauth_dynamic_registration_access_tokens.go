package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	entschema "entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type OAuthDynamicRegistrationAccessTokens struct {
	ent.Schema
}

func (OAuthDynamicRegistrationAccessTokens) Annotations() []entschema.Annotation {
	return []entschema.Annotation{entsql.Annotation{Table: "oauth_dcr_iats"}}
}

func (OAuthDynamicRegistrationAccessTokens) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (OAuthDynamicRegistrationAccessTokens) Fields() []ent.Field {
	return []ent.Field{
		field.String("label").NotEmpty().MaxLen(200).Immutable(),
		field.String("token_hash").NotEmpty().Unique().Immutable().Sensitive(),
		field.String("creator_account_id").GoType(xid.ID{}).Immutable(),
		field.Time("expires_at").Immutable(),
		field.Int("max_registrations").Positive().Immutable(),
		field.Int("registration_count").NonNegative().Default(0),
		field.Time("revoked_at").Optional().Nillable(),
	}
}

func (OAuthDynamicRegistrationAccessTokens) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("creator", Account.Type).Ref("oauth_dcr_iats").Field("creator_account_id").Required().Immutable().Unique(),
		edge.To("clients", OAuthClient.Type),
	}
}
