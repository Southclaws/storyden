package thread_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/utils/integration"
	"github.com/Southclaws/storyden/pkg/resources/seed"
	"github.com/Southclaws/storyden/pkg/services/thread"
)

func TestThreadCreate(t *testing.T) {
	defer integration.Test(t, nil, fx.Invoke(func(
		_ seed.Ready,
		ctx context.Context,
		thread_svc thread.Service,
	) {
		r := require.New(t)
		a := assert.New(t)

		t, err := thread_svc.Create(ctx, "New thread", "hi there", seed.SeedUser_02_User.ID, seed.SeedCategory_01_General.ID, []string{"hello"})
		r.NoError(err)
		r.NotNil(t)

		a.False(t.DeletedAt.IsPresent())
		a.Equal("New thread", t.Title)
		a.Contains(t.Slug, "new-thread")
		a.Equal("hi there", t.Short)
		a.Equal(false, t.Pinned)
		a.Equal(seed.SeedUser_02_User.Name, t.Author.Name)
		// a.Equal([]string{"hello"}, t.Tags) // TODO: upsert tags in resource
		a.Equal(seed.SeedCategory_01_General.Name, t.Category.Name)
		a.Len(t.Posts, 0)
		a.Len(t.Reacts, 0)
	}))
}
