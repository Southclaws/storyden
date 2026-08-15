package robot_toolset_test

import (
	"context"
	"testing"

	"github.com/Southclaws/fault/ftag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot/robot_toolset"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestRepositoryPreventsDeletingAssignedToolset(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		repo *robot_toolset.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			author, err := db.Account.Create().
				SetHandle("toolset-owner").
				SetName("Toolset Owner").
				Save(ctx)
			require.NoError(t, err)

			created, err := repo.Create(
				ctx,
				account.AccountID(author.ID),
				"Research kit",
				"Reusable research tools",
				"Verify claims against sources.",
				[]string{"content_search", "library_search_pages"},
				map[string]any{"icon": "search"},
			)
			require.NoError(t, err)

			other, err := db.Account.Create().
				SetHandle("other-toolset-user").
				SetName("Other Toolset User").
				Save(ctx)
			require.NoError(t, err)
			isAuthor, err := repo.IsAuthor(ctx, created.ID, account.AccountID(author.ID))
			require.NoError(t, err)
			assert.True(t, isAuthor)
			isAuthor, err = repo.IsAuthor(ctx, created.ID, account.AccountID(other.ID))
			require.NoError(t, err)
			assert.False(t, isAuthor)

			stored, err := db.Robot.Create().
				SetName("Research Robot").
				SetDescription("Delegated research specialist").
				SetPlaybook("Research a bounded question.").
				SetModel("mock/test").
				SetToolsets([]string{created.ID.String()}).
				SetAuthorID(author.ID).
				Save(ctx)
			require.NoError(t, err)

			usage, err := repo.UsageCount(ctx, created.ID)
			require.NoError(t, err)
			assert.Equal(t, 1, usage)

			err = repo.Delete(ctx, created.ID)
			require.Error(t, err)
			assert.Equal(t, ftag.AlreadyExists, ftag.Get(err))

			_, err = stored.Update().SetToolsets([]string{}).Save(ctx)
			require.NoError(t, err)
			require.NoError(t, repo.Delete(ctx, created.ID))

			_, err = repo.Get(ctx, created.ID)
			require.Error(t, err)
			assert.Equal(t, ftag.NotFound, ftag.Get(err))
		}))
	}))
}
