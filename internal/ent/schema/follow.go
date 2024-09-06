package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type AccountFollow struct {
	ent.Schema
}

func (AccountFollow) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (AccountFollow) Fields() []ent.Field {
	return []ent.Field{
		field.String("follower_account_id").GoType(xid.ID{}),
		field.String("following_account_id").GoType(xid.ID{}),
	}
}

func (AccountFollow) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("follower", Account.Type).
			Ref("following").
			Field("follower_account_id").
			Unique().
			Required(),

		edge.From("following", Account.Type).
			Ref("followed_by").
			Field("following_account_id").
			Unique().
			Required(),
	}
}

func (AccountFollow) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("follower_account_id", "following_account_id").
			Unique().
			StorageKey("unique_following_pair"),
	}
}
