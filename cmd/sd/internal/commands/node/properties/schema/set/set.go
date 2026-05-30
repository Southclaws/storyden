package set

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

type SetCommand *cobra.Command

func New(store *config.Store) SetCommand {
	longHelp := `# Update Property Schema

Updates the property schema of this node and its siblings.

## How Property Schemas Work

All children of a node use the same schema for properties, resulting in a table-like
structure and behaviour. When you update a node's schema:

- It affects **this node and all its siblings** (nodes with the same parent)
- All siblings will share the same property keys and types
- Values can differ between siblings, but the schema is shared

## Schema Format

Each field is specified as:
~~~
name:type:sort
~~~

Where:
- **name** - The property field name
- **type** - Property type (text, number, boolean, timestamp)
- **sort** - Sort order (asc, desc)

## Type Casting

Property schemas are loosely structured and can automatically cast their values sometimes.
A failed cast will not change data and instead just yield an empty value when reading.
However, changing the schema back to the original type (or a type compatible with what
the type was before changing) will retain the original data upon next read.

This permits easy schema experimentation and undo without data loss.

## Examples

Update schema for a node and its siblings:
~~~bash
sd node properties schema set page-1 status:text:asc priority:number:desc
~~~

This affects page-1 and all its siblings (page-2, page-3, etc if they share a parent).

Multiple properties with different types:
~~~bash
sd node properties schema set my-node title:text:asc count:number:desc completed:boolean:asc
~~~

## See Also

- ` + "`sd node properties schema children`" + ` - Update schema for this node's children
- ` + "`sd node properties set`" + ` - Set property values
`

	command := &cobra.Command{
		Use:   "set <slug> <field>:<type>:<sort>...",
		Short: "Set property schema for a node and its siblings",
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

			result, err := setSchema(cmd.Context(), client.OpenAPI, slug, schema)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Updated property schema for node: %s\n", slug)
			for _, field := range result.Properties {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s (%s) [%s]\n", field.Name, field.Type, field.Sort)
			}

			return nil
		},
	}

	// Setup beautiful markdown help rendering
	help.SetupMarkdownHelp(command)

	return SetCommand(command)
}

func setSchema(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	slug string,
	schema []openapi.PropertySchemaMutableProps,
) (*openapi.NodeUpdatePropertySchemaOK, error) {
	response, err := client.NodeUpdatePropertySchemaWithResponse(ctx, slug, schema)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, schemaSetError(response)
	}

	return response.JSON200, nil
}

func schemaSetError(response *openapi.NodeUpdatePropertySchemaResponse) error {
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
		sortStr := strings.TrimSpace(parts[2])
		sort := strings.ToLower(sortStr)

		if name == "" {
			return nil, fmt.Errorf("field name cannot be empty")
		}

		propType := openapi.PropertyType(typeStr)
		if !isValidPropertyType(propType) {
			return nil, fmt.Errorf("invalid property type %q; must be one of: text, number, boolean, timestamp", typeStr)
		}
		if !isValidSort(sort) {
			return nil, fmt.Errorf("invalid sort %q; must be \"asc\" or \"desc\"", sortStr)
		}

		schema = append(schema, openapi.PropertySchemaMutableProps{
			Name: name,
			Type: propType,
			Sort: sort,
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

func isValidSort(sort string) bool {
	return sort == "asc" || sort == "desc"
}
