package rpc

import (
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
)

type WebSocketHandler struct {
	logger       *slog.Logger
	runner       plugin_runner.Runner
	pluginReader *plugin_reader.Reader
	pluginWriter *plugin_writer.Writer
}

func NewWebSocketHandler(
	logger *slog.Logger,
	runner plugin_runner.Runner,
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

	ctx := r.Context()

	params, err := plugin_auth.ParseConnectionURL(r.URL)
	if err != nil {
		return fault.Wrap(err,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.With("failed to parse connection URL"))
	}

	if err := h.authenticateToken(ctx, *params); err != nil {
		return fault.Wrap(err,
			fctx.With(ctx),
			ftag.With(ftag.Unauthenticated),
			fmsg.With("token authentication failed"))
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		CompressionMode: websocket.CompressionContextTakeover,
	})
	if err != nil {
		return fault.Wrap(err,
			fctx.With(ctx),
			fmsg.With("failed to accept websocket connection"))
	}
	defer conn.CloseNow()

	h.logger.Debug("plugin authenticated and connected", slog.String("plugin_id", params.PluginID.String()))

	session, err := h.runner.GetSession(ctx, params.PluginID)
	if err != nil {
		return fault.Wrap(err,
			fctx.With(ctx),
			fmsg.With("failed to get plugin session"))
	}

	duplex := NewWSDuplex(conn, h.logger)

	err = session.Connect(ctx, duplex)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	select {
	case <-ctx.Done():
		// TODO: Return a done channel from .Connect()
		// case <-done:
	}

	return nil
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
