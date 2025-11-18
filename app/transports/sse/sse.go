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
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/authentication/session"
	storydenagent "github.com/Southclaws/storyden/app/services/semdex/agent"

	"github.com/Southclaws/storyden/app/transports/http/middleware/limiter"
	"github.com/Southclaws/storyden/app/transports/http/middleware/origin"
	"github.com/Southclaws/storyden/app/transports/http/middleware/reqlog"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/httpserver"

	agentpkg "google.golang.org/adk/agent"
	adksession "google.golang.org/adk/session"
	"google.golang.org/genai"
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

	chatAgent storydenagent.Agent,

	mux *http.ServeMux,

	co *origin.Middleware,
	lo *reqlog.Middleware,
	cj *session_cookie.Jar,
	rl *limiter.Middleware,
) {
	if cfg.LanguageModelProvider == "" || chatAgent == nil {
		return
	}

	handler := newChatHandler(logger, chatAgent)

	applied := httpserver.Apply(handler,
		co.WithCORS(),
		lo.WithLogger(),
		cj.WithAuth(),
		rl.WithRequestSizeLimiter(),
		rl.WithRateLimit(),
	)

	lc.Append(fx.StartHook(func() error {
		mux.Handle("/sse/chat", applied)
		return nil
	}))
}

type chatRequest struct {
	ID        string        `json:"id"`
	ThreadID  string        `json:"threadId"`
	SessionID string        `json:"sessionId"`
	Messages  []chatMessage `json:"messages"`
	Data      any           `json:"data"`
}

type chatMessage struct {
	ID       string          `json:"id"`
	Role     string          `json:"role"`
	Parts    []chatPart      `json:"parts"`
	Metadata json.RawMessage `json:"metadata"`
}

type chatPart struct {
	Type   string          `json:"type"`
	Text   string          `json:"text,omitempty"`
	Delta  string          `json:"delta,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
	State  string          `json:"state,omitempty"`
	Source json.RawMessage `json:"source,omitempty"`
}

func newChatHandler(logger *slog.Logger, chatAgent storydenagent.Agent) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()

		accountID, err := session.GetAccountID(ctx)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		var req chatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("sse chat decode", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		prompt := lastUserMessage(req.Messages)
		if prompt == "" {
			http.Error(w, "missing user message", http.StatusBadRequest)
			return
		}

		sessionID := firstNonEmpty(req.SessionID, req.ThreadID, req.ID)
		if sessionID == "" {
			sessionID = fmt.Sprintf("chat-%s", accountID.String())
		}

		userContent := &genai.Content{
			Role:  genai.RoleUser,
			Parts: []*genai.Part{{Text: prompt}},
		}

		stream := chatAgent.Run(ctx, accountID.String(), sessionID, userContent, agentpkg.RunConfig{})

		emitter, err := newStreamEmitter(w)
		if err != nil {
			logger.Error("sse chat flusher", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer emitter.Done()

		responseID := uuid.NewString()
		textID := "text-1"

		if err := emitter.Headers(); err != nil {
			logger.Error("sse chat headers", slog.String("error", err.Error()))
			return
		}

		if err := emitter.Send(map[string]any{"type": "start", "messageId": responseID}); err != nil {
			logger.Debug("sse chat start", slog.String("error", err.Error()))
			return
		}
		if err := emitter.Send(map[string]any{"type": "start-step"}); err != nil {
			return
		}
		if err := emitter.Send(map[string]any{"type": "text-start", "id": textID}); err != nil {
			return
		}

		finishReason := defaultFinishReason
		finalSeen := false

		for event, streamErr := range stream {
			if streamErr != nil {
				_ = emitter.Send(map[string]any{"type": "error", "errorText": streamErr.Error()})
				return
			}

			if ctx.Err() != nil {
				return
			}

			sendTextChunks(event, textID, emitter)

			if event != nil && event.IsFinalResponse() {
				finalSeen = true
				if fr := strings.TrimSpace(string(event.LLMResponse.FinishReason)); fr != "" {
					finishReason = fr
				}
			}
		}

		if !finalSeen {
			finishReason = defaultFinishReason
		}

		_ = emitter.Send(map[string]any{"type": "text-end", "id": textID})
		_ = emitter.Send(map[string]any{"type": "finish-step"})
		_ = emitter.Send(map[string]any{"type": "finish", "finishReason": finishReason})
	})
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

func sendTextChunks(event *adksession.Event, textID string, emitter *streamEmitter) {
	if event == nil || event.LLMResponse.Content == nil {
		return
	}
	for _, part := range event.LLMResponse.Content.Parts {
		if part == nil || strings.TrimSpace(part.Text) == "" {
			continue
		}
		if err := emitter.Send(map[string]any{
			"type":  "text-delta",
			"id":    textID,
			"delta": part.Text,
		}); err != nil {
			return
		}
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

func (s *streamEmitter) Send(payload any) error {
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
