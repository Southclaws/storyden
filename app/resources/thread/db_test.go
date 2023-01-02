package thread_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/internal/integration"
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
				seed.Account_002_Frigg.ID,
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
			a.Len(p.Posts, 1)
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

			a.Len(threads, 3)

			threads, err = repo.List(ctx, time.Now(), 10, thread.WithAuthor(seed.Account_001_Odin.ID))
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

			threads, err := repo.Get(ctx, seed.Post_01_Welcome.ID)
			r.NoError(err)
			r.NotNil(threads)

			a.Equal("Welcome to Storyden!", threads.Title)
			a.Equal("00000000000000000010-welcome-to-storyden", threads.Slug)
			a.Equal(false, threads.Pinned)
			a.False(threads.DeletedAt.IsPresent())
			a.Equal(seed.Category_01_General.ID, threads.Category.ID)

			r.Len(threads.Posts, 10)

			p0 := threads.Posts[0]
			a.Len(p0.Body, 2304)
			a.Equal(seed.Account_001_Odin.ID, p0.Author.ID)

			p1 := threads.Posts[1]
			a.Equal("first 😁", p1.Body)
			a.Equal(seed.Account_004_Loki.ID, p1.Author.ID)

			p2 := threads.Posts[2]
			a.Equal("Nice! One question: what kind of formatting can you use in posts? Is it like the old days with [b]tags[/b] and [color=red]cool stuff[/color] like that?", p2.Body)
			a.Equal(seed.Account_002_Frigg.ID, p2.Author.ID)
		}),
	)
}
