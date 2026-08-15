package robot_session_test

import (
	"context"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	adksession "google.golang.org/adk/v2/session"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestAppendMessagePersistsADKBranchAndIsolationScope(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		repo *robot_session.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			author, err := db.Account.Create().
				SetHandle("session-owner").
				SetName("Session Owner").
				Save(ctx)
			require.NoError(t, err)

			sessionID := robot.SessionID(xid.New())
			_, err = repo.Create(ctx, sessionID, "Delegated task", account.AccountID(author.ID), map[string]any{})
			require.NoError(t, err)

			event := &adksession.Event{
				InvocationID:   "invocation-1",
				Author:         "robot_researcher",
				Branch:         "storyden.robot_researcher",
				IsolationScope: "delegation-call-1",
			}
			err = repo.AppendMessage(
				ctx,
				sessionID,
				opt.NewEmpty[account.AccountID](),
				opt.New(robot.NewBuiltinActor("researcher")),
				event,
			)
			require.NoError(t, err)

			session, _, err := repo.Get(ctx, sessionID, robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 10))
			require.NoError(t, err)
			require.Len(t, session.Messages, 1)

			message := session.Messages[0]
			assert.Equal(t, "storyden.robot_researcher", message.Branch.OrZero())
			assert.Equal(t, "delegation-call-1", message.IsolationScope.OrZero())
			assert.Equal(t, "storyden.robot_researcher", message.Event.Branch)
			assert.Equal(t, "delegation-call-1", message.Event.IsolationScope)
		}))
	}))
}

func TestAppendMessageRollsBackStateDeltaWhenMessageFails(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		repo *robot_session.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			author, err := db.Account.Create().
				SetHandle("state-rollback-owner").
				SetName("State Rollback Owner").
				Save(ctx)
			require.NoError(t, err)

			sessionID := robot.SessionID(xid.New())
			_, err = repo.Create(ctx, sessionID, "Transactional state", account.AccountID(author.ID), map[string]any{
				"active_tools": []string{"thread_get"},
			})
			require.NoError(t, err)

			invalidActor := robot.Actor{
				DatabaseRobotID: opt.New(xid.New()),
				BuiltinRobotID:  opt.New(robot.BuiltinRobotID("denbot")),
			}
			err = repo.AppendMessage(
				ctx,
				sessionID,
				opt.NewEmpty[account.AccountID](),
				opt.New(invalidActor),
				&adksession.Event{
					InvocationID: "invocation-invalid-actor",
					Author:       "denbot",
					Actions: adksession.EventActions{StateDelta: map[string]any{
						"active_tools": []string{"thread_get", "thread_update"},
					}},
				},
			)
			require.Error(t, err)

			session, _, err := repo.Get(ctx, sessionID, robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 10))
			require.NoError(t, err)
			assert.Equal(t, []any{"thread_get"}, session.State["active_tools"])
			assert.Empty(t, session.Messages)
		}))
	}))
}
