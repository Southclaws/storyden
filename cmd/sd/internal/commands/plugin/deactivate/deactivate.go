package deactivate

import (
	"context"

	"github.com/Southclaws/storyden/cmd/sd/internal/output"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

func New(store *config.Store) cligen.PluginDeactivateHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.PluginDeactivateParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}
		plugin, err := plugindev.SetActiveState(ctx, client.OpenAPI, p.PluginInstanceId, openapi.PluginActiveStateInactive)
		if err != nil {
			return err
		}
		return output.JSON(io.Out, plugin)
	}
}
