package robot_session_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	adksession "google.golang.org/adk/v2/session"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestExecutionLeaseSerialisesConcurrentTurns(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		repo *robot_session.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			owner, err := db.Account.Create().
				SetHandle("lease-owner").
				SetName("Lease Owner").
				Save(ctx)
			require.NoError(t, err)

			sessionID := robot.SessionID(xid.New())
			_, err = repo.Create(ctx, sessionID, "Concurrent session", account.AccountID(owner.ID), nil)
			require.NoError(t, err)

			const contenders = 8
			start := make(chan struct{})
			leases := make(chan *robot_session.ExecutionLease, contenders)
			errs := make(chan error, contenders)
			var wg sync.WaitGroup
			for range contenders {
				wg.Add(1)
				go func() {
					defer wg.Done()
					<-start
					lease, err := repo.AcquireExecution(ctx, sessionID, xid.New(), false, time.Minute)
					if err != nil {
						errs <- err
						return
					}
					leases <- lease
				}()
			}
			close(start)
			wg.Wait()
			close(leases)
			close(errs)

			var winner *robot_session.ExecutionLease
			for lease := range leases {
				require.Nil(t, winner, "only one concurrent turn may acquire the session")
				winner = lease
			}
			require.NotNil(t, winner)

			busy := 0
			for err := range errs {
				require.ErrorIs(t, err, robot_session.ErrSessionBusy)
				busy++
			}
			assert.Equal(t, contenders-1, busy)
			require.NoError(t, repo.ReleaseExecution(ctx, *winner))
		}))
	}))
}

func TestExecutionLeaseFencesExpiredWorker(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		repo *robot_session.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			owner, err := db.Account.Create().
				SetHandle("fencing-owner").
				SetName("Fencing Owner").
				Save(ctx)
			require.NoError(t, err)

			sessionID := robot.SessionID(xid.New())
			_, err = repo.Create(ctx, sessionID, "Fenced session", account.AccountID(owner.ID), nil)
			require.NoError(t, err)

			stale, err := repo.AcquireExecution(ctx, sessionID, xid.New(), false, -time.Second)
			require.NoError(t, err)
			current, err := repo.AcquireExecution(ctx, sessionID, xid.New(), false, time.Minute)
			require.NoError(t, err)

			assert.Greater(t, current.Generation, stale.Generation)
			staleCtx := robot_session.WithExecutionLease(ctx, *stale)
			err = repo.AppendMessage(
				staleCtx,
				sessionID,
				opt.NewEmpty[account.AccountID](),
				opt.NewEmpty[robot.Actor](),
				&adksession.Event{InvocationID: "stale-worker"},
			)
			require.ErrorIs(t, err, robot_session.ErrLeaseLost)
			messageCount, err := db.RobotSessionMessage.Query().Count(ctx)
			require.NoError(t, err)
			assert.Zero(t, messageCount)
			require.ErrorIs(t, repo.ReleaseExecution(ctx, *stale), robot_session.ErrLeaseLost)
			require.NoError(t, repo.ReleaseExecution(ctx, *current))
		}))
	}))
}

func TestBlockedSessionRequiresContinuation(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		repo *robot_session.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			owner, err := db.Account.Create().
				SetHandle("blocked-owner").
				SetName("Blocked Owner").
				Save(ctx)
			require.NoError(t, err)

			sessionID := robot.SessionID(xid.New())
			_, err = repo.Create(ctx, sessionID, "Blocked session", account.AccountID(owner.ID), nil)
			require.NoError(t, err)

			blocked, err := repo.AcquireExecution(ctx, sessionID, xid.New(), false, time.Minute)
			require.NoError(t, err)
			leaseCtx := robot_session.WithExecutionLease(ctx, *blocked)
			require.NoError(t, repo.StorePendingClientTools(leaseCtx, sessionID, []string{"confirmation-1"}))
			require.NoError(t, repo.BlockExecution(ctx, *blocked))

			_, err = repo.AcquireExecution(ctx, sessionID, xid.New(), false, time.Minute)
			require.ErrorIs(t, err, robot_session.ErrSessionBusy)

			continued, err := repo.AcquireExecution(ctx, sessionID, xid.New(), true, time.Minute)
			require.NoError(t, err)
			stored, err := db.RobotSession.Get(ctx, xid.ID(sessionID))
			require.NoError(t, err)
			assert.Empty(t, robot.PendingClientToolIDs(stored.State))
			require.NoError(t, repo.ReleaseExecution(ctx, *continued))
		}))
	}))
}

func TestSessionViewsAreIndependentFromSessionCreator(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		repo *robot_session.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			creator, err := db.Account.Create().
				SetHandle("view-creator").
				SetName("View Creator").
				Save(ctx)
			require.NoError(t, err)
			member, err := db.Account.Create().
				SetHandle("view-member").
				SetName("View Member").
				Save(ctx)
			require.NoError(t, err)

			sessionID := robot.SessionID(xid.New())
			_, err = repo.Create(ctx, sessionID, "Shared session", account.AccountID(creator.ID), nil)
			require.NoError(t, err)
			require.NoError(t, repo.EnsureView(ctx, sessionID, account.AccountID(member.ID)))

			stored, err := db.RobotSession.Get(ctx, xid.ID(sessionID))
			require.NoError(t, err)
			storedCreator, err := stored.QueryCreator().Only(ctx)
			require.NoError(t, err)
			assert.Equal(t, creator.ID, storedCreator.ID, "adding a view must not change session creator provenance")

			result, err := repo.List(ctx, pagination.NewPageParams(1, 10), opt.New(account.AccountID(member.ID)))
			require.NoError(t, err)
			require.Len(t, result.Items, 1)
			assert.Equal(t, sessionID.String(), result.Items[0].ID.String())
			assert.Equal(t, creator.ID, xid.ID(result.Items[0].Human.ID))

			require.NoError(t, db.Account.DeleteOneID(member.ID).Exec(ctx))
			_, err = db.RobotSession.Get(ctx, xid.ID(sessionID))
			require.NoError(t, err, "removing a member must not delete the shared session")
		}))
	}))
}
