package deactivate

import (
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

type DeactivateCommand *cobra.Command

func New(store *config.Store) DeactivateCommand {
	command := &cobra.Command{
		Use:   "deactivate <plugin-instance-id>",
		Short: "Stop a supervised plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}
			plugin, err := plugindev.SetActiveState(cmd.Context(), client.OpenAPI, args[0], openapi.PluginActiveStateInactive)
			if err != nil {
				return err
			}
			return output.JSON(cmd.OutOrStdout(), plugin)
		},
	}
	help.SetupMarkdownHelp(command)
	return DeactivateCommand(command)
}
