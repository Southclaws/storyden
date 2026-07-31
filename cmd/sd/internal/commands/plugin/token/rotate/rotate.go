package rotate

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

func New(store *config.Store) cligen.PluginTokenRotateHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.PluginTokenRotateParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}
		token, err := plugindev.CycleToken(ctx, client.OpenAPI, p.PluginInstanceId)
		if err != nil {
			return err
		}
		fmt.Fprintln(io.Out, token)
		return nil
	}
}
