package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type PropertySnapshotEntry struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name"`
	Type  string `json:"type,omitempty"`
	Value string `json:"value"`
	Sort  string `json:"sort,omitempty"`
}

type PropertySnapshot struct {
	Set        bool                    `json:"set"`
	Properties []PropertySnapshotEntry `json:"properties"`
}

type NodeVersion struct {
	ent.Schema
}

func (NodeVersion) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (NodeVersion) Fields() []ent.Field {
	return []ent.Field{
		field.String("node_id").GoType(xid.ID{}),
		field.String("author_id").GoType(xid.ID{}),

		field.Enum("status").Values(VersionStatusValues...).Default(VersionStatusDraft),

		field.String("name"),
		field.String("slug"),
		field.String("description").Optional().Nillable(),
		field.String("content").Optional().Nillable(),
		field.JSON("properties_snapshot", PropertySnapshot{}).Optional(),

		field.JSON("metadata", map[string]any{}).Optional(),
	}
}

func (NodeVersion) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("node_id", "updated_at"),
		index.Fields("node_id", "status", "updated_at"),
	}
}

func (NodeVersion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("node", Node.Type).
			Field("node_id").
			Ref("versions").
			Required().
			Unique().
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.From("author", Account.Type).
			Field("author_id").
			Ref("node_versions").
			Required().
			Unique(),

		edge.To("current_for_nodes", Node.Type),
	}
}
