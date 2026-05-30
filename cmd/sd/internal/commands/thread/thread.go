package thread

import (
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/commands/thread/get"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/thread/list"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

type ThreadCommand *cobra.Command

func New(
	listCommand list.ListCommand,
	getCommand get.GetCommand,
) ThreadCommand {
	command := &cobra.Command{
		Use:   "thread",
		Short: "Work with Storyden threads",
	}

	command.AddCommand((listCommand))
	command.AddCommand((getCommand))

	help.SetupMarkdownHelp(command)

	return ThreadCommand(command)
}
