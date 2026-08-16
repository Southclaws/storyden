package robot_session_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"google.golang.org/adk/v2/model"
	adksession "google.golang.org/adk/v2/session"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestTurnLifecycleUsesOrderedSessionEvents(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		repo *robot_session.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			owner, err := db.Account.Create().
				SetHandle("turn-owner").
				SetName("Turn Owner").
				Save(ctx)
			require.NoError(t, err)

			sessionID := robot.SessionID(xid.New())
			ownerID := account.AccountID(owner.ID)
			_, err = repo.Create(ctx, sessionID, "Queued turn", ownerID, nil)
			require.NoError(t, err)

			turnID := robot.TurnID(xid.New())
			require.NoError(t, repo.EnqueueTurn(ctx, robot_session.EnqueueTurnParams{
				ID:          turnID,
				SessionID:   sessionID,
				InitiatorID: opt.New(ownerID),
				SourceKind:  "user_message",
				RobotRef:    "denbot",
				InputData:   json.RawMessage(`{"text":"hello"}`),
			}))

			active, err := repo.ActiveTurn(ctx, sessionID)
			require.NoError(t, err)
			assert.Equal(t, turnID, active)

			lease, err := repo.AcquireQueuedTurnExecution(ctx, sessionID, xid.ID(turnID), false, time.Minute)
			require.NoError(t, err)
			_, err = repo.AcquireQueuedTurnExecution(ctx, sessionID, xid.ID(turnID), false, time.Minute)
			require.ErrorIs(t, err, robot_session.ErrTurnHandled)
			leaseCtx := robot_session.WithExecutionLease(ctx, *lease)

			require.NoError(t, repo.AppendMessage(leaseCtx, sessionID, opt.New(ownerID), opt.NewEmpty[robot.Actor](), &adksession.Event{
				InvocationID: "invocation-1",
				Author:       "user",
				LLMResponse:  model.LLMResponse{Content: &genai.Content{Role: "user", Parts: []*genai.Part{{Text: "hello"}}}},
			}))
			require.NoError(t, repo.AppendMessage(leaseCtx, sessionID, opt.NewEmpty[account.AccountID](), opt.New(robot.NewBuiltinActor("denbot")), &adksession.Event{
				InvocationID: "invocation-1",
				Author:       "denbot",
				LLMResponse:  model.LLMResponse{Content: &genai.Content{Role: "model", Parts: []*genai.Part{{Text: "hello back"}}}},
			}))

			require.NoError(t, repo.ReleaseExecution(ctx, *lease))
			require.NoError(t, repo.FinishTurn(ctx, sessionID, turnID, robot.TurnStatusCompleted, ""))

			events, next, closed, err := repo.ReadTurnEvents(ctx, sessionID, turnID, 0, 20)
			require.NoError(t, err)
			require.Len(t, events, 4)
			assert.Equal(t, []robot.SessionEventKind{
				robot.SessionEventTurnQueued,
				robot.SessionEventMessage,
				robot.SessionEventMessage,
				robot.SessionEventTurnCompleted,
			}, []robot.SessionEventKind{events[0].Kind, events[1].Kind, events[2].Kind, events[3].Kind})
			assert.Equal(t, events[3].Sequence, next)
			assert.True(t, closed)

			session, _, err := repo.Get(ctx, sessionID, robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 20))
			require.NoError(t, err)
			assert.Len(t, session.Messages, 2, "lifecycle events are not conversation messages")
		}))
	}))
}

func TestQueuedTurnClaimIsIdempotentAndRecoverable(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		repo *robot_session.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			owner, err := db.Account.Create().SetHandle("turn-recovery-owner").SetName("Turn Recovery Owner").Save(ctx)
			require.NoError(t, err)
			sessionID := robot.SessionID(xid.New())
			ownerID := account.AccountID(owner.ID)
			_, err = repo.Create(ctx, sessionID, "Recoverable turn", ownerID, nil)
			require.NoError(t, err)
			turnID := robot.TurnID(xid.New())
			require.NoError(t, repo.EnqueueTurn(ctx, robot_session.EnqueueTurnParams{
				ID: turnID, SessionID: sessionID, InitiatorID: opt.New(ownerID), SourceKind: "user_message", RobotRef: "denbot", InputData: json.RawMessage(`{}`),
			}))

			first, err := repo.AcquireQueuedTurnExecution(ctx, sessionID, xid.ID(turnID), false, -time.Second)
			require.NoError(t, err)
			recoverable, err := repo.RecoverableTurn(ctx, sessionID)
			require.NoError(t, err)
			assert.Equal(t, turnID, recoverable.ID)
			runnable, err := repo.RunnableSessionIDs(ctx, robot.SessionEventPageSize)
			require.NoError(t, err)
			assert.Contains(t, runnable, sessionID)
			recovered, err := repo.AcquireQueuedTurnExecution(ctx, sessionID, xid.ID(turnID), false, time.Minute)
			require.NoError(t, err)
			assert.Greater(t, recovered.Generation, first.Generation)
			require.NoError(t, repo.ReleaseExecution(ctx, *recovered))
			require.NoError(t, repo.FinishTurn(ctx, sessionID, turnID, robot.TurnStatusCompleted, ""))
		}))
	}))
}

func TestQueuedTurnCanBeCancelledBeforeExecution(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		repo *robot_session.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			owner, err := db.Account.Create().SetHandle("turn-cancellation-owner").SetName("Turn Cancellation Owner").Save(ctx)
			require.NoError(t, err)
			sessionID := robot.SessionID(xid.New())
			ownerID := account.AccountID(owner.ID)
			_, err = repo.Create(ctx, sessionID, "Cancelled queued turn", ownerID, nil)
			require.NoError(t, err)
			turnID := robot.TurnID(xid.New())
			require.NoError(t, repo.EnqueueTurn(ctx, robot_session.EnqueueTurnParams{
				ID: turnID, SessionID: sessionID, InitiatorID: opt.New(ownerID), SourceKind: "user_message", RobotRef: "denbot", InputData: json.RawMessage(`{}`),
			}))

			disposition, err := repo.RequestTurnCancellation(ctx, sessionID, turnID)
			require.NoError(t, err)
			assert.True(t, disposition.Finished)
			assert.False(t, disposition.SignalExecution)

			turn, err := repo.NextQueuedTurn(ctx, sessionID)
			assert.ErrorIs(t, err, robot_session.ErrTurnNotFound)
			assert.Nil(t, turn)

			events, _, closed, err := repo.ReadTurnEvents(ctx, sessionID, turnID, 0, 20)
			require.NoError(t, err)
			require.Len(t, events, 2)
			assert.Equal(t, robot.SessionEventTurnQueued, events[0].Kind)
			assert.Equal(t, robot.SessionEventTurnCancelled, events[1].Kind)
			assert.False(t, events[1].ErrorText.Ok(), "cancellation is not a turn error")
			assert.True(t, closed)

			stored, err := db.RobotSessionTurn.Get(ctx, xid.ID(turnID))
			require.NoError(t, err)
			assert.Nil(t, stored.ErrorText)
		}))
	}))
}
