package thread_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/seed"
	thread_repo "github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/app/services/thread"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestThreadCreate(t *testing.T) {
	defer integration.Test(t, nil, fx.Invoke(func(
		_ seed.Ready,
		ctx context.Context,
		thread_svc thread.Service,
	) {
		r := require.New(t)
		a := assert.New(t)

		t, err := thread_svc.Create(ctx,
			"New thread",
			"hi there",
			seed.Account_002_Frigg.ID,
			seed.Category_01_General.ID,
			thread_repo.StatusPublished,
			[]string{"hello"},
			nil)
		r.NoError(err)
		r.NotNil(t)

		a.False(t.DeletedAt.IsPresent())
		a.Equal("New thread", t.Title)
		a.Contains(t.Slug, "new-thread")
		a.Equal("hi there", t.Short)
		a.Equal(false, t.Pinned)
		a.Equal(seed.Account_002_Frigg.Name, t.Author.Name)
		// a.Equal([]string{"hello"}, t.Tags) // TODO: upsert tags in resource
		a.Equal(seed.Category_01_General.Name, t.Category.Name)
		a.Len(t.Posts, 1)
		a.Len(t.Reacts, 0)
	}))
}
