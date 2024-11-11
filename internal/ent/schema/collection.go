package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type Collection struct {
	ent.Schema
}

func (Collection) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}, IndexedAt{}}
}

func (Collection) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("slug"),
		field.String("description").Optional().Nillable(),
		field.String("cover_asset_id").GoType(xid.ID{}).Optional().Nillable(),
		field.Enum("visibility").Values(VisibilityTypes...).Default(VisibilityTypesDraft),
	}
}

func (Collection) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", Account.Type).
			Ref("collections").
			Unique(),

		edge.To("cover_image", Asset.Type).
			Field("cover_asset_id").
			Unique(),

		edge.To("posts", Post.Type).
			Through("collection_posts", CollectionPost.Type),
		edge.To("nodes", Node.Type).
			Through("collection_nodes", CollectionNode.Type),
	}
}
