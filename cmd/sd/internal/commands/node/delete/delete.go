package delete

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

func New(store *config.Store) cligen.NodeDeleteHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.NodeDeleteParams) error {
		ids, err := resolveIDs(p.Slug, p.FromStdin)
		if err != nil {
			return err
		}
		if len(ids) == 0 {
			return fmt.Errorf("no nodes specified; pass slugs as arguments or use --from-stdin")
		}

		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		emit, summary := emitters(cio, p.Format)

		ok := batch.Run(ctx, ids,
			func(ctx context.Context, id string) (string, error) {
				return id, deleteNode(ctx, client.OpenAPI, id, p.Target)
			},
			batch.Options{DryRun: p.DryRun},
			emit,
			summary,
		)
		if !ok {
			return fmt.Errorf("one or more deletions failed")
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

// emitters returns (per-item emit, end-of-run summary). JSONL routes per-item
// results to stdout so they can be piped; plain routes to stderr so stdout
// stays empty for shell composition.
func emitters(cio cligen.IO, format cligen.NodeDeleteFormat) (func(batch.Result), func(int, int)) {
	if format == cligen.NodeDeleteFormatJsonl {
		return batch.JSONLEmitter(cio.Out), nil
	}
	return batch.PlainEmitter(cio.Err), batch.PlainSummary(cio.Err)
}

func deleteNode(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	slug string,
	target string,
) error {
	params := &openapi.NodeDeleteParams{}
	if target != "" {
		targetQuery := openapi.TargetNodeSlugQuery(target)
		params.TargetNode = &targetQuery
	}

	response, err := client.NodeDeleteWithResponse(ctx, slug, params)
	if err != nil {
		return err
	}

	if response.StatusCode() != http.StatusOK {
		return nodeDeleteError(response)
	}

	return nil
}

func nodeDeleteError(response *openapi.NodeDeleteResponse) error {
	if response.StatusCode() == http.StatusNotFound {
		return fmt.Errorf("node not found")
	}

	if response.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("node delete request was not authorised; run sd auth login again")
	}

	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("node delete request failed: %s: %s", response.Status(), body)
	}

	return fmt.Errorf("node delete request failed: %s", response.Status())
}
