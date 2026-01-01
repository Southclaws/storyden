package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

type AuditLog struct {
	ent.Schema
}

func (AuditLog) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (AuditLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("enacted_by_id").
			GoType(xid.ID{}).
			Optional().
			Nillable(),

		field.String("target_id").
			GoType(xid.ID{}).
			Optional().
			Nillable().
			Comment("The ID of the resource relevant to the log entry. This is not a foreign key as reports can refer to a variety of sources, discriminated by the 'target_kind' field."),

		field.String("target_kind").
			Optional().
			Nillable().
			Comment("The datagraph kind of related resource."),

		field.String("type"),

		field.String("error").
			Optional().
			Nillable(),

		field.JSON("metadata", map[string]any{}).
			Optional().
			Comment("Metadata specific to the type of audit log entry."),
	}
}

func (AuditLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("enacted_by", Account.Type).
			Field("enacted_by_id").
			Ref("audit_logs").
			Unique(),
	}
}
