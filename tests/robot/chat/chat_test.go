package chat_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/app/transports/sse"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
	"github.com/Southclaws/storyden/tests/robot"
)

func TestRobotChat(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelSimple),
		sse.Build(),
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
				adminSession := sh.WithSession(adminCtx)

				t.Run("simple_text_response", func(t *testing.T) {
					a := assert.New(t)

					sessionID := xid.New().String()
					stream := doChat(t, root, ts, adminSession, sessionID, "", "hello")
					textDeltas := collectTextDeltas(stream)

					a.NotEmpty(textDeltas)
					a.Equal("hello from mock", strings.Join(textDeltas, ""))
				})

				t.Run("session_messages_cursor_pagination", func(t *testing.T) {
					sessionID := xid.New().String()
					doChat(t, root, ts, adminSession, sessionID, "", "first")
					doChat(t, root, ts, adminSession, sessionID, "", "second")
					doChat(t, root, ts, adminSession, sessionID, "", "third")

					all := tests.AssertRequest(cl.RobotSessionGetWithResponse(root,
						openapi.RobotSessionIDParam(sessionID),
						&openapi.RobotSessionGetParams{},
						adminSession,
					))(t, 200)
					require.NotNil(t, all.JSON200)
					allMessages := all.JSON200.MessageList.Messages
					require.GreaterOrEqual(t, len(allMessages), 6)

					limit := openapi.RobotSessionMessageLimitQuery("2")
					latest := tests.AssertRequest(cl.RobotSessionGetWithResponse(root,
						openapi.RobotSessionIDParam(sessionID),
						&openapi.RobotSessionGetParams{Limit: &limit},
						adminSession,
					))(t, 200)
					require.NotNil(t, latest.JSON200)
					latestMessages := latest.JSON200.MessageList.Messages
					require.Len(t, latestMessages, 2)
					assert.Equal(t, allMessages[len(allMessages)-2].Id, latestMessages[0].Id)
					assert.Equal(t, allMessages[len(allMessages)-1].Id, latestMessages[1].Id)
					require.NotNil(t, latest.JSON200.MessageList.NextBefore)
					assert.Equal(t, latestMessages[0].Id, string(*latest.JSON200.MessageList.NextBefore))

					before := openapi.RobotSessionMessageBeforeQuery(*latest.JSON200.MessageList.NextBefore)
					older := tests.AssertRequest(cl.RobotSessionGetWithResponse(root,
						openapi.RobotSessionIDParam(sessionID),
						&openapi.RobotSessionGetParams{Before: &before, Limit: &limit},
						adminSession,
					))(t, 200)
					require.NotNil(t, older.JSON200)
					olderMessages := older.JSON200.MessageList.Messages
					require.Len(t, olderMessages, 2)
					assert.Equal(t, allMessages[len(allMessages)-4].Id, olderMessages[0].Id)
					assert.Equal(t, allMessages[len(allMessages)-3].Id, olderMessages[1].Id)
				})
			}))
		}),
	)
}
