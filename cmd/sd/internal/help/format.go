package help

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/output"
)

// IsTerminal checks if the given writer is a terminal.
func IsTerminal(w io.Writer) bool {
	return output.IsTerminal(w)
}

// FormatMarkdown renders markdown with glamour if output is a terminal,
// otherwise returns plain text.
func FormatMarkdown(markdown string, out io.Writer) string {
	if !IsTerminal(out) {
		return markdown
	}

	width := getTerminalWidth(out)

	return output.Markdown(markdown, out, width)
}

// FormatHelpMarkdown renders command help markdown with the same terminal
// styling as other markdown output. Keep this separate from command execution:
// pre-rendering help into command Long fields causes ANSI to be escaped later.
func FormatHelpMarkdown(markdown string, out io.Writer) string {
	return FormatMarkdown(markdown, out)
}

// getTerminalWidth returns the terminal width, defaulting to 80 if detection fails.
func getTerminalWidth(w io.Writer) int {
	width := output.TerminalWidth(w, 80)
	// Subtract padding to avoid edge wrapping.
	if width > 10 {
		return width - 2
	}

	return width
}

// SetupMarkdownHelp configures a cobra command to render its Long description
// with beautiful markdown formatting in the terminal.
func SetupMarkdownHelp(cmd *cobra.Command) {
	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		// Render the Long description with glamour
		if c.Long != "" {
			formatted := FormatHelpMarkdown(c.Long, c.OutOrStdout())
			fmt.Fprint(c.OutOrStdout(), formatted)
		}

		// Show usage
		fmt.Fprintf(c.OutOrStdout(), "\nUsage:\n  %s\n", c.UseLine())

		// Show available subcommands
		if c.HasAvailableSubCommands() {
			fmt.Fprint(c.OutOrStdout(), "\nAvailable Commands:\n")
			for _, subcmd := range c.Commands() {
				if !subcmd.IsAvailableCommand() {
					continue
				}
				fmt.Fprintf(c.OutOrStdout(), "  %-15s %s\n", subcmd.Name(), subcmd.Short)
			}
		}

		// Show flags
		if c.HasAvailableFlags() {
			fmt.Fprint(c.OutOrStdout(), "\nFlags:\n")
			fmt.Fprint(c.OutOrStdout(), c.Flags().FlagUsages())
		}

		// Show global flags if any
		if c.HasAvailableInheritedFlags() {
			fmt.Fprint(c.OutOrStdout(), "\nGlobal Flags:\n")
			fmt.Fprint(c.OutOrStdout(), c.InheritedFlags().FlagUsages())
		}

		// Show additional help
		if c.HasHelpSubCommands() {
			fmt.Fprintf(c.OutOrStdout(), "\nUse \"%s [command] --help\" for more information about a command.\n", c.CommandPath())
		}
	})
}
