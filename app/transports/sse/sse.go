package sse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"go.uber.org/fx"

	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	adksession "google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/tool/toolconfirmation"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_querier"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	storydenagent "github.com/Southclaws/storyden/app/services/semdex/robot"
	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry/denbot"
	"github.com/Southclaws/storyden/app/services/semdex/robot/session_coordinator"
	"github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/app/transports/http/middleware/headers"
	"github.com/Southclaws/storyden/app/transports/http/middleware/limiter"
	"github.com/Southclaws/storyden/app/transports/http/middleware/origin"
	"github.com/Southclaws/storyden/app/transports/http/middleware/reqlog"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/app/transports/http/robotprojection"
	"github.com/Southclaws/storyden/internal/infrastructure/httpserver"
	"github.com/Southclaws/storyden/lib/mcp"
)

const (
	contentTypeEventStream = "text/event-stream"
	headerNoCache          = "no-cache"
	headerKeepAlive        = "keep-alive"
	defaultFinishReason    = "stop"
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

	chatAgent *session_coordinator.Coordinator,
	sessionRepo *robot_session.Repository,
	robotQuerier *robot_querier.Querier,
	toolRegistry *tools.Registry,
	mux *http.ServeMux,

	ri *headers.Middleware,
	co *origin.Middleware,
	lo *reqlog.Middleware,
	cj *session_cookie.Jar,
	rl *limiter.Middleware,
) {
	chatHandler := newChatHandler(logger, chatAgent, sessionRepo, robotQuerier, toolRegistry)
	readHandler := newDurableReadHandler(logger, chatAgent, sessionRepo, robotQuerier, toolRegistry)

	apply := func(handler http.Handler) http.Handler {
		return httpserver.Apply(handler,
			ri.WithHeaderContext(),
			co.WithCORS(),
			lo.WithLogger(),
			cj.WithAuth(),
			rl.WithRequestSizeLimiter(),
			rl.WithRateLimit(),
		)
	}

	lc.Append(fx.StartHook(func() error {
		mux.Handle("POST /sse/chat", apply(chatHandler))
		mux.Handle("GET /sse/chat/{sessionID}/stream", apply(newReconnectHandler(logger, sessionRepo)))
		mux.Handle("GET /sse/sessions/{sessionID}/turns/{turnID}", apply(readHandler))
		mux.Handle("HEAD /sse/sessions/{sessionID}/turns/{turnID}", apply(readHandler))
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
	chatAgent *session_coordinator.Coordinator,
	sessionRepo *robot_session.Repository,
	robotQuerier *robot_querier.Querier,
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

		requestedRobotRef := strings.TrimSpace(req.RobotID)
		robotRef := requestedRobotRef
		if robotRef == "" {
			robotRef = denbot.ID
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
		if existingSess != nil {
			rootRobotRef := storydenagent.SessionRootRobotRef(existingSess.State).Or(denbot.ID)
			if requestedRobotRef != "" && requestedRobotRef != rootRobotRef {
				http.Error(w, "robotId can only select the root Robot when starting a session", http.StatusConflict)
				return
			}
			robotRef = rootRobotRef
		}

		reconciliation := reconcilePendingClientTools(req.Messages, readPendingClientTools(existingSessState(existingSess)))
		pendingToolIDs := reconciliation.Pending.IDs
		if interaction, ok := reconciliation.BlockingInteraction.Get(); ok {
			writeChatError(w, interaction)
			return
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

		logger.Debug("sse chat request",
			slog.String("account_id", accountID.String()),
			slog.String("robot_id", robotRef),
			slog.String("session_id", sessionID),
			slog.String("user_message", lastUserMessage(req.Messages)),
			slog.Int("messages", len(req.Messages)),
			slog.Any("init_message", initMessage),
			slog.Any("context", req.Context),
		)

		turnID, err := chatAgent.Start(ctx, robotRef, accountID.String(), sessionID, initMessage, req.Context, storydenagent.RunOptions{
			Mode:      storydenagent.ModeInteractive,
			Source:    storydenagent.SourceInteractiveChat,
			Workspace: workspaceSpec,
		})
		if err != nil {
			logger.Error("Robot turn start", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		streamPath := "/sse/sessions/" + robotSessionID.String() + "/turns/" + turnID.String()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Location", streamPath)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{"streamUrl": streamPath})
	})
}

func newReconnectHandler(logger *slog.Logger, sessions *robot_session.Repository) http.Handler {
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

		sessionID, err := robot.NewSessionID(r.PathValue("sessionID"))
		if err != nil {
			http.Error(w, "invalid session ID: must be a valid xid", http.StatusBadRequest)
			return
		}
		turnID, err := sessions.ResumeTurn(ctx, sessionID, accountID)
		if err != nil {
			if errors.Is(err, robot_session.ErrTurnNotFound) {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			logger.Error("durable stream reconnect", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		allowed, err := sessions.CanReadTurn(ctx, sessionID, turnID, accountID)
		if err != nil {
			logger.Error("durable stream reconnect authorisation", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if !allowed {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		streamPath := "/sse/sessions/" + sessionID.String() + "/turns/" + turnID.String()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Location", streamPath)
		_ = json.NewEncoder(w).Encode(map[string]string{"streamUrl": streamPath})
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

func sendToolCall(ctx context.Context, event *adksession.Event, part *genai.Part, emitter partEmitter, toolRegistry *tools.Registry, logger *slog.Logger) {
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
	emitter partEmitter,
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

func sendToolResult(part *genai.Part, emitter partEmitter, logger *slog.Logger) {
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
