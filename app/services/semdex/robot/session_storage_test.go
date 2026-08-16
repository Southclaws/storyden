package robot_test

import (
	"context"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"google.golang.org/adk/v2/model"
	adksession "google.golang.org/adk/v2/session"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/robot/robot_querier"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	robotservice "github.com/Southclaws/storyden/app/services/semdex/robot"
	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry/denbot"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestSessionStoragePersistsStateDeltaAcrossRequests(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		robotQuerier *robot_querier.Querier,
		sessionRepo *robot_session.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			owner, err := db.Account.Create().
				SetHandle("state-delta-owner").
				SetName("State Delta Owner").
				Save(ctx)
			require.NoError(t, err)

			storage := robotservice.NewSessionStorage(robotQuerier, sessionRepo)
			sessionID := xid.New().String()
			created, err := storage.Create(ctx, &adksession.CreateRequest{
				AppName:   "storyden",
				UserID:    owner.ID.String(),
				SessionID: sessionID,
				State:     map[string]any{"existing": "retained"},
			})
			require.NoError(t, err)

			err = storage.AppendEvent(ctx, created.Session, &adksession.Event{
				InvocationID: "invocation-load-tools",
				Author:       denbot.AgentName,
				Actions: adksession.EventActions{StateDelta: map[string]any{
					"active_tools": []string{"thread_get"},
					"temp:scratch": "discarded",
				}},
			})
			require.NoError(t, err)

			reloaded, err := storage.Get(ctx, &adksession.GetRequest{
				AppName:   "storyden",
				UserID:    owner.ID.String(),
				SessionID: sessionID,
			})
			require.NoError(t, err)

			activeTools, err := reloaded.Session.State().Get("active_tools")
			require.NoError(t, err)
			assert.Equal(t, []any{"thread_get"}, activeTools)

			existing, err := reloaded.Session.State().Get("existing")
			require.NoError(t, err)
			assert.Equal(t, "retained", existing)

			_, err = reloaded.Session.State().Get("temp:scratch")
			assert.ErrorIs(t, err, adksession.ErrStateKeyNotExist)
		}))
	}))
}

func TestSessionStorageExcludesDelegatedEventsFromRootAgentHistory(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		robotQuerier *robot_querier.Querier,
		sessionRepo *robot_session.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			owner, err := db.Account.Create().
				SetHandle("delegation-history-owner").
				SetName("Delegation History Owner").
				Save(ctx)
			require.NoError(t, err)

			storage := robotservice.NewSessionStorage(robotQuerier, sessionRepo)
			sessionID := xid.New().String()
			created, err := storage.Create(ctx, &adksession.CreateRequest{
				AppName:   "storyden",
				UserID:    owner.ID.String(),
				SessionID: sessionID,
			})
			require.NoError(t, err)

			require.NoError(t, storage.AppendEvent(ctx, created.Session, &adksession.Event{
				InvocationID: "root-invocation",
				Author:       denbot.AgentName,
				LLMResponse: model.LLMResponse{Content: genai.NewContentFromText(
					"root response",
					genai.RoleModel,
				)},
			}))
			require.NoError(t, storage.AppendEvent(ctx, created.Session, &adksession.Event{
				InvocationID:   "delegated-invocation",
				Author:         denbot.AgentName,
				Branch:         "robot_specialist",
				IsolationScope: "call_delegate_1",
				LLMResponse: model.LLMResponse{Content: genai.NewContentFromText(
					"private specialist response",
					genai.RoleModel,
				)},
			}))

			reloaded, err := storage.Get(ctx, &adksession.GetRequest{
				AppName:   "storyden",
				UserID:    owner.ID.String(),
				SessionID: sessionID,
			})
			require.NoError(t, err)
			require.Equal(t, 1, reloaded.Session.Events().Len())
			assert.Equal(t, "root response", reloaded.Session.Events().At(0).Content.Parts[0].Text)
			assert.Empty(t, reloaded.Session.Events().At(0).IsolationScope)
		}))
	}))
}
