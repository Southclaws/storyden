package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// React holds the schema definition for the React entity.
type React struct {
    ent.Schema
}

// Fields of React.
func (React) Fields() []ent.Field {
    return []ent.Field{
        field.String("id"),
        field.String("emoji"),
        field.Time("createdAt"),
        field.String("postId"),
        field.String("userId"),
    }
}

// Edges of React.
func (React) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("user", User.Type),
        edge.To("Post", Post.Type),
    }
}
