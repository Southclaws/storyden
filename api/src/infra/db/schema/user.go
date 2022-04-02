package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
    ent.Schema
}

// Fields of User.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.UUID("id", uuid.UUID{}).
            Unique().
            Immutable().
            Default(uuid.New),

        field.String("email"),
        field.String("name"),
        field.String("bio").Optional(),
        field.Bool("admin").Default(false),

        field.Time("createdAt").Default(time.Now()),
        field.Time("updatedAt").Default(time.Now()),
        field.Time("deletedAt").Optional(),
    }
}

// Edges of User.
func (User) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("posts", Post.Type),
        edge.To("reacts", React.Type),
        edge.To("subscriptions", Subscription.Type),
    }
}
