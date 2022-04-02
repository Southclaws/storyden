package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Tag holds the schema definition for the Tag entity.
type Tag struct {
    ent.Schema
}

// Fields of Tag.
func (Tag) Fields() []ent.Field {
    return []ent.Field{
        field.String("id"),
        field.String("name"),
    }
}

// Edges of Tag.
func (Tag) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("posts", Post.Type),
    }
}
