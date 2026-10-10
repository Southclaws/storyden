package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type IdempotencyReceipt struct{ ent.Schema }

func (IdempotencyReceipt) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").MaxLen(64).Immutable(),
		field.String("principal_id").MaxLen(20).Immutable(),
		field.String("operation").MaxLen(64).Immutable(),
		field.String("key_hash").MaxLen(64).Immutable(),
		field.String("claim_token").MaxLen(20).Immutable(),
		field.String("fingerprint").MaxLen(64).Immutable(),
		field.String("response").Optional(),
		field.Time("failed_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("expires_at").Immutable(),
	}
}

func (IdempotencyReceipt) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("principal_id", "operation", "key_hash").Unique(),
		index.Fields("expires_at"),
	}
}
