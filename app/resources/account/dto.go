package account

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/internal/ent"
)

type AccountID xid.ID

func (u AccountID) String() string { return xid.ID(u).String() }

type Account struct {
	ID          AccountID
	Handle      string
	Name        string
	Bio         opt.Optional[string]
	Admin       bool
	ThreadCount int
	PostCount   int
	Interests   []tag.Tag

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

func FromModel(u ent.Account) (o *Account) {
	result := Account{
		ID:     AccountID(u.ID),
		Handle: u.Handle,
		Name:   u.Name,
		Bio:    opt.New(u.Bio),
		Admin:  u.Admin,
		Interests: dt.Map(u.Edges.Tags, func(t *ent.Tag) tag.Tag {
			return tag.Tag{
				ID:   t.ID.String(),
				Name: t.Name,
			}
		}),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: opt.NewPtr(u.DeletedAt),
	}

	return &result
}
