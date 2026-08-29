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
	robotservice "github.com/Southclaws/storyden/app/services/semdex/robot"
	"github.com/Southclaws/storyden/internal/ent"
	ent_robot_session_input "github.com/Southclaws/storyden/internal/ent/robotsessioninput"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestQueuedInputsRemainIndividualWhenMaterialisedAsOneTurn(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		repo *robot_session.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			owner, err := db.Account.Create().SetHandle("input-queue-owner").SetName("Input Queue Owner").Save(ctx)
			require.NoError(t, err)
			ownerID := account.AccountID(owner.ID)
			sessionID := robot.SessionID(xid.New())
			require.NoError(t, createSession(ctx, repo, sessionID, ownerID))

			firstID := robot.InputID(xid.New())
			secondID := robot.InputID(xid.New())
			require.NoError(t, enqueueVisibleInput(ctx, repo, sessionID, ownerID, firstID, "first queued message"))
			require.NoError(t, enqueueVisibleInput(ctx, repo, sessionID, ownerID, firstID, "duplicate retry must be ignored"))
			require.NoError(t, enqueueVisibleInput(ctx, repo, sessionID, ownerID, secondID, "second queued message"))

			queued, err := repo.QueuedInputs(ctx, sessionID, 20)
			require.NoError(t, err)
			require.Len(t, queued, 2)
			assert.Equal(t, []robot.InputID{firstID, secondID}, []robot.InputID{queued[0].ID, queued[1].ID})

			before, _, err := repo.Get(ctx, sessionID, robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 20))
			require.NoError(t, err)
			require.Len(t, before.Messages, 2)
			assert.True(t, before.Messages[0].Queued)
			assert.True(t, before.Messages[1].Queued)

			turnID := robot.TurnID(xid.New())
			require.NoError(t, repo.MaterialiseTurn(ctx, robot_session.MaterialiseTurnParams{
				ID: turnID, SessionID: sessionID, InputIDs: []robot.InputID{firstID, secondID},
				InitiatorID: ownerID, SourceKind: "interactive_chat", RobotRef: "denbot", InputData: json.RawMessage(`{"combined":true}`),
			}))

			claimed, err := db.RobotSessionInput.Query().
				Where(ent_robot_session_input.IDIn(xid.ID(firstID), xid.ID(secondID))).
				Order(ent.Asc(ent_robot_session_input.FieldSequence)).
				All(ctx)
			require.NoError(t, err)
			require.Len(t, claimed, 2)
			for _, input := range claimed {
				assert.Equal(t, ent_robot_session_input.StatusClaimed, input.Status)
				require.NotNil(t, input.TurnID)
				assert.Equal(t, xid.ID(turnID), *input.TurnID)
			}

			lease, err := repo.AcquireQueuedTurnExecution(ctx, sessionID, xid.ID(turnID), false, time.Minute)
			require.NoError(t, err)
			leaseCtx := robot_session.WithExecutionLease(ctx, *lease)
			require.NoError(t, repo.AppendMessage(leaseCtx, sessionID, opt.New(ownerID), opt.NewEmpty[robot.Actor](), textEvent("first queued message\n\nsecond queued message", "user")))
			require.NoError(t, repo.AppendMessage(leaseCtx, sessionID, opt.NewEmpty[account.AccountID](), opt.New(robot.NewBuiltinActor("denbot")), textEvent("one response", "model")))

			after, _, err := repo.Get(ctx, sessionID, robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 20))
			require.NoError(t, err)
			require.Len(t, after.Messages, 3, "the merged runtime input must not become a third user bubble")
			assert.Equal(t, "first queued message", after.Messages[0].Event.Content.Parts[0].Text)
			assert.Equal(t, "second queued message", after.Messages[1].Event.Content.Parts[0].Text)
			assert.Equal(t, "one response", after.Messages[2].Event.Content.Parts[0].Text)
			assert.False(t, after.Messages[0].Queued)
			assert.False(t, after.Messages[1].Queued)

			events, _, err := repo.ReadSessionEvents(ctx, sessionID, 0, 20)
			require.NoError(t, err)
			assert.Equal(t, []robot.SessionEventKind{
				robot.SessionEventInputQueued,
				robot.SessionEventInputQueued,
				robot.SessionEventTurnQueued,
				robot.SessionEventMessage,
				robot.SessionEventMessage,
			}, eventKinds(events))
			require.Len(t, events[2].InputIDs, 2)
			assert.Equal(t, []robot.InputID{firstID, secondID}, events[2].InputIDs)
		}))
	}))
}

func TestEnqueueInputPersistsImportedHistoryAndVisibleInputAtomically(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		repo *robot_session.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			owner, err := db.Account.Create().SetHandle("imported-history-owner").SetName("Imported History Owner").Save(ctx)
			require.NoError(t, err)
			ownerID := account.AccountID(owner.ID)
			sessionID := robot.SessionID(xid.New())
			require.NoError(t, createSession(ctx, repo, sessionID, ownerID))

			inputID := robot.InputID(xid.New())
			history := []*adksession.Event{
				{
					InvocationID: inputID.String(),
					Author:       "user",
					LLMResponse: model.LLMResponse{
						Content: &genai.Content{
							Role: genai.RoleUser,
							Parts: []*genai.Part{
								{
									Text: "first imported message",
									PartMetadata: map[string]any{
										robotservice.MessageSpeakerMetadataKey: "Alice",
									},
								},
							},
						},
					},
				},
				{
					InvocationID: inputID.String(),
					Author:       "imported_assistant",
					LLMResponse: model.LLMResponse{
						Content: genai.NewContentFromText("imported assistant reply", genai.RoleModel),
					},
				},
			}
			require.NoError(t, repo.EnqueueInput(ctx, robot_session.EnqueueInputParams{
				ID:            inputID,
				SessionID:     sessionID,
				AccountID:     ownerID,
				SourceKind:    "plugin_rpc",
				BatchKey:      inputID.String(),
				InputData:     json.RawMessage(`{"turn":1}`),
				HistoryEvents: history,
				VisibleEvent:  opt.New(textEvent("current message", "user")),
			}))

			sess, _, err := repo.Get(ctx, sessionID, robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 20))
			require.NoError(t, err)
			require.Len(t, sess.Messages, 3)
			assert.Equal(t, []string{"first imported message", "imported assistant reply", "current message"}, []string{
				sess.Messages[0].Event.Content.Parts[0].Text,
				sess.Messages[1].Event.Content.Parts[0].Text,
				sess.Messages[2].Event.Content.Parts[0].Text,
			})
			assert.Equal(t, []string{genai.RoleUser, genai.RoleModel, genai.RoleUser}, []string{
				sess.Messages[0].Event.Content.Role,
				sess.Messages[1].Event.Content.Role,
				sess.Messages[2].Event.Content.Role,
			})
			assert.Equal(t, "Alice", sess.Messages[0].Event.Content.Parts[0].PartMetadata[robotservice.MessageSpeakerMetadataKey])
			assert.False(t, sess.Messages[0].Queued)
			assert.False(t, sess.Messages[1].Queued)
			assert.True(t, sess.Messages[2].Queued)

			require.NoError(t, repo.EnqueueInput(ctx, robot_session.EnqueueInputParams{
				ID:            inputID,
				SessionID:     sessionID,
				AccountID:     ownerID,
				SourceKind:    "plugin_rpc",
				BatchKey:      inputID.String(),
				InputData:     json.RawMessage(`{"turn":2}`),
				HistoryEvents: []*adksession.Event{textEvent("duplicate history", "user")},
				VisibleEvent:  opt.New(textEvent("duplicate current message", "user")),
			}))

			afterRetry, _, err := repo.Get(ctx, sessionID, robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 20))
			require.NoError(t, err)
			require.Len(t, afterRetry.Messages, 3)
			assert.Equal(t, "current message", afterRetry.Messages[2].Event.Content.Parts[0].Text)
		}))
	}))
}

func TestQueuedInputsBecomeRunnableAtTheirScheduledTime(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		repo *robot_session.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			owner, err := db.Account.Create().SetHandle("scheduled-input-owner").SetName("Scheduled Input Owner").Save(ctx)
			require.NoError(t, err)
			ownerID := account.AccountID(owner.ID)
			sessionID := robot.SessionID(xid.New())
			require.NoError(t, createSession(ctx, repo, sessionID, ownerID))

			futureID := robot.InputID(xid.New())
			require.NoError(t, repo.EnqueueInput(ctx, robot_session.EnqueueInputParams{
				ID: futureID, SessionID: sessionID, AccountID: ownerID,
				SourceKind: "scheduled", BatchKey: "scheduled", InputData: json.RawMessage(`{}`),
				NotBefore: opt.New(time.Now().Add(time.Hour)),
			}))

			queued, err := repo.QueuedInputs(ctx, sessionID, 20)
			require.NoError(t, err)
			assert.Empty(t, queued)
			runnable, err := repo.RunnableSessionIDs(ctx, 20)
			require.NoError(t, err)
			assert.NotContains(t, runnable, sessionID)

			dueID := robot.InputID(xid.New())
			require.NoError(t, repo.EnqueueInput(ctx, robot_session.EnqueueInputParams{
				ID: dueID, SessionID: sessionID, AccountID: ownerID,
				SourceKind: "scheduled", BatchKey: "scheduled", InputData: json.RawMessage(`{}`),
				NotBefore: opt.New(time.Now().Add(-time.Second)),
			}))

			queued, err = repo.QueuedInputs(ctx, sessionID, 20)
			require.NoError(t, err)
			require.Len(t, queued, 1)
			assert.Equal(t, dueID, queued[0].ID)
			runnable, err = repo.RunnableSessionIDs(ctx, 20)
			require.NoError(t, err)
			assert.Contains(t, runnable, sessionID)
		}))
	}))
}

func createSession(ctx context.Context, repo *robot_session.Repository, sessionID robot.SessionID, ownerID account.AccountID) error {
	_, err := repo.Create(ctx, sessionID, "Queued inputs", ownerID, nil)
	return err
}

func enqueueVisibleInput(
	ctx context.Context,
	repo *robot_session.Repository,
	sessionID robot.SessionID,
	ownerID account.AccountID,
	inputID robot.InputID,
	text string,
) error {
	return repo.EnqueueInput(ctx, robot_session.EnqueueInputParams{
		ID:           inputID,
		SessionID:    sessionID,
		AccountID:    ownerID,
		SourceKind:   "interactive_chat",
		BatchKey:     "same-account-and-runtime",
		InputData:    json.RawMessage(`{}`),
		VisibleEvent: opt.New(textEvent(text, "user")),
	})
}

func textEvent(text, role string) *adksession.Event {
	return &adksession.Event{
		InvocationID: xid.New().String(),
		Author:       role,
		LLMResponse: model.LLMResponse{
			Content: &genai.Content{
				Role: role,
				Parts: []*genai.Part{
					{
						Text: text,
					},
				},
			},
		},
	}
}

func eventKinds(events []robot.SessionEvent) []robot.SessionEventKind {
	result := make([]robot.SessionEventKind, len(events))
	for i, event := range events {
		result[i] = event.Kind
	}
	return result
}
