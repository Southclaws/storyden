package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type Link struct {
	ent.Schema
}

func (Link) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (Link) Fields() []ent.Field {
	return []ent.Field{
		field.String("url").
			Unique().
			Immutable(),
		field.String("slug").
			Unique().
			Immutable(),
		field.String("domain"),
		field.String("title"),
		field.String("description"),
		field.Enum("liveness_status").Values("unknown", "alive", "offline").Default("unknown"),
		field.Time("liveness_checked_at").Optional().Nillable(),
		field.Int("liveness_failures").Default(0),
		field.Enum("scrape_status").Values("pending", "success", "failed").Default("pending"),
		field.Time("scrape_attempted_at").Optional().Nillable(),
		field.Time("scraped_at").Optional().Nillable(),

		field.String("primary_asset_id").
			GoType(xid.ID{}).
			Optional().
			Nillable(),

		field.String("favicon_asset_id").
			GoType(xid.ID{}).
			Optional().
			Nillable(),
	}
}

func (Link) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_at", "id"),
		index.Fields("liveness_checked_at", "id"),
	}
}

func (Link) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posts", Post.Type).
			Comment("Link aggregation posts that have shared this link."),

		edge.To("post_content_references", Post.Type).
			Comment("Posts that reference this link in their content."),

		edge.To("nodes", Node.Type),

		edge.To("node_content_references", Node.Type),

		edge.To("primary_image", Asset.Type).
			Field("primary_asset_id").
			Unique(),

		edge.To("favicon_image", Asset.Type).
			Field("favicon_asset_id").
			Unique(),

		edge.To("assets", Asset.Type),
	}
}
