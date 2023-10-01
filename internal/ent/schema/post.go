package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type Post struct {
	ent.Schema
}

func (Post) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}, DeletedAt{}}
}

func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("first"),

		// parent posts
		field.String("title").Optional(),
		field.String("slug").Optional(),
		field.Bool("pinned").Default(false),

		// child posts
		field.String("root_post_id").GoType(xid.ID{}).Optional(),
		field.String("reply_to_post_id").GoType(xid.ID{}).Optional(),

		// All posts
		field.String("body"),
		field.String("short"),
		field.JSON("metadata", map[string]any{}).
			Optional().
			Comment("Arbitrary metadata used by clients to store domain specific information."),
		field.Enum("status").Values(
			"draft",
			"published",
		).Default("draft"),
		field.String("url").Optional().Nillable(),

		// Edges
		field.String("category_id").GoType(xid.ID{}).Optional(),
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

		edge.To("assets", Asset.Type),

		edge.From("collections", Collection.Type).
			Ref("posts"),
	}
}
