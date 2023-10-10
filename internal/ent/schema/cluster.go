package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type Cluster struct {
	ent.Schema
}

func (Cluster) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (Cluster) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("slug"),
		field.String("image_url").Optional().Nillable(),
		field.String("description"),
		field.String("parent_cluster_id").GoType(xid.ID{}).Optional(),
		field.String("account_id").GoType(xid.ID{}),
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
	}
}
