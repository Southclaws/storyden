package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Server holds the schema definition for the Server entity.
type Server struct {
    ent.Schema
}

// Fields of Server.
func (Server) Fields() []ent.Field {
    return []ent.Field{
        field.String("id"),
        field.String("ip"),
        field.String("hn"),
        field.Int("pc"),
        field.Int("pm"),
        field.String("gm"),
        field.String("la"),
        field.Bool("pa"),
        field.String("vn"),
        field.String("domain").Optional(),
        field.String("description").Optional(),
        field.String("banner").Optional(),
        field.String("userId").Optional(),
        field.Bool("active"),
        field.Time("updatedAt"),
        field.Time("deletedAt").Optional(),
    }
}

// Edges of Server.
func (Server) Edges() []ent.Edge {
    return []ent.Edge{
    edge.To("ru", Rule.Type),
    edge.To("User", User.Type),
    }
}
