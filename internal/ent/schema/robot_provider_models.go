package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type RobotProviderModel struct {
	ent.Schema
}

func (RobotProviderModel) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (RobotProviderModel) Fields() []ent.Field {
	return []ent.Field{
		field.String("provider").
			NotEmpty(),

		field.String("name").
			NotEmpty(),

		field.JSON("raw", map[string]any{}).
			Optional(),

		field.Time("last_seen_at").
			Default(time.Now),
	}
}

func (RobotProviderModel) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("provider", "name").Unique(),
	}
}
