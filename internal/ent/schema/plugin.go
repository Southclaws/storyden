package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/Southclaws/storyden/lib/plugin"
	"github.com/rs/xid"
)

type Plugin struct {
	ent.Schema
}

func (Plugin) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (Plugin) Fields() []ent.Field {
	return []ent.Field{
		field.String("path").
			Unique(),
		field.JSON("manifest", plugin.Manifest{}),
		field.JSON("config", map[string]any{}),

		field.String("active_state"),
		field.Time("active_state_changed_at"),
		field.String("status_message").
			Optional().
			Nillable(),
		field.JSON("status_details", map[string]any{}).
			Optional(),

		field.String("added_by").
			GoType(xid.ID{}),
	}
}

func (Plugin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Required().
			Ref("plugins").
			Field("added_by").
			Unique(),
	}
}
