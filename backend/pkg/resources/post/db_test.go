package post

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/db"
	"github.com/Southclaws/storyden/backend/pkg/resources/user"
)

func TestCreatePost(t *testing.T) {
	// a := assert.New(t)
	r := require.New(t)
	ctx := context.Background()

	d := db.TestDB(t)
	user.Seed(user.New(d))
	repo := New(d)

	p, err := repo.CreatePost(ctx, "body", user.SeedUser_02_User.ID, PostID(uuid.New()), nil)
	r.NoError(err)
	r.NotNil(p)
}
