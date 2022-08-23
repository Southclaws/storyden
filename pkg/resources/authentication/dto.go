package authentication

import (
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/pkg/resources/account"
)

type Service string

type Authentication struct {
	Account    account.Account `json:"account"`
	Service    Service         `json:"service"`
	Identifier string          `json:"identifier"`
	Token      string          `json:"-"`
	Metadata   interface{}     `json:"metadata,omitempty"`
}

func FromModel(m *model.Authentication) *Authentication {
	return &Authentication{
		Account:    account.FromModel(*m.Edges.Account),
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
