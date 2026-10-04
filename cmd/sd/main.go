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
	searchcmd "github.com/Southclaws/storyden/cmd/sd/internal/commands/search"
	storeconfig "github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

func newRootCommand(
	streams cli.Streams,
	store *storeconfig.Store,
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

	root.SilenceUsage = true
	root.SilenceErrors = true

	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		selected, err := cmd.Flags().GetString("context")
		if err != nil {
			return err
		}
		store.SelectedContext = selected
		return nil
	}
	root.SetIn(streams.In)
	root.SetOut(streams.Out)
	root.SetErr(streams.Err)

	help.SetupMarkdownHelp(root)
	carapace.Gen(root)
	searchcmd.Complete(searchCommand)

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

		Build(),
		fx.Invoke(configureDefaultLogger),
		fx.Invoke(cli.Execute),
	)

	if err := app.Start(ctx); err != nil {
		underlying := dig.RootCause(err)
		if cli.IsCommandError(underlying) {
			os.Exit(cli.ExitCode(underlying))
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
