package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type Account struct {
	ent.Schema
}

func (Account) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}, DeletedAt{}, IndexedAt{}}
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
		field.Enum("kind").
			Values("human", "bot").
			Default("human").
			Annotations(
				entsql.Default("human"),
			),
		field.Bool("admin").Default(false),
		field.JSON("links", []ExternalLink{}).Optional(),
		field.JSON("metadata", map[string]any{}).Optional(),

		field.String("invited_by_id").
			GoType(xid.ID{}).
			Optional().
			Nillable(),
	}
}

func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("sessions", Session.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("emails", Email.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("notifications", Notification.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("triggered_notifications", Notification.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("following", AccountFollow.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("followed_by", AccountFollow.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("invitations", Invitation.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.From("invited_by", Invitation.Type).
			Ref("invited").
			Field("invited_by_id").
			Unique(),

		edge.To("posts", Post.Type),
		edge.To("questions", Question.Type),

		edge.To("reacts", React.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("likes", LikePost.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("mentions", MentionProfile.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.From("roles", Role.Type).
			Ref("accounts").
			Through("account_roles", AccountRoles.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

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

		edge.To("events", EventParticipant.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("post_reads", PostRead.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("reports", Report.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("handled_reports", Report.Type).
			Annotations(entsql.OnDelete(entsql.SetNull)),

		edge.To("audit_logs", AuditLog.Type).
			Annotations(entsql.OnDelete(entsql.SetNull)),

		edge.To("robots", Robot.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("robot_sessions", RobotSession.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("robot_messages", RobotSessionMessage.Type).
			Annotations(entsql.OnDelete(entsql.SetNull)),
	}
}
