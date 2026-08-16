package chat_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"

	robottest "github.com/Southclaws/storyden/tests/robot"
)

const (
	mockModelSimple          = "mock/../scripts/robot-chat-simple.yaml"
	mockModelLibraryPageList = "mock/../scripts/robot-chat-library-page-list.yaml"
	mockModelLibraryTool     = "mock/../scripts/robot-chat-library-tool.yaml"
	mockModelToolError       = "mock/../scripts/robot-chat-tool-error.yaml"
	mockModelLLMError        = "mock/../scripts/robot-chat-llm-error.yaml"
	mockModelAck             = "mock/../scripts/robot-chat-ack.yaml"
	mockModelDelayed         = "mock/../scripts/robot-chat-delayed.yaml"

	mockModelLibrarySearchPages = "mock/../scripts/robot-chat-library-search-pages.yaml"
	mockModelContentSearch      = "mock/../scripts/robot-chat-content-search.yaml"
	mockModelThreadSearch       = "mock/../scripts/robot-chat-thread-search.yaml"
	mockModelReplySearch        = "mock/../scripts/robot-chat-reply-search.yaml"
	mockModelPostSearch         = "mock/../scripts/robot-chat-post-search.yaml"
	mockModelMemberSearch       = "mock/../scripts/robot-chat-member-search.yaml"
)

func robotToolsetsPtr(ids ...string) *openapi.RobotToolsetRefList {
	toolsets := openapi.RobotToolsetRefList(ids)
	return &toolsets
}

func robotToolsPtr(names ...string) *openapi.RobotToolNameList {
	tools := openapi.RobotToolNameList(names)
	return &tools
}

func robotIDPtr(robotID string) *string {
	robotID = strings.TrimSpace(robotID)
	if robotID == "" {
		return nil
	}
	return &robotID
}

type fullResponse struct {
	parts       []openapi.StreamPart
	usedLiveSSE bool
}

type delegationStreamData struct {
	StreamID string                 `json:"-"`
	CallID   string                 `json:"callId"`
	Robot    openapi.RobotReference `json:"robot"`
	Request  string                 `json:"request"`
	Status   string                 `json:"status"`
	Messages []openapi.UIMessage    `json:"messages"`
	Error    string                 `json:"error,omitempty"`
}

func doChat(
	t *testing.T,
	ctx context.Context,
	ts *httptest.Server,
	session openapi.RequestEditorFn,
	sessionID, robotID, message string,
) *fullResponse {
	t.Helper()

	var textPart openapi.UIMessagePart
	require.NoError(t, textPart.FromTextUIPart(openapi.TextUIPart{Type: openapi.TextUIPartTypeText, Text: message}))

	body, err := json.Marshal(openapi.RobotChatRequest{
		Id:        sessionID,
		SessionId: &sessionID,
		RobotId:   robotIDPtr(robotID),
		Messages: []openapi.UIMessage{{
			Id:    xid.New().String(),
			Role:  openapi.UIMessageRoleUser,
			Parts: []openapi.UIMessagePart{textPart},
		}},
	})
	require.NoError(t, err)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, ts.URL+"/sse/chat", bytes.NewReader(body))
	require.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")
	require.NoError(t, session(ctx, httpReq))

	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)
	return readDurableChatParts(t, ctx, ts, session, resp)
}

func doChatToolOutput(
	t *testing.T,
	ctx context.Context,
	ts *httptest.Server,
	session openapi.RequestEditorFn,
	sessionID, robotID, toolName, toolCallID string,
	input, output map[string]any,
) *fullResponse {
	t.Helper()

	return doChatToolOutputs(t, ctx, ts, session, sessionID, robotID, []map[string]any{{
		"type":       "tool-" + toolName,
		"state":      "output-available",
		"toolCallId": toolCallID,
		"toolName":   toolName,
		"input":      input,
		"output":     output,
	}})
}

func doChatToolOutputs(
	t *testing.T,
	ctx context.Context,
	ts *httptest.Server,
	session openapi.RequestEditorFn,
	sessionID, robotID string,
	parts []map[string]any,
) *fullResponse {
	t.Helper()

	requestBody := map[string]any{
		"id":        sessionID,
		"sessionId": sessionID,
		"messages": []map[string]any{{
			"id":    xid.New().String(),
			"role":  "assistant",
			"parts": parts,
		}},
	}
	if robotID = strings.TrimSpace(robotID); robotID != "" {
		requestBody["robotId"] = robotID
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, ts.URL+"/sse/chat", bytes.NewReader(body))
	require.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")
	require.NoError(t, session(ctx, httpReq))

	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)
	return readDurableChatParts(t, ctx, ts, session, resp)
}

func readDurableChatParts(
	t *testing.T,
	ctx context.Context,
	ts *httptest.Server,
	session openapi.RequestEditorFn,
	created *http.Response,
) *fullResponse {
	t.Helper()

	location := created.Header.Get("Location")
	require.NotEmpty(t, location)
	read, err := robottest.ReadDurableJSON[openapi.StreamPart](ctx, ts.URL+location, session)
	require.NoError(t, err)
	return &fullResponse{parts: read.Items, usedLiveSSE: read.UsedLiveSSE}
}

func collectToolInputs(ev *fullResponse) []openapi.ToolInputAvailablePart {
	var inputs []openapi.ToolInputAvailablePart
	for _, part := range ev.parts {
		if part.Type != "tool-input-available" {
			continue
		}
		p, err := part.AsToolInputAvailablePart()
		if err == nil {
			inputs = append(inputs, p)
		}
	}
	return inputs
}

func collectSessionToolInputs(messages []openapi.RobotSessionMessage) []openapi.ToolUIPartInputAvailable {
	var inputs []openapi.ToolUIPartInputAvailable
	for _, message := range messages {
		for _, part := range message.Parts {
			if !strings.HasPrefix(string(part.Type), "tool-") {
				continue
			}
			toolPart, err := part.AsToolUIPart()
			if err != nil {
				continue
			}
			input, err := toolPart.AsToolUIPartInputAvailable()
			if err == nil {
				inputs = append(inputs, input)
			}
		}
	}
	return inputs
}

func doChatStatus(
	t *testing.T,
	ctx context.Context,
	ts *httptest.Server,
	session openapi.RequestEditorFn,
	sessionID, message string,
) int {
	return doChatWithRobotStatus(t, ctx, ts, session, sessionID, "", message)
}

func doChatWithRobotStatus(
	t *testing.T,
	ctx context.Context,
	ts *httptest.Server,
	session openapi.RequestEditorFn,
	sessionID, robotID, message string,
) int {
	t.Helper()

	var textPart openapi.UIMessagePart
	require.NoError(t, textPart.FromTextUIPart(openapi.TextUIPart{Type: openapi.TextUIPartTypeText, Text: message}))

	body, err := json.Marshal(openapi.RobotChatRequest{
		Id:        sessionID,
		SessionId: &sessionID,
		RobotId:   robotIDPtr(robotID),
		Messages: []openapi.UIMessage{{
			Id:    xid.New().String(),
			Role:  openapi.UIMessageRoleUser,
			Parts: []openapi.UIMessagePart{textPart},
		}},
	})
	require.NoError(t, err)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, ts.URL+"/sse/chat", bytes.NewReader(body))
	require.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")
	require.NoError(t, session(ctx, httpReq))

	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp.StatusCode
}

func doChatToolOutputsStatus(
	t *testing.T,
	ctx context.Context,
	ts *httptest.Server,
	session openapi.RequestEditorFn,
	sessionID, robotID string,
	parts []map[string]any,
) int {
	t.Helper()

	requestBody := map[string]any{
		"id":        sessionID,
		"sessionId": sessionID,
		"messages": []map[string]any{{
			"id":    xid.New().String(),
			"role":  "assistant",
			"parts": parts,
		}},
	}
	if robotID = strings.TrimSpace(robotID); robotID != "" {
		requestBody["robotId"] = robotID
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, ts.URL+"/sse/chat", bytes.NewReader(body))
	require.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")
	require.NoError(t, session(ctx, httpReq))

	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp.StatusCode
}

func collectTextDeltas(ev *fullResponse) []string {
	var deltas []string
	for _, part := range ev.parts {
		if part.Type != "text-delta" {
			continue
		}
		p, err := part.AsTextDeltaPart()
		if err == nil {
			deltas = append(deltas, p.Delta)
		}
	}
	return deltas
}

func collectReasoningDeltas(ev *fullResponse) []string {
	var deltas []string
	for _, part := range ev.parts {
		if part.Type != "reasoning-delta" {
			continue
		}
		reasoning, err := part.AsReasoningDeltaPart()
		if err == nil {
			deltas = append(deltas, reasoning.Delta)
		}
	}
	return deltas
}

func collectToolCalls(ev *fullResponse) []string {
	var names []string
	for _, part := range ev.parts {
		if part.Type != "tool-input-available" {
			continue
		}
		p, err := part.AsToolInputAvailablePart()
		if err == nil {
			names = append(names, p.ToolName)
		}
	}
	return names
}

func collectToolOutputs(ev *fullResponse) []openapi.ToolOutputAvailablePart {
	var outputs []openapi.ToolOutputAvailablePart
	for _, part := range ev.parts {
		if part.Type != "tool-output-available" {
			continue
		}
		p, err := part.AsToolOutputAvailablePart()
		if err == nil {
			outputs = append(outputs, p)
		}
	}
	return outputs
}

func collectErrorParts(ev *fullResponse) []string {
	var errs []string
	for _, part := range ev.parts {
		if part.Type != "error" {
			continue
		}
		p, err := part.AsErrorPart()
		if err == nil {
			errs = append(errs, p.ErrorText)
		}
	}
	return errs
}

func collectDelegations(ev *fullResponse) []delegationStreamData {
	var delegations []delegationStreamData
	for _, part := range ev.parts {
		if part.Type != "data-delegation" {
			continue
		}
		dataPart, err := part.AsDataPart()
		if err != nil {
			continue
		}
		encoded, err := json.Marshal(dataPart.Data)
		if err != nil {
			continue
		}
		var delegation delegationStreamData
		if err := json.Unmarshal(encoded, &delegation); err == nil {
			if dataPart.Id != nil {
				delegation.StreamID = *dataPart.Id
			}
			delegations = append(delegations, delegation)
		}
	}
	return delegations
}

func writeScript(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
}

func collectPartsOfType(ev *fullResponse, partType string) []openapi.StreamPart {
	var result []openapi.StreamPart
	for _, part := range ev.parts {
		if part.Type == partType {
			result = append(result, part)
		}
	}
	return result
}

func toolOutputResultCount(output openapi.ToolOutputAvailablePart) float64 {
	m, ok := output.Output.(map[string]any)
	if !ok {
		return 0
	}
	v, ok := m["results"].(float64)
	if !ok {
		return 0
	}
	return v
}
