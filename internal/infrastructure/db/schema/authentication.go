package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type Authentication struct {
	ent.Schema
}

func (Authentication) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Default(uuid.New),

		field.String("service").
			NotEmpty().
			Comment("The authentication service name, such as GitHub, Twitter, Discord, etc. Or, 'password' for password auth and 'api_token' for token auth"),

		field.String("identifier").
			Comment("The identifier, usually a user/account ID on some OAuth service or API token name. If it's a password, this is blank."),

		field.String("token").
			NotEmpty().
			Sensitive().
			Comment("The actual authentication token/password/key/etc. If OAuth, it'll be the access_token value, if it's a password, a hash and if it's an api_token type then the API token string."),

		field.JSON("metadata", map[string]interface{}{}).
			Optional().
			Comment("Any necessary metadata specific to the authentication method."),
	}
}

func (Authentication) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Ref("authentication").
			Unique(),
	}
}

func (Authentication) Indexes() []ent.Index {
	return []ent.Index{
		// a given identifier should be unique within the context of a service.
		index.Fields("service", "identifier").
			Unique(),
	}
}

func (Authentication) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.CreateTime{},
	}
}
