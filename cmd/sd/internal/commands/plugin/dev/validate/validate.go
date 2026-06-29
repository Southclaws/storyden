package validate

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

type ValidateCommand *cobra.Command

func New() ValidateCommand {
	var dir string
	var manifestPath string

	command := &cobra.Command{
		Use:   "validate",
		Short: "Validate a plugin manifest and supervised package",
		Long: `# Validate Plugin

Validate ` + "`manifest.yaml`" + ` and the supervised plugin package that would be generated from the current project.
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := plugindev.BuildPackage(cmd.Context(), dir, manifestPath)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Valid plugin package for %s (%d files)\n", pkg.Manifest.ID, len(pkg.Files))
			return nil
		},
	}

	command.Flags().StringVar(&dir, "dir", ".", "Plugin project directory")
	command.Flags().StringVarP(&manifestPath, "manifest", "m", plugindev.ManifestFilename, "Path to plugin manifest YAML")

	help.SetupMarkdownHelp(command)

	return ValidateCommand(command)
}
