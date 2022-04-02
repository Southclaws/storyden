package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Post holds the schema definition for the Post entity.
type Post struct {
    ent.Schema
}

// Fields of Post.
func (Post) Fields() []ent.Field {
    return []ent.Field{
        field.String("id"),
        field.String("title").Optional(),
        field.String("slug").Optional(),
        field.String("body"),
        field.String("short"),
        field.Bool("first"),
        field.Bool("pinned").Default(false),
        field.Time("createdAt"),
        field.Time("updatedAt"),
        field.Time("deletedAt").Optional(),
        field.String("userId"),
        field.String("rootPostId").Optional(),
        field.String("replyPostId").Optional(),
        field.String("categoryId").Optional(),
    }
}

// Edges of Post.
func (Post) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("category", Category.Type),
        edge.To("author", User.Type),

        // edge.To("root", Post.Type),
        edge.To("posts", Post.Type).From("root").Unique(),

        edge.To("replyTo", Post.Type),
        edge.To("replies", Post.Type),
        edge.To("tags", Tag.Type),
        edge.To("reacts", React.Type),
    }
}
