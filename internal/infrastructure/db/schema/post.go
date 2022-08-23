package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Post holds the schema definition for the Post entity.
type Post struct {
	ent.Schema
}

// Fields of Post.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Default(uuid.New),

		field.Bool("first"),

		// parent posts
		field.String("title").Optional(),
		field.String("slug").Optional(),
		field.Bool("pinned").Default(false),

		// child posts
		field.UUID("root_post_id", uuid.UUID{}).Optional(),
		field.UUID("reply_to_post_id", uuid.UUID{}).Optional(),

		// All posts
		field.String("body"),
		field.String("short"),

		field.Time("createdAt").Default(time.Now),
		field.Time("updatedAt").Default(time.Now),
		field.Time("deletedAt").Optional().Nillable(),

		// Edges
		field.UUID("category_id", uuid.UUID{}).Optional(),
	}
}

// Edges of Post.
func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("author", Account.Type).
			Ref("posts").
			Unique().
			Required(),

		edge.From("category", Category.Type).
			Field("category_id").
			Ref("posts").
			Unique().
			Comment("Category is only required for root posts. It should never be added to a child post."),

		edge.From("tags", Tag.Type).
			Ref("posts").
			Comment("Tags are only required for root posts. It should never be added to a child post."),

		edge.To("posts", Post.Type).
			From("root").
			Unique().
			Field("root_post_id").
			Comment("A many-to-many recursive self reference. The root post is the first post in the thread."),

		edge.To("replies", Post.Type).
			From("replyTo").
			Unique().
			Field("reply_to_post_id").
			Comment("A many-to-many recursive self reference. The replyTo post is an optional post that this post is in reply to."),

		edge.To("reacts", React.Type),
	}
}
