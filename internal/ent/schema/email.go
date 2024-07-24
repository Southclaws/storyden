package schema

import (
	"net/mail"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type Email struct {
	ent.Schema
}

func (Email) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (Email) Fields() []ent.Field {
	return []ent.Field{
		field.String("account_id").
			Optional().
			Nillable().
			GoType(xid.ID{}),

		field.String("email_address").
			NotEmpty().
			Immutable().
			Unique().
			MinLen(3).
			MaxLen(254).Validate(func(s string) error {
			_, err := mail.ParseAddress(s)
			return err
		}),

		field.Bool("verified").
			Default(false).
			Annotations(entsql.Default("false")).
			Comment("Whether this email has been verified to be owned by the account via a token send+verify process"),

		field.Bool("is_auth").
			Default(false).
			Annotations(entsql.Default("false")).
			Comment("Whether this email is used for authentication"),
	}
}

func (Email) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Ref("emails").
			Field("account_id").
			Unique(),
	}
}
