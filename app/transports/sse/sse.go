package sse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	adksession "google.golang.org/adk/session"
	"google.golang.org/adk/tool/toolconfirmation"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	storydenagent "github.com/Southclaws/storyden/app/services/semdex/robot"
	"github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/app/transports/http/middleware/headers"
	"github.com/Southclaws/storyden/app/transports/http/middleware/limiter"
	"github.com/Southclaws/storyden/app/transports/http/middleware/origin"
	"github.com/Southclaws/storyden/app/transports/http/middleware/reqlog"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/app/transports/http/robotprojection"
	"github.com/Southclaws/storyden/internal/infrastructure/httpserver"
	"github.com/Southclaws/storyden/lib/mcp"
)

const (
	contentTypeEventStream = "text/event-stream"
	headerNoCache          = "no-cache"
	headerKeepAlive        = "keep-alive"
	defaultFinishReason    = "stop"
	uiMessageStreamVersion = "v1"
)

// Build wires the SSE transport into the application.
func Build() fx.Option {
	return fx.Options(
		fx.Invoke(MountSSE),
	)
}

func MountSSE(
	lc fx.Lifecycle,
	ctx context.Context,
	logger *slog.Logger,

	chatAgent *storydenagent.Agent,
	sessionRepo *robot_session.Repository,
	toolRegistry *tools.Registry,

	mux *http.ServeMux,

	ri *headers.Middleware,
	co *origin.Middleware,
	lo *reqlog.Middleware,
	cj *session_cookie.Jar,
	rl *limiter.Middleware,
) {
	handler := newChatHandler(logger, chatAgent, sessionRepo, toolRegistry)

	applied := httpserver.Apply(handler,
		ri.WithHeaderContext(),
		co.WithCORS(),
		lo.WithLogger(),
		cj.WithAuth(),
		rl.WithRequestSizeLimiter(),
		rl.WithRateLimit(),
	)

	lc.Append(fx.StartHook(func() error {
		mux.Handle("POST /sse/chat", applied)
		return nil
	}))
}

type chatRequest struct {
	ID        string                 `json:"id"`
	ThreadID  string                 `json:"threadId"`
	SessionID string                 `json:"sessionId"`
	RobotID   string                 `json:"robotId,omitempty"`
	Messages  []chatMessage          `json:"messages"`
	Data      any                    `json:"data"`
	Context   *mcp.RobotChatContext  `json:"context,omitempty"`
	Workspace *workspaceMountRequest `json:"workspace,omitempty"`
}

type workspaceMountRequest struct {
	WorkspaceID         string `json:"workspace_id,omitempty"`
	WorkspaceInstanceID string `json:"workspace_instance_id,omitempty"`
}

type chatMessage struct {
	ID       string          `json:"id"`
	Role     string          `json:"role"`
	Parts    []chatPart      `json:"parts"`
	Metadata json.RawMessage `json:"metadata"`
}

type chatPart struct {
	Type       string          `json:"type"`
	Text       string          `json:"text,omitempty"`
	Delta      string          `json:"delta,omitempty"`
	Data       json.RawMessage `json:"data,omitempty"`
	State      string          `json:"state,omitempty"`
	Source     json.RawMessage `json:"source,omitempty"`
	ToolCallId string          `json:"toolCallId,omitempty"`
	ToolName   string          `json:"toolName,omitempty"`
	Input      json.RawMessage `json:"input,omitempty"`
	Output     json.RawMessage `json:"output,omitempty"`
	Approval   *chatApproval   `json:"approval,omitempty"`
}

type chatApproval struct {
	ID       string `json:"id,omitempty"`
	Approved bool   `json:"approved,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

func newChatHandler(
	logger *slog.Logger,
	chatAgent *storydenagent.Agent,
	sessionRepo *robot_session.Repository,
	toolRegistry *tools.Registry,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		accountID, err := session.GetAccountID(ctx)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if err := session.Authorise(ctx, nil, rbac.PermissionUseRobots); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		var req chatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("sse chat decode", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		robotRef := strings.TrimSpace(req.RobotID)
		if robotRef == "" {
			http.Error(w, "robotId is required", http.StatusBadRequest)
			return
		}

		sessionID := firstNonEmpty(req.SessionID, req.ThreadID, req.ID)
		if sessionID == "" {
			sessionID = fmt.Sprintf("chat-%s", accountID.String())
		}
		robotSessionID, err := robot.NewSessionID(sessionID)
		if err != nil {
			http.Error(w, "invalid session ID: must be a valid xid", http.StatusBadRequest)
			return
		}

		existingSess, _, sessionErr := sessionRepo.Get(ctx, robotSessionID, robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 1))
		if sessionErr != nil {
			if ftag.Get(sessionErr) == ftag.NotFound {
				sessionErr = nil
			}
		}
		if sessionErr != nil {
			http.Error(w, "failed to retrieve session: "+sessionErr.Error(), http.StatusInternalServerError)
			return
		}

		reconciliation := reconcilePendingClientTools(req.Messages, readPendingClientTools(existingSessState(existingSess)))
		pendingToolIDs := reconciliation.Pending.IDs
		if interaction, ok := reconciliation.BlockingInteraction.Get(); ok {
			writeChatError(w, interaction)
			return
		}
		if recovery, ok := reconciliation.StaleRobotSwitch.Get(); ok {
			logger.Warn("recovering stale pending robot switch",
				slog.String("session_id", sessionID),
				slog.String("tool_call_id", recovery.ToolCallID),
				slog.String("robot_id", recovery.RobotID))

			if err := persistClientToolResult(ctx, sessionRepo, robotSessionID, accountID, robotSwitchToolResultContent(recovery.ToolCallID, recovery.RobotID)); err != nil {
				logger.Error("failed to persist recovered robot switch tool result",
					slog.String("error", err.Error()),
					slog.String("robot_id", recovery.RobotID))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			state := clearPendingClientTools(existingSessState(existingSess))
			if err := updateCurrentRobotID(ctx, sessionRepo, robotSessionID, state, opt.New(recovery.RobotID)); err != nil {
				logger.Error("failed to recover current robot after stale switch",
					slog.String("error", err.Error()),
					slog.String("robot_id", recovery.RobotID))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			robotRef = recovery.RobotID
			pendingToolIDs = nil
		} else if ownerRobotID, ok := reconciliation.OwnerRobotID.Get(); ok {
			robotRef = ownerRobotID
		}

		robotID := opt.NewEmpty[xid.ID]()
		if id, err := xid.FromString(robotRef); err == nil {
			robotID = opt.New(id)
		}

		workspaceSpec, err := workspaceMountSpecFromRequest(req.Workspace)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		initMessage, err := getLastMessage(req.Messages, pendingToolIDs, logger)
		if err != nil {
			logger.Error("sse chat convert message", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if len(pendingToolIDs) > 0 && existingSess != nil {
			state := clearPendingClientTools(existingSess.State)
			if err := sessionRepo.UpdateState(ctx, robotSessionID, state); err != nil {
				logger.Error("failed to clear pending tool IDs", slog.String("error", err.Error()))
			}
		}

		nextCurrentRobotID := getRobotSwitchTargetID(req.Messages, pendingToolIDs)
		if targetRobotID, ok := nextCurrentRobotID.Get(); ok {
			if len(pendingToolIDs) > 0 {
				if err := persistClientToolResult(ctx, sessionRepo, robotSessionID, accountID, initMessage); err != nil {
					logger.Error("failed to persist robot switch tool result",
						slog.String("error", err.Error()),
						slog.String("robot_id", targetRobotID))
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				sess, _, err := sessionRepo.Get(ctx, robotSessionID, robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 1))
				if err != nil {
					logger.Error("failed to retrieve session after robot switch",
						slog.String("error", err.Error()),
						slog.String("robot_id", targetRobotID))
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				if err := updateCurrentRobotID(ctx, sessionRepo, robotSessionID, sess.State, opt.New(targetRobotID)); err != nil {
					logger.Error("failed to update current robot after switch",
						slog.String("error", err.Error()),
						slog.String("robot_id", targetRobotID))
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			}

			if err := finishEmptyStream(w); err != nil {
				logger.Error("failed to finish robot switch stream",
					slog.String("error", err.Error()),
					slog.String("robot_id", targetRobotID))
				return
			}

			logger.Info("robot switch completed without resuming previous robot",
				slog.String("session_id", sessionID),
				slog.String("robot_id", targetRobotID),
				slog.Bool("pending", len(pendingToolIDs) > 0))
			return
		}

		logger.Debug("sse chat request",
			slog.String("account_id", accountID.String()),
			slog.String("robot_id", robotRef),
			slog.String("session_id", sessionID),
			slog.String("user_message", lastUserMessage(req.Messages)),
			slog.Int("messages", len(req.Messages)),
			slog.Any("init_message", initMessage),
			slog.Any("context", req.Context),
		)

		stream := chatAgent.Run(ctx, robotID, accountID.String(), sessionID, initMessage, req.Context, storydenagent.RunOptions{
			Mode:      storydenagent.ModeInteractive,
			Source:    storydenagent.SourceInteractiveChat,
			RobotID:   opt.New(robotRef),
			Workspace: workspaceSpec,
		})

		emitter, err := newStreamEmitter(w)
		if err != nil {
			logger.Error("sse chat flusher", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer emitter.Done()

		responseID := uuid.NewString()
		textID := uuid.NewString()

		if err := emitter.Headers(); err != nil {
			logger.Error("sse chat headers", slog.String("error", err.Error()))
			return
		}

		startPart := openapi.StreamPart{}
		err = startPart.FromStartPart(openapi.StartPart{MessageId: responseID})
		if err != nil {
			logger.Error("failed to create start part", slog.String("error", err.Error()))
			return
		}

		if err := emitter.Send(startPart); err != nil {
			logger.Debug("sse chat start", slog.String("error", err.Error()))
			return
		}

		finishReason := defaultFinishReason
		finalSeen := false
		eventCount := 0

		for event, streamErr := range stream {
			eventCount++

			if streamErr != nil {
				if errors.Is(streamErr, context.Canceled) || errors.Is(ctx.Err(), context.Canceled) {
					logger.Debug("sse stream cancelled", slog.Int("event_num", eventCount))
					return
				}

				humanReadable := streamErrorText(streamErr)

				logger.Error("sse stream error",
					slog.Int("event_num", eventCount),
					slog.String("error", streamErr.Error()),
					slog.String("error_message", humanReadable),
				)

				errorPart := openapi.StreamPart{}
				err = errorPart.FromErrorPart(openapi.ErrorPart{
					ErrorText: humanReadable,
				})
				_ = emitter.Send(errorPart)
				return
			}

			if ctx.Err() != nil {
				return
			}

			if event != nil {
				logger.Info("sse event received",
					slog.Int("event_num", eventCount),
					slog.String("author", event.Author),
					slog.String("branch", event.Branch),
					slog.Bool("is_final", event.IsFinalResponse()),
					slog.Bool("has_content", event.LLMResponse.Content != nil),
					slog.Bool("is_partial", event.LLMResponse.Partial),
					slog.Bool("turn_complete", event.LLMResponse.TurnComplete),
					slog.String("finish_reason", string(event.LLMResponse.FinishReason)),
				)
			}

			if event != nil && event.LLMResponse.Content != nil {
				// Check for tool calls first
				for _, part := range event.LLMResponse.Content.Parts {
					if part == nil {
						continue
					}

					if part.FunctionCall != nil {
						if part.FunctionCall.Name == toolconfirmation.FunctionCallName {
							sendToolConfirmationCall(ctx, event, part, emitter, sessionRepo, robotSessionID, opt.New(robotRef), toolRegistry, logger)
						} else if toolRequiresConfirmation(ctx, toolRegistry, part.FunctionCall.Name) {
							continue
						} else {
							sendToolCall(ctx, event, part, emitter, toolRegistry, logger)
						}
					}

					if part.FunctionResponse != nil {
						if event.Author == "user" {
							continue
						}
						if isClientSidePending(part.FunctionResponse.Response) || toolRequiresConfirmation(ctx, toolRegistry, part.FunctionResponse.Name) {
							continue
						}
						sendToolResult(part, emitter, logger)
					}
				}

				// WORKAROUND: Check for client-side tool markers from BeforeToolCallback
				//
				// When a client-side tool is called, the
				// interceptClientSideTools callback in agent.go returns a marker
				// {"_client_side_pending": true} instead of executing the tool.
				//
				// We detect this marker here and immediately end the SSE stream WITHOUT
				// sending a finish event. This prevents the LLM from seeing or responding
				// to the marker.
				//
				// The frontend will then:
				// 1. Execute the tool on the client side
				// 2. POST back the real result
				// 3. Backend continues the agent with the real result
				//
				// See: agent.go:interceptClientSideTools
				for _, part := range event.LLMResponse.Content.Parts {
					if part != nil && part.FunctionResponse != nil {
						if isClientSidePending(part.FunctionResponse.Response) {
							logger.Info("client-side tool pending, ending stream to wait for client result",
								slog.String("tool_call_id", part.FunctionResponse.ID))

							if err := storePendingToolID(ctx, sessionRepo, robotSessionID, part.FunctionResponse.ID, opt.New(robotRef)); err != nil {
								logger.Error("failed to store pending tool ID",
									slog.String("error", err.Error()),
									slog.String("tool_call_id", part.FunctionResponse.ID))
							}

							finishPart := openapi.StreamPart{}
							_ = finishPart.FromFinishMessagePart(openapi.FinishMessagePart{})
							_ = emitter.Send(finishPart)

							return
						}
					}
				}

				if hasPendingConfirmation(event) {
					finishPart := openapi.StreamPart{}
					_ = finishPart.FromFinishMessagePart(openapi.FinishMessagePart{})
					_ = emitter.Send(finishPart)

					return
				}

				sendPresentationChunks(event, textID, emitter)
			}

			if event != nil && event.IsFinalResponse() {
				finalSeen = true
				if fr := strings.TrimSpace(string(event.LLMResponse.FinishReason)); fr != "" {
					finishReason = strings.ToLower(fr)
				}
			}
		}

		logger.Info("sse stream complete",
			slog.Int("total_events", eventCount),
			slog.Bool("final_seen", finalSeen),
			slog.String("finish_reason", finishReason))

		if !finalSeen {
			finishReason = defaultFinishReason
		}

		sess, _, err := sessionRepo.Get(ctx, robotSessionID, robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 1))
		if err == nil {
			if id, ok := nextCurrentRobotID.Get(); ok {
				if err := updateCurrentRobotID(ctx, sessionRepo, robotSessionID, sess.State, opt.New(id)); err != nil {
					logger.Error("failed to update current robot",
						slog.String("error", err.Error()),
						slog.String("robot_id", id))
				}
			} else {
				if err := updateCurrentRobotID(ctx, sessionRepo, robotSessionID, sess.State, opt.New(robotRef)); err != nil {
					logger.Error("failed to update current robot",
						slog.String("error", err.Error()),
						slog.String("robot_id", robotRef))
				}
			}

			dataPart := openapi.StreamPart{}
			dataPart.FromDataPart(openapi.DataPart{
				Data: sess.Name,
			})
			dataPart.Type = "data-session_name"
			_ = emitter.Send(dataPart)
		}

		finishPart := openapi.StreamPart{}
		err = finishPart.FromFinishMessagePart(openapi.FinishMessagePart{})
		_ = emitter.Send(finishPart)
	})
}

func getLastMessage(messages []chatMessage, pendingToolIDs []string, logger *slog.Logger) (*genai.Content, error) {
	if len(messages) == 0 {
		return nil, fmt.Errorf("no messages provided")
	}

	lastMessage := messages[len(messages)-1]

	content := &genai.Content{
		Role:  lastMessage.Role,
		Parts: []*genai.Part{},
	}

	pendingSet := make(map[string]bool)
	for _, id := range pendingToolIDs {
		pendingSet[id] = true
	}

	switch strings.ToLower(lastMessage.Role) {
	case "user":
		for _, part := range lastMessage.Parts {
			switch part.Type {
			case "text":
				if part.Text != "" {
					content.Parts = append(content.Parts, &genai.Part{Text: part.Text})
				}
			}
		}

	case "assistant":
		content.Role = "user"

		for _, part := range lastMessage.Parts {
			if strings.HasPrefix(part.Type, "tool-") && part.State == "approval-responded" {
				if part.ToolCallId == "" {
					return nil, fmt.Errorf("tool approval missing toolCallId: type=%s", part.Type)
				}
				if part.Approval == nil {
					return nil, fmt.Errorf("tool approval missing approval payload: type=%s, toolCallId=%s", part.Type, part.ToolCallId)
				}

				approvalID := part.Approval.ID
				if approvalID == "" {
					approvalID = part.ToolCallId
				}

				if len(pendingSet) > 0 && !pendingSet[approvalID] && !pendingSet[part.ToolCallId] {
					logger.Info("skipping tool approval not in pending list",
						slog.String("tool_call_id", part.ToolCallId),
						slog.String("approval_id", approvalID),
						slog.String("tool_name", part.ToolName))
					continue
				}

				content.Parts = append(content.Parts, &genai.Part{
					FunctionResponse: &genai.FunctionResponse{
						ID:   approvalID,
						Name: toolconfirmation.FunctionCallName,
						Response: map[string]any{
							"confirmed": part.Approval.Approved,
						},
					},
				})

				logger.Info("tool approval received from frontend",
					slog.String("tool_call_id", part.ToolCallId),
					slog.String("approval_id", approvalID),
					slog.String("tool_name", part.ToolName),
					slog.Bool("approved", part.Approval.Approved))

				continue
			}

			if strings.HasPrefix(part.Type, "tool-") && part.State == "output-available" {
				if part.ToolCallId == "" {
					return nil, fmt.Errorf("tool result missing toolCallId: type=%s", part.Type)
				}

				if len(pendingSet) > 0 && !pendingSet[part.ToolCallId] {
					logger.Info("skipping tool result not in pending list",
						slog.String("tool_call_id", part.ToolCallId),
						slog.String("tool_name", part.ToolName))
					continue
				}

				output, err := resolveToolOutput(part)
				if err != nil {
					return nil, fmt.Errorf("failed to parse tool output for %s: %w", part.ToolCallId, err)
				}

				toolName := part.ToolName
				if toolName == "" && strings.HasPrefix(part.Type, "tool-") {
					toolName = strings.TrimPrefix(part.Type, "tool-")
				}
				if isStorydenConfirmationOutput(part.Output) {
					toolName = toolconfirmation.FunctionCallName
				}

				if toolName == "" {
					return nil, fmt.Errorf("tool result missing tool name: type=%s, toolCallId=%s", part.Type, part.ToolCallId)
				}

				content.Parts = append(content.Parts, &genai.Part{
					FunctionResponse: &genai.FunctionResponse{
						ID:       part.ToolCallId,
						Name:     toolName,
						Response: output,
					},
				})

				logger.Info("tool result received from frontend",
					slog.String("tool_call_id", part.ToolCallId),
					slog.String("tool_name", toolName),
					slog.Any("output", output))
			}
		}
	}

	if len(content.Parts) == 0 {
		return nil, fmt.Errorf("user message has no content")
	}

	return content, nil
}

func resolveToolOutput(part chatPart) (map[string]any, error) {
	var output map[string]any
	if err := json.Unmarshal(part.Output, &output); err != nil {
		return nil, err
	}

	confirmation, ok := output["_storyden_confirmation"].(map[string]any)
	if !ok {
		return output, nil
	}

	approved, _ := confirmation["approved"].(bool)
	return map[string]any{"confirmed": approved}, nil
}

func isStorydenConfirmationOutput(output json.RawMessage) bool {
	var payload map[string]any
	if err := json.Unmarshal(output, &payload); err != nil {
		return false
	}
	_, ok := payload["_storyden_confirmation"].(map[string]any)
	return ok
}

func updateCurrentRobotID(ctx context.Context, sessionRepo *robot_session.Repository, sessionID robot.SessionID, state map[string]any, robotID opt.Optional[string]) error {
	if state == nil {
		state = make(map[string]any)
	}

	if id, ok := robotID.Get(); ok {
		state["current_robot_id"] = id
	} else {
		delete(state, "current_robot_id")
	}

	return sessionRepo.UpdateState(ctx, sessionID, state)
}

func lastUserMessage(messages []chatMessage) string {
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		if !strings.EqualFold(msg.Role, "user") {
			continue
		}

		text := extractTextFromParts(msg.Parts)
		if text != "" {
			return text
		}
	}
	return ""
}

func extractTextFromParts(parts []chatPart) string {
	var b strings.Builder
	for _, part := range parts {
		if part.Type != "text" && part.Type != "reasoning" {
			continue
		}
		fragment := part.Text
		if fragment == "" {
			fragment = part.Delta
		}
		fragment = strings.TrimSpace(fragment)
		if fragment == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString(fragment)
	}
	return b.String()
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func workspaceMountSpecFromRequest(req *workspaceMountRequest) (opt.Optional[storydenagent.WorkspaceMountSpec], error) {
	if req == nil {
		return opt.NewEmpty[storydenagent.WorkspaceMountSpec](), nil
	}

	hasWorkspaceID := strings.TrimSpace(req.WorkspaceID) != ""
	hasInstanceID := strings.TrimSpace(req.WorkspaceInstanceID) != ""
	if hasWorkspaceID == hasInstanceID {
		return opt.NewEmpty[storydenagent.WorkspaceMountSpec](), fmt.Errorf("provide exactly one workspace_id or workspace_instance_id")
	}

	if hasWorkspaceID {
		id, err := robot.NewWorkspaceID(req.WorkspaceID)
		if err != nil {
			return opt.NewEmpty[storydenagent.WorkspaceMountSpec](), err
		}
		return opt.New(storydenagent.WorkspaceMountSpec{
			WorkspaceID: opt.New(id),
			Metadata:    map[string]any{},
		}), nil
	}

	id, err := robot.NewWorkspaceInstanceID(req.WorkspaceInstanceID)
	if err != nil {
		return opt.NewEmpty[storydenagent.WorkspaceMountSpec](), err
	}
	return opt.New(storydenagent.WorkspaceMountSpec{
		WorkspaceInstanceID: opt.New(id),
		Metadata:            map[string]any{},
	}), nil
}

// isClientSidePending checks if a tool response is the special marker from
// interceptClientSideTools (see agent.go) indicating that a client-side tool
// was called and we should wait for the real result from the frontend.
func isClientSidePending(response map[string]any) bool {
	if response == nil {
		return false
	}
	pending, ok := response["_client_side_pending"].(bool)
	return ok && pending
}

func toolRequiresConfirmation(ctx context.Context, toolRegistry *tools.Registry, toolName string) bool {
	if toolRegistry == nil || toolName == "" {
		return false
	}
	tool, err := toolRegistry.GetTool(ctx, toolName)
	return err == nil && tool.Definition.RequiresConfirmation
}

func hasPendingConfirmation(event *adksession.Event) bool {
	if event == nil || event.LLMResponse.Content == nil {
		return false
	}
	for _, part := range event.LLMResponse.Content.Parts {
		if part == nil || part.FunctionCall == nil {
			continue
		}
		if part.FunctionCall.Name == toolconfirmation.FunctionCallName {
			return true
		}
	}
	return false
}

func sendPresentationChunks(event *adksession.Event, fallbackTextID string, emitter *streamEmitter) {
	for _, streamPart := range robotprojection.PresentationStreamParts(event, fallbackTextID) {
		_ = emitter.Send(streamPart)
	}
}

func streamErrorText(err error) string {
	raw := strings.TrimSpace(err.Error())
	issue := strings.TrimSpace(fmsg.GetIssue(err))

	if issue == "" {
		return raw
	}

	if raw == "" || raw == issue {
		return issue
	}

	return fmt.Sprintf("%s (%s)", issue, raw)
}

func sendToolCall(ctx context.Context, event *adksession.Event, part *genai.Part, emitter *streamEmitter, toolRegistry *tools.Registry, logger *slog.Logger) {
	fc := part.FunctionCall
	if fc == nil {
		return
	}

	toolCallId := fc.ID
	toolName := fc.Name

	logger.Info("tool call detected",
		slog.String("tool_call_id", toolCallId),
		slog.String("tool_name", toolName),
		slog.Any("args", fc.Args),
		slog.Any("long_running_ids", event.LongRunningToolIDs),
	)

	metadata := robotprojection.ToolMetadataFromRegistry(ctx, toolRegistry)(toolName)
	for _, streamPart := range robotprojection.FunctionCallStreamPartsWithMetadata(fc, metadata) {
		_ = emitter.Send(streamPart)
	}
}

func sendToolConfirmationCall(
	ctx context.Context,
	event *adksession.Event,
	part *genai.Part,
	emitter *streamEmitter,
	sessionRepo *robot_session.Repository,
	robotSessionID robot.SessionID,
	robotID opt.Optional[string],
	toolRegistry *tools.Registry,
	logger *slog.Logger,
) {
	fc := part.FunctionCall
	if fc == nil {
		return
	}

	original, err := toolconfirmation.OriginalCallFrom(fc)
	if err != nil {
		logger.Error("failed to parse tool confirmation call",
			slog.String("tool_call_id", fc.ID),
			slog.String("error", err.Error()))
		return
	}

	logger.Info("tool confirmation requested",
		slog.String("confirmation_call_id", fc.ID),
		slog.String("tool_call_id", original.ID),
		slog.String("tool_name", original.Name),
		slog.Any("args", original.Args),
		slog.Any("long_running_ids", event.LongRunningToolIDs),
	)

	if err := storePendingToolID(ctx, sessionRepo, robotSessionID, fc.ID, robotID); err != nil {
		logger.Error("failed to store pending confirmation ID",
			slog.String("error", err.Error()),
			slog.String("tool_call_id", fc.ID))
	}

	confirmationPart := &genai.Part{
		FunctionCall: &genai.FunctionCall{
			ID:   fc.ID,
			Name: original.Name,
			Args: original.Args,
		},
	}
	sendToolCall(ctx, event, confirmationPart, emitter, toolRegistry, logger)
	_ = emitter.Send(robotprojection.ToolApprovalRequestStreamPart(fc.ID, fc.ID))
}

func sendToolResult(part *genai.Part, emitter *streamEmitter, logger *slog.Logger) {
	fr := part.FunctionResponse
	if fr == nil {
		return
	}

	toolCallId := fr.ID
	toolName := fr.Name

	logger.Info("tool result detected",
		slog.String("tool_call_id", toolCallId),
		slog.String("tool_name", toolName),
		slog.Any("response", fr.Response),
	)

	if streamPart, ok := robotprojection.FunctionResponseStreamPart(fr); ok {
		_ = emitter.Send(streamPart)
	}
}

type streamEmitter struct {
	w              http.ResponseWriter
	flusher        http.Flusher
	once           sync.Once
	mu             sync.Mutex
	headersWritten bool
}

func newStreamEmitter(w http.ResponseWriter) (*streamEmitter, error) {
	flusher, ok := GetFlusher(w)
	if !ok {
		return nil, errors.New("streaming unsupported")
	}
	return &streamEmitter{w: w, flusher: flusher}, nil
}

func GetFlusher(w http.ResponseWriter) (http.Flusher, bool) {
	for {
		if f, ok := w.(http.Flusher); ok {
			return f, true
		}
		// Try to unwrap
		if unwrapper, ok := w.(interface{ Unwrap() http.ResponseWriter }); ok {
			w = unwrapper.Unwrap()
		} else {
			return nil, false
		}
	}
}

func (s *streamEmitter) Headers() error {
	s.once.Do(func() {
		header := s.w.Header()
		header.Set("Content-Type", contentTypeEventStream)
		header.Set("Cache-Control", headerNoCache)
		header.Set("Connection", headerKeepAlive)
		header.Set("X-Accel-Buffering", "no")
		header.Set("X-Vercel-AI-UI-Message-Stream", uiMessageStreamVersion)
		s.headersWritten = true
	})
	return nil
}

func (s *streamEmitter) Send(payload openapi.StreamPart) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.headersWritten {
		if err := s.Headers(); err != nil {
			return err
		}
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if _, err := s.w.Write([]byte("data: ")); err != nil {
		return err
	}
	if _, err := s.w.Write(data); err != nil {
		return err
	}
	if _, err := s.w.Write([]byte("\n\n")); err != nil {
		return err
	}
	s.flusher.Flush()
	return nil
}

func (s *streamEmitter) Done() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.headersWritten {
		s.Headers()
	}
	_, _ = s.w.Write([]byte("data: [DONE]\n\n"))
	s.flusher.Flush()
}
