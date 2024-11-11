package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Event struct {
	ent.Schema
}

func (Event) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}, DeletedAt{}, IndexedAt{}}
}

func (Event) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),

		field.String("slug").
			Unique(),

		field.String("description").
			Optional().
			Nillable(),

		field.Time("start_time"),

		field.Time("end_time"),

		field.String("participation_policy"),

		field.Enum("visibility").
			Values(VisibilityTypes...).
			Default(VisibilityTypesDraft),

		field.String("location_type").
			Optional().
			Nillable(),

		field.String("location_name").
			Optional().
			Nillable(),

		field.String("location_address").
			Optional().
			Nillable(),

		field.Float("location_latitude").
			Optional().
			Nillable(),

		field.Float("location_longitude").
			Optional().
			Nillable(),

		field.String("location_url").
			Optional().
			Nillable(),

		field.Int("capacity").
			Optional().
			Nillable(),

		field.JSON("metadata", map[string]any{}).
			Optional(),
	}
}

func (Event) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("slug").Unique(),
	}
}

func (Event) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("participants", EventParticipant.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.From("thread", Post.Type).
			Ref("event").
			Unique().
			Required(),

		edge.From("primary_image", Asset.Type).
			Ref("event").
			Unique(),
	}
}
