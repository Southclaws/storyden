package node

import (
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/assets"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/children"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/create"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/delete"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/get"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/list"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/meta"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/move"
	nodeopen "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/open"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/search"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/tree"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/update"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/node/visibility"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

type NodeCommand *cobra.Command

func New(
	listCommand list.ListCommand,
	treeCommand tree.TreeCommand,
	getCommand get.GetCommand,
	createCommand create.CreateCommand,
	updateCommand update.UpdateCommand,
	deleteCommand delete.DeleteCommand,
	moveCommand move.MoveCommand,
	metaCommand meta.MetaCommand,
	assetsCommand assets.AssetsCommand,
	visibilityCommand visibility.VisibilityCommand,
	childrenCommand children.ChildrenCommand,
	propertiesCommand properties.PropertiesCommand,
	openCommand nodeopen.OpenCommand,
	searchCommand search.SearchCommand,
) NodeCommand {
	command := &cobra.Command{
		Use:   "node",
		Short: "Work with Storyden nodes (pages)",
	}

	command.AddCommand(listCommand)
	command.AddCommand(treeCommand)
	command.AddCommand(getCommand)
	command.AddCommand(createCommand)
	command.AddCommand(updateCommand)
	command.AddCommand(deleteCommand)
	command.AddCommand(moveCommand)
	command.AddCommand(metaCommand)
	command.AddCommand(assetsCommand)
	command.AddCommand(visibilityCommand)
	command.AddCommand(childrenCommand)
	command.AddCommand(propertiesCommand)
	command.AddCommand(openCommand)
	command.AddCommand(searchCommand)

	help.SetupMarkdownHelp(command)

	return NodeCommand(command)
}
