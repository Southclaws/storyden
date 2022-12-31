package account_test

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestCreateUser(t *testing.T) {
	defer integration.Test(t, nil, fx.Invoke(func(ctx context.Context, repo account.Repository) {
		r := require.New(t)
		a := assert.New(t)

		u, err := repo.Create(ctx, "southclaws")
		r.NoError(err)
		r.NotNil(u)

		a.Equal("southclaws", u.Name)

		u1, err := repo.GetByID(ctx, u.ID)
		r.NoError(err)
		a.NotNil(u1)

		a.Equal("southclaws", u1.Name)

		// Duplicate email address should fail.
		u2, err := repo.Create(ctx, "southclaws")
		r.Error(err)
		a.Nil(u2)
	}))
}

func TestGetByID(t *testing.T) {
	defer integration.Test(t, nil, fx.Invoke(func(ctx context.Context, repo account.Repository) {
		r := require.New(t)
		a := assert.New(t)

		acc, err := repo.GetByID(ctx, seed.Account_001_Odin.ID)
		r.NoError(err)
		r.NotNil(acc)

		a.Equal(seed.Account_001_Odin.Name, acc.Name)
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

			handles := lo.Map(u, func(t *account.Account, i int) string { return t.Handle })

			a.Contains(handles, seed.Account_001_Odin.Handle)
			a.Contains(handles, seed.Account_002_Frigg.Handle)
		}))
}
