package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

type TUICommand *cobra.Command

func New(store *config.Store) TUICommand {
	command := &cobra.Command{
		Use:   "tui",
		Short: "Explore Storyden content in an interactive TUI",
		Long: `# Storyden TUI

Open a k9s-style explorer for Storyden content.

Use it for interactive browsing. List commands stay plain and script-friendly.

## Navigation

- **n/t** - switch nodes or threads
- **↑/↓** or **j/k** - move selection
- **←/→** or **h/l** - change pages
- **Enter** - open the selected node or thread
- **c** - drill into selected node children
- **Backspace** - return from the open view or go back to the previous node level
- **r** - refresh current view
- **q** or **Esc** - quit
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			nodes, err := fetchNodes(cmd.Context(), client.OpenAPI, 1)
			if err != nil {
				return err
			}
			threads, err := fetchThreads(cmd.Context(), client.OpenAPI, 1)
			if err != nil {
				return err
			}

			program := tea.NewProgram(
				newModel(cmd.Context(), client.OpenAPI, cmd.OutOrStdout(), nodes, threads),
				tea.WithContext(cmd.Context()),
				tea.WithInput(cmd.InOrStdin()),
				tea.WithOutput(cmd.OutOrStdout()),
			)

			_, err = program.Run()
			return err
		},
	}

	help.SetupMarkdownHelp(command)

	return TUICommand(command)
}
