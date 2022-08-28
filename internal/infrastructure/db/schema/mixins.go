package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/rs/xid"
)

type Identifier struct{ mixin.Schema }

func (Identifier) Fields() []ent.Field {
	return []ent.Field{
		field.Bytes("id").
			MaxLen(20).
			NotEmpty().
			Immutable().
			GoType(xid.ID{}).
			SchemaType(map[string]string{
				dialect.MySQL:    "binary(12)",
				dialect.Postgres: "bytea",
			}).
			DefaultFunc(func() xid.ID { return xid.New() }),
	}
}

type CreatedAt struct{ mixin.Schema }

func (CreatedAt) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

type UpdatedAt struct{ mixin.Schema }

func (UpdatedAt) Fields() []ent.Field {
	return []ent.Field{
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

type DeletedAt struct{ mixin.Schema }

func (DeletedAt) Fields() []ent.Field {
	return []ent.Field{
		field.Time("deleted_at").
			Optional().
			Nillable(),
	}
}
