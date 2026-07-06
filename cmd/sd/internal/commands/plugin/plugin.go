package plugin

import (
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/activate"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/deactivate"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/delete"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/get"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/list"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/logs"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/token"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

type PluginCommand *cobra.Command

func New(
	devCommand dev.DevCommand,
	listCommand list.ListCommand,
	getCommand get.GetCommand,
	deleteCommand delete.DeleteCommand,
	activateCommand activate.ActivateCommand,
	deactivateCommand deactivate.DeactivateCommand,
	logsCommand logs.LogsCommand,
	tokenCommand token.TokenCommand,
) PluginCommand {
	command := &cobra.Command{
		Use:   "plugin",
		Short: "Build, develop, and manage Storyden plugins",
		Long: `# Plugin Commands

Create local plugin projects, run plugins against the current Storyden instance, and manage installed plugins.

## Examples

Create a minimal plugin project:
~~~bash
sd plugin dev new my-plugin
~~~

Run a local external plugin against the current instance:
~~~bash
sd plugin dev run
~~~

List plugins installed on the current instance:
~~~bash
sd plugin list
~~~

Stream supervised plugin logs:
~~~bash
sd plugin logs <plugin-instance-id>
~~~

`,
	}

	command.AddCommand(devCommand)
	command.AddCommand(listCommand)
	command.AddCommand(getCommand)
	command.AddCommand(deleteCommand)
	command.AddCommand(activateCommand)
	command.AddCommand(deactivateCommand)
	command.AddCommand(logsCommand)
	command.AddCommand(tokenCommand)

	help.SetupMarkdownHelp(command)

	return PluginCommand(command)
}
