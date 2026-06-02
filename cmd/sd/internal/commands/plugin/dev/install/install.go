package install

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/pluginapi"
)

type InstallCommand *cobra.Command

func New(store *config.Store) InstallCommand {
	var dir string
	var manifestPath string

	command := &cobra.Command{
		Use:   "install",
		Short: "Install or update a local plugin as supervised",
		Long: `# Install Plugin

Package the local plugin project and install it as a supervised plugin on the current Storyden instance.

If a supervised plugin with the same manifest ID already exists, its package is replaced.
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := pluginapi.BuildPackage(cmd.Context(), dir, manifestPath)
			if err != nil {
				return err
			}

			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			plugin, updated, err := pluginapi.InstallSupervisedPackage(cmd.Context(), client.OpenAPI, pkg)
			if err != nil {
				return err
			}

			action := "Installed"
			if updated {
				action = "Updated"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s supervised plugin %s (%s)\n", action, plugin.Name, plugin.Id)
			return nil
		},
	}

	command.Flags().StringVar(&dir, "dir", ".", "Plugin project directory")
	command.Flags().StringVarP(&manifestPath, "manifest", "m", pluginapi.ManifestFilename, "Path to plugin manifest YAML")

	help.SetupMarkdownHelp(command)

	return InstallCommand(command)
}
