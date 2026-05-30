package create

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	domainvisibility "github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

type CreateCommand *cobra.Command

func New(store *config.Store) CreateCommand {
	var name string
	var slug string
	var description string
	var content string
	var contentFile string
	var markdown bool
	var parent string
	var visibility string
	var tags []string
	var url string
	var hideChildTree bool

	command := &cobra.Command{
		Use:   "create",
		Short: "Create a new node",
		Long: `# Create a New Node

Create a new node (page, document, wiki entry) in Storyden.

Nodes are the fundamental content unit in Storyden. They can be organized hierarchically using parent-child relationships, and each node can contain rich-text HTML content, metadata, and custom properties.

Content is stored as HTML. By default, whatever you pass to ` + "`--content`" + ` or ` + "`--content-file`" + ` is parsed and sanitised as HTML, so ` + "`# Title`" + ` would be stored as literal text rather than a heading. Pass ` + "`--markdown`" + ` to author content as Markdown instead; it is converted to HTML before sending.

## Examples

Create a simple node:
~~~bash
sd node create --name "My Page"
~~~

Create a node with inline content (HTML):
~~~bash
sd node create --name "Quick Note" --content "<h1>Title</h1><p>Some content here</p>"
~~~

Create a node with content from a file (the file must contain HTML):
~~~bash
sd node create --name "Documentation" --content-file page.html
~~~

Create a node with Markdown content (converted to HTML):
~~~bash
sd node create --name "Notes" --markdown --content-file notes.md
~~~

Create a child node under a parent:
~~~bash
sd node create --name "Getting Started" --parent docs
~~~

Create a node that links to an external URL:
~~~bash
sd node create --name "Anthropic" --url https://anthropic.com
~~~

Create a node and hide its children from tree views:
~~~bash
sd node create --name "Internal" --hide-child-tree
~~~

Create a draft node with content from stdin:
~~~bash
echo "<p>Hello World</p>" | sd node create --name "Draft" --content-file - --visibility draft
~~~

Create a fully-specified node:
~~~bash
sd node create \
  --name "Tutorial" \
  --slug custom-tutorial-slug \
  --description "A comprehensive tutorial" \
  --content-file tutorial.md \
  --parent guides \
  --visibility published \
  --tags "beginner,tutorial,guide"
~~~

## About Slugs

Slugs are URL-friendly identifiers. If you don't provide one, it's auto-generated from the name. For example, "My Page" becomes "my-page".

## About Visibility

Choose visibility based on your workflow:
- **draft** - Work in progress, only you can see it
- **review** - Ready for moderator review
- **published** - Public and searchable (default for most content)
- **unlisted** - Public via direct link but not in listings
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if err := validateVisibility(visibility); err != nil {
				return err
			}

			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			finalContent, err := readContent(content, contentFile, cmd.InOrStdin())
			if err != nil {
				return err
			}

			finalContent, err = contentToHTML(finalContent, markdown)
			if err != nil {
				return err
			}

			node, err := createNode(cmd.Context(), client.OpenAPI, openapi.NodeInitialProps{
				Name:          name,
				Slug:          stringPtr(slug),
				Description:   stringPtr(description),
				Content:       stringPtr(finalContent),
				Parent:        stringPtr(parent),
				Visibility:    visibilityPtr(visibility),
				Tags:          tagsPtr(tags),
				Url:           stringPtr(url),
				HideChildTree: boolPtr(hideChildTree),
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Created node: %s (slug: %s)\n", node.Name, node.Slug)

			return nil
		},
	}

	command.Flags().StringVar(&name, "name", "", "Node name (required)")
	command.Flags().StringVar(&slug, "slug", "", "Node slug (auto-generated if not provided)")
	command.Flags().StringVar(&description, "description", "", "Node description")
	command.Flags().StringVar(&content, "content", "", "Node content as HTML (or Markdown with --markdown)")
	command.Flags().StringVar(&contentFile, "content-file", "", "Read content from file (use - for stdin)")
	command.Flags().BoolVar(&markdown, "markdown", false, "Treat content as Markdown and convert it to HTML")
	command.Flags().StringVar(&parent, "parent", "", "Parent node slug")
	command.Flags().StringVar(&visibility, "visibility", "", "Visibility: draft, review, published, unlisted")
	command.Flags().StringSliceVar(&tags, "tags", nil, "Tags (comma-separated)")
	command.Flags().StringVar(&url, "url", "", "External URL to attach as the node's link")
	command.Flags().BoolVar(&hideChildTree, "hide-child-tree", false, "Hide this node's children from tree views")

	help.SetupMarkdownHelp(command)

	return CreateCommand(command)
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

// contentToHTML converts Markdown content to HTML when the --markdown flag is
// set, reusing the same converter the backend uses so output matches what the
// server would produce. Without the flag, content passes through unchanged
// because the API already expects HTML.
func contentToHTML(content string, markdown bool) (string, error) {
	if !markdown || content == "" {
		return content, nil
	}

	rt, err := datagraph.NewRichTextFromMarkdown(content)
	if err != nil {
		return "", fmt.Errorf("failed to convert markdown content: %w", err)
	}

	return rt.HTML(), nil
}

func readContent(content string, contentFile string, stdin io.Reader) (string, error) {
	if content != "" && contentFile != "" {
		return "", fmt.Errorf("cannot specify both --content and --content-file")
	}

	if content != "" {
		return content, nil
	}

	if contentFile == "" {
		return "", nil
	}

	if contentFile == "-" {
		bytes, err := io.ReadAll(stdin)
		if err != nil {
			return "", fmt.Errorf("failed to read from stdin: %w", err)
		}

		return string(bytes), nil
	}

	bytes, err := os.ReadFile(contentFile)
	if err != nil {
		return "", fmt.Errorf("failed to read content file: %w", err)
	}

	return string(bytes), nil
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

func validateVisibility(visibility string) error {
	if visibility == "" {
		return nil
	}

	if _, err := domainvisibility.NewVisibility(visibility); err != nil {
		return fmt.Errorf("invalid --visibility: %s", visibility)
	}

	return nil
}

func tagsPtr(tags []string) *[]string {
	if len(tags) == 0 {
		return nil
	}

	return &tags
}
