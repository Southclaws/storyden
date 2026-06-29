package pluginnew

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type NewCommand *cobra.Command

func New() NewCommand {
	var template string
	var commandName string
	var author string
	var force bool

	command := &cobra.Command{
		Use:   "new <directory>",
		Short: "Create a minimal plugin project",
		Long: `# New Plugin

Create a minimal Storyden plugin directory containing a ` + "`manifest.yaml`" + ` file.

A future template workflow can clone richer starters. For now, ` + "`--template`" + ` is reserved and reports that templates are not available yet.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(template) != "" {
				return fmt.Errorf("--template is reserved for template repositories but is not implemented yet")
			}

			dir := args[0]
			name := filepath.Base(filepath.Clean(dir))
			id := plugindev.Slugify(name)
			if id == "" {
				return fmt.Errorf("could not derive a valid plugin id from %q; pass a named directory", dir)
			}

			if commandName == "" {
				commandName = "./" + name
			}
			if author == "" {
				author = plugindev.DefaultAuthor()
			}

			manifest := rpc.Manifest{
				ID:          id,
				Name:        plugindev.Titleize(name),
				Author:      plugindev.Slugify(author),
				Description: "Describe what this plugin does.",
				Version:     "0.1.0",
				Command:     commandName,
				Args:        []string{},
			}

			return plugindev.WriteNewManifest(cmd.OutOrStdout(), dir, manifest, force)
		},
	}

	command.Flags().StringVar(&template, "template", "", "Template repository URL (reserved; not implemented yet)")
	command.Flags().StringVar(&commandName, "command", "", "Command used to run the plugin during development")
	command.Flags().StringVar(&author, "author", "", "Plugin author slug")
	command.Flags().BoolVar(&force, "force", false, "Overwrite an existing manifest.yaml")

	help.SetupMarkdownHelp(command)

	return NewCommand(command)
}
