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
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
)

type MetaCommand *cobra.Command

func New(store *config.Store) MetaCommand {
	command := &cobra.Command{
		Use:   "meta",
		Short: "Get or set raw node metadata",
		Long: `# Node Metadata

Read and replace a node's raw JSON metadata.

Metadata is intentionally treated as schema-less JSON. Client-specific metadata is not interpreted by the CLI.
`,
	}

	command.AddCommand(newGetCommand(store))
	command.AddCommand(newSetCommand(store))

	help.SetupMarkdownHelp(command)

	return MetaCommand(command)
}

func newGetCommand(store *config.Store) *cobra.Command {
	command := &cobra.Command{
		Use:   "get <slug>",
		Short: "Get node metadata as JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			node, err := nodeapi.Fetch(cmd.Context(), client.OpenAPI, args[0])
			if err != nil {
				return err
			}

			meta := node.Meta
			if meta == nil {
				meta = openapi.Metadata{}
			}

			return output.JSON(cmd.OutOrStdout(), meta)
		},
	}

	help.SetupMarkdownHelp(command)

	return command
}

func newSetCommand(store *config.Store) *cobra.Command {
	var file string

	command := &cobra.Command{
		Use:   "set <slug> [json]",
		Short: "Replace node metadata with a JSON object",
		Long: `# Set Node Metadata

Replace a node's raw metadata with a JSON object.

## Examples

Set inline JSON:
~~~bash
sd node meta set docs '{"source":"import"}'
~~~

Set JSON from a file:
~~~bash
sd node meta set docs --file meta.json
~~~

Read JSON from stdin:
~~~bash
cat meta.json | sd node meta set docs --file -
~~~
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 || len(args) > 2 {
				return fmt.Errorf("accepts 1 or 2 arg(s), received %d", len(args))
			}
			if len(args) == 2 && file != "" {
				return fmt.Errorf("cannot specify both inline JSON and --file")
			}
			if len(args) == 1 && file == "" {
				return fmt.Errorf("provide metadata JSON as an argument or with --file")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			input := ""
			if len(args) == 2 {
				input = args[1]
			}

			meta, err := readMetadata(input, file, cmd.InOrStdin())
			if err != nil {
				return err
			}

			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			node, err := setMetadata(cmd.Context(), client.OpenAPI, args[0], meta)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Updated metadata for node: %s (slug: %s)\n", node.Name, node.Slug)
			return nil
		},
	}

	command.Flags().StringVar(&file, "file", "", "Read metadata JSON from file (use - for stdin)")
	help.SetupMarkdownHelp(command)

	return command
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
