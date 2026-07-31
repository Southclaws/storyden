package pluginnew

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func New() cligen.PluginDevNewHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.PluginDevNewParams) error {
		dir := p.Directory
		name := filepath.Base(filepath.Clean(dir))
		id := plugindev.Slugify(name)
		if id == "" {
			return fmt.Errorf("could not derive a valid plugin id from %q; pass a named directory", dir)
		}

		commandName := p.Command
		if commandName == "" {
			commandName = "./" + name
		}
		author := p.Author
		if author == "" {
			author = plugindev.DefaultAuthor()
		}

		manifest := rpc.Manifest{
			ID:          id,
			Name:        plugindev.Titleize(name),
			Author:      plugindev.Slugify(author),
			Description: "Describe what this plugin does.",
			Version:     "0.1.0",
			Command:     commandName,
			Args:        []string{},
		}

		return plugindev.WriteNewManifest(io.Out, dir, manifest, p.Force)
	}
}
