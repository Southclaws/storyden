package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type OAuthRegistrationApproval struct{ ent.Schema }

func (OAuthRegistrationApproval) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (OAuthRegistrationApproval) Fields() []ent.Field {
	return []ent.Field{
		field.String("registration_code_hash").NotEmpty().Unique().Immutable().Sensitive(),
		field.String("verification_code_hash").NotEmpty().Unique().Immutable().Sensitive(),
		field.String("verification_code_display").NotEmpty().Immutable(),
		field.JSON("metadata", json.RawMessage{}).Immutable(),
		field.Time("expires_at").Immutable(),
		field.Time("next_poll_at"),
		field.Int("poll_interval_seconds").Positive().Default(5),
		field.String("approved_by_account_id").GoType(xid.ID{}).Optional().Nillable(),
		field.Time("approved_at").Optional().Nillable(),
		field.Time("denied_at").Optional().Nillable(),
		field.Time("cancelled_at").Optional().Nillable(),
		field.Time("consumed_at").Optional().Nillable(),
	}
}

func (OAuthRegistrationApproval) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("approved_by_account", Account.Type).Ref("oauth_registration_approvals").Field("approved_by_account_id").Unique(),
	}
}
