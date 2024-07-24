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

type ExternalLink struct {
	Text string
	URL  string
}

func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.String("handle").Unique().NotEmpty(),
		field.String("name").NotEmpty(),
		field.String("bio").Optional(),
		field.Bool("admin").Default(false),
		field.JSON("links", []ExternalLink{}).Optional(),
		field.JSON("metadata", map[string]any{}).Optional(),
	}
}

func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("emails", Email.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

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

		edge.To("nodes", Node.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)), // TODO: Don't cascade but do something more clever

		edge.To("assets", Asset.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)), // TODO: Don't cascade but do something more clever
	}
}
