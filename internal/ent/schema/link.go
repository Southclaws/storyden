package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Link struct {
	ent.Schema
}

func (Link) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (Link) Fields() []ent.Field {
	return []ent.Field{
		field.String("url").
			Unique().
			Immutable(),
		field.String("slug").
			Unique().
			Immutable(),
		field.String("domain"),
		field.String("title"),
		field.String("description"),
	}
}

func (Link) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posts", Post.Type),
		edge.To("clusters", Cluster.Type),
		edge.To("items", Item.Type),
		edge.To("assets", Asset.Type),
	}
}
