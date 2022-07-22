package user

import (
	"time"

	"4d63.com/optional"
	"github.com/google/uuid"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/backend/internal/utils"
)

type UserID uuid.UUID

func (u UserID) String() string { return uuid.UUID(u).String() }

type User struct {
	ID          UserID                    `json:"id"`
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

const Role = "User"

func (u *User) GetRole() string { return Role }

func FromModel(u model.User) (o User) {
	result := User{
		ID:        UserID(u.ID),
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

func FromModelPublic(u model.User) (o User) {
	m := FromModel(u)
	m.Email = ""

	return m
}
