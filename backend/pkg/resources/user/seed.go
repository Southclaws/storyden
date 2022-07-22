package user

import (
	"context"
	"fmt"

	"4d63.com/optional"
	"github.com/google/uuid"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model"
)

var (
	SeedUser_01_Admin = User{
		ID:    UserID(uuid.MustParse("00000000-0000-0000-0000-000000000000")),
		Email: "tim@storyd.en",
		Name:  "TimManTheTinMan",
		Bio:   optional.Of("I run this place"),
		Admin: true,
	}

	SeedUser_02_User = User{
		ID:    UserID(uuid.MustParse("00000000-0000-0000-0000-000000000000")),
		Email: "tam@storyd.en",
		Name:  "IDontLikeTom",
		Bio:   optional.Of("I'm just called Tam"),
	}
)

func NewLocalWithSeed() Repository {
	m := NewLocal()
	Seed(m)
	return m
}

func NewWithSeed(db *model.Client) Repository {
	m := New(db)
	Seed(m)
	return m
}

func Seed(r Repository) {
	ctx := context.Background()

	var u *User

	u, _ = r.CreateUser(ctx, SeedUser_01_Admin.Email, SeedUser_01_Admin.Name)
	SeedUser_01_Admin.ID = u.ID

	u, _ = r.CreateUser(ctx, SeedUser_02_User.Email, SeedUser_02_User.Name)
	SeedUser_02_User.ID = u.ID

	fmt.Println("created seed users", SeedUser_01_Admin.ID, SeedUser_02_User.ID)
}
