package chat_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
	"github.com/Southclaws/storyden/tests/robot"
)

func TestRobotRunningTurnCanBeCancelled(t *testing.T) {
	t.Parallel()

	integration.Test(
		t,
		&config.Config{LanguageModelProvider: "mock"},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelAck),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			settingsRepo *settings.SettingsRepository,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				auth := sh.WithSession(adminCtx)
				scriptName := "robot-chat-cancellation-" + xid.New().String() + ".yaml"
				scriptPath := filepath.Join("..", "scripts", scriptName)
				writeScript(t, scriptPath, `steps:
  - match:
      contains: "cancel this turn"
    respond:
      delay_ms: 5000
      text: "This response should never be stored."
      finish: "stop"
  - match:
      any: true
    respond:
      text: "The next turn still runs."
      finish: "stop"
`)
				defer os.Remove(scriptPath)
				require.NoError(t, robot.SetRobotSettings(root, settingsRepo, "mock/../scripts/"+scriptName))

				sessionID := xid.New().String()
				accepted := enqueueChatMessage(t, root, ts, auth, sessionID, "", "cancel this turn")
				var turnID string
				_, err := robot.ReadDurableJSONUntil(
					root,
					ts.URL+accepted.location,
					auth,
					"-1",
					func(event openapi.RobotSessionStreamEvent) bool {
						if event.EventKind != openapi.TurnQueued || event.TurnId == nil {
							return false
						}
						turnID = string(*event.TurnId)
						return true
					},
				)
				require.NoError(t, err)
				require.NotEmpty(t, turnID)

				require.Eventually(t, func() bool {
					snapshot, err := cl.RobotSessionGetWithResponse(
						root,
						openapi.RobotSessionIDParam(sessionID),
						&openapi.RobotSessionGetParams{},
						auth,
					)
					return err == nil && snapshot.JSON200 != nil && snapshot.JSON200.ActiveTurnId != nil && string(*snapshot.JSON200.ActiveTurnId) == turnID
				}, time.Second, 10*time.Millisecond, "the session snapshot should identify its active turn")

				cancelled := tests.AssertRequest(cl.RobotSessionTurnCancelWithResponse(
					root,
					openapi.RobotSessionIDParam(sessionID),
					openapi.RobotTurnIDParam(turnID),
					auth,
				))(t, http.StatusAccepted)
				require.NotNil(t, cancelled)

				cancelCtx, cancel := context.WithTimeout(root, 3*time.Second)
				defer cancel()
				read, err := robot.ReadDurableJSONUntil(
					cancelCtx,
					ts.URL+accepted.location,
					auth,
					"-1",
					func(event openapi.RobotSessionStreamEvent) bool {
						return event.EventKind == openapi.TurnCancelled && event.TurnId != nil && string(*event.TurnId) == turnID
					},
				)
				require.NoError(t, err)
				assert.Condition(t, func() bool {
					for _, event := range read.Items {
						if event.EventKind == openapi.TurnCancelled && event.TurnId != nil && string(*event.TurnId) == turnID {
							for _, part := range event.Parts {
								assert.NotEqual(t, "error", part.Type, "cancellation should not be projected as an error")
							}
							return true
						}
					}
					return false
				}, "session stream should expose the cancelled turn")

				followUp := doChat(t, root, ts, auth, sessionID, "", "continue after cancellation")
				assert.Equal(t, "The next turn still runs.", strings.Join(collectTextDeltas(followUp), ""))
			}))
		}),
	)
}
