package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type Asset struct {
	ent.Schema
}

func (Asset) Mixin() []ent.Mixin {
	return []ent.Mixin{CreatedAt{}, UpdatedAt{}}
}

func (Asset) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			NotEmpty().
			Immutable().
			Unique(),

		field.String("url"),
		field.String("mimetype"),
		field.Int("width"),
		field.Int("height"),

		// Edges
		field.String("account_id").GoType(xid.ID{}),
	}
}

func (Asset) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("posts", Post.Type).
			Ref("assets"),

		edge.From("owner", Account.Type).
			Field("account_id").
			Ref("assets").
			Unique().
			Required(),
	}
}
