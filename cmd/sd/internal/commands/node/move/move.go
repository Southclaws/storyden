package move

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

func New(store *config.Store) cligen.NodeMoveHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.NodeMoveParams) error {
		ids, err := resolveIDs(p.Slug, p.FromStdin)
		if err != nil {
			return err
		}
		if len(ids) == 0 {
			return fmt.Errorf("no nodes specified; pass slugs as arguments or use --from-stdin")
		}
		if (p.Before != "" || p.After != "") && len(ids) > 1 {
			return fmt.Errorf("--before/--after only make sense with a single node, not %d", len(ids))
		}

		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		emit, summary := emitters(cio, p.Format)

		ok := batch.Run(ctx, ids,
			func(ctx context.Context, id string) (string, error) {
				props := buildProps(p.Parent, p.Before, p.After, p.ToRoot)
				node, err := moveNode(ctx, client.OpenAPI, id, props)
				if err != nil {
					return "", err
				}
				return string(node.Name), nil
			},
			batch.Options{DryRun: p.DryRun},
			emit,
			summary,
		)
		if !ok {
			return fmt.Errorf("one or more moves failed")
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

func emitters(cio cligen.IO, format cligen.NodeMoveFormat) (func(batch.Result), func(int, int)) {
	if format == cligen.NodeMoveFormatJsonl {
		return batch.JSONLEmitter(cio.Out), nil
	}
	return batch.PlainEmitter(cio.Err), batch.PlainSummary(cio.Err)
}

func buildProps(parent, before, after string, toRoot bool) openapi.NodePositionMutableProps {
	props := openapi.NodePositionMutableProps{}
	if parent != "" {
		props.Parent.Set(parent)
	} else if toRoot {
		props.Parent.Set("")
		props.Parent.SetNull()
	}
	if before != "" {
		props.Before = &before
	} else if after != "" {
		props.After = &after
	}
	return props
}

func moveNode(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	slug string,
	props openapi.NodePositionMutableProps,
) (*openapi.NodeWithChildren, error) {
	response, err := client.NodeUpdatePositionWithResponse(ctx, slug, props)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, nodeMoveError(response)
	}

	return response.JSON200, nil
}

func nodeMoveError(response *openapi.NodeUpdatePositionResponse) error {
	if response.StatusCode() == http.StatusNotFound {
		return fmt.Errorf("node not found")
	}

	if response.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("node move request was not authorised; run sd auth login again")
	}

	if response.StatusCode() == http.StatusBadRequest {
		body := strings.TrimSpace(string(response.Body))
		if body != "" {
			return fmt.Errorf("invalid move request: %s", body)
		}

		return fmt.Errorf("invalid move request")
	}

	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("node move request failed: %s: %s", response.Status(), body)
	}

	return fmt.Errorf("node move request failed: %s", response.Status())
}
