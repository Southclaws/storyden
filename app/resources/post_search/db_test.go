package post_search_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/post_search"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/internal/integration"
)

func Test_database_Search(t *testing.T) {
	defer integration.Test(t, nil, fx.Invoke(
		func(
			_ seed.Ready,
			ctx context.Context,
			repo post_search.Repository,
		) {
			r := require.New(t)
			a := assert.New(t)

			empty, err := repo.Search(ctx)
			r.NoError(err)
			r.NotNil(empty)

			a.Empty(empty)

			odins, err := repo.Search(ctx, post_search.WithAuthorHandle("odin"))
			r.NoError(err)
			r.NotNil(odins)

			a.NotEmpty(odins)
			a.Len(odins, 3)

			lokis, err := repo.Search(ctx, post_search.WithAuthorHandle("loki"))
			r.NoError(err)
			r.NotNil(lokis)

			a.NotEmpty(lokis)
			a.Len(lokis, 4)

			storydens, err := repo.Search(ctx, post_search.WithBodyContains("storyden"))
			r.NoError(err)
			r.NotNil(storydens)

			a.NotEmpty(storydens)
			a.Len(storydens, 2)

			welcomes, err := repo.Search(ctx, post_search.WithTitleContains("welcome"))
			r.NoError(err)
			r.NotNil(welcomes)

			a.NotEmpty(welcomes)
			a.Len(welcomes, 1)

			odinwelcomes, err := repo.Search(ctx,
				post_search.WithAuthorHandle("odin"),
				post_search.WithTitleContains("welcome"),
			)
			r.NoError(err)
			r.NotNil(odinwelcomes)

			a.NotEmpty(odinwelcomes)
			a.Len(odinwelcomes, 1)

			odinthreads, err := repo.Search(ctx,
				post_search.WithKinds(post_search.KindThread),
				post_search.WithAuthorHandle("odin"),
			)
			r.NoError(err)
			r.NotNil(odinthreads)

			a.NotEmpty(odinthreads)
			a.Len(odinthreads, 2)

			odinposts, err := repo.Search(ctx,
				post_search.WithKinds(post_search.KindPost),
				post_search.WithAuthorHandle("odin"),
			)
			r.NoError(err)
			r.NotNil(odinposts)

			a.NotEmpty(odinposts)
			a.Len(odinposts, 1)

			odinthreadsandposts, err := repo.Search(ctx,
				post_search.WithKinds(post_search.KindPost, post_search.KindThread),
				post_search.WithAuthorHandle("odin"),
			)
			r.NoError(err)
			r.NotNil(odinthreadsandposts)

			a.NotEmpty(odinthreadsandposts)
			a.Len(odinthreadsandposts, 3)
		}),
	)
}
