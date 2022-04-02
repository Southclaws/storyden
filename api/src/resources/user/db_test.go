package user

import (
	"context"
	"testing"

	"github.com/Southclaws/storyden/api/src/infra/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	r := New(db)

	u, err := r.CreateUser(ctx, SeedUser_01_Admin.Email, SeedUser_01_Admin.Name)
	require.NoError(t, err)

	assert.Equal(t, SeedUser_01_Admin.Email, u.Email)
	assert.Equal(t, SeedUser_01_Admin.Name, u.Name)
}
