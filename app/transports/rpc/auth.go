package rpc

import (
	"context"

	"github.com/Southclaws/storyden/app/services/plugin/plugin_auth"
)

func (h *WebSocketHandler) authenticateToken(ctx context.Context, params plugin_auth.ConnectionParams) error {
	secret, err := h.pluginReader.GetAuthSecret(ctx, params.PluginID)
	if err != nil {
		return err
	}

	// Cycle the secret immediately. Plugin cannot reconnect.
	if _, err := h.pluginWriter.CycleAuthSecret(ctx, params.PluginID); err != nil {
		return err
	}

	decryptedPluginID, err := plugin_auth.OpenToken(params.Token, secret)
	if err != nil {
		return err
	}

	if decryptedPluginID != params.PluginID {
		return plugin_auth.ErrInvalidToken
	}

	return nil
}
