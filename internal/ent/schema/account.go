package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Account struct {
	ent.Schema
}

func (Account) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}, DeletedAt{}}
}

func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.String("handle").Unique().NotEmpty(),
		field.String("name").NotEmpty(),
		field.String("bio").Optional(),
		field.Bool("admin").Default(false),
	}
}

func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posts", Post.Type),

		edge.To("reacts", React.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.From("roles", Role.Type).
			Ref("accounts"),

		edge.To("authentication", Authentication.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("tags", Tag.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("collections", Collection.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("clusters", Cluster.Type).
			Annotations(entsql.OnDelete(entsql.SetNull)),

		edge.To("items", Item.Type).
			Annotations(entsql.OnDelete(entsql.SetNull)),

		edge.To("assets", Asset.Type).
			Annotations(entsql.OnDelete(entsql.SetNull)),
	}
}
