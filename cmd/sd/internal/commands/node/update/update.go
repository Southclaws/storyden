package update

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
)

type UpdateCommand *cobra.Command

func New(store *config.Store) UpdateCommand {
	var name string
	var slug string
	var description string
	var content string
	var contentFile string
	var markdown bool
	var url string
	var clearURL bool
	var hideChildTree bool
	var jsonInput string
	var tags []string

	command := &cobra.Command{
		Use:   "update <slug>",
		Short: "Update an existing node",
		Long: `# Update a Node

Update an existing node's properties. Only the fields you specify will be changed - everything else stays the same.

Content is stored as HTML. By default ` + "`--content`" + ` and ` + "`--content-file`" + ` are treated as HTML; pass ` + "`--markdown`" + ` to author them as Markdown, which is converted to HTML before sending.

## Examples

Update just the name:
~~~bash
sd node update my-page --name "New Title"
~~~

Update content from a file (HTML):
~~~bash
sd node update my-page --content-file updated.html
~~~

Update content from stdin (HTML):
~~~bash
echo "<p>New Content</p>" | sd node update my-page --content-file -
~~~

Update content authored as Markdown (converted to HTML):
~~~bash
sd node update my-page --markdown --content-file updated.md
~~~

Update multiple fields:
~~~bash
sd node update my-page --name "Better Title" --description "Improved description"
~~~

Change the slug (careful - breaks existing URLs):
~~~bash
sd node update old-slug --slug new-slug
~~~

Set an external URL:
~~~bash
sd node update my-page --url https://example.com
~~~

Hide a node's children from tree views:
~~~bash
sd node update my-page --hide-child-tree
~~~

Update with a raw JSON body from stdin:
~~~bash
jq '.node_data' file.json | sd node update my-page --json -
~~~

The positional slug chooses the node to update. In JSON mode, the input object is sent as the update body and can include API fields such as "name", "description", "content", "url", "hide_child_tree", "meta", "asset_ids", or "primary_image_asset_id".
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeSlug := args[0]

			props, err := buildMutableProps(cmd, mutablePropsInput{
				name:          name,
				slug:          slug,
				description:   description,
				content:       content,
				contentFile:   contentFile,
				markdown:      markdown,
				url:           url,
				clearURL:      clearURL,
				hideChildTree: hideChildTree,
				jsonInput:     jsonInput,
				tags:          tags,
			})
			if err != nil {
				return err
			}

			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			node, err := nodeapi.Update(cmd.Context(), client.OpenAPI, nodeSlug, props)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Updated node: %s (slug: %s)\n", node.Name, node.Slug)

			return nil
		},
	}

	command.Flags().StringVar(&name, "name", "", "New node name")
	command.Flags().StringVar(&slug, "slug", "", "New node slug")
	command.Flags().StringVar(&description, "description", "", "New node description")
	command.Flags().StringVar(&content, "content", "", "New node content as HTML (or Markdown with --markdown)")
	command.Flags().StringVar(&contentFile, "content-file", "", "Read content from file (use - for stdin)")
	command.Flags().BoolVar(&markdown, "markdown", false, "Treat content as Markdown and convert it to HTML")
	command.Flags().StringVar(&url, "url", "", "External URL for the node")
	command.Flags().BoolVar(&clearURL, "clear-url", false, "Clear the node's external URL")
	command.Flags().BoolVar(&hideChildTree, "hide-child-tree", false, "Hide this node's children from tree views")
	command.Flags().StringVar(&jsonInput, "json", "", "Read raw node update JSON from file (use - for stdin)")
	command.Flags().StringSliceVar(&tags, "tags", nil, "New tags (comma-separated)")

	help.SetupMarkdownHelp(command)

	return UpdateCommand(command)
}

type mutablePropsInput struct {
	name          string
	slug          string
	description   string
	content       string
	contentFile   string
	markdown      bool
	url           string
	clearURL      bool
	hideChildTree bool
	jsonInput     string
	tags          []string
}

func buildMutableProps(cmd *cobra.Command, input mutablePropsInput) (openapi.NodeMutableProps, error) {
	if input.jsonInput != "" {
		if changed := changedFieldFlags(cmd); len(changed) > 0 {
			return openapi.NodeMutableProps{}, fmt.Errorf("cannot combine --json with %s", strings.Join(changed, ", "))
		}

		return readJSONProps(input.jsonInput, cmd.InOrStdin())
	}

	if input.url != "" && input.clearURL {
		return openapi.NodeMutableProps{}, fmt.Errorf("cannot specify both --url and --clear-url")
	}

	finalContent, err := readContent(input.content, input.contentFile, cmd.InOrStdin())
	if err != nil {
		return openapi.NodeMutableProps{}, err
	}

	finalContent, err = contentToHTML(finalContent, input.markdown)
	if err != nil {
		return openapi.NodeMutableProps{}, err
	}

	props := openapi.NodeMutableProps{
		Name:        stringPtr(input.name),
		Slug:        stringPtr(input.slug),
		Description: stringPtr(input.description),
		Content:     stringPtr(finalContent),
	}

	if len(input.tags) > 0 {
		props.Tags = &input.tags
	}
	if cmd.Flags().Changed("url") {
		props.Url.Set(input.url)
	}
	if input.clearURL {
		props.Url.SetNull()
	}
	if cmd.Flags().Changed("hide-child-tree") {
		props.HideChildTree = &input.hideChildTree
	}

	return props, nil
}

func changedFieldFlags(cmd *cobra.Command) []string {
	names := []string{
		"name",
		"slug",
		"description",
		"content",
		"content-file",
		"url",
		"clear-url",
		"hide-child-tree",
		"tags",
	}

	var changed []string
	for _, name := range names {
		if cmd.Flags().Changed(name) {
			changed = append(changed, "--"+name)
		}
	}

	return changed
}

func readJSONProps(source string, stdin io.Reader) (openapi.NodeMutableProps, error) {
	var data []byte

	if source == "-" {
		bytes, err := io.ReadAll(stdin)
		if err != nil {
			return openapi.NodeMutableProps{}, fmt.Errorf("failed to read JSON from stdin: %w", err)
		}
		data = bytes
	} else {
		bytes, err := os.ReadFile(source)
		if err != nil {
			return openapi.NodeMutableProps{}, fmt.Errorf("failed to read JSON file: %w", err)
		}
		data = bytes
	}

	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err != nil {
		return openapi.NodeMutableProps{}, fmt.Errorf("invalid node update JSON: %w", err)
	}
	if object == nil {
		return openapi.NodeMutableProps{}, fmt.Errorf("node update JSON must be an object")
	}

	var props openapi.NodeMutableProps
	if err := json.Unmarshal(data, &props); err != nil {
		return openapi.NodeMutableProps{}, fmt.Errorf("invalid node update JSON: %w", err)
	}

	return props, nil
}

// contentToHTML converts Markdown content to HTML when the --markdown flag is
// set, reusing the same converter the backend uses. Without the flag, content
// passes through unchanged because the API already expects HTML.
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
