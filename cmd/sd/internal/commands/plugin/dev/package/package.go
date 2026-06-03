package packagecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/pluginapi"
)

type PackageCommand *cobra.Command

func New() PackageCommand {
	var dir string
	var manifestPath string
	var outputPath string
	var force bool

	command := &cobra.Command{
		Use:   "package",
		Short: "Create a supervised plugin package zip",
		Long: `# Package Plugin

Create a zip archive for supervised plugin distribution.

The package includes a generated ` + "`manifest.json`" + ` from ` + "`manifest.yaml`" + ` and the project files in the plugin directory.
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := outputPath
			if path == "" {
				mf, err := pluginapi.ReadProjectManifest(dir, manifestPath)
				if err != nil {
					return err
				}
				path = pluginapi.DefaultPackagePath(dir, mf.Manifest)
			}

			excludes := []string{path}
			pkg, err := pluginapi.BuildPackage(cmd.Context(), dir, manifestPath, excludes...)
			if err != nil {
				return err
			}

			if err := pluginapi.WritePackageFile(path, pkg, force); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Created plugin package %s (%d files)\n", path, len(pkg.Files))
			return nil
		},
	}

	command.Flags().StringVar(&dir, "dir", ".", "Plugin project directory")
	command.Flags().StringVarP(&manifestPath, "manifest", "m", pluginapi.ManifestFilename, "Path to plugin manifest YAML")
	command.Flags().StringVarP(&outputPath, "output", "o", "", "Output zip path")
	command.Flags().BoolVar(&force, "force", false, "Overwrite an existing package")

	help.SetupMarkdownHelp(command)

	return PackageCommand(command)
}
