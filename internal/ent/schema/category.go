package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

// Category holds the schema definition for the Category entity.
type Category struct {
	ent.Schema
}

func (Category) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

// Fields of Category.
func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.String("slug").Unique(),
		field.String("description").Default("(No description)"),
		field.String("colour").Default("#8577ce"),
		field.Int("sort").Default(-1),
		field.Bool("admin").Default(false),
		field.String("parent_category_id").
			GoType(xid.ID{}).
			Optional(),
		field.String("cover_image_asset_id").
			GoType(xid.ID{}).
			Optional().
			Nillable(),
		field.JSON("metadata", map[string]any{}).
			Optional().
			Comment("Arbitrary metadata used by clients to store domain specific information."),
	}
}

// Edges of Category.
func (Category) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posts", Post.Type),
		edge.To("children", Category.Type).
			From("parent").
			Unique().
			Field("parent_category_id").
			Comment("Optional recursive self reference to the parent category."),
		edge.To("cover_image", Asset.Type).
			Field("cover_image_asset_id").
			Unique(),
	}
}
