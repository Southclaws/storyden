package help

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/output"
)

// IsTerminal checks if the given writer is a terminal.
func IsTerminal(w io.Writer) bool {
	return output.IsTerminal(w)
}

// FormatMarkdown renders markdown with glamour if output is a terminal,
// otherwise preserves raw Markdown.
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

// SetupMarkdownHelp renders the whole help document as Markdown in a terminal,
// preserving the same Markdown document for agents reading a pipe.
// Install on each child too: Fang replaces the root help function at execution
// time, so inherited help would otherwise differ for bare command groups.
func SetupMarkdownHelp(cmd *cobra.Command) {
	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		var doc strings.Builder
		if c.Long != "" {
			doc.WriteString(c.Long)
		} else {
			fmt.Fprintf(&doc, "# %s\n\n%s\n", c.CommandPath(), c.Short)
		}
		if c.Example != "" && !strings.Contains(c.Long, "## Examples") {
			fmt.Fprintf(&doc, "\n## Examples\n\n~~~sh\n%s\n~~~\n", strings.TrimSpace(c.Example))
		}
		fmt.Fprintf(&doc, "\n## Usage\n\n~~~text\n%s\n~~~\n", c.UseLine())
		if c.HasAvailableSubCommands() {
			doc.WriteString("\n## Commands\n\n")
			for _, subcmd := range c.Commands() {
				if subcmd.IsAvailableCommand() {
					fmt.Fprintf(&doc, "- `%s`: %s\n", subcmd.Name(), subcmd.Short)
				}
			}
		}
		if flags := c.LocalFlags(); flags.HasAvailableFlags() {
			fmt.Fprintf(&doc, "\n## Flags\n\n~~~text\n%s~~~\n", flags.FlagUsages())
		}
		if c.HasAvailableInheritedFlags() {
			fmt.Fprintf(&doc, "\n## Global flags\n\n~~~text\n%s~~~\n", c.InheritedFlags().FlagUsages())
		}
		if c.HasAvailableSubCommands() {
			fmt.Fprintf(&doc, "\nUse `%s COMMAND --help` for command details.\n", c.CommandPath())
		}
		fmt.Fprint(c.OutOrStdout(), FormatHelpMarkdown(doc.String(), c.OutOrStdout()))
	})
	for _, child := range cmd.Commands() {
		SetupMarkdownHelp(child)
	}
}
