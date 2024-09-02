package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type LikePost struct {
	ent.Schema
}

func (LikePost) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (LikePost) Fields() []ent.Field {
	return []ent.Field{
		field.String("account_id").GoType(xid.ID{}),
		field.String("post_id").GoType(xid.ID{}),
	}
}

func (LikePost) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Ref("likes").
			Field("account_id").
			Unique().
			Required(),

		edge.From("Post", Post.Type).
			Ref("likes").
			Field("post_id").
			Unique().
			Required(),
	}
}

func (LikePost) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("account_id", "post_id").
			Unique().
			StorageKey("unique_like_post"),
	}
}
