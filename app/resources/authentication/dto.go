package authentication

import (
	"github.com/Southclaws/fault"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

type Service string

type Authentication struct {
	Account    account.Account
	Service    Service
	Identifier string
	Token      string
	Metadata   interface{}
}

func FromModel(m *ent.Authentication) (*Authentication, error) {
	accEdge, err := m.Edges.AccountOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	acc, err := account.FromModel(accEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Authentication{
		Account:    *acc,
		Service:    Service(m.Service),
		Identifier: m.Identifier,
		Token:      m.Token,
		Metadata:   m.Metadata,
	}, nil
}
