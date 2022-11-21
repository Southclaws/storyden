package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type React struct {
	ent.Schema
}

func (React) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (React) Fields() []ent.Field {
	return []ent.Field{
		field.String("emoji"),
	}
}

func (React) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("account", Account.Type).Unique(),
		edge.To("Post", Post.Type).Unique(),
	}
}
