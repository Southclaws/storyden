package create

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
	commandcontent "github.com/Southclaws/storyden/cmd/sd/internal/content"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
)

func New(store *config.Store) cligen.NodeCreateHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.NodeCreateParams) error {
		if p.Name == "" {
			return fmt.Errorf("--name is required")
		}
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		finalContent, err := commandcontent.Read(p.Content, p.ContentFile, cio.In)
		if err != nil {
			return err
		}

		finalContent, err = commandcontent.ToHTML(finalContent, p.Markdown)
		if err != nil {
			return err
		}

		node, err := createNode(ctx, client.OpenAPI, openapi.NodeInitialProps{
			Name:          p.Name,
			Slug:          stringPtr(p.Slug),
			Description:   stringPtr(p.Description),
			Content:       stringPtr(finalContent),
			Parent:        stringPtr(p.Parent),
			Visibility:    visibilityPtr(string(p.Visibility)),
			Tags:          tagsPtr(p.Tags),
			Url:           stringPtr(p.Url),
			HideChildTree: boolPtr(p.HideChildTree),
		})
		if err != nil {
			return err
		}

		if p.Output == "json" {
			return output.JSON(cio.Out, node)
		}

		fmt.Fprintf(cio.Out, "Created node: %s (slug: %s)\n", node.Name, node.Slug)

		return nil
	}
}

func createNode(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	props openapi.NodeInitialProps,
) (*openapi.Node, error) {
	response, err := client.NodeCreateWithResponse(ctx, props)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, nodeCreateError(response)
	}

	return response.JSON200, nil
}

func nodeCreateError(response *openapi.NodeCreateResponse) error {
	if response.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("node create request was not authorised; run sd auth login again")
	}

	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("node create request failed: %s: %s", response.Status(), body)
	}

	return fmt.Errorf("node create request failed: %s", response.Status())
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}

// boolPtr returns nil when the flag is left at its default false so the field
// is omitted from the request body; only sends true.
func boolPtr(b bool) *bool {
	if !b {
		return nil
	}
	return &b
}

func visibilityPtr(v string) *openapi.Visibility {
	if v == "" {
		return nil
	}

	vis := openapi.Visibility(v)

	return &vis
}

func tagsPtr(tags []string) *[]string {
	if len(tags) == 0 {
		return nil
	}

	return &tags
}
