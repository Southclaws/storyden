package path

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
)

func New(store *config.Store) cligen.ConfigPathHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.ConfigPathParams) error {
		fmt.Fprintln(io.Out, store.Path())
		return nil
	}
}
