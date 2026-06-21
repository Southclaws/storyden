package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type OAuthRemoteAuthorisationFlow struct {
	ent.Schema
}

func (OAuthRemoteAuthorisationFlow) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (OAuthRemoteAuthorisationFlow) Fields() []ent.Field {
	return []ent.Field{
		field.String("connection_id").GoType(xid.ID{}).Immutable(),
		field.String("state_hash").Unique().NotEmpty().Immutable(),
		field.String("pkce_verifier").NotEmpty().Sensitive().Immutable(),
		field.String("redirect_uri").NotEmpty().Immutable(),
		field.Time("expires_at").Immutable(),
		field.Time("consumed_at").Optional().Nillable(),
	}
}

func (OAuthRemoteAuthorisationFlow) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("connection", OAuthRemoteConnection.Type).
			Ref("authorisation_flows").
			Field("connection_id").
			Required().
			Immutable().
			Unique(),
	}
}
