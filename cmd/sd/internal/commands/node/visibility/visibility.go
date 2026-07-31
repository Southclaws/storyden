package visibility

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/batch"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
)

func New(store *config.Store) cligen.NodeVisibilityHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.NodeVisibilityParams) error {
		ids, err := resolveIDs(p.Slug, p.FromStdin)
		if err != nil {
			return err
		}
		if len(ids) == 0 {
			return fmt.Errorf("no nodes specified; pass slugs as positional args or use --from-stdin")
		}

		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		emit, summary := emitters(cio, p.Format)

		ok := batch.Run(ctx, ids,
			func(ctx context.Context, id string) (string, error) {
				node, err := updateVisibility(ctx, client.OpenAPI, id, string(p.Visibility))
				if err != nil {
					return "", err
				}
				return fmt.Sprintf("%s → %s", node.Name, node.Visibility), nil
			},
			batch.Options{DryRun: p.DryRun},
			emit,
			summary,
		)
		if !ok {
			return fmt.Errorf("one or more visibility updates failed")
		}
		return nil
	}
}

func resolveIDs(args []string, fromStdin bool) ([]string, error) {
	if fromStdin {
		ids, err := batch.ReadIdentifiers(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("read stdin: %w", err)
		}
		return append(ids, args...), nil
	}
	return args, nil
}

func emitters(cio cligen.IO, format cligen.NodeVisibilityFormat) (func(batch.Result), func(int, int)) {
	if format == cligen.NodeVisibilityFormatJsonl {
		return batch.JSONLEmitter(cio.Out), nil
	}
	return batch.PlainEmitter(cio.Err), batch.PlainSummary(cio.Err)
}

func updateVisibility(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	slug string,
	visibility string,
) (*openapi.NodeWithChildren, error) {
	vis := openapi.Visibility(visibility)

	response, err := client.NodeUpdateVisibilityWithResponse(ctx, slug, openapi.VisibilityMutationProps{
		Visibility: vis,
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, visibilityUpdateError(response)
	}

	return response.JSON200, nil
}

func visibilityUpdateError(response *openapi.NodeUpdateVisibilityResponse) error {
	if response.StatusCode() == http.StatusNotFound {
		return fmt.Errorf("node not found")
	}

	if response.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("visibility update request was not authorised; run sd auth login again")
	}

	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("visibility update request failed: %s: %s", response.Status(), body)
	}

	return fmt.Errorf("visibility update request failed: %s", response.Status())
}
