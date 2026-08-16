package chat_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	ent_robot_session "github.com/Southclaws/storyden/internal/ent/robotsession"
	ent_robot_session_input "github.com/Southclaws/storyden/internal/ent/robotsessioninput"
	ent_robot_session_turn "github.com/Southclaws/storyden/internal/ent/robotsessionturn"
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

func TestRobotQueuedMessagesShareOneNextTurn(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{LanguageModelProvider: "mock"},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelQueue),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			db *ent.Client,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				auth := sh.WithSession(adminCtx)
				sessionID := xid.New().String()
				sessionXID, err := xid.FromString(sessionID)
				require.NoError(t, err)

				first := enqueueChatMessage(t, root, ts, auth, sessionID, "", "first message starts the active turn")
				require.Eventually(t, func() bool {
					session, err := db.RobotSession.Query().Where(ent_robot_session.IDEQ(sessionXID)).Only(root)
					return err == nil && session.ExecutionStatus == ent_robot_session.ExecutionStatusRunning
				}, 5*time.Second, 25*time.Millisecond)

				second := enqueueChatMessage(t, root, ts, auth, sessionID, "", "second message waits in the queue")
				third := enqueueChatMessage(t, root, ts, auth, sessionID, "", "third message joins the same queue batch")
				require.Eventually(t, func() bool {
					count, err := db.RobotSessionInput.Query().Where(
						ent_robot_session_input.SessionIDEQ(sessionXID),
						ent_robot_session_input.StatusEQ(ent_robot_session_input.StatusQueued),
					).Count(root)
					return err == nil && count == 2
				}, 5*time.Second, 25*time.Millisecond)

				completed := 0
				read, err := robot.ReadDurableJSONUntil(
					root,
					ts.URL+first.location,
					auth,
					"-1",
					func(event openapi.RobotSessionStreamEvent) bool {
						if event.EventKind == openapi.TurnCompleted {
							completed++
						}
						return completed == 2
					},
				)
				require.NoError(t, err)
				assert.True(t, read.UsedLiveSSE)

				turns, err := db.RobotSessionTurn.Query().
					Where(ent_robot_session_turn.SessionIDEQ(sessionXID)).
					WithInputs().
					Order(ent.Asc(ent_robot_session_turn.FieldCreatedAt)).
					All(root)
				require.NoError(t, err)
				require.Len(t, turns, 2)
				assert.Len(t, turns[0].Edges.Inputs, 1)
				assert.Len(t, turns[1].Edges.Inputs, 2)

				snapshot, err := cl.RobotSessionGetWithResponse(root, openapi.RobotSessionIDParam(sessionID), &openapi.RobotSessionGetParams{}, auth)
				require.NoError(t, err)
				require.NotNil(t, snapshot.JSON200)
				userMessageIDs := make([]string, 0, 3)
				assistantMessages := 0
				for _, message := range snapshot.JSON200.MessageList.Messages {
					switch message.Role {
					case openapi.RobotSessionMessageRoleUser:
						userMessageIDs = append(userMessageIDs, string(message.Id))
					case openapi.RobotSessionMessageRoleAssistant:
						assistantMessages++
					}
				}
				assert.ElementsMatch(t, []string{
					string(first.reference.MessageId),
					string(second.reference.MessageId),
					string(third.reference.MessageId),
				}, userMessageIDs, "queued inputs remain three visible messages")
				assert.Equal(t, 2, assistantMessages, "the two queued inputs produce one shared next response")
			}))
		}),
	)
}
