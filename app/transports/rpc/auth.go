package rpc

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_auth"
)

func (h *WebSocketHandler) authenticateToken(ctx context.Context, params plugin_auth.ConnectionParams) (plugin.InstallationID, error) {
	if plugin_auth.IsExternalToken(params.Token) {
		rec, err := h.pluginReader.GetByExternalToken(ctx, params.Token)
		if err != nil {
			return plugin.InstallationID{}, plugin_auth.ErrInvalidToken
		}

		if suppliedID, ok := params.PluginID.Get(); ok && suppliedID != rec.InstallationID {
			return plugin.InstallationID{}, plugin_auth.ErrInvalidToken
		}

		return rec.InstallationID, nil
	}

	pluginID, ok := params.PluginID.Get()
	if !ok {
		return plugin.InstallationID{}, plugin_auth.ErrInvalidToken
	}

	secret, err := h.pluginReader.GetAuthSecret(ctx, pluginID)
	if err != nil {
		return plugin.InstallationID{}, err
	}

	// Cycle the secret immediately for supervised plugin connection tokens.
	if _, err := h.pluginWriter.CycleAuthSecret(ctx, pluginID); err != nil {
		return plugin.InstallationID{}, err
	}

	decryptedPluginID, err := plugin_auth.OpenToken(params.Token, secret)
	if err != nil {
		return plugin.InstallationID{}, err
	}

	if decryptedPluginID != pluginID {
		return plugin.InstallationID{}, plugin_auth.ErrInvalidToken
	}

	return pluginID, nil
}
