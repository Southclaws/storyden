package create

import (
	"context"
	"fmt"
	"strings"

	domainvisibility "github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	commandcontent "github.com/Southclaws/storyden/cmd/sd/internal/content"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
	"github.com/Southclaws/storyden/cmd/sd/internal/threadapi"
	"github.com/spf13/cobra"
)

func New(store *config.Store) cligen.ThreadCreateHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.ThreadCreateParams) error {
		format := string(p.Output)
		visibility := string(p.Visibility)

		if strings.TrimSpace(p.Title) == "" {
			return fmt.Errorf("--title is required")
		}
		if format != "plain" && format != "json" {
			return fmt.Errorf("--output must be plain or json")
		}
		if _, err := domainvisibility.NewVisibility(visibility); err != nil {
			return fmt.Errorf("invalid --visibility: %s", visibility)
		}
		body, err := commandcontent.Read(p.Content, p.ContentFile, cio.In)
		if err != nil {
			return err
		}
		body, err = commandcontent.ToHTML(body, p.Markdown)
		if err != nil {
			return err
		}
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}
		vis := openapi.Visibility(visibility)
		props := openapi.ThreadInitialProps{Title: p.Title, Body: &body, Visibility: &vis}
		if p.Category != "" {
			props.Category = &p.Category
		}
		response, err := threadapi.Create(ctx, client.OpenAPI, props)
		if err != nil {
			return err
		}
		if format == "json" {
			return output.JSON(cio.Out, response)
		}
		_, err = fmt.Fprintf(cio.Out, "Created thread: %s (%s)\n", response.Title, response.Id)
		return err
	}
}
