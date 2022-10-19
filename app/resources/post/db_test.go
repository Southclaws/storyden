package post_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/internal/utils/integration"
)

func TestCreate(t *testing.T) {
	defer integration.Test(t, nil, fx.Invoke(
		func(
			_ seed.Ready,
			ctx context.Context,
			repo post.Repository,
		) {
			r := require.New(t)
			a := assert.New(t)

			p, err := repo.Create(ctx,
				"My cool post in reply to your thread",
				seed.Account_002.ID,
				seed.Post_01.ID,
				nil, nil)
			r.NoError(err)
			r.NotNil(p)

			a.Equal("My cool post in reply to your thread", p.Body)
			a.WithinDuration(p.CreatedAt, time.Now(), time.Second*5)
			a.WithinDuration(p.UpdatedAt, time.Now(), time.Second*5)
			a.False(p.DeletedAt.IsPresent())
		}),
	)
}
