package download

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

func New(store *config.Store) cligen.PluginDevDownloadHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.PluginDevDownloadParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		data, err := plugindev.DownloadPackage(ctx, client.OpenAPI, p.PluginInstanceId)
		if err != nil {
			return err
		}

		if p.NoUnzip {
			filename := plugindev.PackageFilename(ctx, data, p.PluginInstanceId)
			target := filepath.Join(p.Dir, filename)
			if err := plugindev.WritePackageFile(target, &plugindev.PackageArchive{Bytes: data}, p.Force); err != nil {
				return err
			}
			fmt.Fprintf(io.Out, "Downloaded plugin package: %s\n", target)
			return nil
		}

		result, err := plugindev.ExtractPackageArchive(data, p.Dir, p.Force)
		if err != nil {
			return err
		}
		fmt.Fprintf(io.Out, "Extracted plugin package to %s (%d files)\n", p.Dir, len(result.Files))
		return nil
	}
}
