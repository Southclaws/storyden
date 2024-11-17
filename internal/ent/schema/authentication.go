package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type Authentication struct {
	ent.Schema
}

func (Authentication) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (Authentication) Fields() []ent.Field {
	return []ent.Field{
		field.String("service").
			NotEmpty().
			Comment("The authentication service name, such as GitHub, Twitter, Discord, etc. Or, 'password' for password auth and 'api_token' for token auth"),

		field.String("token_type").
			NotEmpty().
			Comment("The type of secret/token used by the service to secure the authentication record."),

		field.String("identifier").
			Comment("The identifier, usually a user/account ID on some OAuth service or API token name."),

		field.String("token").
			NotEmpty().
			Sensitive().
			Comment("The actual authentication token/password/key/etc. If OAuth, it'll be the access_token value, if it's a password, a hash and if it's an api_token type then the API token string."),

		field.String("name").
			Optional().
			Nillable().
			Comment("A human-readable name for the authentication method. For WebAuthn, this may be the device OS or nickname."),

		field.JSON("metadata", map[string]interface{}{}).
			Optional().
			Comment("Any necessary metadata specific to the authentication method."),

		field.String("account_authentication").GoType(xid.ID{}),

		field.String("email_address_record_id").
			GoType(xid.ID{}).
			Optional().
			Nillable().
			NotEmpty(),
	}
}

func (Authentication) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Ref("authentication").
			Field("account_authentication").
			Required().
			Unique(),

		// Only one auth method may be linked to an email address.
		edge.From("email_address", Email.Type).
			Field("email_address_record_id").
			Ref("authentication_record").
			Unique(),
	}
}

func (Authentication) Indexes() []ent.Index {
	return []ent.Index{
		// Each pair of service and identifier can only exist once.
		index.Fields("service", "identifier", "account_authentication").
			Unique(),

		// Each pair of token type and identifier can only exist once.
		index.Fields("token_type", "identifier", "account_authentication").
			Unique(),
	}
}
