package bindings

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	adksession "google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/tool/toolconfirmation"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	storydenagent "github.com/Southclaws/storyden/app/services/semdex/robot"
	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry/denbot"
	"github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/app/transports/http/robotprojection"
)

const (
	contentTypeEventStream = "text/event-stream"
	headerNoCache          = "no-cache"
	headerKeepAlive        = "keep-alive"
	defaultFinishReason    = "stop"
)

type robotSessionCreateResponseFunc func(http.ResponseWriter)

func (f robotSessionCreateResponseFunc) VisitRobotSessionCreateResponse(w http.ResponseWriter) error {
	f(w)
	return nil
}

type robotSessionStreamResponseFunc func(http.ResponseWriter)

func (f robotSessionStreamResponseFunc) VisitRobotSessionStreamResponse(w http.ResponseWriter) error {
	f(w)
	return nil
}

type robotSessionStreamHeadResponseFunc func(http.ResponseWriter)

func (f robotSessionStreamHeadResponseFunc) VisitRobotSessionStreamHeadResponse(w http.ResponseWriter) error {
	f(w)
	return nil
}

type chatRequest struct {
	ID        string                    `json:"id"`
	ThreadID  string                    `json:"threadId"`
	SessionID string                    `json:"sessionId"`
	RobotID   string                    `json:"robotId,omitempty"`
	Messages  []chatMessage             `json:"messages"`
	Data      any                       `json:"data"`
	Context   *openapi.RobotChatContext `json:"context,omitempty"`
	Workspace *workspaceMountRequest    `json:"workspace,omitempty"`
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

func (r *Robots) RobotSessionCreate(ctx context.Context, request openapi.RobotSessionCreateRequestObject) (openapi.RobotSessionCreateResponseObject, error) {
	return robotSessionCreateResponseFunc(func(w http.ResponseWriter) {
		r.createSessionTurn(ctx, w, request.Body)
	}), nil
}

func (r *Robots) createSessionTurn(ctx context.Context, w http.ResponseWriter, body *openapi.RobotSessionCreateJSONRequestBody) {
	if body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		r.logger.Error("Robot session request encode", slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var req chatRequest
	if err := json.Unmarshal(encoded, &req); err != nil {
		r.logger.Error("Robot session request decode", slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	normalizeChatPartTypes(req.Messages)

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
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

	existingSess, _, sessionErr := r.sessionRepo.Get(ctx, robotSessionID, robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 1))
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

	initMessage, err := getLastMessage(req.Messages, pendingToolIDs, r.logger)
	if err != nil {
		r.logger.Error("Robot session convert message", slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	invocationContext := invocationContextFromRequest(req.Context)
	lastMessageID, err := xid.FromString(req.Messages[len(req.Messages)-1].ID)
	if err != nil {
		http.Error(w, "message ID must be a valid xid", http.StatusBadRequest)
		return
	}

	r.logger.Debug(
		"Robot session turn request",
		slog.String("account_id", accountID.String()),
		slog.String("robot_id", robotRef),
		slog.String("session_id", sessionID),
		slog.String("user_message", lastUserMessage(req.Messages)),
		slog.Int("messages", len(req.Messages)),
		slog.Any("init_message", initMessage),
		slog.Any("invocation_context", invocationContext),
	)

	inputID, err := r.coordinator.Enqueue(ctx, robot.InputID(lastMessageID), robotRef, accountID.String(), sessionID, initMessage, invocationContext, storydenagent.RunOptions{
		Mode:      storydenagent.ModeInteractive,
		Source:    storydenagent.SourceInteractiveChat,
		Workspace: workspaceSpec,
	})
	if err != nil {
		r.logger.Error("Robot turn start", slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	streamPath := "/api/robots/sessions/" + robotSessionID.String() + "/stream"

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Location", streamPath)
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(openapi.RobotSessionInputReference{
		StreamUrl: streamPath,
		SessionId: openapi.Identifier(robotSessionID.String()),
		MessageId: openapi.Identifier(inputID.String()),
	})
}

func invocationContextFromRequest(context *openapi.RobotChatContext) storydenagent.InvocationContext {
	if context == nil {
		return nil
	}

	result := storydenagent.InvocationContext{}
	if context.DatagraphItem != nil {
		result[storydenagent.InvocationContextKeyDatagraphItem] = context.DatagraphItem
	}
	if context.PageType != nil {
		result[storydenagent.InvocationContextKeyPageType] = *context.PageType
	}
	return result
}

func (r *Robots) RobotSessionStream(ctx context.Context, request openapi.RobotSessionStreamRequestObject) (openapi.RobotSessionStreamResponseObject, error) {
	return robotSessionStreamResponseFunc(func(w http.ResponseWriter) {
		offsetValue := ""
		if request.Params.Offset != nil {
			offsetValue = string(*request.Params.Offset)
		}
		live := ""
		if request.Params.Live != nil {
			live = string(*request.Params.Live)
		}
		r.serveSessionEvents(ctx, w, string(request.SessionId), offsetValue, live, false)
	}), nil
}

func (r *Robots) RobotSessionStreamHead(ctx context.Context, request openapi.RobotSessionStreamHeadRequestObject) (openapi.RobotSessionStreamHeadResponseObject, error) {
	return robotSessionStreamHeadResponseFunc(func(w http.ResponseWriter) {
		offsetValue := ""
		if request.Params.Offset != nil {
			offsetValue = string(*request.Params.Offset)
		}
		r.serveSessionEvents(ctx, w, string(request.SessionId), offsetValue, "", true)
	}), nil
}

func (r *Robots) serveSessionEvents(ctx context.Context, w http.ResponseWriter, sessionValue, offsetValue, live string, head bool) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	sessionID, err := robot.NewSessionID(sessionValue)
	if err != nil {
		http.Error(w, "invalid session ID: must be a valid xid", http.StatusBadRequest)
		return
	}
	allowed, err := r.sessionRepo.CanReadSession(ctx, sessionID, accountID)
	if err != nil {
		r.logger.Error("Robot session stream authorisation", slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	offset, err := parseStreamOffset(offsetValue)
	if err != nil {
		http.Error(w, "invalid stream offset", http.StatusBadRequest)
		return
	}
	if err := r.sessionRepo.AcknowledgeSessionEvents(ctx, sessionID, accountID, offset); err != nil {
		r.logger.Error("Robot session stream acknowledgement", slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if head {
		serveSessionCatchUp(ctx, w, r.logger, r.sessionRepo, r.robotQuerier, r.tools, sessionID, offset, true)
		return
	}
	switch live {
	case "":
		serveSessionCatchUp(ctx, w, r.logger, r.sessionRepo, r.robotQuerier, r.tools, sessionID, offset, false)
	case "sse":
		serveSessionLive(ctx, w, r.logger, r.coordinator, r.sessionRepo, r.robotQuerier, r.tools, sessionID, offset)
	default:
		http.Error(w, "unsupported stream live mode", http.StatusBadRequest)
	}
}

func normalizeChatPartTypes(messages []chatMessage) {
	for messageIndex := range messages {
		for partIndex := range messages[messageIndex].Parts {
			part := &messages[messageIndex].Parts[partIndex]
			if part.Type != "" {
				continue
			}
			if part.ToolName != "" {
				part.Type = "tool-" + part.ToolName
				continue
			}
			if part.Text != "" {
				part.Type = "text"
			}
		}
	}
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
					logger.Debug("skipping tool approval not in pending list",
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

				logger.Debug("tool approval received from frontend",
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
					logger.Debug("skipping tool result not in pending list",
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

				logger.Debug("tool result received from frontend",
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

func sendToolCall(ctx context.Context, event *adksession.Event, part *genai.Part, emitter partEmitter, toolRegistry *tools.Registry, logger *slog.Logger) {
	fc := part.FunctionCall
	if fc == nil {
		return
	}

	toolCallId := fc.ID
	toolName := fc.Name

	logger.Debug(
		"tool call detected",
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

	logger.Debug(
		"tool confirmation requested",
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

	logger.Debug(
		"tool result detected",
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
