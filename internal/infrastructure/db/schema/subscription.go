package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Subscription struct {
	ent.Schema
}

func (Subscription) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (Subscription) Fields() []ent.Field {
	return []ent.Field{
		field.String("refers_type").NotEmpty(),
		field.String("refers_to").NotEmpty(),
	}
}

func (Subscription) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("account", Account.Type).Unique(),
		edge.To("notifications", Notification.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
