package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	charmLog "charm.land/log/v2"
	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/cmd/sd/internal/cli"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/auth"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/auth/login"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/auth/remove"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/auth/switcher"
	configcmd "github.com/Southclaws/storyden/cmd/sd/internal/commands/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/config/path"
	nodecmd "github.com/Southclaws/storyden/cmd/sd/internal/commands/node"
	nodeassets "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/assets"
	nodechildren "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/children"
	nodecreate "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/create"
	nodedelete "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/delete"
	nodeget "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/get"
	nodelist "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/list"
	nodemeta "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/meta"
	nodemove "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/move"
	nodeopen "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/open"
	nodeproperties "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties"
	propertiesget "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/get"
	propertiesschema "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/schema"
	schemachildren "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/schema/children"
	schemaget "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/schema/get"
	schemaset "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/schema/set"
	propertiesset "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/set"
	nodesearch "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/search"
	nodetree "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/tree"
	nodeupdate "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/update"
	nodevisibility "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/visibility"
	plugincmd "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin"
	pluginactivate "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/activate"
	plugindeactivate "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/deactivate"
	plugindelete "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/delete"
	plugindev "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev"
	plugindevinstall "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/install"
	plugindevnew "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/new"
	plugindevpackage "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/package"
	plugindevrun "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/run"
	plugindevvalidate "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/validate"
	pluginget "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/get"
	pluginlist "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/list"
	pluginlogs "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/logs"
	plugintoken "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/token"
	plugintokenrotate "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/token/rotate"
	searchcmd "github.com/Southclaws/storyden/cmd/sd/internal/commands/search"
	threadcmd "github.com/Southclaws/storyden/cmd/sd/internal/commands/thread"
	threadget "github.com/Southclaws/storyden/cmd/sd/internal/commands/thread/get"
	threadlist "github.com/Southclaws/storyden/cmd/sd/internal/commands/thread/list"
	tuicmd "github.com/Southclaws/storyden/cmd/sd/internal/commands/tui"
	storeconfig "github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

func newRootCommand(
	streams cli.Streams,
	authCommand auth.AuthCommand,
	configCommand configcmd.ConfigCommand,
	threadCommand threadcmd.ThreadCommand,
	nodeCommand nodecmd.NodeCommand,
	pluginCommand plugincmd.PluginCommand,
	searchCommand searchcmd.SearchCommand,
	tuiCommand tuicmd.TUICommand,
) *cobra.Command {
	root := &cobra.Command{
		Use:   "sd",
		Short: "Storyden CLI",
		Long: `# Storyden CLI

The **sd** command-line tool provides a powerful interface for working with Storyden instances.

## Getting Started

To read and write information from a Storyden instance, authenticate using its public web/API address:
~~~bash
sd auth login https://your-instance.com
~~~

## Configuration

The CLI stores authentication and context configuration in:
- Windows: ` + "`%APPDATA%/storyden/config.yaml`" + `
- macOS: ` + "`~/Library/Application Support/storyden/config.yaml`" + `
- Linux: ` + "`~/.config/storyden/config.yaml`" + ` (or ` + "`$XDG_CONFIG_HOME/storyden/config.yaml`" + `)

View your config file location:
~~~bash
sd config path
~~~

## Multiple Instances

You can authenticate with multiple Storyden instances and switch between them:
~~~bash
sd auth login https://instance1.com
sd auth login https://instance2.com
sd auth switch
~~~
`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.SetIn(streams.In)
	root.SetOut(streams.Out)
	root.SetErr(streams.Err)

	root.AddCommand((*cobra.Command)(authCommand))
	root.AddCommand((*cobra.Command)(configCommand))
	root.AddCommand((*cobra.Command)(threadCommand))
	root.AddCommand((*cobra.Command)(nodeCommand))
	root.AddCommand((*cobra.Command)(pluginCommand))
	root.AddCommand((*cobra.Command)(searchCommand))
	root.AddCommand((*cobra.Command)(tuiCommand))

	help.SetupMarkdownHelp(root)
	carapace.Gen(root)

	return root
}

func newLogger(streams cli.Streams) *slog.Logger {
	return slog.New(charmLog.New(streams.Err))
}

func configureDefaultLogger(logger *slog.Logger) {
	slog.SetDefault(logger)
}

func main() {
	ctx, cf := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cf()

	app := fx.New(
		fx.NopLogger,

		fx.Provide(func() context.Context { return ctx }),

		fx.Provide(
			storeconfig.NewStore,
			cli.NewStreams,
			newLogger,
			login.New,
			remove.New,
			switcher.New,
			auth.New,
			path.New,
			configcmd.New,
			threadlist.New,
			threadget.New,
			threadcmd.New,
			tuicmd.New,
			nodelist.New,
			nodetree.New,
			nodeget.New,
			nodecreate.New,
			nodeupdate.New,
			nodedelete.New,
			nodemove.New,
			nodeopen.New,
			nodesearch.New,
			nodemeta.New,
			nodeassets.New,
			nodevisibility.New,
			nodechildren.New,
			propertiesget.New,
			propertiesset.New,
			schemaget.New,
			schemaset.New,
			schemachildren.New,
			propertiesschema.New,
			nodeproperties.New,
			nodecmd.New,
			plugindevnew.New,
			plugindevrun.New,
			plugindevpackage.New,
			plugindevvalidate.New,
			plugindevinstall.New,
			plugindev.New,
			pluginlist.New,
			pluginget.New,
			plugindelete.New,
			pluginactivate.New,
			plugindeactivate.New,
			pluginlogs.New,
			plugintokenrotate.New,
			plugintoken.New,
			plugincmd.New,
			searchcmd.New,
			newRootCommand,
		),
		fx.Invoke(configureDefaultLogger),
		fx.Invoke(cli.Execute),
	)

	if err := app.Start(ctx); err != nil {
		underlying := dig.RootCause(err)
		if cli.IsCommandError(underlying) {
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, underlying)
		os.Exit(1)
	}

	stopCtx, stop := context.WithTimeout(context.Background(), time.Second*5)
	defer stop()

	if err := app.Stop(stopCtx); err != nil {
		slog.Error("fatal error occurred", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
