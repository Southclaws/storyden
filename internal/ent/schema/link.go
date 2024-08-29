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
		edge.To("posts", Post.Type).
			Comment("Link aggregation posts that have shared this link."),

		edge.To("post_content_references", Post.Type).
			Comment("Posts that reference this link in their content."),

		edge.To("nodes", Node.Type),

		edge.To("node_content_references", Node.Type),

		edge.To("assets", Asset.Type),
	}
}
