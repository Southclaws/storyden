package activate

import (
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

type ActivateCommand *cobra.Command

func New(store *config.Store) ActivateCommand {
	command := &cobra.Command{
		Use:   "activate <plugin-instance-id>",
		Short: "Start a supervised plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}
			plugin, err := plugindev.SetActiveState(cmd.Context(), client.OpenAPI, args[0], openapi.PluginActiveStateActive)
			if err != nil {
				return err
			}
			return output.JSON(cmd.OutOrStdout(), plugin)
		},
	}
	help.SetupMarkdownHelp(command)
	return ActivateCommand(command)
}
