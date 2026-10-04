package reply

import (
	"context"
	"fmt"
	"strings"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	commandcontent "github.com/Southclaws/storyden/cmd/sd/internal/content"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
	"github.com/Southclaws/storyden/cmd/sd/internal/threadapi"
	"github.com/spf13/cobra"
)

func New(store *config.Store) cligen.ThreadReplyHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.ThreadReplyParams) error {
		format := string(p.Output)

		if format != "plain" && format != "json" {
			return fmt.Errorf("--output must be plain or json")
		}
		body, err := commandcontent.Read(p.Content, p.ContentFile, cio.In)
		if err != nil {
			return err
		}
		if strings.TrimSpace(body) == "" {
			return fmt.Errorf("reply content is required")
		}
		body, err = commandcontent.ToHTML(body, p.Markdown)
		if err != nil {
			return err
		}
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}
		props := openapi.ReplyInitialProps{Body: body}
		if p.ReplyTo != "" {
			props.ReplyTo = &p.ReplyTo
		}
		response, err := threadapi.Reply(ctx, client.OpenAPI, p.ThreadId, props)
		if err != nil {
			return err
		}
		if format == "json" {
			return output.JSON(cio.Out, response)
		}
		_, err = fmt.Fprintf(cio.Out, "Created reply: %s\n", response.Id)
		return err
	}
}
