package packagecmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

func New() cligen.PluginDevPackageHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.PluginDevPackageParams) error {
		path := p.OutputFile
		if path == "" {
			mf, err := plugindev.ReadProjectManifest(p.Dir, p.Manifest)
			if err != nil {
				return err
			}
			path = plugindev.DefaultPackagePath(p.Dir, mf.Manifest)
		}

		excludes := []string{path}
		pkg, err := plugindev.BuildPackage(ctx, p.Dir, p.Manifest, excludes...)
		if err != nil {
			return err
		}

		if err := plugindev.WritePackageFile(path, pkg, p.Force); err != nil {
			return err
		}

		fmt.Fprintf(io.Out, "Created plugin package %s (%d files)\n", path, len(pkg.Files))
		return nil
	}
}
