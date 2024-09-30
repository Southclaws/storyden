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

type CollectionNode struct {
	ent.Schema
}

func (CollectionNode) Mixin() []ent.Mixin {
	return []ent.Mixin{CreatedAt{}}
}

func (CollectionNode) Fields() []ent.Field {
	return []ent.Field{
		field.String("collection_id").
			MaxLen(20).
			NotEmpty().
			GoType(xid.ID{}).
			DefaultFunc(func() xid.ID { return xid.New() }),

		field.String("node_id").
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

func (CollectionNode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("collection", Collection.Type).
			Unique().
			Required().
			Field("collection_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("node", Node.Type).
			Unique().
			Required().
			Field("node_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (CollectionNode) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("collection_id", "node_id").
			Unique().
			StorageKey("unique_collection_node"),
	}
}

func (CollectionNode) Annotations() []schema.Annotation {
	return []schema.Annotation{
		field.ID("collection_id", "node_id"),
	}
}
