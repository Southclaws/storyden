package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Tag struct {
	ent.Schema
}

func (Tag) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (Tag) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique().
			Immutable(),
	}
}

func (Tag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posts", Post.Type).Through("post_tags", TagPost.Type),
		edge.To("nodes", Node.Type).Through("node_tags", TagNode.Type),
		edge.From("accounts", Account.Type).
			Ref("tags"),
	}
}
