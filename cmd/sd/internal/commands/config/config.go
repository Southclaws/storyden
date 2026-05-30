package configcmd

import (
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/commands/config/path"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

type ConfigCommand *cobra.Command

func New(pathCommand path.PathCommand) ConfigCommand {
	command := &cobra.Command{
		Use:   "config",
		Short: "Inspect and manage sd configuration",
	}

	command.AddCommand((pathCommand))

	help.SetupMarkdownHelp(command)

	return ConfigCommand(command)
}
