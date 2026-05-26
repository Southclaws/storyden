package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type OAuthRefreshToken struct {
	ent.Schema
}

func (OAuthRefreshToken) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (OAuthRefreshToken) Fields() []ent.Field {
	return []ent.Field{
		field.String("client_id").GoType(xid.ID{}).Immutable(),
		field.String("account_id").GoType(xid.ID{}).Immutable(),
		field.String("token_hash").NotEmpty().Immutable(),
		field.String("scope"),
		field.Time("expires_at").Immutable(),
		field.Time("revoked_at").Optional().Nillable(),
		field.String("replaced_by_token_id").GoType(xid.ID{}).Optional().Nillable(),
		field.Time("last_used_at").Optional().Nillable(),
	}
}

func (OAuthRefreshToken) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("client", OAuthClient.Type).
			Ref("refresh_tokens").
			Field("client_id").
			Required().
			Immutable().
			Unique(),

		edge.From("account", Account.Type).
			Ref("oauth_refresh_tokens").
			Field("account_id").
			Required().
			Immutable().
			Unique(),

		edge.To("replaced_by", OAuthRefreshToken.Type).
			From("replaces").
			Field("replaced_by_token_id").
			Unique().
			Annotations(entsql.OnDelete(entsql.SetNull)),
	}
}

func (OAuthRefreshToken) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("token_hash").Unique(),
	}
}
