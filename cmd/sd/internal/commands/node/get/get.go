package get

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/render"
)

func New(store *config.Store) cligen.NodeGetHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.NodeGetParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		node, err := nodeapi.Fetch(ctx, client.OpenAPI, p.Slug)
		if err != nil {
			return err
		}

		switch p.Output {
		case cligen.NodeGetOutputJson:
			return render.NodeWithAncestorsJSON(io.Out, node)
		case cligen.NodeGetOutputYaml:
			return render.NodeWithAncestorsYAML(io.Out, node)
		default:
			return render.NodeWithAncestorsMarkdown(io.Out, node)
		}
	}
}
