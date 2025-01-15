package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type TagPost struct {
	ent.Schema
}

func (TagPost) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}}
}

func (TagPost) Fields() []ent.Field {
	return []ent.Field{
		field.String("tag_id").GoType(xid.ID{}),
		field.String("post_id").GoType(xid.ID{}),
	}
}

func (TagPost) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tag", Tag.Type).
			Field("tag_id").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("post", Post.Type).
			Field("post_id").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (TagPost) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("post_id"),
	}
}
