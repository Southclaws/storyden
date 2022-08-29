package thread_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/utils/integration"
	"github.com/Southclaws/storyden/pkg/resources/seed"
	"github.com/Southclaws/storyden/pkg/resources/thread"
)

func TestCreate(t *testing.T) {
	defer integration.Test(t, nil, fx.Invoke(
		func(
			_ seed.Ready,
			ctx context.Context,
			repo thread.Repository,
		) {
			r := require.New(t)
			a := assert.New(t)

			p, err := repo.Create(ctx,
				"A Super Nice Thread",
				"Lorem ipsum",
				seed.Account_002.ID,
				seed.Category_01_General.ID,
				[]string{})
			r.NoError(err)
			r.NotNil(p)

			a.Equal("A Super Nice Thread", p.Title)
			a.Contains(p.Slug, "a-super-nice-thread")
			a.Equal(false, p.Pinned)
			a.WithinDuration(p.CreatedAt, time.Now(), time.Second*5)
			a.WithinDuration(p.UpdatedAt, time.Now(), time.Second*5)
			a.False(p.DeletedAt.IsPresent())
			a.Equal(seed.Category_01_General.ID, p.Category.ID)
			a.Len(p.Posts, 0)
		}),
	)
}

func TestList(t *testing.T) {
	defer integration.Test(t, nil, fx.Invoke(
		func(
			_ seed.Ready,
			ctx context.Context,
			repo thread.Repository,
		) {
			r := require.New(t)
			a := assert.New(t)

			threads, err := repo.List(ctx, time.Now(), 10)
			r.NoError(err)
			r.NotNil(threads)

			a.Len(threads, 2)
		}),
	)
}

func TestGet(t *testing.T) {
	defer integration.Test(t, nil, fx.Invoke(
		func(
			_ seed.Ready,
			ctx context.Context,
			repo thread.Repository,
		) {
			r := require.New(t)
			a := assert.New(t)

			threads, err := repo.Get(ctx, seed.Post_01.ID)
			r.NoError(err)
			r.NotNil(threads)

			a.Equal("Hello world!", threads.Title)
			a.Contains(threads.Slug, "hello-world")
			a.Equal(false, threads.Pinned)
			a.False(threads.DeletedAt.IsPresent())
			a.Equal(seed.Category_01_General.ID, threads.Category.ID)

			r.Len(threads.Posts, 2)

			p0 := threads.Posts[0]
			a.Equal("First reply", p0.Body)
			a.Equal(seed.Account_003.ID, p0.Author.ID)

			p1 := threads.Posts[1]
			a.Equal("Second reply", p1.Body)
			a.Equal(seed.Account_004.ID, p1.Author.ID)
		}),
	)
}
