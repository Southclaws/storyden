package authentication

import (
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/backend/pkg/resources/user"
	"github.com/samber/lo"
)

type Service string

type Authentication struct {
	User       user.User   `json:"user"`
	Service    Service     `json:"service"`
	Identifier string      `json:"identifier"`
	Token      string      `json:"-"`
	Metadata   interface{} `json:"metadata,omitempty"`
}

func FromModel(m *model.Authentication) *Authentication {
	return &Authentication{
		User:       user.FromModel(*m.Edges.User),
		Service:    Service(m.Service),
		Identifier: m.Identifier,
		Token:      m.Token,
		Metadata:   m.Metadata,
	}
}

func FromModelMany(m []*model.Authentication) []Authentication {
	return lo.Map(m, func(t *model.Authentication, i int) Authentication {
		return *FromModel(t)
	})
}
