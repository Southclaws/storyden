package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/rs/xid"
)

type Property struct {
	ent.Schema
}

func (Property) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}}
}

func (Property) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("type"),
		field.String("value"),
		field.String("node_id").GoType(xid.ID{}).Optional(),
	}
}

func (Property) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("node", Node.Type).
			Field("node_id").
			Ref("properties").
			Unique(),
	}
}

func (Property) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("node_id", "name").Unique(),
		index.Fields("name"),
		index.Fields("node_id"),
	}
}
