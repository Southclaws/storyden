package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type PushSubscription struct {
	ent.Schema
}

func (PushSubscription) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (PushSubscription) Fields() []ent.Field {
	return []ent.Field{
		field.String("account_id").GoType(xid.ID{}).Immutable(),
		field.String("endpoint").NotEmpty().Unique(),
		field.String("p256dh").NotEmpty(),
		field.String("auth").NotEmpty().Sensitive(),
	}
}

func (PushSubscription) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Ref("push_subscriptions").
			Field("account_id").
			Required().
			Unique().
			Immutable(),
	}
}
