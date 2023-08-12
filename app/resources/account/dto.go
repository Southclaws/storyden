package account

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
)

type AccountID xid.ID

func (u AccountID) String() string { return xid.ID(u).String() }

type Account struct {
	ID     AccountID
	Handle string
	Name   string
	Bio    opt.Optional[string]
	Admin  bool
	Auths  []string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt opt.Optional[time.Time]
}

// Name is the role/resource name.
const Name = "Account"

func (a *Account) GetRole() string {
	if a.Admin {
		return "everyone"
	}

	return "owner"
}

func (*Account) GetResourceName() string { return Name }

func FromModel(a *ent.Account) (*Account, error) {
	auths := dt.Map(a.Edges.Authentication, func(a *ent.Authentication) string {
		return a.Service
	})

	return &Account{
		ID:     AccountID(a.ID),
		Handle: a.Handle,
		Name:   a.Name,
		Bio:    opt.New(a.Bio),
		Admin:  a.Admin,
		Auths:  auths,

		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		DeletedAt: opt.NewPtr(a.DeletedAt),
	}, nil
}
