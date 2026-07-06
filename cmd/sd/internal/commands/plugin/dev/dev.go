package dev

import (
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/download"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/install"
	pluginnew "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/new"
	pluginpackage "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/package"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/run"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/symbols"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/validate"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

type DevCommand *cobra.Command

func New(
	newCommand pluginnew.NewCommand,
	runCommand run.RunCommand,
	packageCommand pluginpackage.PackageCommand,
	validateCommand validate.ValidateCommand,
	installCommand install.InstallCommand,
	downloadCommand download.DownloadCommand,
	symbolsCommand symbols.SymbolsCommand,
) DevCommand {
	command := &cobra.Command{
		Use:   "dev",
		Short: "Create, run, package, validate, and install plugin projects",
		Long: `# Plugin Development

Create and iterate on local Storyden plugin projects.

## Examples

Create a minimal plugin project:
~~~bash
sd plugin dev new my-plugin
~~~

Run a local external plugin:
~~~bash
sd plugin dev run
~~~

Package and install as a supervised plugin:
~~~bash
sd plugin dev install
~~~

Download and unpack an installed supervised plugin package:
~~~bash
sd plugin dev download <plugin-instance-id>
~~~
`,
	}

	command.AddCommand(newCommand)
	command.AddCommand(runCommand)
	command.AddCommand(packageCommand)
	command.AddCommand(validateCommand)
	command.AddCommand(installCommand)
	command.AddCommand(downloadCommand)
	command.AddCommand(symbolsCommand)

	help.SetupMarkdownHelp(command)

	return DevCommand(command)
}
