package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type React struct {
	ent.Schema
}

func (React) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (React) Fields() []ent.Field {
	return []ent.Field{
		field.String("account_id").GoType(xid.ID{}),
		field.String("post_id").GoType(xid.ID{}),
		field.String("emoji"),
	}
}

func (React) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Ref("reacts").
			Field("account_id").
			Unique().
			Required(),

		edge.From("Post", Post.Type).
			Ref("reacts").
			Field("post_id").
			Unique().
			Required(),
	}
}
