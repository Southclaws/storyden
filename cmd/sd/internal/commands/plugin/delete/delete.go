package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

func New(store *config.Store) cligen.PluginDeleteHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.PluginDeleteParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}
		if err := plugindev.DeletePlugin(ctx, client.OpenAPI, p.PluginInstanceId); err != nil {
			return err
		}
		fmt.Fprintf(io.Out, "Deleted plugin %s\n", p.PluginInstanceId)
		return nil
	}
}
