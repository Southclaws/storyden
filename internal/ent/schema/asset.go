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
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (Asset) Fields() []ent.Field {
	return []ent.Field{
		field.String("url"),
		field.String("mimetype"),
		field.Int("width"),
		field.Int("height"),

		// Edges
		field.String("post_id").GoType(xid.ID{}),
	}
}

// Edges of Asset.
func (Asset) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("post", Post.Type).
			Field("post_id").
			Ref("assets").
			Unique().
			Required(),
	}
}
