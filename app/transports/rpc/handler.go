package rpc

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/coder/websocket"

	"github.com/Southclaws/storyden/app/resources/plugin/plugin_reader"
	"github.com/Southclaws/storyden/app/resources/plugin/plugin_writer"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_auth"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/duplex"
)

type WebSocketHandler struct {
	logger       *slog.Logger
	runner       plugin_runner.Host
	pluginReader *plugin_reader.Reader
	pluginWriter *plugin_writer.Writer
}

func NewWebSocketHandler(
	logger *slog.Logger,
	runner plugin_runner.Host,
	pluginReader *plugin_reader.Reader,
	pluginWriter *plugin_writer.Writer,
) *WebSocketHandler {
	return &WebSocketHandler{
		logger:       logger,
		runner:       runner,
		pluginReader: pluginReader,
		pluginWriter: pluginWriter,
	}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	if err := h.handleWebSocket(w, r); err != nil {
		h.logger.Error("websocket handler error", "error", err)

		switch ftag.Get(err) {
		case ftag.InvalidArgument:
			http.Error(w, "invalid connection parameters", http.StatusBadRequest)

		case ftag.Unauthenticated:
			http.Error(w, "unauthenticated", http.StatusUnauthorized)

		case ftag.PermissionDenied:
			http.Error(w, "permission denied", http.StatusForbidden)

		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func (h *WebSocketHandler) handleWebSocket(w http.ResponseWriter, r *http.Request) error {
	w, ok := GetFlusher(w)
	if !ok {
		panic("websocket handler requires http.Flusher")
	}

	requestCtx := r.Context()

	params, err := plugin_auth.ParseConnectionURL(r.URL)
	if err != nil {
		return fault.Wrap(err,
			fctx.With(requestCtx),
			ftag.With(ftag.InvalidArgument),
			fmsg.With("failed to parse connection URL"))
	}

	pluginID, err := h.authenticateToken(requestCtx, *params)
	if err != nil {
		return fault.Wrap(err,
			fctx.With(requestCtx),
			ftag.With(ftag.Unauthenticated),
			fmsg.With("token authentication failed"))
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		CompressionMode: websocket.CompressionContextTakeover,
	})
	if err != nil {
		return fault.Wrap(err,
			fctx.With(requestCtx),
			fmsg.With("failed to accept websocket connection"))
	}
	defer func() {
		_ = conn.Close(websocket.StatusNormalClosure, "closed")
	}()

	logger := h.logger.With(slog.String("plugin_id", pluginID.String()))
	logger.Debug("plugin authenticated and connected")
	logger.Debug("starting plugin websocket session")

	duplex := NewWSDuplex(conn, h.logger)
	// Preserve request-scoped values/metadata but detach from short-lived
	// cancellation semantics of upgraded HTTP requests.
	sessionCtx, sessionCancel := context.WithCancel(context.WithoutCancel(requestCtx))
	defer sessionCancel()

	err = h.runner.Connect(sessionCtx, pluginID, duplex)
	if err != nil {
		if isExpectedConnectDisconnect(err) {
			logger.Debug("plugin websocket disconnected",
				slog.Any("error", err),
				slog.String("close_status", websocket.CloseStatus(err).String()))
			return nil
		}

		logger.Warn("plugin connection rejected",
			slog.Any("error", err))
		_ = conn.Close(websocket.StatusPolicyViolation, "connection rejected")
		return nil
	}

	logger.Debug("plugin websocket session ended")
	return nil
}

func isExpectedConnectDisconnect(err error) bool {
	if err == nil {
		return true
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, io.EOF) {
		return true
	}
	if duplex.IsExpectedDisconnect(err) {
		return true
	}

	switch websocket.CloseStatus(err) {
	case websocket.StatusNormalClosure, websocket.StatusGoingAway, websocket.StatusNoStatusRcvd:
		return true
	default:
		return false
	}
}

func GetFlusher(w http.ResponseWriter) (http.ResponseWriter, bool) {
	for {
		if _, ok := w.(http.Flusher); ok {
			return w, true
		}
		// Try to unwrap
		if unwrapper, ok := w.(interface{ Unwrap() http.ResponseWriter }); ok {
			w = unwrapper.Unwrap()
		} else {
			return nil, false
		}
	}
}
