package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Setting struct {
	ent.Schema
}

func (Setting) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),

		field.String("value"),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}
