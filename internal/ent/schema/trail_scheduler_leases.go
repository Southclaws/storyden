package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type TrailSchedulerLease struct {
	ent.Schema
}

func (TrailSchedulerLease) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Default(1).Immutable(),
		field.String("lease_token").Optional().Nillable(),
		field.Time("lease_expires_at").Optional().Nillable(),
	}
}
