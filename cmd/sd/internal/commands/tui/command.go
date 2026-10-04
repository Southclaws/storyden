package tui

import (
	"context"

	tea "charm.land/bubbletea/v2"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
)

func New(store *config.Store) cligen.TuiHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.TuiParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		nodes, err := fetchNodes(ctx, client.OpenAPI, 1)
		if err != nil {
			return err
		}
		threads, err := fetchThreads(ctx, client.OpenAPI, 1)
		if err != nil {
			return err
		}

		program := tea.NewProgram(
			newModel(ctx, client.OpenAPI, io.Out, nodes, threads),
			tea.WithContext(ctx),
			tea.WithInput(io.In),
			tea.WithOutput(io.Out),
		)

		_, err = program.Run()
		return err
	}
}
