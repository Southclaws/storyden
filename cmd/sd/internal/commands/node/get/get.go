package get

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/render"
)

type GetCommand *cobra.Command

const (
	formatJSON     = "json"
	formatMarkdown = "markdown"
	formatYAML     = "yaml"
)

func New(store *config.Store) GetCommand {
	var format string

	command := &cobra.Command{
		Use:     "get <slug>",
		Aliases: []string{"show"},
		Short:   "Get a node by its slug",
		Long: `# Get a Node

Retrieve and display a node's content and metadata.

The default output renders the node's markdown content beautifully in your terminal with metadata shown at the top.

## Examples

View a node:
~~~bash
sd node get my-page
~~~

Get node as JSON for processing:
~~~bash
sd node get my-page --format json | jq '.content'
~~~

Get node as YAML:
~~~bash
sd node get my-page --format yaml
~~~

View and pipe to another tool:
~~~bash
sd node get docs --format json | jq -r '.name'
~~~

Extract just the content:
~~~bash
sd node get my-page --format json | jq -r '.content' > content.md
~~~
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			slug := args[0]

			if err := validateFormat(format); err != nil {
				return err
			}

			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			node, err := nodeapi.Fetch(cmd.Context(), client.OpenAPI, slug)
			if err != nil {
				return err
			}

			switch format {
			case formatJSON:
				return render.NodeJSON(cmd.OutOrStdout(), node)
			case formatMarkdown:
				return render.NodeMarkdown(cmd.OutOrStdout(), node)
			case formatYAML:
				return render.NodeYAML(cmd.OutOrStdout(), node)
			default:
				return fmt.Errorf("unsupported format %q", format)
			}
		},
	}

	command.Flags().StringVarP(&format, "format", "f", formatMarkdown, "Output format (markdown, json, yaml)")

	help.SetupMarkdownHelp(command)

	return GetCommand(command)
}

func validateFormat(format string) error {
	switch format {
	case formatJSON, formatMarkdown, formatYAML:
		return nil
	default:
		return fmt.Errorf("--format must be one of: json, markdown, yaml")
	}
}
