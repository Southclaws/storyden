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

type CollectionPost struct {
	ent.Schema
}

func (CollectionPost) Mixin() []ent.Mixin {
	return []ent.Mixin{CreatedAt{}}
}

func (CollectionPost) Fields() []ent.Field {
	return []ent.Field{
		field.String("collection_id").
			MaxLen(20).
			NotEmpty().
			GoType(xid.ID{}).
			DefaultFunc(func() xid.ID { return xid.New() }),

		field.String("post_id").
			MaxLen(20).
			NotEmpty().
			GoType(xid.ID{}).
			DefaultFunc(func() xid.ID { return xid.New() }),

		field.String("membership_type").
			Default("normal").
			Annotations(
				entsql.Default("normal"),
			),
	}
}

func (CollectionPost) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("collection", Collection.Type).
			Unique().
			Required().
			Field("collection_id"),

		edge.To("post", Post.Type).
			Unique().
			Required().
			Field("post_id"),
	}
}

func (CollectionPost) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("collection_id", "post_id").
			Unique().
			StorageKey("unique_collection_post"),
	}
}

func (CollectionPost) Annotations() []schema.Annotation {
	return []schema.Annotation{
		field.ID("collection_id", "post_id"),
	}
}
