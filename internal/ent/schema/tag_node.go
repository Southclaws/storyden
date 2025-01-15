package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type TagNode struct {
	ent.Schema
}

func (TagNode) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}}
}

func (TagNode) Fields() []ent.Field {
	return []ent.Field{
		field.String("tag_id").GoType(xid.ID{}),
		field.String("node_id").GoType(xid.ID{}),
	}
}

func (TagNode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tag", Tag.Type).
			Field("tag_id").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("node", Node.Type).
			Field("node_id").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (TagNode) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("node_id"),
	}
}
