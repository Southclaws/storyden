package get

import (
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
	"github.com/Southclaws/storyden/cmd/sd/internal/pluginapi"
)

type GetCommand *cobra.Command

func New(store *config.Store) GetCommand {
	command := &cobra.Command{
		Use:     "get <plugin-instance-id>",
		Aliases: []string{"status"},
		Short:   "Get plugin status and manifest information",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}
			plugin, err := pluginapi.GetPlugin(cmd.Context(), client.OpenAPI, args[0])
			if err != nil {
				return err
			}
			return output.JSON(cmd.OutOrStdout(), plugin)
		},
	}
	help.SetupMarkdownHelp(command)
	return GetCommand(command)
}
