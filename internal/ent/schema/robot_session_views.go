package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

// RobotSessionView is one account's view of a shared Robot session. It is not
// ownership or an execution lease. Future durable stream cursors belong here.
type RobotSessionView struct {
	ent.Schema
}

func (RobotSessionView) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (RobotSessionView) Fields() []ent.Field {
	return []ent.Field{
		field.String("session_id").GoType(xid.ID{}),
		field.String("account_id").GoType(xid.ID{}),
		field.Time("last_accessed_at").Default(time.Now),
	}
}

func (RobotSessionView) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("session_id", "account_id").Unique(),
		index.Fields("account_id", "last_accessed_at"),
	}
}

func (RobotSessionView) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("session", RobotSession.Type).
			Field("session_id").
			Ref("views").
			Unique().
			Required(),

		edge.From("account", Account.Type).
			Field("account_id").
			Ref("robot_session_views").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
