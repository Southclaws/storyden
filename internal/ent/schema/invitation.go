package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type Invitation struct {
	ent.Schema
}

func (Invitation) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}, DeletedAt{}}
}

func (Invitation) Fields() []ent.Field {
	return []ent.Field{
		field.String("message").
			Optional().
			Nillable(),

		field.String("creator_account_id").
			GoType(xid.ID{}),
	}
}

func (Invitation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("creator", Account.Type).
			Ref("invitations").
			Field("creator_account_id").
			Unique().
			Required(),

		edge.To("invited", Account.Type),
	}
}
