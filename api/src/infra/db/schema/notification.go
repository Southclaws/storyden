package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// Notification holds the schema definition for the Notification entity.
type Notification struct {
	ent.Schema
}

// Fields of Notification.
func (Notification) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Default(uuid.New),

		field.String("title"),
		field.String("description"),
		field.String("link"),
		field.Bool("read"),

		// BUG: ent does not do this automatically.
		field.Time("create_time").Default(time.Now),
	}
}

// Edges of Notification.
func (Notification) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("subscription", Subscription.Type),
	}
}

func (Notification) Mixins() []ent.Mixin {
	return []ent.Mixin{
		mixin.CreateTime{},
	}
}
