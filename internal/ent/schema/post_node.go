package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type PostNode struct {
	ent.Schema
}

func (PostNode) Mixin() []ent.Mixin {
	return []ent.Mixin{CreatedAt{}}
}

func (PostNode) Fields() []ent.Field {
	return []ent.Field{
		field.String("node_id").
			MaxLen(20).
			NotEmpty().
			GoType(xid.ID{}).
			DefaultFunc(func() xid.ID { return xid.New() }),

		field.String("post_id").
			MaxLen(20).
			NotEmpty().
			GoType(xid.ID{}).
			DefaultFunc(func() xid.ID { return xid.New() }),
	}
}

func (PostNode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("node", Node.Type).
			Unique().
			Required().
			Field("node_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("post", Post.Type).
			Unique().
			Required().
			Field("post_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (PostNode) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("node_id", "post_id").
			Unique().
			StorageKey("unique_post_node"),

		index.Fields("node_id", "created_at").
			StorageKey("idx_post_node_paginate"),
	}
}

func (PostNode) Annotations() []schema.Annotation {
	return []schema.Annotation{
		field.ID("node_id", "post_id"),
	}
}
