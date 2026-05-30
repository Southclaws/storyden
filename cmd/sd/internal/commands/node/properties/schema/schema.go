package schema

import (
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/schema/children"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/schema/get"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/schema/set"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

type SchemaCommand *cobra.Command

func New(
	getCommand get.GetCommand,
	setCommand set.SetCommand,
	childrenCommand children.ChildrenCommand,
) SchemaCommand {
	command := &cobra.Command{
		Use:   "schema",
		Short: "Manage node property schemas",
	}

	command.AddCommand((getCommand))
	command.AddCommand((setCommand))
	command.AddCommand((childrenCommand))

	help.SetupMarkdownHelp(command)

	return SchemaCommand(command)
}
