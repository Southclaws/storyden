package children

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

type ChildrenCommand *cobra.Command

func New(store *config.Store) ChildrenCommand {
	longHelp := `# Update Children Property Schema

Updates the property schema of the children of this node.

## How Property Schemas Work

All children of a node use the same schema for properties, resulting in a table-like
structure and behaviour. This means:

- All **sibling nodes** share the same property schema
- When you update a child's schema, it affects **all its siblings**
- Setting a children schema defines what properties the **children will have**

## Schema Format

Each field is specified as:
~~~
name:type:sort
~~~

Where:
- **name** - The property field name
- **type** - Property type (text, number, boolean, timestamp)
- **sort** - Sort order (asc, desc)

## Examples

Define schema for children of 'docs':
~~~bash
sd node properties schema children docs status:text:asc priority:number:desc
~~~

Now when you create children of 'docs', they will all have these properties:
~~~bash
sd node create --parent docs --name "Page 1"
sd node properties get page-1  # Shows status and priority (initially empty)
sd node properties set page-1 status=draft priority=1
~~~

Multiple properties:
~~~bash
sd node properties schema children tutorials difficulty:text:asc duration:number:asc hands-on:boolean:desc
~~~

## See Also

- ` + "`sd node properties schema set`" + ` - Update schema for this node and its siblings
- ` + "`sd node properties set`" + ` - Set property values
`

	command := &cobra.Command{
		Use:   "children <slug> <field>:<type>:<sort>...",
		Short: "Set property schema for a node's children",
		Long:  longHelp,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			slug := args[0]
			fields := args[1:]

			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			schema, err := parseSchema(fields)
			if err != nil {
				return err
			}

			result, err := setChildrenSchema(cmd.Context(), client.OpenAPI, slug, schema)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Updated children property schema for node: %s\n", slug)
			for _, field := range result.Properties {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s (%s) [%s]\n", field.Name, field.Type, field.Sort)
			}

			return nil
		},
	}

	// Setup beautiful markdown help rendering
	help.SetupMarkdownHelp(command)

	return ChildrenCommand(command)
}

func setChildrenSchema(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	slug string,
	schema []openapi.PropertySchemaMutableProps,
) (*openapi.NodeUpdatePropertySchemaOK, error) {
	response, err := client.NodeUpdateChildrenPropertySchemaWithResponse(ctx, slug, schema)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, schemaSetError(response)
	}

	return response.JSON200, nil
}

func schemaSetError(response *openapi.NodeUpdateChildrenPropertySchemaResponse) error {
	if response.StatusCode() == http.StatusNotFound {
		return fmt.Errorf("node not found")
	}

	if response.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("schema update request was not authorised; run sd auth login again")
	}

	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("schema update request failed: %s: %s", response.Status(), body)
	}

	return fmt.Errorf("schema update request failed: %s", response.Status())
}

func parseSchema(fields []string) ([]openapi.PropertySchemaMutableProps, error) {
	schema := make([]openapi.PropertySchemaMutableProps, 0, len(fields))

	for _, field := range fields {
		parts := strings.Split(field, ":")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid schema field format %q; expected name:type:sort", field)
		}

		name := strings.TrimSpace(parts[0])
		typeStr := strings.TrimSpace(parts[1])
		sortStr := strings.ToLower(strings.TrimSpace(parts[2]))

		if name == "" {
			return nil, fmt.Errorf("field name cannot be empty")
		}

		propType := openapi.PropertyType(typeStr)
		if !isValidPropertyType(propType) {
			return nil, fmt.Errorf("invalid property type %q; must be one of: text, number, boolean, timestamp", typeStr)
		}
		if sortStr != "asc" && sortStr != "desc" {
			return nil, fmt.Errorf("invalid sort %q; must be \"asc\" or \"desc\"", parts[2])
		}

		schema = append(schema, openapi.PropertySchemaMutableProps{
			Name: name,
			Type: propType,
			Sort: sortStr,
		})
	}

	return schema, nil
}

func isValidPropertyType(t openapi.PropertyType) bool {
	switch t {
	case openapi.PropertyTypeText, openapi.PropertyTypeNumber,
		openapi.PropertyTypeBoolean, openapi.PropertyTypeTimestamp:
		return true
	default:
		return false
	}
}
