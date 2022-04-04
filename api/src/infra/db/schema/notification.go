package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Notification holds the schema definition for the Notification entity.
type Notification struct {
	ent.Schema
}

// Fields of Notification.
func (Notification) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.String("title"),
		field.String("description"),
		field.String("link"),
		field.Bool("read"),
		field.Time("createdAt"),
		field.String("subscriptionId"),
	}
}

// Edges of Notification.
func (Notification) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("subscription", Subscription.Type),
	}
}
