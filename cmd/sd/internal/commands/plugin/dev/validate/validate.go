package validate

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

func New() cligen.PluginDevValidateHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.PluginDevValidateParams) error {
		pkg, err := plugindev.BuildPackage(ctx, p.Dir, p.Manifest)
		if err != nil {
			return err
		}

		fmt.Fprintf(io.Out, "Valid plugin package for %s (%d files)\n", pkg.Manifest.ID, len(pkg.Files))
		return nil
	}
}
