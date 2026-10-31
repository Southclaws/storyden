package model_debug_dump

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/rs/xid"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/model"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
)

type Writer struct {
	enabled bool
	store   object.Storer
	logger  *slog.Logger
}

type snapshot struct {
	sessionID    string
	invocationID string
	agentName    string
	request      *model.LLMRequest
}

func New(cfg config.Config, store object.Storer, logger *slog.Logger) *Writer {
	return &Writer{enabled: cfg.RobotModelDebugDump, store: store, logger: logger}
}

func (w *Writer) BeforeModel(ctx agent.Context, request *model.LLMRequest) (*model.LLMResponse, error) {
	if !w.enabled {
		return nil, nil
	}
	err := w.write(ctx, snapshot{
		sessionID:    ctx.SessionID(),
		invocationID: ctx.InvocationID(),
		agentName:    ctx.AgentName(),
		request:      request,
	})
	if err != nil {
		w.logger.Warn("failed to write Robot model debug dump",
			slog.String("session", ctx.SessionID()),
			slog.String("invocation", ctx.InvocationID()),
			slog.Any("error", err),
		)
	}
	return nil, nil
}

func (w *Writer) write(ctx context.Context, input snapshot) error {
	now := time.Now().UTC()
	content, err := render(input, now)
	if err != nil {
		return err
	}
	session := strings.ReplaceAll(url.PathEscape(input.sessionID), ".", "%2E")
	if session == "" {
		return fmt.Errorf("model debug dump requires a session ID")
	}
	key := path.Join("debug/model-requests", session, now.Format("20060102T150405.000000000Z")+"-"+xid.New().String()+".md")
	if err := w.store.Write(ctx, key, strings.NewReader(content), int64(len(content))); err != nil {
		return err
	}
	w.logger.Debug("wrote Robot model debug dump", slog.String("key", key))
	return nil
}
