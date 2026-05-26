package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type OAuthDeviceAuthorisation struct {
	ent.Schema
}

func (OAuthDeviceAuthorisation) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (OAuthDeviceAuthorisation) Fields() []ent.Field {
	return []ent.Field{
		field.String("client_id").GoType(xid.ID{}).Immutable(),
		field.String("device_code_hash").NotEmpty().Immutable(),
		field.String("user_code_hash").NotEmpty().Immutable(),
		field.String("user_code_display").NotEmpty().Immutable(),
		field.String("scope"),
		field.Time("expires_at").Immutable(),
		field.Int("poll_interval_seconds").Default(5),
		field.Time("last_polled_at").Optional().Nillable(),
		field.String("claimed_by_account_id").GoType(xid.ID{}).Optional().Nillable(),
		field.String("approved_by_account_id").GoType(xid.ID{}).Optional().Nillable(),
		field.Time("approved_at").Optional().Nillable(),
		field.Time("denied_at").Optional().Nillable(),
		field.Time("consumed_at").Optional().Nillable(),
	}
}

func (OAuthDeviceAuthorisation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("client", OAuthClient.Type).
			Ref("device_authorisations").
			Field("client_id").
			Required().
			Immutable().
			Unique(),

		edge.From("claimed_by_account", Account.Type).
			Ref("claimed_oauth_device_authorisations").
			Field("claimed_by_account_id").
			Unique(),

		edge.From("approved_by_account", Account.Type).
			Ref("approved_oauth_device_authorisations").
			Field("approved_by_account_id").
			Unique(),
	}
}

func (OAuthDeviceAuthorisation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("device_code_hash").Unique(),
		index.Fields("user_code_hash"),
	}
}
