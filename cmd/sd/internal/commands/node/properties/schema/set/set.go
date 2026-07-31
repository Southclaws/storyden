package set

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
)

func New(store *config.Store) cligen.NodePropertiesSchemaSetHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.NodePropertiesSchemaSetParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		schema, err := parseSchema(p.Token)
		if err != nil {
			return err
		}

		result, err := setSchema(ctx, client.OpenAPI, p.Slug, schema)
		if err != nil {
			return err
		}

		fmt.Fprintf(io.Out, "Updated property schema for node: %s\n", p.Slug)
		for _, field := range result.Properties {
			fmt.Fprintf(io.Out, "  %s (%s) [%s]\n", field.Name, field.Type, field.Sort)
		}

		return nil
	}
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
