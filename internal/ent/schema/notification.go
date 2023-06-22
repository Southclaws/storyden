package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Notification struct {
	ent.Schema
}

func (Notification) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (Notification) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.String("description"),
		field.String("link"),
		field.Bool("read"),
	}
}

// func (Notification) Edges() []ent.Edge {
// 	return []ent.Edge{
// 		edge.To("subscription", Subscription.Type).
// 			Annotations(entsql.Annotation{
// 				OnDelete: entsql.Cascade,
// 			}).
// 			Unique(),
// 	}
// }
