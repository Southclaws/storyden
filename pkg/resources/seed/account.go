package seed

import (
	"context"
	"fmt"

	"4d63.com/optional"
	"github.com/google/uuid"

	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/pkg/resources/account"
)

var (
	SeedUser_01_Admin = account.Account{
		ID:    account.AccountID(uuid.MustParse("00000000-0000-0000-0000-000000000000")),
		Email: "tim@storyd.en",
		Name:  "TimManTheTinMan",
		Bio:   optional.Of("I run this place"),
		Admin: true,
	}

	SeedUser_02_User = account.Account{
		ID:    account.AccountID(uuid.MustParse("00000000-0000-0000-0000-000000000000")),
		Email: "tam@storyd.en",
		Name:  "IDontLikeTom",
		Bio:   optional.Of("I'm just called Tam"),
	}
)

func users(r account.Repository) {
	ctx := context.Background()

	var u *account.Account

	u = utils.Must(r.Create(ctx, SeedUser_01_Admin.Email, SeedUser_01_Admin.Name))
	SeedUser_01_Admin.ID = u.ID

	u = utils.Must(r.Create(ctx, SeedUser_02_User.Email, SeedUser_02_User.Name))
	SeedUser_02_User.ID = u.ID

	fmt.Println("created seed users")
}
