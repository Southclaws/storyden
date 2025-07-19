package authentication

import (
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

var ErrExpired = fault.New("authentication expired", ftag.With(ftag.PermissionDenied))

type ID = xid.ID

type Authentication struct {
	ID         ID
	Created    time.Time
	Expires    opt.Optional[time.Time]
	Account    account.Account
	Service    Service
	Type       TokenType
	Identifier string
	Token      string
	Name       opt.Optional[string]
	Disabled   bool
	Metadata   interface{}
}

func (a Authentication) IsExpired() bool {
	exp, ok := a.Expires.Get()
	if !ok {
		return false
	}
	return exp.Before(time.Now())
}

func (a Authentication) CheckExpired() error {
	exp, ok := a.Expires.Get()
	if !ok {
		return nil
	}
	if exp.After(time.Now()) {
		return nil
	}
	return fault.Wrap(ErrExpired)
}

func FromModel(m *ent.Authentication) (*Authentication, error) {
	accEdge, err := m.Edges.AccountOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	acc, err := account.MapRef(accEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	tokenType, err := NewTokenType(m.TokenType)
	if err != nil {
		return nil, err
	}

	service, err := NewService(m.Service)
	if err != nil {
		return nil, err
	}

	return &Authentication{
		ID:         ID(m.ID),
		Created:    m.CreatedAt,
		Expires:    opt.NewPtr(m.ExpiresAt),
		Account:    *acc,
		Service:    service,
		Type:       tokenType,
		Identifier: m.Identifier,
		Token:      m.Token,
		Name:       opt.NewPtr(m.Name),
		Disabled:   m.Disabled,
		Metadata:   m.Metadata,
	}, nil
}
