package authentication

import (
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

type ID = xid.ID

type Service string

type Authentication struct {
	ID         ID
	Created    time.Time
	Account    account.Account
	Service    Service
	Identifier string
	Token      string
	Name       opt.Optional[string]
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
		ID:         ID(m.ID),
		Created:    m.CreatedAt,
		Account:    *acc,
		Service:    Service(m.Service),
		Identifier: m.Identifier,
		Token:      m.Token,
		Name:       opt.NewPtr(m.Name),
		Metadata:   m.Metadata,
	}, nil
}
