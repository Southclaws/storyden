package schema

import (
	"encoding/json"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type BackgroundJob struct {
	ent.Schema
}

func (BackgroundJob) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (BackgroundJob) Fields() []ent.Field {
	return []ent.Field{
		field.String("kind"),
		field.String("key").NotEmpty(),
		field.JSON("payload", json.RawMessage{}),
		field.Int64("revision").Default(0),
		field.Enum("status").Values("pending", "processing", "skipped", "failed", "complete").Default("pending"),
		field.Int("attempts").Default(0),
		field.String("reason").Default(""),
		field.String("lease_token").Default(""),
		field.Time("lease_until").Default(time.Now),
		field.Time("dispatch_until").Default(time.Now),
		field.Time("available_at").Default(time.Now),
	}
}

func (BackgroundJob) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("kind", "status", "available_at", "dispatch_until"),
		index.Fields("kind", "key").Unique(),
	}
}
