package delete

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/pluginapi"
)

type DeleteCommand *cobra.Command

func New(store *config.Store) DeleteCommand {
	command := &cobra.Command{
		Use:   "delete <plugin-instance-id>",
		Short: "Delete a plugin from the current instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}
			if err := pluginapi.DeletePlugin(cmd.Context(), client.OpenAPI, args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Deleted plugin %s\n", args[0])
			return nil
		},
	}
	help.SetupMarkdownHelp(command)
	return DeleteCommand(command)
}
