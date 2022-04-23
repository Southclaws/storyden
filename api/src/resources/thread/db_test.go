package thread

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/api/src/infra/db"
	"github.com/Southclaws/storyden/api/src/resources/user"
)

func TestCreatePost(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

	d := db.TestDB(t)
	user.Seed(user.New(d))
	repo := New(d)

	p, err := repo.CreateThread(ctx, "title", "body", user.SeedUser_02_User.ID, "", []string{})
	r.NoError(err)
	r.NotNil(p)
}
