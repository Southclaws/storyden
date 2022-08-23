package account_test

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/utils/integration"
	"github.com/Southclaws/storyden/pkg/resources/account"
	"github.com/Southclaws/storyden/pkg/resources/seed"
)

func TestCreateUser(t *testing.T) {
	defer integration.Test(t, nil, fx.Invoke(func(ctx context.Context, repo account.Repository) {
		r := require.New(t)
		a := assert.New(t)

		u, err := repo.Create(ctx, seed.SeedUser_01_Admin.Email, seed.SeedUser_01_Admin.Name)
		r.NoError(err)
		r.NotNil(u)

		a.Equal(seed.SeedUser_01_Admin.Email, u.Email)
		a.Equal(seed.SeedUser_01_Admin.Name, u.Name)

		u1, err := repo.GetByID(ctx, u.ID)
		r.NoError(err)
		a.NotNil(u1)

		a.Equal(seed.SeedUser_01_Admin.Email, u1.Email)
		a.Equal(seed.SeedUser_01_Admin.Name, u1.Name)

		// Duplicate email address should fail.
		u2, err := repo.Create(ctx, seed.SeedUser_01_Admin.Email, seed.SeedUser_01_Admin.Name)
		r.Error(err)
		a.Nil(u2)
	}))
}

func TestGetByID(t *testing.T) {
	defer integration.Test(t, nil, fx.Invoke(func(ctx context.Context, repo account.Repository) {
		r := require.New(t)
		a := assert.New(t)

		none, err := repo.GetByID(ctx, seed.SeedUser_01_Admin.ID)
		r.NoError(err)
		a.Nil(none)

		u, err := repo.Create(ctx, seed.SeedUser_01_Admin.Email, seed.SeedUser_01_Admin.Name)
		r.NoError(err)

		u, err = repo.GetByID(ctx, u.ID)
		r.NoError(err)
		a.NotNil(u)
	}))
}

func TestGetByEmail(t *testing.T) {
	defer integration.Test(t, nil, fx.Invoke(func(ctx context.Context, repo account.Repository) {
		r := require.New(t)
		a := assert.New(t)

		none, ok, err := repo.LookupByEmail(ctx, seed.SeedUser_01_Admin.Email)
		r.NoError(err)
		r.False(ok)
		a.Nil(none)

		u, err := repo.Create(ctx, seed.SeedUser_01_Admin.Email, seed.SeedUser_01_Admin.Name)
		r.NoError(err)

		u, ok, err = repo.LookupByEmail(ctx, seed.SeedUser_01_Admin.Email)
		r.NoError(err)
		r.True(ok)
		a.NotNil(u)
	}))
}

func TestGetAll(t *testing.T) {
	defer integration.Test(t, nil, fx.Invoke(
		func(
			_ seed.Ready,
			ctx context.Context,
			repo account.Repository,
		) {
			r := require.New(t)
			a := assert.New(t)

			u, err := repo.List(ctx, "asc", 10, 0)
			r.NoError(err)
			a.NotNil(u)

			emails := lo.Map(u, func(t account.Account, i int) string { return t.Email })

			a.Contains(emails, seed.SeedUser_01_Admin.Email)
			a.Contains(emails, seed.SeedUser_02_User.Email)
		}))
}
