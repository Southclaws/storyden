package account

import (
	"time"

	"4d63.com/optional"
	"github.com/Southclaws/dt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
)

type AccountID xid.ID

func (u AccountID) String() string { return xid.ID(u).String() }

type Account struct {
	ID          AccountID
	Handle      string
	Name        string
	Bio         optional.Optional[string]
	Admin       bool
	ThreadCount int
	PostCount   int
	Interests   []tag.Tag

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt optional.Optional[time.Time]
}

// Name is the role/resource name.
const Name = "Account"

func (*Account) GetRole() string { return Name }

func (*Account) GetResourceName() string { return Name }

func FromModel(u model.Account) (o *Account) {
	result := Account{
		ID:     AccountID(u.ID),
		Handle: u.Handle,
		Name:   u.Name,
		Bio:    optional.Of(u.Bio),
		Admin:  u.Admin,
		Interests: dt.Map(u.Edges.Tags, func(t *model.Tag) tag.Tag {
			return tag.Tag{
				ID:   t.ID.String(),
				Name: t.Name,
			}
		}),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: optional.OfPtr(u.DeletedAt),
	}

	return &result
}
