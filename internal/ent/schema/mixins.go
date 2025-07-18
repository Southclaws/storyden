package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/rs/xid"
)

type Identifier struct{ mixin.Schema }

func (Identifier) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			MaxLen(20).
			NotEmpty().
			Immutable().
			GoType(xid.ID{}).
			DefaultFunc(func() xid.ID { return xid.New() }),
	}
}

type CreatedAt struct{ mixin.Schema }

func (CreatedAt) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			Annotations(
				entsql.Default("CURRENT_TIMESTAMP"),
			).
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

type UpdatedAtManual struct{ mixin.Schema }

func (UpdatedAtManual) Fields() []ent.Field {
	return []ent.Field{
		field.Time("updated_at"),
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

type ExpiresAt struct{ mixin.Schema }

func (ExpiresAt) Fields() []ent.Field {
	return []ent.Field{
		field.Time("expires_at").
			Optional().
			Nillable(),
	}
}

type IndexedAt struct{ mixin.Schema }

func (IndexedAt) Fields() []ent.Field {
	return []ent.Field{
		field.Time("indexed_at").
			Optional().
			Nillable(),
	}
}
