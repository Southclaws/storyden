package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Rule holds the schema definition for the Rule entity.
type Rule struct {
    ent.Schema
}

// Fields of Rule.
func (Rule) Fields() []ent.Field {
    return []ent.Field{
        field.Int("id"),
        field.String("name"),
        field.String("value"),
        field.String("serverId").Optional(),
    }
}

// Edges of Rule.
func (Rule) Edges() []ent.Edge {
    return []ent.Edge{
    edge.To("Server", Server.Type),
    }
}
