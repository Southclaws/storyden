package chat_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	robot_tests "github.com/Southclaws/storyden/tests/robot"
)

func TestRobotWebMCPToolBlocksReplaysAndResumes(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{LanguageModelProvider: "mock"},
		e2e.Setup(),
		robot_tests.WithRobotSettings(mockModelAck),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			sessionRepo *robot_session.Repository,
			settingsRepo *settings.SettingsRepository,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				auth := sh.WithSession(adminCtx)
				scriptName := "robot-chat-webmcp-" + xid.New().String() + ".yaml"
				scriptPath := filepath.Join("..", "scripts", scriptName)
				writeScript(t, scriptPath, `steps:
  - match:
      contains: "add a browser block"
    respond:
      tool_calls:
        - id: call_webmcp_add_1
          name: library_page_block_add
          args:
            block: assets
            after_block: directory
  - match:
      tool_result: library_page_block_add
    respond:
      text: "Browser layout result received."
      finish: "stop"
`)
				defer os.Remove(scriptPath)
				require.NoError(t, robot_tests.SetRobotSettings(root, settingsRepo, "mock/../scripts/"+scriptName))

				pageType := "library_page"
				toolTitle := "Add Library Page Block"
				readOnly := false
				sessionID := xid.New().String()
				clientID := "browser-" + xid.New().String()
				accepted := enqueueWebMCPChatMessage(t, root, ts, auth, openapi.RobotChatRequest{
					Id:        sessionID,
					SessionId: &sessionID,
					Context:   &openapi.RobotChatContext{PageType: &pageType},
					ClientTools: &openapi.RobotClientToolContext{
						ClientId: clientID,
						Tools: []openapi.RobotClientTool{{
							Name:        "library_page_block_add",
							Title:       &toolTitle,
							Description: "Add a hidden block to the current Library page.",
							InputSchema: map[string]any{
								"type":     "object",
								"required": []any{"block"},
								"properties": map[string]any{
									"block":       map[string]any{"type": "string"},
									"after_block": map[string]any{"type": "string"},
								},
							},
							Annotations: &openapi.RobotClientToolAnnotations{ReadOnlyHint: &readOnly},
						}},
					},
				}, "add a browser block")

				first := readDurableChatParts(t, root, ts, auth, accepted)
				inputs := collectToolInputs(first)
				require.Len(t, inputs, 1)
				assert.Equal(t, "library_page_block_add", string(inputs[0].ToolName))
				assert.Equal(t, map[string]any{"block": "assets", "after_block": "directory"}, inputs[0].Input)
				require.NotNil(t, inputs[0].Dynamic)
				assert.True(t, *inputs[0].Dynamic)
				require.NotNil(t, inputs[0].ProviderMetadata)
				assertWebMCPMetadata(t, *inputs[0].ProviderMetadata, clientID, pageType)
				assert.Empty(t, collectToolOutputs(first), "the pending server marker must not be projected as a completed tool result")

				snapshot, err := cl.RobotSessionGetWithResponse(root, openapi.RobotSessionIDParam(sessionID), &openapi.RobotSessionGetParams{}, auth)
				require.NoError(t, err)
				require.NotNil(t, snapshot.JSON200)
				require.NotNil(t, snapshot.JSON200.ActiveTurnId, "blocked client tool keeps the turn active for reconnection")
				persisted := findDynamicToolCall(t, snapshot.JSON200.MessageList.Messages, "call_webmcp_add_1")
				assert.Equal(t, "library_page_block_add", persisted.ToolName)
				require.NotNil(t, persisted.CallProviderMetadata)
				assertWebMCPMetadata(t, *persisted.CallProviderMetadata, clientID, pageType)

				parsedSessionID, err := robot.NewSessionID(sessionID)
				require.NoError(t, err)
				blockedSession, _, err := sessionRepo.Get(root, parsedSessionID, robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 50))
				require.NoError(t, err)
				assert.Equal(t, []string{"call_webmcp_add_1"}, robot.PendingClientToolIDs(blockedSession.State))

				assert.Equal(t, http.StatusConflict, doChatStatus(t, root, ts, auth, sessionID, "do not skip the pending browser tool"))

				second := doChatToolOutputs(t, root, ts, auth, sessionID, "", []map[string]any{{
					"type":       "dynamic-tool",
					"state":      "output-available",
					"toolCallId": "call_webmcp_add_1",
					"toolName":   "library_page_block_add",
					"input":      map[string]any{"block": "assets", "after_block": "directory"},
					"output": map[string]any{
						"message": "Added the assets block after the directory block.",
						"layout": map[string]any{"blocks": []any{
							map[string]any{"type": "title"},
							map[string]any{"type": "directory", "layout": "table"},
							map[string]any{"type": "assets", "layout": "grid", "grid_size": float64(3)},
						}},
					},
				}})

				assert.Empty(t, collectErrorParts(second))
				assert.Empty(t, collectToolOutputs(second), "the browser-supplied result is input to the resumed turn, not a new UI output")
				assert.Equal(t, "Browser layout result received.", strings.Join(collectTextDeltas(second), ""))

				resumedSession, _, err := sessionRepo.Get(root, parsedSessionID, robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 50))
				require.NoError(t, err)
				assert.Empty(t, robot.PendingClientToolIDs(resumedSession.State))
				_, hasActiveTurn := resumedSession.ActiveTurnID.Get()
				assert.False(t, hasActiveTurn, "resumed client tool turn must complete")

				matchingResponses := 0
				for _, message := range resumedSession.Messages {
					if message.Event.LLMResponse.Content == nil {
						continue
					}
					for _, part := range message.Event.LLMResponse.Content.Parts {
						if part == nil || part.FunctionResponse == nil || part.FunctionResponse.ID != "call_webmcp_add_1" {
							continue
						}
						matchingResponses++
						assert.Equal(t, "library_page_block_add", part.FunctionResponse.Name)
						assert.Equal(t, "Added the assets block after the directory block.", part.FunctionResponse.Response["message"])
						assert.False(t, isPendingClientToolResponse(part.FunctionResponse.Response))
						assert.Equal(t, "user", message.Event.Author)
					}
				}
				assert.Equal(t, 1, matchingResponses, "the durable history must contain exactly the real browser result")
			}))
		}),
	)
}

func enqueueWebMCPChatMessage(
	t *testing.T,
	ctx context.Context,
	ts *httptest.Server,
	session openapi.RequestEditorFn,
	request openapi.RobotChatRequest,
	message string,
) acceptedInput {
	t.Helper()

	var textPart openapi.UIMessagePart
	require.NoError(t, textPart.FromTextUIPart(openapi.TextUIPart{Type: openapi.TextUIPartTypeText, Text: message}))
	request.Messages = []openapi.UIMessage{{
		Id:    xid.New().String(),
		Role:  openapi.UIMessageRoleUser,
		Parts: []openapi.UIMessagePart{textPart},
	}}
	body, err := json.Marshal(request)
	require.NoError(t, err)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, ts.URL+"/api/robots/sessions", bytes.NewReader(body))
	require.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")
	require.NoError(t, session(ctx, httpReq))

	response, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, http.StatusAccepted, response.StatusCode)

	return decodeAcceptedInput(t, response)
}

func findDynamicToolCall(t *testing.T, messages []openapi.RobotSessionMessage, callID string) openapi.ToolUIPartInputAvailable {
	t.Helper()

	for _, message := range messages {
		for _, part := range message.Parts {
			if part.Type != "dynamic-tool" {
				continue
			}
			toolPart, err := part.AsToolUIPart()
			if err != nil {
				continue
			}
			input, err := toolPart.AsToolUIPartInputAvailable()
			if err == nil && input.ToolCallId == callID {
				return input
			}
		}
	}

	t.Fatalf("dynamic tool call %q not found in session snapshot", callID)
	return openapi.ToolUIPartInputAvailable{}
}

func assertWebMCPMetadata(t *testing.T, metadata openapi.ArbitraryData, clientID, pageType string) {
	t.Helper()

	metadataMap, ok := metadata.(map[string]any)
	require.True(t, ok)
	storyden, ok := metadataMap["storyden"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "webmcp", storyden["source"])
	assert.Equal(t, clientID, storyden["client_id"])
	scope, ok := storyden["scope"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, pageType, scope["page_type"])
}
