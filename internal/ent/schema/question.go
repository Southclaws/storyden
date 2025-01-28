package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type Question struct {
	ent.Schema
}

func (Question) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, IndexedAt{}}
}

func (Question) Fields() []ent.Field {
	return []ent.Field{
		field.String("slug").Unique(),
		field.String("query"),
		field.String("result"),

		field.JSON("metadata", map[string]any{}).Optional(),

		field.String("account_id").GoType(xid.ID{}).Optional(),
		field.String("parent_question_id").GoType(xid.ID{}).Optional(),
	}
}

func (Question) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("author", Account.Type).
			Field("account_id").
			Ref("questions").
			Unique(),

		edge.To("parent_question", Question.Type).
			From("parent").
			Unique().
			Field("parent_question_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
