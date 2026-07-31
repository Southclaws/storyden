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
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/auth/login"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/auth/remove"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/auth/switcher"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/config/path"
	infocmd "github.com/Southclaws/storyden/cmd/sd/internal/commands/info"
	nodeassets "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/assets"
	nodechildren "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/children"
	nodecreate "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/create"
	nodedelete "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/delete"
	nodeget "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/get"
	nodelist "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/list"
	nodemeta "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/meta"
	nodemove "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/move"
	nodeopen "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/open"
	propertiesget "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/get"
	schemachildren "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/schema/children"
	schemaget "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/schema/get"
	schemaset "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/schema/set"
	propertiesset "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/set"
	nodesearch "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/search"
	nodetree "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/tree"
	nodeupdate "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/update"
	nodevisibility "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/visibility"
	pluginactivate "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/activate"
	plugindeactivate "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/deactivate"
	plugindelete "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/delete"
	plugindevdownload "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/download"
	plugindevinstall "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/install"
	plugindevnew "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/new"
	plugindevpackage "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/package"
	plugindevrun "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/run"
	plugindevsymbols "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/symbols"
	plugindevvalidate "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/validate"
	pluginget "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/get"
	pluginlist "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/list"
	pluginlogs "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/logs"
	plugintokenrotate "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/token/rotate"
	searchcmd "github.com/Southclaws/storyden/cmd/sd/internal/commands/search"
	threadget "github.com/Southclaws/storyden/cmd/sd/internal/commands/thread/get"
	threadlist "github.com/Southclaws/storyden/cmd/sd/internal/commands/thread/list"
	tuicmd "github.com/Southclaws/storyden/cmd/sd/internal/commands/tui"
	storeconfig "github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

func newRootCommand(
	streams cli.Streams,
	authCommand cligen.AuthCommand,
	configCommand cligen.ConfigCommand,
	infoCommand cligen.InfoCommand,
	searchCommand cligen.SearchCommand,
	nodeCommand cligen.NodeCommand,
	pluginCommand cligen.PluginCommand,
	threadCommand cligen.ThreadCommand,
	tuiCommand cligen.TuiCommand,
) *cobra.Command {
	root := cligen.NewRootCommand(
		authCommand,
		configCommand,
		infoCommand,
		searchCommand,
		nodeCommand,
		pluginCommand,
		threadCommand,
		tuiCommand,
	)

	root.Long = `# Storyden CLI

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

## Instance Information

Agents and scripts can inspect top-line information about the current authenticated instance:
~~~bash
sd info
~~~

## Multiple Instances

You can authenticate with multiple Storyden instances and switch between them:
~~~bash
sd auth login https://instance1.com
sd auth login https://instance2.com
sd auth switch
~~~
`
	root.SilenceUsage = true
	root.SilenceErrors = true

	root.SetIn(streams.In)
	root.SetOut(streams.Out)
	root.SetErr(streams.Err)

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

			// auth
			login.New,
			remove.New,
			switcher.New,
			cligen.NewAuthCommand,

			// config
			path.New,
			cligen.NewConfigCommand,

			// info
			infocmd.New,
			infocmd.NewMetadata,
			cligen.NewInfoCommand,

			// search
			searchcmd.New,
			cligen.NewSearchCommand,

			// node
			nodelist.New,
			nodetree.New,
			nodeget.New,
			nodecreate.New,
			nodeupdate.New,
			nodedelete.New,
			nodemove.New,
			nodeopen.New,
			nodesearch.New,
			nodemeta.NewGet,
			nodemeta.NewSet,
			nodeassets.NewUpload,
			nodeassets.NewAdd,
			nodeassets.NewRemove,
			nodeassets.NewDownload,
			nodeassets.NewPrimarySet,
			nodeassets.NewPrimaryClear,
			nodeassets.NewPrimaryDownload,
			nodevisibility.New,
			nodechildren.New,
			propertiesget.New,
			propertiesset.New,
			schemaget.New,
			schemaset.New,
			schemachildren.New,
			cligen.NewNodeCommand,

			// plugin
			plugindevnew.New,
			plugindevrun.New,
			plugindevpackage.New,
			plugindevvalidate.New,
			plugindevinstall.New,
			plugindevdownload.New,
			plugindevsymbols.NewPackages,
			plugindevsymbols.NewPackage,
			plugindevsymbols.NewDetail,
			plugindevsymbols.NewSearch,
			pluginlist.New,
			pluginget.New,
			plugindelete.New,
			pluginactivate.New,
			plugindeactivate.New,
			pluginlogs.New,
			plugintokenrotate.New,
			cligen.NewPluginCommand,

			// thread
			threadlist.New,
			threadget.New,
			cligen.NewThreadCommand,

			// tui
			tuicmd.New,
			cligen.NewTuiCommand,

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
