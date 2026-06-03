package rotate

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/pluginapi"
)

type RotateCommand *cobra.Command

func New(store *config.Store) RotateCommand {
	command := &cobra.Command{
		Use:   "rotate <plugin-instance-id>",
		Short: "Regenerate and print an external plugin RPC token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}
			token, err := pluginapi.CycleToken(cmd.Context(), client.OpenAPI, args[0])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), token)
			return nil
		},
	}
	help.SetupMarkdownHelp(command)
	return RotateCommand(command)
}
