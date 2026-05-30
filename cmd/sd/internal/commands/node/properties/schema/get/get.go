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
		Short: "Get property schema for a node and its siblings",
		Long: `# Get Property Schema

View the property schema that applies to this node and all its siblings.

Property schemas define the structure (field names, types, sort order) for properties. All sibling nodes share the same schema.

## Examples

View schema:
~~~bash
sd node properties schema get my-page
~~~

Get as JSON:
~~~bash
sd node properties schema get my-page --format json
~~~

Get as YAML:
~~~bash
sd node properties schema get my-page --format yaml
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

			schema := node.ChildPropertySchema

			switch format {
			case formatJSON:
				return output.JSON(cmd.OutOrStdout(), schema)
			case formatYAML:
				return renderYAML(cmd.OutOrStdout(), schema)
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

func renderYAML(out io.Writer, schema []openapi.PropertySchema) error {
	type schemaField struct {
		Name string `yaml:"name"`
		Type string `yaml:"type"`
		Sort string `yaml:"sort"`
	}
	payload := struct {
		Schema []schemaField `yaml:"schema"`
	}{
		Schema: make([]schemaField, 0, len(schema)),
	}

	for _, field := range schema {
		payload.Schema = append(payload.Schema, schemaField{
			Name: field.Name,
			Type: string(field.Type),
			Sort: field.Sort,
		})
	}

	return output.YAML(out, payload)
}
