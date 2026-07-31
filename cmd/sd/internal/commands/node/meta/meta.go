package meta

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
)

func NewGet(store *config.Store) cligen.NodeMetaGetHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.NodeMetaGetParams) (cligen.Metadata, error) {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return nil, err
		}

		node, err := nodeapi.Fetch(ctx, client.OpenAPI, p.Slug)
		if err != nil {
			return nil, err
		}

		if node.Meta == nil {
			return openapi.Metadata{}, nil
		}

		return node.Meta, nil
	}
}

func NewSet(store *config.Store) cligen.NodeMetaSetHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.NodeMetaSetParams) error {
		if p.Json != "" && p.File != "" {
			return fmt.Errorf("cannot specify both inline JSON and --file")
		}
		if p.Json == "" && p.File == "" {
			return fmt.Errorf("provide metadata JSON as an argument or with --file")
		}

		meta, err := readMetadata(p.Json, p.File, cio.In)
		if err != nil {
			return err
		}

		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		node, err := setMetadata(ctx, client.OpenAPI, p.Slug, meta)
		if err != nil {
			return err
		}

		fmt.Fprintf(cio.Out, "Updated metadata for node: %s (slug: %s)\n", node.Name, node.Slug)
		return nil
	}
}

func readMetadata(input string, file string, stdin io.Reader) (openapi.Metadata, error) {
	if input != "" && file != "" {
		return nil, fmt.Errorf("cannot specify both inline JSON and --file")
	}

	var data []byte

	switch {
	case input != "":
		data = []byte(input)
	case file == "-":
		bytes, err := io.ReadAll(stdin)
		if err != nil {
			return nil, fmt.Errorf("failed to read metadata from stdin: %w", err)
		}
		data = bytes
	case file != "":
		bytes, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read metadata file: %w", err)
		}
		data = bytes
	default:
		return nil, fmt.Errorf("provide metadata JSON as an argument or with --file")
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, fmt.Errorf("invalid metadata JSON: %w", err)
	}
	if parsed == nil {
		return nil, fmt.Errorf("metadata must be a JSON object")
	}

	return openapi.Metadata(parsed), nil
}

func setMetadata(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	slug string,
	metadata openapi.Metadata,
) (*openapi.NodeWithChildren, error) {
	node, err := nodeapi.Update(ctx, client, slug, openapi.NodeMutableProps{
		Meta: &metadata,
	})
	if err != nil {
		return nil, err
	}

	return node, nil
}
