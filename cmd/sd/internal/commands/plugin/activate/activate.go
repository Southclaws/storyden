package activate

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligenconv"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

func New(store *config.Store) cligen.PluginActivateHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.PluginActivateParams) (cligen.Plugin, error) {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return cligen.Plugin{}, err
		}
		plugin, err := plugindev.SetActiveState(ctx, client.OpenAPI, p.PluginInstanceId, openapi.PluginActiveStateActive)
		if err != nil {
			return cligen.Plugin{}, err
		}
		return cligenconv.Convert[cligen.Plugin](plugin)
	}
}
