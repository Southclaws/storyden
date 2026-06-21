package chat_test

import (
	"context"
	"net/http"
	"net/http/httptest"
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

func TestRobotSessionGetSharedAcrossRobotUsers(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelAck),
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
				ownerCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				ownerSession := sh.WithSession(ownerCtx)

				attackerCtx, _ := e2e.WithAccount(root, aw, seed.Account_002_Frigg)
				attackerSession := sh.WithSession(attackerCtx)

				sessionID := xid.New().String()
				doChat(t, root, ts, ownerSession, sessionID, "", "my secret message")

				t.Run("owner_can_read_own_session", func(t *testing.T) {
					resp := tests.AssertRequest(cl.RobotSessionGetWithResponse(root,
						openapi.RobotSessionIDParam(sessionID),
						&openapi.RobotSessionGetParams{},
						ownerSession,
					))(t, http.StatusOK)
					require.NotNil(t, resp.JSON200)
				})

				t.Run("other_robot_user_can_read_session", func(t *testing.T) {
					resp := tests.AssertRequest(cl.RobotSessionGetWithResponse(root,
						openapi.RobotSessionIDParam(sessionID),
						&openapi.RobotSessionGetParams{},
						attackerSession,
					))(t, http.StatusOK)
					require.NotNil(t, resp.JSON200)
				})
			}))
		}),
	)
}

func TestRobotSessionsListNoFilterReturnsAllSessions(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelAck),
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
				ownerCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				ownerSession := sh.WithSession(ownerCtx)

				otherCtx, _ := e2e.WithAccount(root, aw, seed.Account_002_Frigg)
				otherSession := sh.WithSession(otherCtx)

				ownerSessionID := xid.New().String()
				doChat(t, root, ts, ownerSession, ownerSessionID, "", "odin's private thought")

				otherSessionID := xid.New().String()
				doChat(t, root, ts, otherSession, otherSessionID, "", "frigg's private thought")

				t.Run("list_without_filter_exposes_other_users_sessions", func(t *testing.T) {
					a := assert.New(t)
					resp := tests.AssertRequest(cl.RobotSessionsListWithResponse(root,
						&openapi.RobotSessionsListParams{},
						otherSession,
					))(t, http.StatusOK)
					require.NotNil(t, resp.JSON200)

					var foundOwnerSession bool
					for _, s := range resp.JSON200.Sessions {
						if string(s.Id) == ownerSessionID {
							foundOwnerSession = true
							break
						}
					}

					a.True(foundOwnerSession)
				})
			}))
		}),
	)
}

func TestSSEChatCanContinueSessionOwnedByAnotherRobotUser(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelAck),
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
				ownerCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				ownerSession := sh.WithSession(ownerCtx)

				attackerCtx, _ := e2e.WithAccount(root, aw, seed.Account_002_Frigg)
				attackerSession := sh.WithSession(attackerCtx)

				sharedSessionID := xid.New().String()

				doChat(t, root, ts, ownerSession, sharedSessionID, "", "initial turn")

				t.Run("other_robot_user_can_post_to_owners_session", func(t *testing.T) {
					a := assert.New(t)

					doChat(t, root, ts, attackerSession, sharedSessionID, "", "attacker message")

					limit := openapi.RobotSessionMessageLimitQuery("10")
					resp := tests.AssertRequest(cl.RobotSessionGetWithResponse(root,
						openapi.RobotSessionIDParam(sharedSessionID),
						&openapi.RobotSessionGetParams{Limit: &limit},
						attackerSession,
					))(t, http.StatusOK)
					require.NotNil(t, resp.JSON200)

					a.GreaterOrEqual(resp.JSON200.MessageList.Results, 2)
				})
			}))
		}),
	)
}
