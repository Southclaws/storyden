package user

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Southclaws/storyden/api/src/infra/db"
	"github.com/Southclaws/storyden/api/src/utils"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func implementations(t *testing.T, seed bool) []utils.ImplConstructor[Repository] {
	if seed {
		return []utils.ImplConstructor[Repository]{
			func() Repository { return NewWithSeed(db.TestDB(t)) },
			func() Repository { return NewMockWithSeed() },
		}
	} else {
		return []utils.ImplConstructor[Repository]{
			func() Repository { return New(db.TestDB(t)) },
			func() Repository { return NewMock() },
		}
	}
}

func TestCreateUser(t *testing.T) {
	utils.TestAll(t, implementations(t, false),
		func(t *testing.T, r Repository) {
			ctx := context.Background()

			u, err := r.CreateUser(ctx, SeedUser_01_Admin.Email, SeedUser_01_Admin.Name)
			require.NoError(t, err)
			require.NotNil(t, u)

			assert.Equal(t, SeedUser_01_Admin.Email, u.Email)
			assert.Equal(t, SeedUser_01_Admin.Name, u.Name)

			u1, err := r.GetUser(ctx, u.ID, false)
			require.NoError(t, err)
			assert.NotNil(t, u1)

			assert.Equal(t, SeedUser_01_Admin.Email, u1.Email)
			assert.Equal(t, SeedUser_01_Admin.Name, u1.Name)

			// Duplicate email address should fail.
			u2, err := r.CreateUser(ctx, SeedUser_01_Admin.Email, SeedUser_01_Admin.Name)
			require.Error(t, err)
			assert.Nil(t, u2)
		})
}

func TestGetByID(t *testing.T) {
	utils.TestAll(t, implementations(t, false),
		func(t *testing.T, r Repository) {
			ctx := context.Background()

			none, err := r.GetUser(ctx, SeedUser_01_Admin.ID, false)
			require.NoError(t, err)
			assert.Nil(t, none)

			u, err := r.CreateUser(ctx, SeedUser_01_Admin.Email, SeedUser_01_Admin.Name)
			require.NoError(t, err)

			u, err = r.GetUser(ctx, u.ID, false)
			require.NoError(t, err)
			assert.NotNil(t, u)
		})
}

func TestGetByEmail(t *testing.T) {
	utils.TestAll(t, implementations(t, false),
		func(t *testing.T, r Repository) {
			ctx := context.Background()

			none, err := r.GetUserByEmail(ctx, SeedUser_01_Admin.Email, false)
			require.NoError(t, err)
			assert.Nil(t, none)

			u, err := r.CreateUser(ctx, SeedUser_01_Admin.Email, SeedUser_01_Admin.Name)
			require.NoError(t, err)

			u, err = r.GetUserByEmail(ctx, SeedUser_01_Admin.Email, false)
			require.NoError(t, err)
			assert.NotNil(t, u)
		})
}

func TestGetAll(t *testing.T) {
	utils.TestAll(t, implementations(t, true),
		func(t *testing.T, r Repository) {
			ctx := context.Background()

			u, err := r.GetUsers(ctx, "asc", 10, 0, false)
			require.NoError(t, err)
			assert.NotNil(t, u)

			emails := lo.Map(u, func(t User, i int) string { return t.Email })

			assert.Contains(t, emails, SeedUser_01_Admin.Email)
			assert.Contains(t, emails, SeedUser_02_User.Email)
		})
}

func TestUpdateUser(t *testing.T) {
	utils.TestAll(t, implementations(t, true),
		func(t *testing.T, r Repository) {
			ctx := context.Background()

			fmt.Println("BEFORE GET", SeedUser_02_User.ID)

			before, err := r.GetUser(ctx, SeedUser_02_User.ID, false)
			fmt.Println(before, err, SeedUser_02_User.ID)
			require.NoError(t, err)
			assert.NotNil(t, before)

			after, err := r.UpdateUser(ctx, SeedUser_02_User.ID, utils.Ref("timmy@storyd.en"), nil, nil)
			require.NoError(t, err)
			require.NotNil(t, after)

			assert.Equal(t, "timmy@storyd.en", after.Email)
		})
}

func TestSetAdmin(t *testing.T) {
	utils.TestAll(t, implementations(t, true),
		func(t *testing.T, r Repository) {
			ctx := context.Background()

			err := r.SetAdmin(ctx, SeedUser_02_User.ID, true)
			require.NoError(t, err)

			after, err := r.GetUser(ctx, SeedUser_02_User.ID, false)
			require.NoError(t, err)
			require.NotNil(t, after)
			assert.True(t, after.Admin)
		})
}

func TestBan(t *testing.T) {
	utils.TestAll(t, implementations(t, true),
		func(t *testing.T, r Repository) {
			ctx := context.Background()

			u, err := r.Ban(ctx, SeedUser_02_User.ID)
			require.NoError(t, err)
			require.NotNil(t, u)

			after, err := r.GetUser(ctx, SeedUser_02_User.ID, false)
			require.NoError(t, err)
			require.NotNil(t, after)

			assert.True(t, after.DeletedAt.IsPresent())
			assert.WithinDuration(t, time.Now(), after.DeletedAt.ElseZero(), time.Second)
		})
}

func TestUnban(t *testing.T) {
	utils.TestAll(t, implementations(t, true),
		func(t *testing.T, r Repository) {
			ctx := context.Background()

			u1, err := r.Ban(ctx, SeedUser_02_User.ID)
			require.NoError(t, err)
			require.NotNil(t, u1)

			u2, err := r.GetUser(ctx, SeedUser_02_User.ID, false)
			require.NoError(t, err)
			require.NotNil(t, u2)

			assert.True(t, u2.DeletedAt.IsPresent())
			assert.WithinDuration(t, time.Now(), u2.DeletedAt.ElseZero(), time.Second)

			u3, err := r.Unban(ctx, SeedUser_02_User.ID)
			require.NoError(t, err)
			require.NotNil(t, u3)

			assert.False(t, u3.DeletedAt.IsPresent())
		})
}
