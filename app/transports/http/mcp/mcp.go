package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/google/uuid"
	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_traversal"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/queue"
)

type Handler http.Handler

func Build() fx.Option {
	return fx.Provide(newHandler)
}

type SessionID = uuid.UUID

type mcp struct {
	logger     *zap.Logger
	address    url.URL
	qf         *queue.QueueFactory
	sessionMap *xsync.MapOf[SessionID, session]

	// resource queriers
	nodeLister  node_traversal.Repository
	nodeQuerier *node_querier.Querier
}

type session struct {
	id    SessionID
	topic pubsub.Topic[Message]
}

func newHandler(
	cfg config.Config,
	logger *zap.Logger,
	qf *queue.QueueFactory,
	nodeLister node_traversal.Repository,
	nodeQuerier *node_querier.Querier,
) Handler {
	mux := http.NewServeMux()

	mcp := &mcp{
		logger:      logger,
		address:     cfg.PublicAPIAddress,
		qf:          qf,
		sessionMap:  xsync.NewMapOf[SessionID, session](),
		nodeLister:  nodeLister,
		nodeQuerier: nodeQuerier,
	}

	mux.HandleFunc("/sse", mcp.sseHandler)
	mux.HandleFunc("/messages", mcp.jsonRPCHandler)

	return mux
}

type Message FixedJSONRPCResponse

func (m *mcp) sseHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := SessionID(uuid.New())
	topicName := fmt.Sprintf("mcp_%s", sessionID.String())
	topic := queue.NewNamed[Message](m.qf, topicName)

	session := session{
		id:    sessionID,
		topic: topic,
	}

	m.sessionMap.Store(sessionID, session)

	mc, err := topic.Subscribe(context.TODO())
	if err != nil {
		http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	endpoint := m.address.JoinPath("/mcp/messages").String()

	fmt.Fprintf(w, "event: endpoint\ndata: %s?sessionId=%s\n\n", endpoint, sessionID)
	w.(http.Flusher).Flush()

	for msg := range mc {
		b, err := json.Marshal(msg.Payload)
		if err != nil {
			msg.Ack()
			http.Error(w, "Failed to marshal", http.StatusInternalServerError)
			return
		}

		m.logger.Debug("sending response",
			zap.String("session_id", sessionID.String()),
			zap.Int("id", int(msg.Payload.ID)),
			zap.Any("message", msg.Payload),
		)

		fmt.Fprintf(w, "event: message\ndata: %s\n\n", string(b))
		w.(http.Flusher).Flush()
		msg.Ack()
	}

	m.sessionMap.Delete(sessionID)
}

func (m *mcp) jsonRPCHandler(w http.ResponseWriter, r *http.Request) {
	rawSessionID := r.URL.Query().Get("sessionId")
	id, err := uuid.Parse(rawSessionID)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	sessionID := SessionID(id)

	session, ok := m.sessionMap.Load(sessionID)
	if !ok {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	reusableBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}

	var req JSONRPCRequest
	if err := json.Unmarshal(reusableBody, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	res := FixedJSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.Id,
	}

	fmt.Printf("received method request from %s with %d method: %s\n", sessionID.String(), req.Id, req.Method)

	switch req.Method {

	case "initialize":
		res.Result = m.buildInitializeResponse()

	case "notifications/initialized":
		w.WriteHeader(http.StatusAccepted)
		return

	case "ping":
		res.Result = map[string]string{}

	case "completion/complete":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "logging/setLevel":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "notifications/cancelled":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "notifications/message":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "notifications/progress":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "notifications/prompts/list_changed":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "notifications/resources/list_changed":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "notifications/resources/updated":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "notifications/roots/list_changed":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "notifications/tools/list_changed":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "prompts/get":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "prompts/list":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "resources/list":
		resources, err := m.getResources(r.Context())
		if err != nil {
			http.Error(w, "Failed to list resources", http.StatusInternalServerError)
			return
		}
		res.Result = ListResourcesResult{
			Resources: resources,
		}

	case "resources/read":

		var req ReadResourceRequest
		if err := json.Unmarshal(reusableBody, &req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		nu, err := url.Parse(req.Params.Uri)
		if err != nil {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}
		_, nid := path.Split(nu.Path)

		n, err := m.nodeQuerier.Get(r.Context(), library.NewKey(nid))
		if err != nil {
			http.Error(w, "Failed to list resources", http.StatusInternalServerError)
			return
		}
		res.Result = ReadResourceResult{
			Contents: []any{
				TextResourceContents{
					Text:     n.Content.OrZero().HTML(),
					MimeType: opt.New("text/html").Ptr(),
					Uri:      req.Params.Uri,
				},
			},
		}

	case "resources/subscribe":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "resources/templates/list":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "resources/unsubscribe":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "roots/list":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "sampling/createMessage":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "tools/call":
		w.WriteHeader(http.StatusNotImplemented)
		return

	case "tools/list":
		w.WriteHeader(http.StatusNotImplemented)
		return

	default:
		http.Error(w, "Method not found", http.StatusNotFound)
		return
	}

	err = session.topic.Publish(r.Context(), Message(res))
	if err != nil {
		http.Error(w, "Failed to publish", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (m *mcp) buildInitializeResponse() InitializeResult {
	return InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: ServerCapabilities{
			Prompts:   &ServerCapabilitiesPrompts{},
			Resources: &ServerCapabilitiesResources{},
			Tools:     &ServerCapabilitiesTools{},
		},
		ServerInfo: Implementation{
			Name:    "Storyden",
			Version: "rolling",
		},
	}
}

func (m *mcp) getResources(ctx context.Context) ([]Resource, error) {
	nodes, err := m.nodeLister.Subtree(ctx, nil, true)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	resources := dt.Map(nodes, func(n *library.Node) Resource {
		return Resource{
			Annotations: &ResourceAnnotations{
				Audience: []Role{
					RoleAssistant,
					RoleUser,
				},
			},
			Description: n.Description.Ptr(),
			MimeType:    opt.New("application/json").Ptr(),
			Name:        n.Name,
			Uri:         m.address.JoinPath("/nodes/" + n.Mark.Slug()).String(),
		}
	})

	return resources, nil
}
