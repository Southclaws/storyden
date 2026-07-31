package install

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

func New(store *config.Store) cligen.PluginDevInstallHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.PluginDevInstallParams) error {
		pkg, err := plugindev.BuildPackage(ctx, p.Dir, p.Manifest)
		if err != nil {
			return err
		}

		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		plugin, updated, err := plugindev.InstallSupervisedPackage(ctx, client.OpenAPI, pkg)
		if err != nil {
			return err
		}

		action := "Installed"
		if updated {
			action = "Updated"
		}
		fmt.Fprintf(io.Out, "%s supervised plugin %s (%s)\n", action, plugin.Name, plugin.Id)
		return nil
	}
}
