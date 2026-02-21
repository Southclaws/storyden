package storyden

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coder/websocket"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_auth"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Plugin struct {
	rpcURL    *url.URL
	logger    *slog.Logger
	startTime time.Time
	mode      plugin.PluginMode

	pending  *xsync.Map[xid.ID, chan *rpc.PluginToHostResponse]
	handlers *xsync.Map[string, EventHandler]

	configureHandlerMu sync.RWMutex
	configureHandler   ConfigureHandler

	ctx    context.Context
	cancel context.CancelFunc

	runStarted   atomic.Bool
	shutdownOnce sync.Once
	shutdownErr  error

	connStateMu sync.RWMutex
	connState   *sdkConnectionState

	loopMu sync.Mutex
	loop   chan struct{}
}

type sdkConnectionState struct {
	conn     *websocket.Conn
	outbound chan outboundWrite
	done     <-chan struct{}
}

type outboundWrite struct {
	data   []byte
	result chan error
}

type EventHandler func(context.Context, rpc.EventPayload) error
type ConfigureHandler func(context.Context, map[string]any) error

const (
	initialReconnectWait = 250 * time.Millisecond
	maxReconnectWait     = 10 * time.Second
	defaultRPCTimeout    = 30 * time.Second
)

func New(ctx context.Context) (*Plugin, error) {
	urlString := os.Getenv("STORYDEN_RPC_URL")
	if urlString == "" {
		return nil, fmt.Errorf("STORYDEN_RPC_URL environment variable is not set")
	}
	rpcURL, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	pluginCtx, cancel := context.WithCancel(ctx)

	p := &Plugin{
		rpcURL: rpcURL,
		mode:   modeFromRPCURL(rpcURL),

		pending:  xsync.NewMap[xid.ID, chan *rpc.PluginToHostResponse](),
		handlers: xsync.NewMap[string, EventHandler](),

		ctx:    pluginCtx,
		cancel: cancel,
		logger: slog.Default(),
	}

	return p, nil
}

func (p *Plugin) On(eventType string, handler EventHandler) {
	p.handlers.Store(eventType, handler)
	p.logger.Debug("register handler", slog.String("eventType", eventType))
}

func (p *Plugin) OnConfigure(handler ConfigureHandler) {
	p.configureHandlerMu.Lock()
	p.configureHandler = handler
	p.configureHandlerMu.Unlock()
	p.logger.Debug("register configure handler")
}

// Run connects the plugin to the host and starts the WebSocket RPC read/write loops.
func (p *Plugin) Run(ctx context.Context) error {
	p.runStarted.Store(true)

	retryWait := initialReconnectWait

	for {
		conn, resp, err := websocket.Dial(ctx, p.rpcURL.String(), nil)
		if err != nil {
			if p.mode == plugin.PluginModeExternal && shouldRetryDial(err, resp) && p.ctx.Err() == nil && ctx.Err() == nil {
				p.logger.Warn("failed to connect, retrying",
					slog.String("rpc_url", p.rpcEndpointURL()),
					slog.String("error", p.sanitizeError(err)),
					slog.Duration("retry_in", retryWait))
				select {
				case <-ctx.Done():
					return p.Shutdown()
				case <-p.ctx.Done():
					return nil
				case <-time.After(retryWait):
				}

				retryWait *= 2
				if retryWait > maxReconnectWait {
					retryWait = maxReconnectWait
				}
				continue
			}
			return fmt.Errorf("failed to connect to %s: %s", p.rpcEndpointURL(), p.sanitizeError(err))
		}

		retryWait = initialReconnectWait

		connCtx, connCancel := context.WithCancel(p.ctx)
		outbound := make(chan outboundWrite)

		p.setConnState(&sdkConnectionState{
			conn:     conn,
			outbound: outbound,
			done:     connCtx.Done(),
		})

		if p.startTime.IsZero() {
			p.startTime = time.Now()
		}

		done := make(chan struct{})
		errCh := make(chan error, 2)

		p.loopMu.Lock()
		p.loop = done
		p.loopMu.Unlock()

		go p.readLoop(connCtx, conn, errCh)
		go p.writeLoop(connCtx, conn, outbound, errCh)

		var disconnectErr error
		receivedDisconnectErr := false
		var readErr error
		select {
		case <-ctx.Done():
			p.clearConnState()
			connCancel()
			conn.CloseNow()
			readErr = p.waitConnectionLoops(errCh, receivedDisconnectErr, disconnectErr)
			close(done)
			return p.Shutdown()
		case <-p.ctx.Done():
			p.clearConnState()
			connCancel()
			conn.CloseNow()
			_ = p.waitConnectionLoops(errCh, receivedDisconnectErr, disconnectErr)
			close(done)
			return nil
		case disconnectErr = <-errCh:
			receivedDisconnectErr = true
		}

		p.clearConnState()
		connCancel()
		conn.CloseNow()

		readErr = p.waitConnectionLoops(errCh, receivedDisconnectErr, disconnectErr)
		close(done)

		p.loopMu.Lock()
		if p.loop == done {
			p.loop = nil
		}
		p.loopMu.Unlock()

		p.clearPending()

		if p.mode == plugin.PluginModeSupervised {
			return nil
		}

		if !shouldRetryDisconnect(readErr) {
			return nil
		}

		if ctx.Err() != nil || p.ctx.Err() != nil {
			return nil
		}

		p.logger.Warn("connection dropped, reconnecting",
			slog.String("rpc_url", p.rpcEndpointURL()),
			slog.String("error", p.sanitizeError(readErr)),
			slog.Duration("retry_in", retryWait))

		select {
		case <-ctx.Done():
			return p.Shutdown()
		case <-p.ctx.Done():
			return nil
		case <-time.After(retryWait):
		}

		retryWait *= 2
		if retryWait > maxReconnectWait {
			retryWait = maxReconnectWait
		}
	}
}

func (p *Plugin) rpcEndpointURL() string {
	if p.rpcURL == nil {
		return ""
	}

	u := *p.rpcURL
	u.RawQuery = ""
	u.Fragment = ""
	return u.String()
}

func (p *Plugin) rpcDialURLWithQuery() string {
	if p.rpcURL == nil {
		return ""
	}

	u := *p.rpcURL
	switch u.Scheme {
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
	}

	return u.String()
}

func (p *Plugin) rpcDialEndpointURL() string {
	if p.rpcURL == nil {
		return ""
	}

	u := *p.rpcURL
	switch u.Scheme {
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
	}
	u.RawQuery = ""
	u.Fragment = ""

	return u.String()
}

func (p *Plugin) sanitizeError(err error) string {
	if err == nil {
		return ""
	}

	message := err.Error()

	if wsWithQuery := p.rpcURL.String(); wsWithQuery != "" {
		message = strings.ReplaceAll(message, wsWithQuery, p.rpcEndpointURL())
	}
	if httpWithQuery := p.rpcDialURLWithQuery(); httpWithQuery != "" {
		message = strings.ReplaceAll(message, httpWithQuery, p.rpcDialEndpointURL())
	}

	return message
}

func (p *Plugin) Shutdown() error {
	p.shutdownOnce.Do(func() {
		p.logger.Info("plugin shutting down")

		p.cancel()
		p.clearPending()

		conn := p.clearConnState()
		if conn != nil {
			conn.CloseNow()
		}

		// If Run/readLoop never started, there is nothing to wait for.
		if !p.runStarted.Load() {
			return
		}

		p.loopMu.Lock()
		loop := p.loop
		p.loopMu.Unlock()

		if loop == nil {
			return
		}

		select {
		case <-loop:
		case <-time.After(5 * time.Second):
			if p.shutdownErr == nil {
				p.shutdownErr = fmt.Errorf("timeout waiting for read loop to stop")
			}
		}
	})

	return p.shutdownErr
}

func (p *Plugin) clearPending() {
	p.pending.Range(func(id xid.ID, ch chan *rpc.PluginToHostResponse) bool {
		select {
		case ch <- nil:
		default:
		}
		p.pending.Delete(id)
		return true
	})
}

func (p *Plugin) readLoop(ctx context.Context, conn *websocket.Conn, errCh chan<- error) {
	var loopErr error
	defer func() {
		errCh <- loopErr
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, message, err := conn.Read(ctx)
		if err != nil {
			if ctx.Err() != nil || p.ctx.Err() != nil {
				return
			}

			loopErr = err
			return
		}

		if err := p.handleMessage(message); err != nil {
			p.logger.Error("failed to handle message",
				slog.String("error", err.Error()))
		}
	}
}

func (p *Plugin) writeLoop(ctx context.Context, conn *websocket.Conn, outbound <-chan outboundWrite, errCh chan<- error) {
	var loopErr error
	defer func() {
		errCh <- loopErr
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case req, ok := <-outbound:
			if !ok {
				return
			}

			if err := conn.Write(ctx, websocket.MessageText, req.data); err != nil {
				loopErr = err
				req.complete(fmt.Errorf("failed to write message: %w", err))
				return
			}

			req.complete(nil)
		}
	}
}

func (w outboundWrite) complete(err error) {
	select {
	case w.result <- err:
	default:
	}
}

func (p *Plugin) waitConnectionLoops(errCh <-chan error, hasFirst bool, firstErr error) error {
	err := firstErr
	if !hasFirst {
		err = <-errCh
	}
	other := <-errCh
	if err == nil {
		err = other
	}
	return err
}

func (p *Plugin) setConnState(state *sdkConnectionState) {
	p.connStateMu.Lock()
	p.connState = state
	p.connStateMu.Unlock()
}

func (p *Plugin) clearConnState() *websocket.Conn {
	p.connStateMu.Lock()
	defer p.connStateMu.Unlock()

	if p.connState == nil {
		return nil
	}

	conn := p.connState.conn
	p.connState = nil
	return conn
}

func (p *Plugin) getConnState() *sdkConnectionState {
	p.connStateMu.RLock()
	defer p.connStateMu.RUnlock()
	return p.connState
}

func (p *Plugin) getConn() *websocket.Conn {
	state := p.getConnState()
	if state == nil {
		return nil
	}
	return state.conn
}

func (p *Plugin) enqueueWrite(ctx context.Context, data []byte) error {
	state := p.getConnState()
	if state == nil {
		return fmt.Errorf("connection closed")
	}

	req := outboundWrite{
		data:   data,
		result: make(chan error, 1),
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-p.ctx.Done():
		return fmt.Errorf("plugin shutting down")
	case <-state.done:
		return fmt.Errorf("connection closed")
	case state.outbound <- req:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-p.ctx.Done():
		return fmt.Errorf("plugin shutting down")
	case <-state.done:
		return fmt.Errorf("connection closed")
	case err := <-req.result:
		return err
	}
}

func (p *Plugin) Send(ctx context.Context, payload rpc.PluginToHostRequestUnion) (rpc.PluginToHostResponseUnionUnion, error) {
	id, data, err := marshalRequestWithGeneratedID(payload)
	if err != nil {
		return nil, err
	}

	reqCtx, cancel := withDefaultTimeout(ctx)
	defer cancel()

	respch := make(chan *rpc.PluginToHostResponse, 1)
	p.pending.Store(id, respch)
	defer p.pending.Delete(id)

	if err := p.enqueueWrite(reqCtx, data); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	var response *rpc.PluginToHostResponse
	select {
	case <-reqCtx.Done():
		return nil, reqCtx.Err()
	case <-p.ctx.Done():
		return nil, fmt.Errorf("plugin shutting down")
	case response = <-respch:
		if response == nil {
			return nil, fmt.Errorf("connection closed")
		}
	}

	if rpcErr, ok := response.Error.Get(); ok {
		if msg, ok := rpcErr.Message.Get(); ok {
			return nil, fmt.Errorf("rpc error: %s", msg)
		}
		return nil, fmt.Errorf("rpc error")
	}

	return response.Result.PluginToHostResponseUnionUnion, nil
}

func withDefaultTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, defaultRPCTimeout)
}

func marshalRequestWithGeneratedID(payload rpc.PluginToHostRequestUnion) (xid.ID, []byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return xid.NilID(), nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	var body map[string]json.RawMessage
	if err := json.Unmarshal(data, &body); err != nil {
		return xid.NilID(), nil, fmt.Errorf("failed to decode request object: %w", err)
	}

	if len(body) == 0 {
		return xid.NilID(), nil, fmt.Errorf("request payload cannot be empty")
	}

	id := xid.New()
	body["id"] = json.RawMessage(strconv.Quote(id.String()))

	data, err = json.Marshal(body)
	if err != nil {
		return xid.NilID(), nil, fmt.Errorf("failed to re-encode request: %w", err)
	}

	return id, data, nil
}

func modeFromRPCURL(rpcURL *url.URL) plugin.PluginMode {
	token := rpcURL.Query().Get("token")
	if strings.HasPrefix(token, plugin_auth.ExternalTokenPrefix) {
		return plugin.PluginModeExternal
	}
	return plugin.PluginModeSupervised
}

func shouldRetryDial(err error, resp *http.Response) bool {
	if resp != nil && resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return false
	}

	status := websocket.CloseStatus(err)
	if status == websocket.StatusPolicyViolation || status == websocket.StatusUnsupportedData {
		return false
	}

	return true
}

func shouldRetryDisconnect(err error) bool {
	if err == nil {
		return false
	}

	status := websocket.CloseStatus(err)
	switch status {
	case -1:
		return true
	case websocket.StatusGoingAway:
		return true
	case websocket.StatusAbnormalClosure:
		return true
	case websocket.StatusInternalError:
		return true
	case websocket.StatusServiceRestart:
		return true
	case websocket.StatusTryAgainLater:
		return true
	default:
		return false
	}
}
