package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// Subscription holds the schema definition for the Subscription entity.
type Subscription struct {
	ent.Schema
}

// Fields of Subscription.
func (Subscription) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Default(uuid.New),

		field.String("refers_type").NotEmpty(),
		field.String("refers_to").NotEmpty(),

		field.Time("delete_time").Optional(),

		// BUG: ent does not do this automatically.
		field.Time("create_time").Default(time.Now),
		field.Time("update_time").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of Subscription.
func (Subscription) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type),
		edge.To("notifications", Notification.Type),
	}
}

func (Subscription) Mixins() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
