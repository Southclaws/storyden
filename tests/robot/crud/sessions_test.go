package crud_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/app/transports/sse"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
	"github.com/Southclaws/storyden/tests/robot"
)

// The mock model path resolves relative to the test package directory:
// tests/robot/crud/../scripts/robot-chat-simple.yaml
const mockModel = "mock/../scripts/robot-chat-simple.yaml"

// startSession sends a single chat message via the SSE endpoint, drains the
// stream, and returns the session ID. It creates a session as a side-effect
// without asserting anything about the LLM response content.
func startSession(t *testing.T, ctx context.Context, ts *httptest.Server, session openapi.RequestEditorFn, robotID string) string {
	t.Helper()

	sessionID := xid.New().String()

	var textPart openapi.UIMessagePart
	require.NoError(t, textPart.FromTextUIPart(openapi.TextUIPart{Type: openapi.Text, Text: "hello"}))

	var robotIDPtr *string
	if robotID != "" {
		robotIDPtr = &robotID
	}

	body, err := json.Marshal(openapi.RobotChatRequest{
		Id:        sessionID,
		SessionId: &sessionID,
		RobotId:   robotIDPtr,
		Messages: []openapi.UIMessage{{
			Id:    xid.New().String(),
			Role:  openapi.UIMessageRoleUser,
			Parts: []openapi.UIMessagePart{textPart},
		}},
	})
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ts.URL+"/sse/chat", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, session(ctx, req))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Drain stream so the session is fully persisted before we query it.
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		if strings.TrimPrefix(scanner.Text(), "data: ") == "[DONE]" {
			break
		}
	}
	require.NoError(t, scanner.Err())

	return sessionID
}

func TestRobotChatSSERequiresAuthWhenRobotsDisabled(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		sse.Build(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
		) {
			lc.Append(fx.StartHook(func() {
				req, err := http.NewRequestWithContext(root, http.MethodPost, ts.URL+"/sse/chat", strings.NewReader(`{}`))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")

				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
			}))
		}),
	)
}

// TestRobotSessionsVisibility verifies the shared-ownership model described in
// the OpenAPI spec:
//
//	"sessions, messages and usage is not considered hidden to other accounts
//	with the usage permission. Robots are intended as administrative or
//	moderation tools to be shared among the team rather than private assistants."
func TestRobotSessionsVisibility(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		robot.WithRobotSettings(mockModel),
		sse.Build(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			roles *role_repo.Repository,
			assignments *role_assign.Assignment,
		) {
			lc.Append(fx.StartHook(func() {
				// User A: session owner.
				userACtx, userAAcc := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				userASession := sh.WithSession(userACtx)
				grantRobotPerms(t, root, roles, assignments, userAAcc.ID, rbac.PermissionUseRobots, rbac.PermissionManageRobots)

				// User B: different account, USE_ROBOTS only.
				userBCtx, userBAcc := e2e.WithAccount(root, aw, seed.Account_005_Þórr)
				userBSession := sh.WithSession(userBCtx)
				grantRobotPerms(t, root, roles, assignments, userBAcc.ID, rbac.PermissionUseRobots)
				_ = userBAcc

				// No-perm user.
				nopermCtx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
				nopermSession := sh.WithSession(nopermCtx)

				// Create a robot for User A to chat with.
				rb := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "visibility-test-robot-" + uuid.NewString(),
						Description: "Robot for session visibility tests",
						Playbook:    "Answer everything with 'hello from mock'.",
						Model:       robotModel(mockModel),
					},
					userASession,
				))(t, http.StatusOK)

				// User A starts a session via the chat endpoint — the real API path.
				sessionID := startSession(t, root, ts, userASession, string(rb.JSON200.Id))
				sessionIDParam := openapi.RobotSessionIDParam(sessionID)
				userAIDParam := openapi.AccountIDQueryParam(userAAcc.ID.String())

				t.Run("owner_can_list_own_sessions", func(t *testing.T) {
					a := assert.New(t)

					list := tests.AssertRequest(cl.RobotSessionsListWithResponse(root,
						&openapi.RobotSessionsListParams{},
						userASession,
					))(t, http.StatusOK)

					found := false
					for _, s := range list.JSON200.Sessions {
						if s.Id == sessionIDParam {
							found = true
							a.Equal(userAAcc.ID.String(), s.CreatedBy.Id)
						}
					}
					a.True(found, "owner must see their own session in the default list")
				})

				t.Run("user_b_can_list_user_a_sessions_via_account_filter", func(t *testing.T) {
					a := assert.New(t)

					list := tests.AssertRequest(cl.RobotSessionsListWithResponse(root,
						&openapi.RobotSessionsListParams{AccountId: &userAIDParam},
						userBSession,
					))(t, http.StatusOK)

					found := false
					for _, s := range list.JSON200.Sessions {
						if s.Id == sessionIDParam {
							found = true
							a.Equal(userAAcc.ID.String(), s.CreatedBy.Id)
						}
					}
					a.True(found, "any USE_ROBOTS user must be able to list another account's sessions")
				})

				t.Run("user_b_default_list_includes_user_a_sessions", func(t *testing.T) {
					a := assert.New(t)

					list := tests.AssertRequest(cl.RobotSessionsListWithResponse(root,
						&openapi.RobotSessionsListParams{},
						userBSession,
					))(t, http.StatusOK)

					found := false
					for _, s := range list.JSON200.Sessions {
						if s.Id == sessionIDParam {
							found = true
							a.Equal(userAAcc.ID.String(), s.CreatedBy.Id)
						}
					}
					a.True(found, "unfiltered Robot session list should include team sessions")
				})

				// Desired behaviour per spec: any member with USE_ROBOTS can retrieve
				// any session by ID regardless of authorship.
				t.Run("user_b_can_get_user_a_session_by_id", func(t *testing.T) {
					a := assert.New(t)

					get := tests.AssertRequest(cl.RobotSessionGetWithResponse(root,
						sessionIDParam,
						&openapi.RobotSessionGetParams{},
						userBSession,
					))(t, http.StatusOK)

					a.Equal(string(sessionIDParam), string(get.JSON200.Id))
					a.Equal(userAAcc.ID.String(), get.JSON200.CreatedBy.Id)
				})

				t.Run("owner_can_get_own_session_by_id", func(t *testing.T) {
					a := assert.New(t)

					get := tests.AssertRequest(cl.RobotSessionGetWithResponse(root,
						sessionIDParam,
						&openapi.RobotSessionGetParams{},
						userASession,
					))(t, http.StatusOK)

					a.Equal(string(sessionIDParam), string(get.JSON200.Id))
					a.Equal(userAAcc.ID.String(), get.JSON200.CreatedBy.Id)
				})

				t.Run("noperm_cannot_list_sessions", func(t *testing.T) {
					r := require.New(t)

					resp, err := cl.RobotSessionsListWithResponse(root,
						&openapi.RobotSessionsListParams{},
						nopermSession,
					)
					r.NoError(err)
					r.Equal(http.StatusForbidden, resp.StatusCode())
				})

				t.Run("noperm_cannot_get_session_by_id", func(t *testing.T) {
					r := require.New(t)

					resp, err := cl.RobotSessionGetWithResponse(root,
						sessionIDParam,
						&openapi.RobotSessionGetParams{},
						nopermSession,
					)
					r.NoError(err)
					r.Equal(http.StatusForbidden, resp.StatusCode())
				})
			}))
		}),
	)
}
