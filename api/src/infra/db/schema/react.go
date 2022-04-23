package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// React holds the schema definition for the React entity.
type React struct {
	ent.Schema
}

// Fields of React.
func (React) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Default(uuid.New),

		field.String("emoji"),

		field.Time("createdAt").Default(time.Now),
	}
}

// Edges of React.
func (React) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).Unique(),
		edge.To("Post", Post.Type).Unique(),
	}
}
