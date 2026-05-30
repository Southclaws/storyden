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
	longHelp := `# Set Property Values

Set property values on a node.

## How Properties Work

Properties are key-value pairs attached to nodes. All sibling nodes share the same
property schema (defined keys and types). The CLI automatically handles property IDs
for you when updating existing properties.

## Syntax

**For NEW properties** (not in schema yet):
~~~
name:type=value
~~~

**For EXISTING properties** (already in schema):
~~~
name=value
~~~

The CLI fetches existing properties and infers types automatically.

## Available Types

- **text** - String values
- **number** - Numeric values
- **boolean** - true/false values
- **timestamp** - ISO 8601 date/time values

## Examples

Create new properties (type required):
~~~bash
sd node properties set my-node status:text=draft priority:number=1
~~~

Update existing properties (type auto-detected):
~~~bash
sd node properties set my-node status=published priority=2
~~~

Mix new and existing properties:
~~~bash
sd node properties set my-node status=active new_field:text=hello
~~~

## See Also

- ` + "`sd node properties schema set`" + ` - Update property schema for this node and its siblings
- ` + "`sd node properties get`" + ` - View current property values
`

	command := &cobra.Command{
		Use:   "set <slug> <property>=<value>...",
		Short: "Set property values for a node",
		Long:  longHelp,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			slug := args[0]
			properties := args[1:]

			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			// First, fetch existing properties to get their FIDs
			existingProps, err := fetchNodeProperties(cmd.Context(), client.OpenAPI, slug)
			if err != nil {
				return err
			}

			propMutations, err := parseProperties(properties, existingProps)
			if err != nil {
				return err
			}

			result, err := setProperties(cmd.Context(), client.OpenAPI, slug, propMutations)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Updated properties for node: %s\n", slug)
			for _, prop := range result.Properties {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s (%s): %v\n", prop.Name, prop.Type, prop.Value)
			}

			return nil
		},
	}

	// Setup beautiful markdown help rendering
	help.SetupMarkdownHelp(command)

	return SetCommand(command)
}

func setProperties(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	slug string,
	properties []openapi.PropertyMutation,
) (*openapi.NodeUpdatePropertiesOK, error) {
	props := openapi.PropertyMutableProps{
		Properties: properties,
	}
	response, err := client.NodeUpdatePropertiesWithResponse(ctx, slug, props)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, propertiesSetError(response)
	}

	return response.JSON200, nil
}

func propertiesSetError(response *openapi.NodeUpdatePropertiesResponse) error {
	if response.StatusCode() == http.StatusNotFound {
		return fmt.Errorf("node not found")
	}

	if response.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("properties update request was not authorised; run sd auth login again")
	}

	if response.StatusCode() == http.StatusBadRequest {
		body := strings.TrimSpace(string(response.Body))

		// Check for common error patterns and provide helpful messages
		if strings.Contains(body, "missing type") || strings.Contains(body, "no type on new field") {
			return fmt.Errorf(`property type required for new properties

When creating a new property, you must specify its type:
  sd node properties set <slug> name:type=value

Example:
  sd node properties set my-node status:text=draft priority:number=1

Available types: text, number, boolean, timestamp

Original error: %s`, body)
		}

		if strings.Contains(body, "cannot remove and add the same property") || strings.Contains(body, "field ID") {
			return fmt.Errorf(`property ID conflict detected

This usually means the property exists but wasn't detected. The CLI should
handle this automatically. Please report this as a bug.

Original error: %s`, body)
		}

		if body != "" {
			return fmt.Errorf("invalid properties: %s", body)
		}

		return fmt.Errorf("invalid properties")
	}

	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("properties update request failed: %s: %s", response.Status(), body)
	}

	return fmt.Errorf("properties update request failed: %s", response.Status())
}

func fetchNodeProperties(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	slug string,
) ([]openapi.Property, error) {
	response, err := client.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{})
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, propertiesFetchError(response)
	}

	return response.JSON200.Properties, nil
}

func propertiesFetchError(response *openapi.NodeGetResponse) error {
	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("failed to fetch node properties: %s: %s", response.Status(), body)
	}

	return fmt.Errorf("failed to fetch node properties: %s", response.Status())
}

func parseProperties(properties []string, existing []openapi.Property) ([]openapi.PropertyMutation, error) {
	mutations := make([]openapi.PropertyMutation, 0, len(properties))

	// Create a map of existing properties by name for quick lookup
	existingMap := make(map[string]openapi.Property)
	for _, prop := range existing {
		existingMap[prop.Name] = prop
	}

	for _, prop := range properties {
		parts := strings.SplitN(prop, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid property format %q; expected name=value or name:type=value", prop)
		}

		nameAndType := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Check if type is specified (name:type)
		var name string
		var propType *openapi.PropertyType

		if strings.Contains(nameAndType, ":") {
			typeParts := strings.SplitN(nameAndType, ":", 2)
			name = strings.TrimSpace(typeParts[0])
			typeStr := strings.TrimSpace(typeParts[1])

			if typeStr == "" {
				return nil, fmt.Errorf("property type cannot be empty; expected name:type=value")
			}
			pt := openapi.PropertyType(typeStr)
			if !isValidPropertyType(pt) {
				return nil, fmt.Errorf("invalid property type %q; must be one of: text, number, boolean, timestamp", typeStr)
			}
			propType = &pt
		} else {
			name = nameAndType
		}

		if name == "" {
			return nil, fmt.Errorf("property name cannot be empty")
		}

		mutation := openapi.PropertyMutation{
			Name:  name,
			Type:  propType,
			Value: value,
		}

		// If the property exists, include its FID and type for update
		if existingProp, exists := existingMap[name]; exists {
			mutation.Fid = &existingProp.Fid
			// If no type was explicitly specified, use the existing type
			if propType == nil {
				mutation.Type = &existingProp.Type
			}
		}

		mutations = append(mutations, mutation)
	}

	return mutations, nil
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
