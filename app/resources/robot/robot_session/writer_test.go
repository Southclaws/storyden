package robot_session_test

import (
	"context"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"google.golang.org/adk/v2/model"
	adksession "google.golang.org/adk/v2/session"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
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

func TestAppendMessageLinksMultipleAssetsWithoutLosingContentOrder(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		repo *robot_session.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			owner, err := db.Account.Create().
				SetHandle("message-assets-owner").
				SetName("Message Assets Owner").
				Save(ctx)
			require.NoError(t, err)

			firstID := asset.AssetID(xid.New())
			secondID := asset.AssetID(xid.New())
			for id, name := range map[asset.AssetID]string{firstID: "first.png", secondID: "second.webp"} {
				_, err := db.Asset.Create().
					SetID(id).
					SetAccountID(owner.ID).
					SetFilename(asset.NewExistingFilename(id, name).String()).
					SetMimeType("image/png").
					SetSize(1).
					Save(ctx)
				require.NoError(t, err)
			}

			sessionID := robot.SessionID(xid.New())
			_, err = repo.Create(ctx, sessionID, "Image references", account.AccountID(owner.ID), map[string]any{})
			require.NoError(t, err)

			event := &adksession.Event{
				InvocationID: "invocation-with-images",
				Author:       "user",
				LLMResponse: model.LLMResponse{
					Content: &genai.Content{Role: genai.RoleUser, Parts: []*genai.Part{
						{Text: "compare"},
						robot.NewImageAssetPart(firstID),
						robot.NewImageAssetPart(secondID),
						robot.NewImageAssetPart(firstID),
					}},
				},
			}
			err = repo.AppendMessage(ctx, sessionID, opt.New(account.AccountID(owner.ID)), opt.NewEmpty[robot.Actor](), event)
			require.NoError(t, err)

			stored, err := db.RobotSessionMessage.Query().WithAssets().Only(ctx)
			require.NoError(t, err)
			require.Len(t, stored.Edges.Assets, 2, "the relation is a deduplicated asset reference index")
			require.Len(t, stored.EventData.Content.Parts, 4, "the event remains the ordered model content source")
			assert.Equal(t, firstID.String(), stored.EventData.Content.Parts[1].PartMetadata[robot.ImageAssetIDMetadataKey])
			assert.Equal(t, secondID.String(), stored.EventData.Content.Parts[2].PartMetadata[robot.ImageAssetIDMetadataKey])
			assert.Equal(t, firstID.String(), stored.EventData.Content.Parts[3].PartMetadata[robot.ImageAssetIDMetadataKey])

			require.NoError(t, db.Asset.DeleteOneID(firstID).Exec(ctx))
			stored, err = db.RobotSessionMessage.Query().WithAssets().Only(ctx)
			require.NoError(t, err, "deleting an asset must preserve the message")
			require.Len(t, stored.Edges.Assets, 1, "deleting an asset removes only its join-row reference")
			assert.Equal(t, secondID, stored.Edges.Assets[0].ID)

			require.NoError(t, db.RobotSession.DeleteOneID(xid.ID(sessionID)).Exec(ctx))
			remaining, err := db.Asset.Get(ctx, secondID)
			require.NoError(t, err, "deleting a session must not delete its assets")
			count, err := remaining.QueryRobotMessages().Count(ctx)
			require.NoError(t, err)
			assert.Zero(t, count, "deleting a message removes only the join-row reference")
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
