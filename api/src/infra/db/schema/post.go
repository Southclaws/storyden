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
		edge.From("author", User.Type).
			Ref("posts").
			Unique().
			Required(),

		edge.From("category", Category.Type).
			Ref("posts").
			Comment("Category is only required for root posts. It should never be added to a child post."),

		edge.From("tags", Tag.Type).
			Ref("posts").
			Comment("Tagss are only required for root posts. It should never be added to a child post."),

		edge.To("posts", Post.Type).
			From("root").
			Comment("A many-to-many recursive self reference. The root post is the first post in the thread."),

		edge.To("replyTo", Post.Type).
			From("replies").
			Comment("A many-to-many recursive self reference. The replyTo post is an optional post that this post is in reply to."),

		edge.To("reacts", React.Type),
	}
}
