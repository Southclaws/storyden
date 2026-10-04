package create

import (
	"fmt"
	"strings"

	domainvisibility "github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	commandcontent "github.com/Southclaws/storyden/cmd/sd/internal/content"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
	"github.com/Southclaws/storyden/cmd/sd/internal/threadapi"
	"github.com/spf13/cobra"
)

type CreateCommand *cobra.Command

func New(store *config.Store) CreateCommand {
	var title, content, file, visibility, category, format string
	var markdown bool
	cmd := &cobra.Command{Use: "create", Short: "Create a thread as the selected account", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(title) == "" {
				return fmt.Errorf("--title is required")
			}
			if format != "plain" && format != "json" {
				return fmt.Errorf("--format must be plain or json")
			}
			if _, err := domainvisibility.NewVisibility(visibility); err != nil {
				return fmt.Errorf("invalid --visibility: %s", visibility)
			}
			body, err := commandcontent.Read(content, file, cmd.InOrStdin())
			if err != nil {
				return err
			}
			body, err = commandcontent.ToHTML(body, markdown)
			if err != nil {
				return err
			}
			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}
			vis := openapi.Visibility(visibility)
			props := openapi.ThreadInitialProps{Title: title, Body: &body, Visibility: &vis}
			if category != "" {
				props.Category = &category
			}
			response, err := threadapi.Create(cmd.Context(), client.OpenAPI, props)
			if err != nil {
				return err
			}
			if format == "json" {
				return output.JSON(cmd.OutOrStdout(), response)
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Created thread: %s (%s)\n", response.Title, response.Id)
			return err
		}}
	cmd.Flags().StringVar(&title, "title", "", "Thread title (required)")
	cmd.Flags().StringVar(&content, "content", "", "Content as HTML (or Markdown with --markdown)")
	cmd.Flags().StringVar(&file, "content-file", "", "Read content from file (use - for stdin)")
	cmd.Flags().BoolVar(&markdown, "markdown", false, "Convert Markdown content to HTML")
	cmd.Flags().StringVar(&visibility, "visibility", "published", "Visibility: draft, review, published, unlisted")
	cmd.Flags().StringVar(&category, "category", "", "Category ID")
	cmd.Flags().StringVar(&format, "format", "plain", "Output format: plain, json")
	help.SetupMarkdownHelp(cmd)
	return CreateCommand(cmd)
}
