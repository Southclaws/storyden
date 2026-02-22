package rpc_handler

import (
	"context"
	"log/slog"
	"net/url"
	"sync"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/authentication/access_key"
	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/resources/plugin/plugin_reader"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Handler struct {
	installationID plugin.InstallationID
	logger         *slog.Logger
	manifest       *plugin.Validated
	apiBaseURL     url.URL
	accountQuerier *account_querier.Querier
	accountWriter  *account_writer.Writer
	accessKeys     *access_key.Repository
	pluginReader   *plugin_reader.Reader

	mu           sync.Mutex
	cachedAccess *rpc.RPCResponseAccessGetResult
}

func New(
	logger *slog.Logger,
	installationID plugin.InstallationID,
	manifest *plugin.Validated,
	apiBaseURL url.URL,
	accountQuerier *account_querier.Querier,
	accountWriter *account_writer.Writer,
	accessKeys *access_key.Repository,
	pluginReader *plugin_reader.Reader,
) *Handler {
	return &Handler{
		installationID: installationID,
		logger:         logger,
		manifest:       manifest,
		apiBaseURL:     apiBaseURL,
		accountQuerier: accountQuerier,
		accountWriter:  accountWriter,
		accessKeys:     accessKeys,
		pluginReader:   pluginReader,
	}
}

func (h *Handler) Handle(ctx context.Context, req rpc.PluginToHostRequestUnion) (*rpc.PluginToHostResponse, error) {
	switch v := req.(type) {
	case *rpc.RPCRequestAccessGet:
		result, err := h.handleAccessGet(ctx, v)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return &rpc.PluginToHostResponse{
			ID:      v.ID,
			Jsonrpc: "2.0",
			Result: rpc.PluginToHostResponseUnion{
				PluginToHostResponseUnionUnion: result,
			},
		}, nil

	case *rpc.RPCRequestGetConfig:
		result, err := h.handleGetConfig(ctx, v)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return &rpc.PluginToHostResponse{
			ID:      v.ID,
			Jsonrpc: "2.0",
			Result: rpc.PluginToHostResponseUnion{
				PluginToHostResponseUnionUnion: result,
			},
		}, nil

	default:
		h.logger.WarnContext(ctx, "unhandled RPC method type",
			slog.String("type", req.PluginToHostRequestType()),
		)

		return nil, fault.New("unhandled RPC method type")
	}
}

func (h *Handler) OnDisconnect() {
	h.mu.Lock()
	h.cachedAccess = nil
	h.mu.Unlock()
}
