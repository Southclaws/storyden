package get

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/render"
	"github.com/Southclaws/storyden/cmd/sd/internal/threadapi"
)

func New(store *config.Store) cligen.ThreadGetHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.ThreadGetParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		thread, err := threadapi.Fetch(ctx, client.OpenAPI, p.ThreadMark)
		if err != nil {
			return err
		}

		switch p.Output {
		case cligen.ThreadGetOutputJson:
			return render.ThreadJSON(io.Out, thread)
		case cligen.ThreadGetOutputYaml:
			return render.ThreadYAML(io.Out, thread)
		default:
			return render.ThreadMarkdown(io.Out, thread)
		}
	}
}
