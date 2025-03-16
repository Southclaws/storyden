package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/google/uuid"
	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_traversal"
	"github.com/Southclaws/storyden/app/transports/http/mcp/mcp_schema"
	"github.com/Southclaws/storyden/app/transports/http/mcp/resources"
	"github.com/Southclaws/storyden/app/transports/http/mcp/tools"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/queue"
)

type Handler = http.Handler

func Build() fx.Option {
	return fx.Options(
		fx.Provide(newHandler),
		fx.Provide(
			resources.New,
			tools.New,
		),
	)
}

var errNotImplemented = fault.New("not implemented", ftag.With(ftag.NotFound))

type SessionID = uuid.UUID

type mcp struct {
	logger      *zap.Logger
	config      config.Config
	sessions    *xsync.MapOf[SessionID, *session]
	qf          *queue.QueueFactory
	nodeLister  node_traversal.Repository
	nodeQuerier *node_querier.Querier

	// providers
	toolsProvider     *tools.Provider
	resourcesProvider *resources.Provider
}

type session struct {
	id    SessionID
	topic pubsub.Topic[Message]
}

func newHandler(
	logger *zap.Logger,
	config config.Config,
	qf *queue.QueueFactory,
	nodeLister node_traversal.Repository,
	nodeQuerier *node_querier.Querier,
	toolsProvider *tools.Provider,
	resourcesProvider *resources.Provider,
) Handler {
	mux := http.NewServeMux()

	mcp := &mcp{
		logger:      logger,
		config:      config,
		sessions:    xsync.NewMapOf[SessionID, *session](),
		qf:          qf,
		nodeLister:  nodeLister,
		nodeQuerier: nodeQuerier,

		toolsProvider:     toolsProvider,
		resourcesProvider: resourcesProvider,
	}

	mux.HandleFunc("/sse", mcp.sseHandler)
	mux.HandleFunc("/messages", mcp.jsonRPCHandler)

	return mux
}

type Message mcp_schema.FixedJSONRPCResponse

func (m *mcp) sseHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := SessionID(uuid.New())
	topicName := fmt.Sprintf("mcp_%s", sessionID.String())
	topic := queue.NewNamed[Message](m.qf, topicName)

	session := &session{
		id:    sessionID,
		topic: topic,
	}

	m.sessions.Store(sessionID, session)

	mc, err := topic.Subscribe(context.TODO())
	if err != nil {
		http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	endpoint := m.config.PublicAPIAddress.JoinPath("/mcp/messages").String()

	fmt.Fprintf(w, "event: endpoint\ndata: %s?sessionId=%s\n\n", endpoint, sessionID)
	w.(http.Flusher).Flush()

	for {
		select {
		case <-r.Context().Done():
			m.sessions.Delete(sessionID)
			return

		case msg, ok := <-mc:
			if !ok {
				return
			}

			b, err := json.Marshal(msg.Payload)
			if err != nil {
				msg.Ack()
				http.Error(w, "Failed to marshal", http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "event: message\ndata: %s\n\n", string(b))
			w.(http.Flusher).Flush()
			msg.Ack()
		}
	}

}

func (m *mcp) jsonRPCHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("sessionId"))
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	sessionID := SessionID(id)

	session, ok := m.sessions.Load(sessionID)
	if !ok {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	reusableBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var req mcp_schema.JSONRPCRequest
	if err := json.Unmarshal(reusableBody, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("received method request from %s with %d method: %s\n", sessionID.String(), req.Id, req.Method)

	result, err := m.requestHandler(r.Context(), req.Method, reusableBody)
	if err != nil {
		// TODO: handle error properly lol
		http.Error(w, "Failed to handle request", http.StatusInternalServerError)
		return
	}

	res := mcp_schema.FixedJSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.Id,
		Result:  result,
	}

	err = session.topic.Publish(r.Context(), Message(res))
	if err != nil {
		http.Error(w, "Failed to publish", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (m *mcp) requestHandler(ctx context.Context, method string, body []byte) (any, error) {
	switch method {
	case "initialize":
		return m.buildInitializeResponse(), nil

	case "notifications/initialized":
		return nil, nil

	case "ping":
		return map[string]string{}, nil

	case "completion/complete":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "logging/setLevel":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "notifications/cancelled":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "notifications/message":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "notifications/progress":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "notifications/prompts/list_changed":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "notifications/resources/list_changed":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "notifications/resources/updated":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "notifications/roots/list_changed":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "notifications/tools/list_changed":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "prompts/get":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "prompts/list":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "resources/list":
		r, err := m.resourcesProvider.ListResources(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		return r, nil

	case "resources/read":
		var req mcp_schema.ReadResourceRequest
		if err := json.Unmarshal(body, &req); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		return m.resourcesProvider.ReadResource(ctx, req)

	case "resources/subscribe":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "resources/templates/list":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "resources/unsubscribe":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "roots/list":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "sampling/createMessage":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "tools/call":
		return nil, fault.Wrap(errNotImplemented, fctx.With(ctx))

	case "tools/list":
		var req mcp_schema.ListToolsRequest
		if err := json.Unmarshal(body, &req); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		return m.toolsProvider.ListTools(ctx, req)

	default:
		err := fault.Newf("method not found: %s", method)
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
}

func (m *mcp) buildInitializeResponse() mcp_schema.InitializeResult {
	return mcp_schema.InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: mcp_schema.ServerCapabilities{
			Prompts:   &mcp_schema.ServerCapabilitiesPrompts{},
			Resources: &mcp_schema.ServerCapabilitiesResources{},
			Tools:     &mcp_schema.ServerCapabilitiesTools{},
		},
		ServerInfo: mcp_schema.Implementation{
			Name:    "Storyden",
			Version: "rolling",
		},
	}
}
