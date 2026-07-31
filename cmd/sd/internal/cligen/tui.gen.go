package cligen

import (
	"context"
	"github.com/spf13/cobra"
)

type TuiParams struct {
}

type TuiHandler func(ctx context.Context, cmd *cobra.Command, io IO, p TuiParams) error

func NewTuiCommand(tui TuiHandler) TuiCommand {
	cmd := &cobra.Command{
		Use:   "tui",
		Short: "Explore Storyden content in an interactive TUI.",
		Long: `Launches an interactive terminal explorer. Has no flags or arguments
and no machine-readable output — it requires a real terminal.
`,
		Args: rangeArgs(0, 0),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		p := TuiParams{}

		return tui(cmd.Context(), cmd, newIO(cmd), p)
	}

	return TuiCommand(cmd)
}

type TuiCommand *cobra.Command
