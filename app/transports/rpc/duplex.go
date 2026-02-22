package rpc

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/coder/websocket"

	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/duplex"
)

type wsDuplex struct {
	conn *websocket.Conn
	log  *slog.Logger
	mu   sync.Mutex
}

func NewWSDuplex(conn *websocket.Conn, log *slog.Logger) duplex.Duplex {
	return &wsDuplex{conn: conn, log: log}
}

func (d *wsDuplex) Send(ctx context.Context, b []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.conn.Write(ctx, websocket.MessageText, b); err != nil {
		return mapDuplexError(err)
	}
	return nil
}

func (d *wsDuplex) Recv(ctx context.Context) ([]byte, error) {
	for {
		mt, b, err := d.conn.Read(ctx)
		if err != nil {
			return nil, mapDuplexError(err)
		}
		if mt == websocket.MessageText {
			return b, nil
		}
	}
}

func (d *wsDuplex) Close(cause error) error {
	code, reason := causeToWebSocketClose(cause)
	return d.conn.Close(code, reason)
}

func mapDuplexError(err error) error {
	var closeErr websocket.CloseError
	if errors.As(err, &closeErr) {
		return duplex.NewError(
			websocketStatusToErrorKind(closeErr.Code),
			closeErr.Reason,
		)
	}

	switch websocket.CloseStatus(err) {
	case websocket.StatusNormalClosure:
		return duplex.NewError(duplex.ErrClosed, "")
	case websocket.StatusGoingAway:
		return duplex.NewError(duplex.ErrClosed, "")
	case websocket.StatusNoStatusRcvd:
		return duplex.NewError(duplex.ErrClosed, "")
	case websocket.StatusPolicyViolation:
		return duplex.NewError(duplex.ErrRejected, "")
	case websocket.StatusInternalError:
		return duplex.NewError(duplex.ErrFailed, "")
	case websocket.StatusServiceRestart:
		return duplex.NewError(duplex.ErrUnavailable, "")
	case websocket.StatusTryAgainLater:
		return duplex.NewError(duplex.ErrUnavailable, "")
	}

	return err
}

func websocketStatusToErrorKind(code websocket.StatusCode) error {
	switch code {
	case websocket.StatusNormalClosure:
		return duplex.ErrClosed
	case websocket.StatusGoingAway:
		return duplex.ErrClosed
	case websocket.StatusNoStatusRcvd:
		return duplex.ErrClosed
	case websocket.StatusPolicyViolation:
		return duplex.ErrRejected
	case websocket.StatusServiceRestart:
		return duplex.ErrUnavailable
	case websocket.StatusTryAgainLater:
		return duplex.ErrUnavailable
	default:
		return duplex.ErrFailed
	}
}

func causeToWebSocketClose(cause error) (websocket.StatusCode, string) {
	if cause == nil || errors.Is(cause, duplex.ErrClosed) {
		return websocket.StatusNormalClosure, "closed"
	}

	if errors.Is(cause, duplex.ErrRejected) {
		return websocket.StatusPolicyViolation, errorReason(cause, "rejected")
	}

	if errors.Is(cause, duplex.ErrUnavailable) {
		return websocket.StatusTryAgainLater, errorReason(cause, "unavailable")
	}

	return websocket.StatusInternalError, errorReason(cause, "internal error")
}

func errorReason(err error, fallback string) string {
	var de duplex.Error
	if errors.As(err, &de) && de.Reason != "" {
		return de.Reason
	}

	if err == nil {
		return fallback
	}

	msg := err.Error()
	if msg == "" {
		return fallback
	}

	return msg
}
