package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type Cluster struct {
	ent.Schema
}

func (Cluster) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}, DeletedAt{}}
}

func (Cluster) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("slug").Unique(),
		field.String("image_url").Optional().Nillable(),
		field.String("description"),
		field.String("content").Optional().Nillable(),
		field.String("parent_cluster_id").GoType(xid.ID{}).Optional(),
		field.String("account_id").GoType(xid.ID{}),
		field.Any("properties").Optional(),
	}
}

func (Cluster) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("slug"),
	}
}

func (Cluster) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", Account.Type).
			Ref("clusters").
			Field("account_id").
			Unique().
			Required(),

		edge.To("clusters", Cluster.Type).
			From("parent").
			Unique().
			Field("parent_cluster_id").
			Comment("A many-to-many recursive self reference. The parent cluster, if any."),

		edge.To("items", Item.Type),

		edge.To("assets", Asset.Type),

		edge.From("tags", Tag.Type).
			Ref("clusters"),

		edge.From("links", Link.Type).
			Ref("clusters"),
	}
}
