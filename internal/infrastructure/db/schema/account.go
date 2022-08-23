package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Account struct {
	ent.Schema
}

func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Default(uuid.New),

		field.String("email").Unique(),
		field.String("name").NotEmpty(),
		field.String("bio").Optional(),
		field.Bool("admin").Default(false),

		field.Time("createdAt").Default(time.Now),
		field.Time("updatedAt").Default(time.Now),
		field.Time("deletedAt").Optional(),
	}
}

func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posts", Post.Type),

		edge.To("reacts", React.Type),

		edge.To("subscriptions", Subscription.Type),

		edge.To("authentication", Authentication.Type),
	}
}
