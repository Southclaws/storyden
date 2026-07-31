package cligen

import (
	"github.com/spf13/cobra"
)

func NewRootCommand(
	auth AuthCommand,
	config ConfigCommand,
	info InfoCommand,
	search SearchCommand,
	node NodeCommand,
	plugin PluginCommand,
	thread ThreadCommand,
	tui TuiCommand,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sd",
		Short: "Command-line interface for working with Storyden instances.",
		Long:  "The **sd** command-line tool provides an interface for working with Storyden\ninstances. Authenticate against an instance's public web/API address with\n`sd auth login`, then read and write nodes, threads, and other content.\n",
	}

	cmd.AddCommand(
		(*cobra.Command)(auth),
		(*cobra.Command)(config),
		(*cobra.Command)(info),
		(*cobra.Command)(search),
		(*cobra.Command)(node),
		(*cobra.Command)(plugin),
		(*cobra.Command)(thread),
		(*cobra.Command)(tui),
	)

	return cmd
}
