package robot_session_test

import (
	"context"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/tests"
)

func TestRunnableSessionLimitAppliesAfterSessionDeduplication(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("skipping global runnable session ordering assertion on shared postgres database")
	}

	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		repo *robot_session.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			owner, err := db.Account.Create().SetHandle("runnable-session-" + xid.New().String()).SetName("Runnable Session Owner").Save(ctx)
			require.NoError(t, err)
			ownerID := account.AccountID(owner.ID)
			secondSessionID := robot.SessionID(xid.New())
			firstSessionID := robot.SessionID(xid.New())
			require.NoError(t, createSession(ctx, repo, firstSessionID, ownerID))
			require.NoError(t, createSession(ctx, repo, secondSessionID, ownerID))

			_, err = db.RobotSession.UpdateOneID(xid.ID(secondSessionID)).
				SetNextEventSequence(100).
				Save(ctx)
			require.NoError(t, err)

			require.NoError(t, enqueueVisibleInput(ctx, repo, firstSessionID, ownerID, robot.InputID(xid.New()), "first"))
			require.NoError(t, enqueueVisibleInput(ctx, repo, firstSessionID, ownerID, robot.InputID(xid.New()), "second"))
			require.NoError(t, enqueueVisibleInput(ctx, repo, secondSessionID, ownerID, robot.InputID(xid.New()), "third"))

			runnable, err := repo.RunnableSessionIDs(ctx, 2)
			require.NoError(t, err)
			assert.Equal(t, []robot.SessionID{firstSessionID, secondSessionID}, runnable)
		}))
	}))
}
