package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type RobotSession struct {
	ent.Schema
}

func (RobotSession) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (RobotSession) Fields() []ent.Field {
	return []ent.Field{
		field.String("account_id").
			GoType(xid.ID{}).
			Comment("UserID (account ID) from ADK Session"),

		field.JSON("state", map[string]any{}).
			Optional().
			Comment("Session state from ADK"),
	}
}

func (RobotSession) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", Account.Type).
			Field("account_id").
			Ref("robot_sessions").
			Unique().
			Required(),

		edge.To("messages", RobotSessionMessage.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
