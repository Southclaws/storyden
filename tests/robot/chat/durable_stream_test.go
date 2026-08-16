package chat_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests/robot"
)

func TestRobotDurableStreamOfficialClient(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{LanguageModelProvider: "mock"},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelDelayed),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				sessionID := xid.New().String()
				stream := doChat(t, root, ts, sh.WithSession(adminCtx), sessionID, "", "exercise durable live tail")

				assert.True(t, stream.usedLiveSSE, "official client must advance from catch-up to the SSE live tail")
				assert.Equal(t, "Durable live tail completed.", strings.Join(collectTextDeltas(stream), ""))

				snapshot, err := cl.RobotSessionGetWithResponse(root, openapi.RobotSessionIDParam(sessionID), &openapi.RobotSessionGetParams{}, sh.WithSession(adminCtx))
				assert.NoError(t, err)
				if !assert.NotNil(t, snapshot.JSON200) {
					return
				}

				readCtx, cancel := context.WithTimeout(root, 10*time.Second)
				defer cancel()
				readResult := make(chan *robot.DurableJSONRead[openapi.RobotSessionStreamEvent], 1)
				readError := make(chan error, 1)
				go func() {
					result, err := robot.ReadDurableJSONUntil(
						readCtx,
						ts.URL+"/api/robots/sessions/"+sessionID+"/stream",
						sh.WithSession(adminCtx),
						snapshot.JSON200.StreamOffset,
						func(event openapi.RobotSessionStreamEvent) bool {
							return event.EventKind == openapi.TurnCompleted
						},
					)
					if err != nil {
						readError <- err
						return
					}
					readResult <- result
				}()

				doChat(t, root, ts, sh.WithSession(adminCtx), sessionID, "", "exercise durable live tail")

				select {
				case err := <-readError:
					assert.NoError(t, err)
				case result := <-readResult:
					assert.True(t, result.UsedLiveSSE, "session reader must remain attached while idle")
					for _, event := range result.Items {
						assert.NotNil(t, event.Parts, "session event parts must always be an array")
					}
					assert.Condition(t, func() bool {
						for _, event := range result.Items {
							if event.Message != nil && event.Message.Role == openapi.RobotSessionMessageRoleUser {
								return true
							}
						}
						return false
					}, "session reader must receive the initiating member message")
				case <-readCtx.Done():
					assert.Fail(t, "session reader timed out")
				}
			}))
		}),
	)
}
