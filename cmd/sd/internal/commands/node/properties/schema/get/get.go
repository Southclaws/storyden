package get

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligenconv"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
)

func New(store *config.Store) cligen.NodePropertiesSchemaGetHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.NodePropertiesSchemaGetParams) (cligen.PropertySchemaList, error) {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return nil, err
		}

		node, err := nodeapi.Fetch(ctx, client.OpenAPI, p.Slug)
		if err != nil {
			return nil, err
		}

		// Same full shape for both json and yaml, matching properties get's
		// parity fix. codegen's generated RunE encodes the result for both.
		return cligenconv.Convert[cligen.PropertySchemaList](node.ChildPropertySchema)
	}
}
