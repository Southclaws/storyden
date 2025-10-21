package storyden

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func (p *Plugin) handleMessage(message []byte) error {
	var req rpc.HostToPluginRequest

	if err := json.Unmarshal(message, &req); err != nil {
		var base rpc.JsonRpcResponse
		if err := json.Unmarshal(message, &base); err == nil && !base.ID.IsNil() {
			return p.handleResponse(message, base.ID)
		}
		return fmt.Errorf("invalid message format: %w", err)
	}

	if req.HostToPluginRequestUnion == nil {
		return fmt.Errorf("empty request")
	}

	switch r := req.HostToPluginRequestUnion.(type) {
	case *rpc.RPCRequestConfigure:
		return p.handleConfigure(*r)
	case *rpc.RPCRequestEvent:
		return p.handleEvent(*r)
	case *rpc.RPCRequestPing:
		return p.handlePing(*r)
	default:
		return fmt.Errorf("unknown request type: %T", r)
	}
}

func (p *Plugin) handleResponse(message []byte, id xid.ID) error {
	var resp rpc.PluginToHostResponse
	if err := json.Unmarshal(message, &resp); err != nil {
		return fmt.Errorf("invalid response format: %w", err)
	}

	ch, ok := p.pending.LoadAndDelete(id)

	if ok {
		select {
		case ch <- &resp:
		case <-time.After(time.Second):
			p.logger.Warn("timed out sending response",
				slog.String("id", id.String()))
		}
	} else {
		p.logger.Warn("received response for unknown pending request id",
			slog.String("id", id.String()))
	}

	return nil
}

func (p *Plugin) handleConfigure(req rpc.RPCRequestConfigure) error {
	p.configureHandlerMu.RLock()
	handler := p.configureHandler
	p.configureHandlerMu.RUnlock()

	go func() {
		ok := true
		if handler != nil {
			if err := handler(p.ctx, req.Params); err != nil {
				ok = false
				p.logger.Error("configure handler error",
					slog.String("error", err.Error()))
			}
		}

		if err := p.sendResponse(req.ID, rpc.HostToPluginResponseUnion{
			HostToPluginResponseUnionUnion: &rpc.RPCResponseConfigure{
				Method: opt.New("configure"),
				Ok:     ok,
			},
		}); err != nil {
			p.logger.Error("failed to send configure response",
				slog.String("error", err.Error()))
		}
	}()

	return nil
}

func (p *Plugin) handleEvent(req rpc.RPCRequestEvent) error {
	eventType := req.Params.EventPayloadType()

	handler, ok := p.handlers.Load(eventType)

	if !ok {
		p.logger.Warn("no handler for event",
			slog.String("event_type", eventType))

		return p.sendResponse(req.ID, rpc.HostToPluginResponseUnion{
			HostToPluginResponseUnionUnion: &rpc.RPCResponseEvent{
				Method: opt.New("event"),
				Ok:     true,
			},
		})
	}

	go func() {
		if err := handler(p.ctx, req.Params); err != nil {
			p.logger.Error("event handler error",
				slog.String("event_type", eventType),
				slog.String("error", err.Error()))

			if sendErr := p.sendErrorResponse(req.ID, -32000, fmt.Sprintf("handler error: %v", err)); sendErr != nil {
				p.logger.Error("failed to send event error response",
					slog.String("event_type", eventType),
					slog.String("error", sendErr.Error()))
			}
			return
		}

		if err := p.sendResponse(req.ID, rpc.HostToPluginResponseUnion{
			HostToPluginResponseUnionUnion: &rpc.RPCResponseEvent{
				Method: opt.New("event"),
				Ok:     true,
			},
		}); err != nil {
			p.logger.Error("failed to send event response",
				slog.String("event_type", eventType),
				slog.String("error", err.Error()))
		}
	}()

	return nil
}

func (p *Plugin) handlePing(req rpc.RPCRequestPing) error {
	uptime := time.Since(p.startTime).Seconds()

	return p.sendResponse(req.ID, rpc.HostToPluginResponseUnion{
		HostToPluginResponseUnionUnion: &rpc.RPCResponsePing{
			Method:        opt.New("ping"),
			Pong:          true,
			UptimeSeconds: opt.New(uptime),
			Status:        opt.New("healthy"),
		},
	})
}

func (p *Plugin) sendResponse(id xid.ID, result rpc.HostToPluginResponseUnion) error {
	resp := rpc.HostToPluginResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Result:  result,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	return p.enqueueWrite(p.ctx, data)
}

func (p *Plugin) sendErrorResponse(id xid.ID, code int, message string) error {
	data, err := json.Marshal(rpc.JsonRpcResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Error: opt.New(rpc.JsonRpcResponseError{
			Code:    opt.New(code),
			Message: opt.New(message),
		}),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal error response: %w", err)
	}

	return p.enqueueWrite(p.ctx, data)
}
