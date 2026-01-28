package rpc

import (
	"context"
	"log/slog"
	"sync"

	"github.com/coder/websocket"

	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
)

type wsDuplex struct {
	conn *websocket.Conn
	log  *slog.Logger
	mu   sync.Mutex
}

func NewWSDuplex(conn *websocket.Conn, log *slog.Logger) plugin_runner.Duplex {
	return &wsDuplex{conn: conn, log: log}
}

func (d *wsDuplex) Send(ctx context.Context, b []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.conn.Write(ctx, websocket.MessageText, b)
}

func (d *wsDuplex) Recv(ctx context.Context) ([]byte, error) {
	for {
		mt, b, err := d.conn.Read(ctx)
		if err != nil {
			return nil, err
		}
		if mt == websocket.MessageText {
			return b, nil
		}
	}
}

func (d *wsDuplex) Close() error {
	return d.conn.Close(websocket.StatusNormalClosure, "closing")
}
