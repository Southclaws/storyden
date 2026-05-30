package get

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
)

type GetCommand *cobra.Command

const (
	formatJSON = "json"
	formatYAML = "yaml"
)

func New(store *config.Store) GetCommand {
	var format string

	command := &cobra.Command{
		Use:   "get <slug>",
		Short: "Get property values for a node",
		Long: `# Get Property Values

View all property values for a node.

Properties are custom key-value pairs that store structured metadata. All sibling nodes share the same property schema, creating a table-like structure.

## Examples

View properties:
~~~bash
sd node properties get my-page
~~~

Get as JSON:
~~~bash
sd node properties get my-page --format json
~~~

Get as YAML:
~~~bash
sd node properties get my-page --format yaml
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
				return output.JSON(cmd.OutOrStdout(), node.Properties)
			case formatYAML:
				return renderYAML(cmd.OutOrStdout(), node.Properties)
			default:
				return fmt.Errorf("unsupported format %q", format)
			}
		},
	}

	command.Flags().StringVarP(&format, "format", "f", formatJSON, "Output format: json, yaml")

	help.SetupMarkdownHelp(command)

	return GetCommand(command)
}

func validateFormat(format string) error {
	switch format {
	case formatJSON, formatYAML:
		return nil
	default:
		return fmt.Errorf("--format must be one of: json, yaml")
	}
}

func renderYAML(out io.Writer, properties []openapi.Property) error {
	payload := yamlPropertiesPayload{
		Properties: yamlProperties(properties),
	}

	return output.YAML(out, payload)
}

type yamlPropertiesPayload struct {
	Properties []yamlProperty `yaml:"properties"`
}

type yamlProperty struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

func yamlProperties(properties []openapi.Property) []yamlProperty {
	out := make([]yamlProperty, len(properties))
	for i, prop := range properties {
		out[i] = yamlProperty{
			Name:  string(prop.Name),
			Type:  string(prop.Type),
			Value: string(prop.Value),
		}
	}

	return out
}
