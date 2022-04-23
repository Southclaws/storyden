package thread

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/api/src/infra/db"
	"github.com/Southclaws/storyden/api/src/resources/category"
	"github.com/Southclaws/storyden/api/src/resources/user"
	"github.com/Southclaws/storyden/api/src/utils"
)

func implementations(t *testing.T, seed bool) []utils.ImplConstructor[Repository] {
	if seed {
		return []utils.ImplConstructor[Repository]{
			func() Repository { return NewWithSeed(db.TestDB(t)) },
			func() Repository { return NewLocalWithSeed() },
		}
	} else {
		return []utils.ImplConstructor[Repository]{
			func() Repository { return New(db.TestDB(t)) },
			func() Repository { return NewLocal() },
		}
	}
}

func TestCreatePost(t *testing.T) {
	utils.TestAll(t, implementations(t, false), func(t1 *testing.T, repo Repository) {
		r := require.New(t)
		ctx := context.Background()

		d := db.TestDB(t)
		category.Seed(category.New(d))
		user.Seed(user.New(d))

		p, err := repo.CreateThread(ctx, "title", "body", user.SeedUser_02_User.ID, category.SeedCategory_01_General.ID, []string{})
		r.NoError(err)
		r.NotNil(p)
	})
}
