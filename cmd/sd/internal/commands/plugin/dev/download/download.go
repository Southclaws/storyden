package download

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

type DownloadCommand *cobra.Command

func New(store *config.Store) DownloadCommand {
	var dir string
	var noUnzip bool
	var force bool

	command := &cobra.Command{
		Use:   "download <plugin-instance-id>",
		Short: "Download an installed plugin package",
		Long: `# Download Plugin Package

Download the package archive for an installed supervised plugin. By default the archive is unzipped into the current directory.

## Examples

Download and unzip into the current directory:
~~~bash
sd plugin dev download d8cg...
~~~

Download and unzip into a specific directory:
~~~bash
sd plugin dev download d8cg... --dir ./plugin-source
~~~

Download the zip without extracting it:
~~~bash
sd plugin dev download d8cg... --no-unzip
~~~
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			data, err := plugindev.DownloadPackage(cmd.Context(), client.OpenAPI, args[0])
			if err != nil {
				return err
			}

			if noUnzip {
				filename := plugindev.PackageFilename(cmd.Context(), data, args[0])
				target := filepath.Join(dir, filename)
				if err := plugindev.WritePackageFile(target, &plugindev.PackageArchive{Bytes: data}, force); err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Downloaded plugin package: %s\n", target)
				return nil
			}

			result, err := plugindev.ExtractPackageArchive(data, dir, force)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Extracted plugin package to %s (%d files)\n", dir, len(result.Files))
			return nil
		},
	}

	command.Flags().StringVar(&dir, "dir", ".", "Directory to unzip into or save the zip in")
	command.Flags().BoolVar(&noUnzip, "no-unzip", false, "Save the package zip without extracting it")
	command.Flags().BoolVar(&force, "force", false, "Overwrite existing files")

	help.SetupMarkdownHelp(command)

	return DownloadCommand(command)
}
