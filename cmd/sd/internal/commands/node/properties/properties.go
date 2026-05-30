package properties

import (
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/get"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/schema"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/set"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

type PropertiesCommand *cobra.Command

func New(
	getCommand get.GetCommand,
	setCommand set.SetCommand,
	schemaCommand schema.SchemaCommand,
) PropertiesCommand {
	command := &cobra.Command{
		Use:   "properties",
		Short: "Manage node properties",
	}

	command.AddCommand((getCommand))
	command.AddCommand((setCommand))
	command.AddCommand((schemaCommand))

	help.SetupMarkdownHelp(command)

	return PropertiesCommand(command)
}
