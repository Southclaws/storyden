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

	"github.com/Southclaws/opt"
	adksession "google.golang.org/adk/session"
	"google.golang.org/genai"

	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	storydenagent "github.com/Southclaws/storyden/app/services/semdex/robot"
	"github.com/Southclaws/storyden/app/transports/http/middleware/limiter"
	"github.com/Southclaws/storyden/app/transports/http/middleware/origin"
	"github.com/Southclaws/storyden/app/transports/http/middleware/reqlog"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/httpserver"
	"github.com/Southclaws/storyden/mcp"
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
	cfg config.Config,

	chatAgent *storydenagent.Agent,
	sessionRepo *robot_session.Repository,

	mux *http.ServeMux,

	co *origin.Middleware,
	lo *reqlog.Middleware,
	cj *session_cookie.Jar,
	rl *limiter.Middleware,
) {
	if cfg.LanguageModelProvider == "" {
		return
	}

	handler := newChatHandler(logger, chatAgent, sessionRepo)

	applied := httpserver.Apply(handler,
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
	ID        string                `json:"id"`
	ThreadID  string                `json:"threadId"`
	SessionID string                `json:"sessionId"`
	RobotID   string                `json:"robotId,omitempty"`
	Messages  []chatMessage         `json:"messages"`
	Data      any                   `json:"data"`
	Context   *mcp.RobotChatContext `json:"context,omitempty"`
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
	Output     json.RawMessage `json:"output,omitempty"`
}

func newChatHandler(logger *slog.Logger, chatAgent *storydenagent.Agent, sessionRepo *robot_session.Repository) http.Handler {
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

		robotID, err := opt.MapErr(opt.NewIf(req.RobotID, func(s string) bool {
			return s != ""
		}), xid.FromString)
		if err != nil {
			http.Error(w, "invalid robot ID: must be a valid xid", http.StatusBadRequest)
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

		existingSess, _, sessionErr := sessionRepo.Get(ctx, robotSessionID, pagination.NewPageParams(0, 0))

		pendingToolIDs := []string{}
		if sessionErr == nil && existingSess.State != nil {
			if existing, ok := existingSess.State["pending_client_tools"]; ok {
				if ids, ok := existing.([]interface{}); ok {
					for _, id := range ids {
						if s, ok := id.(string); ok {
							pendingToolIDs = append(pendingToolIDs, s)
						}
					}
				}
			}
		}

		initMessage, err := getLastMessage(req.Messages, pendingToolIDs, logger)
		if err != nil {
			logger.Error("sse chat convert message", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if len(pendingToolIDs) > 0 && existingSess != nil {
			state := existingSess.State
			if state == nil {
				state = make(map[string]any)
			}
			state["pending_client_tools"] = []string{}
			if err := sessionRepo.UpdateState(ctx, robotSessionID, state); err != nil {
				logger.Error("failed to clear pending tool IDs", slog.String("error", err.Error()))
			}
		}

		logger.Debug("sse chat request",
			slog.String("account_id", accountID.String()),
			slog.String("robot_id", robotID.String()),
			slog.String("session_id", sessionID),
			slog.String("user_message", lastUserMessage(req.Messages)),
			slog.Int("messages", len(req.Messages)),
			slog.Any("init_message", initMessage),
			slog.Any("context", req.Context),
		)

		stream := chatAgent.Run(ctx, robotID, accountID.String(), sessionID, initMessage, req.Context)

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
		textStarted := false

		for event, streamErr := range stream {
			eventCount++

			if streamErr != nil {
				humanReadable := fmsg.GetIssue(streamErr)
				if humanReadable == "" {
					humanReadable = "An unknown error occurred."
				}

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

					// Tool call invocations
					if part.FunctionCall != nil {
						sendToolCall(event, part, emitter, logger)
					}

					// Tool results
					if part.FunctionResponse != nil {
						sendToolResult(event, part, emitter, logger)
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

							if err := storePendingToolID(ctx, sessionRepo, robotSessionID, part.FunctionResponse.ID); err != nil {
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

				// Then handle text content
				hasText := false
				for _, part := range event.LLMResponse.Content.Parts {
					if part != nil && strings.TrimSpace(part.Text) != "" {
						hasText = true
						break
					}
				}

				if hasText && !textStarted {
					textStartPart := openapi.StreamPart{}
					err = textStartPart.FromTextStartPart(openapi.TextStartPart{
						Id: textID,
					})
					_ = emitter.Send(textStartPart)
					textStarted = true
				}

				if textStarted {
					sendTextChunks(event, textID, emitter)
				}
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
			slog.Bool("text_started", textStarted),
			slog.String("finish_reason", finishReason))

		if !finalSeen {
			finishReason = defaultFinishReason
		}

		if textStarted {
			textEndPart := openapi.StreamPart{}
			err = textEndPart.FromTextEndPart(openapi.TextEndPart{
				Id: textID,
			})
			_ = emitter.Send(textEndPart)
		}

		sess, _, err := sessionRepo.Get(ctx, robotSessionID, pagination.NewPageParams(0, 0))
		if err == nil {
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
		for _, part := range lastMessage.Parts {
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

				var output map[string]any
				if err := json.Unmarshal(part.Output, &output); err != nil {
					return nil, fmt.Errorf("failed to parse tool output for %s: %w", part.ToolCallId, err)
				}

				toolName := part.ToolName
				if toolName == "" && strings.HasPrefix(part.Type, "tool-") {
					toolName = strings.TrimPrefix(part.Type, "tool-")
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

func sendTextChunks(event *adksession.Event, textID string, emitter *streamEmitter) {
	if event == nil || event.LLMResponse.Content == nil {
		return
	}
	for _, part := range event.LLMResponse.Content.Parts {
		if part == nil || strings.TrimSpace(part.Text) == "" {
			continue
		}
		textDeltaPart := openapi.StreamPart{}
		if err := textDeltaPart.FromTextDeltaPart(openapi.TextDeltaPart{
			Id:    textID,
			Delta: part.Text,
		}); err != nil {
			return
		}

		if err := emitter.Send(textDeltaPart); err != nil {
			return
		}
	}
}

func sendToolCall(event *adksession.Event, part *genai.Part, emitter *streamEmitter, logger *slog.Logger) {
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

	// Emit tool-input-start using generated type
	toolInputStartPart := openapi.StreamPart{}
	_ = toolInputStartPart.FromToolInputStartPart(openapi.ToolInputStartPart{
		ToolCallId: toolCallId,
		ToolName:   toolName,
	})
	_ = emitter.Send(toolInputStartPart)

	// Emit tool-input-delta with arguments (serialize args as JSON string)
	argsJSON, err := json.Marshal(fc.Args)
	if err == nil {
		toolInputDeltaPart := openapi.StreamPart{}
		_ = toolInputDeltaPart.FromToolInputDeltaPart(openapi.ToolInputDeltaPart{
			ToolCallId:     toolCallId,
			InputTextDelta: string(argsJSON),
		})
		_ = emitter.Send(toolInputDeltaPart)
	}

	// Emit tool-input-available to signal tool is ready for execution
	toolInputAvailablePart := openapi.StreamPart{}
	_ = toolInputAvailablePart.FromToolInputAvailablePart(openapi.ToolInputAvailablePart{
		ToolCallId: toolCallId,
		ToolName:   toolName,
		Input:      fc.Args,
	})
	_ = emitter.Send(toolInputAvailablePart)
}

func sendToolResult(event *adksession.Event, part *genai.Part, emitter *streamEmitter, logger *slog.Logger) {
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

	// Emit tool-output-available with the result
	toolOutputAvailablePart := openapi.StreamPart{}
	_ = toolOutputAvailablePart.FromToolOutputAvailablePart(openapi.ToolOutputAvailablePart{
		ToolCallId: toolCallId,
		Output:     fr.Response,
	})
	_ = emitter.Send(toolOutputAvailablePart)
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

func storePendingToolID(ctx context.Context, sessionRepo *robot_session.Repository, sessionID robot.SessionID, toolCallID string) error {
	sess, _, err := sessionRepo.Get(ctx, sessionID, pagination.NewPageParams(0, 0))
	if err != nil {
		return err
	}

	state := sess.State
	if state == nil {
		state = make(map[string]any)
	}

	var pendingIDs []string
	if existing, ok := state["pending_client_tools"]; ok {
		if ids, ok := existing.([]interface{}); ok {
			for _, id := range ids {
				if s, ok := id.(string); ok {
					pendingIDs = append(pendingIDs, s)
				}
			}
		}
	}

	pendingIDs = append(pendingIDs, toolCallID)
	state["pending_client_tools"] = pendingIDs

	return sessionRepo.UpdateState(ctx, sessionID, state)
}
