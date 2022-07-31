package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/internal/utils"
	"github.com/Southclaws/storyden/backend/internal/utils/bdd"
	"github.com/Southclaws/storyden/backend/pkg/resources/seed"
	"github.com/Southclaws/storyden/backend/pkg/resources/user"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	bdd.Test(t, nil, fx.Invoke(func(repo user.Repository) {
		r := require.New(t)
		a := assert.New(t)

		u, err := repo.CreateUser(ctx, seed.SeedUser_01_Admin.Email, seed.SeedUser_01_Admin.Name)
		r.NoError(err)
		r.NotNil(u)

		a.Equal(seed.SeedUser_01_Admin.Email, u.Email)
		a.Equal(seed.SeedUser_01_Admin.Name, u.Name)

		u1, err := repo.GetUser(ctx, u.ID, false)
		r.NoError(err)
		a.NotNil(u1)

		a.Equal(seed.SeedUser_01_Admin.Email, u1.Email)
		a.Equal(seed.SeedUser_01_Admin.Name, u1.Name)

		// Duplicate email address should fail.
		u2, err := repo.CreateUser(ctx, seed.SeedUser_01_Admin.Email, seed.SeedUser_01_Admin.Name)
		r.Error(err)
		a.Nil(u2)
	}))
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()

	bdd.Test(t, nil, fx.Invoke(func(repo user.Repository) {
		r := require.New(t)
		a := assert.New(t)

		none, err := repo.GetUser(ctx, seed.SeedUser_01_Admin.ID, false)
		r.NoError(err)
		a.Nil(none)

		u, err := repo.CreateUser(ctx, seed.SeedUser_01_Admin.Email, seed.SeedUser_01_Admin.Name)
		r.NoError(err)

		u, err = repo.GetUser(ctx, u.ID, false)
		r.NoError(err)
		a.NotNil(u)
	}))
}

func TestGetByEmail(t *testing.T) {
	ctx := context.Background()

	bdd.Test(t, nil, fx.Invoke(func(repo user.Repository) {
		r := require.New(t)
		a := assert.New(t)

		none, err := repo.GetUserByEmail(ctx, seed.SeedUser_01_Admin.Email, false)
		r.NoError(err)
		a.Nil(none)

		u, err := repo.CreateUser(ctx, seed.SeedUser_01_Admin.Email, seed.SeedUser_01_Admin.Name)
		r.NoError(err)

		u, err = repo.GetUserByEmail(ctx, seed.SeedUser_01_Admin.Email, false)
		r.NoError(err)
		a.NotNil(u)
	}))
}

func TestGetAll(t *testing.T) {
	ctx := context.Background()

	bdd.Test(t, nil, fx.Invoke(
		func(
			_ seed.Ready,
			repo user.Repository,
		) {
			r := require.New(t)
			a := assert.New(t)

			u, err := repo.GetUsers(ctx, "asc", 10, 0, false)
			r.NoError(err)
			a.NotNil(u)

			emails := lo.Map(u, func(t user.User, i int) string { return t.Email })

			a.Contains(emails, seed.SeedUser_01_Admin.Email)
			a.Contains(emails, seed.SeedUser_02_User.Email)
		}))
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()

	bdd.Test(t, nil, fx.Invoke(func(_ seed.Ready,
		repo user.Repository,
	) {
		r := require.New(t)
		a := assert.New(t)

		before, err := repo.GetUser(ctx, seed.SeedUser_02_User.ID, false)
		r.NoError(err)
		a.NotNil(before)

		after, err := repo.UpdateUser(ctx, seed.SeedUser_02_User.ID, utils.Ref("timmy@storyd.en"), nil, nil)
		r.NoError(err)
		r.NotNil(after)

		a.Equal("timmy@storyd.en", after.Email)
	}))
}

func TestSetAdmin(t *testing.T) {
	ctx := context.Background()

	bdd.Test(t, nil, fx.Invoke(func(
		_ seed.Ready,
		repo user.Repository,
	) {
		r := require.New(t)
		a := assert.New(t)

		err := repo.SetAdmin(ctx, seed.SeedUser_02_User.ID, true)
		r.NoError(err)

		after, err := repo.GetUser(ctx, seed.SeedUser_02_User.ID, false)
		r.NoError(err)
		r.NotNil(after)
		a.True(after.Admin)
	}))
}

func TestBan(t *testing.T) {
	ctx := context.Background()

	bdd.Test(t, nil, fx.Invoke(func(
		_ seed.Ready,
		repo user.Repository,
	) {
		r := require.New(t)
		a := assert.New(t)

		u, err := repo.Ban(ctx, seed.SeedUser_02_User.ID)
		r.NoError(err)
		r.NotNil(u)

		after, err := repo.GetUser(ctx, seed.SeedUser_02_User.ID, false)
		r.NoError(err)
		r.NotNil(after)

		a.True(after.DeletedAt.IsPresent())
		a.WithinDuration(time.Now(), after.DeletedAt.ElseZero(), time.Second)
	}))
}

func TestUnban(t *testing.T) {
	ctx := context.Background()

	bdd.Test(t, nil, fx.Invoke(func(
		_ seed.Ready,
		repo user.Repository,
	) {
		r := require.New(t)
		a := assert.New(t)

		u1, err := repo.Ban(ctx, seed.SeedUser_02_User.ID)
		r.NoError(err)
		r.NotNil(u1)

		u2, err := repo.GetUser(ctx, seed.SeedUser_02_User.ID, false)
		r.NoError(err)
		r.NotNil(u2)

		a.True(u2.DeletedAt.IsPresent())
		a.WithinDuration(time.Now(), u2.DeletedAt.ElseZero(), time.Second)

		u3, err := repo.Unban(ctx, seed.SeedUser_02_User.ID)
		r.NoError(err)
		r.NotNil(u3)

		a.False(u3.DeletedAt.IsPresent())
	}))
}
