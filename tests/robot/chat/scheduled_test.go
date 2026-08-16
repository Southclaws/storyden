package chat_test

import (
	"context"
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
	"github.com/Southclaws/storyden/tests/robot"
)

func TestRobotChecksBackAfterScheduledDelay(t *testing.T) {
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
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			settingsRepo *settings.SettingsRepository,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				auth := sh.WithSession(adminCtx)
				scriptName := "robot-chat-check-back-" + xid.New().String() + ".yaml"
				scriptPath := filepath.Join("..", "scripts", scriptName)
				writeScript(t, scriptPath, `steps:
  - match:
      contains: "scheduled check is now due"
    respond:
      text: "The scheduled check resumed in a new turn."
      finish: "stop"
  - match:
      tool_result: check_back_later
      tool_result_status: pending
    respond:
      text: "I scheduled the check and will return later."
      finish: "stop"
  - match:
      contains: "check this later"
    respond:
      tool_calls:
        - id: call_check_back_1
          name: check_back_later
          args:
            duration_seconds: 1
            task: "Check the current session state now."
`)
				defer os.Remove(scriptPath)
				require.NoError(t, robot.SetRobotSettings(root, settingsRepo, "mock/../scripts/"+scriptName))

				sessionID := xid.New().String()
				scheduledAt := time.Now()
				accepted := enqueueChatMessage(t, root, ts, auth, sessionID, "", "check this later")
				readCtx, cancel := context.WithTimeout(root, 20*time.Second)
				defer cancel()
				completedTurns := 0
				read, err := robot.ReadDurableJSONUntil(
					readCtx,
					ts.URL+accepted.location,
					auth,
					"-1",
					func(event openapi.RobotSessionStreamEvent) bool {
						if event.EventKind == openapi.TurnCompleted {
							completedTurns++
						}
						return completedTurns == 2
					},
				)
				require.NoError(t, err)
				assert.GreaterOrEqual(t, time.Since(scheduledAt), 900*time.Millisecond)

				stream := &fullResponse{}
				for _, event := range read.Items {
					stream.parts = append(stream.parts, event.Parts...)
				}
				text := strings.Join(collectTextDeltas(stream), "")
				assert.Contains(t, text, "I scheduled the check and will return later.")
				assert.Contains(t, text, "The scheduled check resumed in a new turn.")
			}))
		}),
	)
}
