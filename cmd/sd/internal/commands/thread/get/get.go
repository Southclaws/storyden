package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligenconv"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/render"
	"github.com/Southclaws/storyden/cmd/sd/internal/threadapi"
)

func New(store *config.Store) cligen.ThreadGetHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.ThreadGetParams) (cligen.Thread, error) {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return cligen.Thread{}, err
		}

		thread, err := threadapi.Fetch(ctx, client.OpenAPI, p.ThreadMark)
		if err != nil {
			return cligen.Thread{}, err
		}

		result, err := cligenconv.Convert[cligen.Thread](thread)
		if err != nil {
			return cligen.Thread{}, err
		}

		switch p.Format {
		case cligen.ThreadGetFormatMarkdown:
			if err := render.ThreadMarkdown(io.Out, thread); err != nil {
				return cligen.Thread{}, err
			}
		case cligen.ThreadGetFormatJson, cligen.ThreadGetFormatYaml:
			// codegen's generated RunE encodes result for these formats.
		default:
			return cligen.Thread{}, fmt.Errorf("unsupported format %q", p.Format)
		}

		return result, nil
	}
}
