package get

import (
	"context"
	"io"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
)

func New(store *config.Store) cligen.NodePropertiesSchemaGetHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.NodePropertiesSchemaGetParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		node, err := nodeapi.Fetch(ctx, client.OpenAPI, p.Slug)
		if err != nil {
			return err
		}

		if string(p.Output) == "yaml" {
			return renderYAML(io.Out, node.ChildPropertySchema)
		}
		return output.JSON(io.Out, node.ChildPropertySchema)
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
