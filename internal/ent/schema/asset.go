package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type Asset struct {
	ent.Schema
}

func (Asset) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (Asset) Fields() []ent.Field {
	return []ent.Field{
		field.String("filename"),

		field.Int("size").
			Annotations(entsql.Default("0")),

		field.String("mime_type").
			Annotations(entsql.Default("application/octet-stream")),

		field.JSON("metadata", map[string]any{}).Optional(),

		// Edges
		field.String("account_id").GoType(xid.ID{}),
		field.String("parent_asset_id").GoType(xid.ID{}).Optional().Nillable(),
	}
}

func (Asset) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("filename"),
	}
}

func (Asset) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("posts", Post.Type).
			Ref("assets"),

		edge.From("nodes", Node.Type).
			Ref("assets"),

		edge.From("links", Link.Type).
			Ref("assets"),

		edge.From("owner", Account.Type).
			Field("account_id").
			Ref("assets").
			Unique().
			Required(),

		edge.To("assets", Asset.Type).
			From("parent").
			Unique().
			Field("parent_asset_id"),

		edge.To("event", Event.Type),
	}
}
