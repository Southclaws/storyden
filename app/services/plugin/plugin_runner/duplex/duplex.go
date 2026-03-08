package duplex

import "context"

// Duplex represents a bidirectional communication channel between the Host and
// an authenticated and connected plugin. This is implemented by a WebSocket.
type Duplex interface {
	Send(ctx context.Context, b []byte) error
	Recv(ctx context.Context) ([]byte, error)
	Close(cause error) error
}
