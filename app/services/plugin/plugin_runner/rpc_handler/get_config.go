package rpc_handler

import (
	"context"

	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func (h *Handler) handleGetConfig(ctx context.Context, req *rpc.RPCRequestGetConfig) (rpc.RPCResponseGetConfig, error) {
	config, err := h.pluginReader.GetConfig(ctx, h.installationID)
	if err != nil {
		return rpc.RPCResponseGetConfig{}, err
	}

	if params, ok := req.Params.Get(); ok && len(params.Keys) > 0 {
		filtered := map[string]any{}
		for _, key := range params.Keys {
			value, exists := config[key]
			if !exists {
				continue
			}
			filtered[key] = value
		}
		config = filtered
	}

	return rpc.RPCResponseGetConfig{
		Method: "get_config",
		Config: config,
	}, nil
}
