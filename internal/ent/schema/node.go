package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type Node struct {
	ent.Schema
}

func (Node) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}, DeletedAt{}}
}

func (Node) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("slug").Unique(),
		field.String("description").Optional().Nillable(),
		field.String("content").Optional().Nillable(),
		field.String("parent_node_id").GoType(xid.ID{}).Optional(),
		field.String("account_id").GoType(xid.ID{}),
		field.Enum("visibility").Values(VisibilityTypes...).Default(VisibilityTypesDraft),
		field.Any("properties").Optional(),
	}
}

func (Node) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("slug"),
	}
}

func (Node) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", Account.Type).
			Ref("nodes").
			Field("account_id").
			Unique().
			Required(),

		edge.To("nodes", Node.Type).
			From("parent").
			Unique().
			Field("parent_node_id").
			Comment("A many-to-many recursive self reference. The parent node, if any."),

		edge.To("assets", Asset.Type),

		edge.From("tags", Tag.Type).
			Ref("nodes"),

		edge.From("links", Link.Type).
			Ref("nodes"),

		edge.From("collections", Collection.Type).
			Ref("nodes"),
	}
}
