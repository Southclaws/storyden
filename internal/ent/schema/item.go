package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type Item struct {
	ent.Schema
}

func (Item) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}, DeletedAt{}}
}

func (Item) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("slug"),
		field.String("image_url").Optional().Nillable(),
		field.String("description"),
		field.String("account_id").GoType(xid.ID{}),
		field.Any("properties").Optional(),
	}
}

func (Item) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", Account.Type).
			Ref("items").
			Field("account_id").
			Unique().
			Required(),

		edge.From("clusters", Cluster.Type).
			Ref("items"),

		edge.To("assets", Asset.Type),

		edge.From("tags", Tag.Type).
			Ref("items"),
	}
}
