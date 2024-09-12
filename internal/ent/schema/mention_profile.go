package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type MentionProfile struct {
	ent.Schema
}

func (MentionProfile) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (MentionProfile) Fields() []ent.Field {
	return []ent.Field{
		field.String("account_id").GoType(xid.ID{}),
		field.String("post_id").GoType(xid.ID{}),
	}
}

func (MentionProfile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Ref("mentions").
			Field("account_id").
			Unique().
			Required(),

		edge.From("Post", Post.Type).
			Ref("mentions").
			Field("post_id").
			Unique().
			Required(),
	}
}

func (MentionProfile) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("account_id", "post_id").
			Unique().
			StorageKey("unique_mentions_post"),
	}
}
