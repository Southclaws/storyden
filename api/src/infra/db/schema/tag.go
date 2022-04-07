package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Tag holds the schema definition for the Tag entity.
type Tag struct {
	ent.Schema
}

// Fields of Tag.
func (Tag) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Unique().
			Immutable().
			Default(uuid.New),

		field.String("name").
			Unique().
			Immutable(),
	}
}

// Edges of Tag.
func (Tag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posts", Post.Type),
	}
}
