package reply

import (
	"fmt"
	"strings"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	commandcontent "github.com/Southclaws/storyden/cmd/sd/internal/content"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
	"github.com/Southclaws/storyden/cmd/sd/internal/threadapi"
	"github.com/spf13/cobra"
)

type ReplyCommand *cobra.Command

func New(store *config.Store) ReplyCommand {
	var content, file, replyTo, format string
	var markdown bool
	cmd := &cobra.Command{Use: "reply THREAD_ID", Short: "Reply to a thread as the selected account", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if format != "plain" && format != "json" {
				return fmt.Errorf("--format must be plain or json")
			}
			body, err := commandcontent.Read(content, file, cmd.InOrStdin())
			if err != nil {
				return err
			}
			if strings.TrimSpace(body) == "" {
				return fmt.Errorf("reply content is required")
			}
			body, err = commandcontent.ToHTML(body, markdown)
			if err != nil {
				return err
			}
			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}
			props := openapi.ReplyInitialProps{Body: body}
			if replyTo != "" {
				props.ReplyTo = &replyTo
			}
			response, err := threadapi.Reply(cmd.Context(), client.OpenAPI, args[0], props)
			if err != nil {
				return err
			}
			if format == "json" {
				return output.JSON(cmd.OutOrStdout(), response)
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Created reply: %s\n", response.Id)
			return err
		}}
	cmd.Flags().StringVar(&content, "content", "", "Content as HTML (or Markdown with --markdown)")
	cmd.Flags().StringVar(&file, "content-file", "", "Read content from file (use - for stdin)")
	cmd.Flags().BoolVar(&markdown, "markdown", false, "Convert Markdown content to HTML")
	cmd.Flags().StringVar(&replyTo, "reply-to", "", "Post ID to reply to within the thread")
	cmd.Flags().StringVar(&format, "format", "plain", "Output format: plain, json")
	help.SetupMarkdownHelp(cmd)
	return ReplyCommand(cmd)
}
