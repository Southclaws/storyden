package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/Southclaws/lexorank"
	"github.com/rs/xid"
)

type Node struct {
	ent.Schema
}

func (Node) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}, DeletedAt{}, IndexedAt{}}
}

func (Node) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("slug").Unique(),
		field.String("description").Optional().Nillable(),
		field.String("content").Optional().Nillable(),
		field.String("parent_node_id").GoType(xid.ID{}).Optional(),
		field.String("account_id").GoType(xid.ID{}),
		field.String("property_schema_id").GoType(xid.ID{}).Optional().Nillable(),
		field.String("primary_asset_id").GoType(xid.ID{}).Optional().Nillable(),
		field.String("link_id").GoType(xid.ID{}).Optional(),
		field.Enum("visibility").Values(VisibilityTypes...).Default(VisibilityTypesDraft),
		field.String("sort").GoType(lexorank.Key{}).DefaultFunc(func() lexorank.Key {
			return lexorank.Top
		}),
		field.JSON("metadata", map[string]any{}).Optional(),
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

		edge.To("primary_image", Asset.Type).
			Field("primary_asset_id").
			Unique(),

		edge.To("assets", Asset.Type),

		edge.From("tags", Tag.Type).
			Ref("nodes"),

		edge.To("properties", Property.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("property_schema", PropertySchema.Type).
			Field("property_schema_id").
			Ref("node").
			Unique().
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.From("link", Link.Type).
			Field("link_id").
			Ref("nodes").
			Unique(),

		edge.From("content_links", Link.Type).
			Ref("node_content_references"),

		edge.From("collections", Collection.Type).
			Ref("nodes").
			Through("collection_nodes", CollectionNode.Type),
	}
}
