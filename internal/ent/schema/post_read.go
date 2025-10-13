package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type PostRead struct {
	ent.Schema
}

func (PostRead) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}}
}

func (PostRead) Fields() []ent.Field {
	return []ent.Field{
		field.String("root_post_id").GoType(xid.ID{}),
		field.String("account_id").GoType(xid.ID{}),
		field.Time("last_seen_at"),
	}
}

func (PostRead) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("root_post", Post.Type).
			Ref("post_reads").
			Field("root_post_id").
			Unique().
			Required(),

		edge.From("account", Account.Type).
			Ref("post_reads").
			Field("account_id").
			Unique().
			Required(),
	}
}

func (PostRead) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("root_post_id", "account_id").
			Unique().
			StorageKey("unique_post_read"),
	}
}
