package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type Session struct {
	ent.Schema
}

func (Session) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.String("account_id").
			Immutable().
			GoType(xid.ID{}).
			NotEmpty(),

		field.Time("expires_at").
			Immutable(),

		field.Time("revoked_at").
			Optional().
			Nillable(),
	}
}

func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Ref("sessions").
			Field("account_id").
			Required().
			Immutable().
			Unique(),
	}
}
