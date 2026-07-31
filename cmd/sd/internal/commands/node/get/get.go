package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligenconv"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/render"
)

func New(store *config.Store) cligen.NodeGetHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.NodeGetParams) (cligen.NodeWithChildren, error) {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return cligen.NodeWithChildren{}, err
		}

		node, err := nodeapi.Fetch(ctx, client.OpenAPI, p.Slug)
		if err != nil {
			return cligen.NodeWithChildren{}, err
		}

		result, err := cligenconv.Convert[cligen.NodeWithChildren](node)
		if err != nil {
			return cligen.NodeWithChildren{}, err
		}

		switch p.Format {
		case cligen.NodeGetFormatMarkdown:
			if err := render.NodeMarkdown(io.Out, node); err != nil {
				return cligen.NodeWithChildren{}, err
			}
		case cligen.NodeGetFormatJson, cligen.NodeGetFormatYaml:
			// codegen's generated RunE encodes result for these formats.
		default:
			return cligen.NodeWithChildren{}, fmt.Errorf("unsupported format %q", p.Format)
		}

		return result, nil
	}
}
