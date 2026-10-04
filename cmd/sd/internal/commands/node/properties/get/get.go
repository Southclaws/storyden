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

func New(store *config.Store) cligen.NodePropertiesGetHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.NodePropertiesGetParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		node, err := nodeapi.Fetch(ctx, client.OpenAPI, p.Slug)
		if err != nil {
			return err
		}

		if string(p.Output) == "yaml" {
			return renderYAML(io.Out, node.Properties)
		}
		return output.JSON(io.Out, node.Properties)
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
