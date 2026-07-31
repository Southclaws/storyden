package cligen

import (
	"context"
	"github.com/spf13/cobra"
)

type ConfigPathParams struct {
}

type ConfigPathHandler func(ctx context.Context, cmd *cobra.Command, io IO, p ConfigPathParams) error

func newConfigPathCommand(configPath ConfigPathHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "path",
		Short: "Print the sd config file path.",
		Args:  rangeArgs(0, 0),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		p := ConfigPathParams{}

		return configPath(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

func NewConfigCommand(
	configPath ConfigPathHandler,
) ConfigCommand {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Inspect and manage sd configuration.",
	}

	cmd.AddCommand(
		newConfigPathCommand(configPath),
	)

	return ConfigCommand(cmd)
}

type ConfigCommand *cobra.Command
