package account

import (
	"time"

	"4d63.com/optional"
	"github.com/google/uuid"

	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/internal/utils"
)

type AccountID uuid.UUID

func (u AccountID) String() string { return uuid.UUID(u).String() }

type Account struct {
	ID          AccountID                 `json:"id"`
	Email       string                    `json:"email"`
	Name        string                    `json:"name"`
	Bio         optional.Optional[string] `json:"bio"`
	Admin       bool                      `json:"admin"`
	ThreadCount int                       `json:"threadCount"`
	PostCount   int                       `json:"postCount"`

	CreatedAt time.Time                    `json:"createdAt"`
	UpdatedAt time.Time                    `json:"updatedAt"`
	DeletedAt optional.Optional[time.Time] `json:"deletedAt"`
}

// Name is the role/resource name.
const Name = "Account"

func (*Account) GetRole() string { return Name }

func (*Account) GetResourceName() string { return Name }

func FromModel(u model.Account) (o Account) {
	result := Account{
		ID:        AccountID(u.ID),
		Email:     u.Email,
		Name:      u.Name,
		Bio:       optional.Of(u.Bio),
		Admin:     u.Admin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: utils.OptionalZero(u.DeletedAt),
	}

	return result
}
