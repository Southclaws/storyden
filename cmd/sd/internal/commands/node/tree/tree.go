package tree

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/spf13/cobra"

	domainvisibility "github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/filter"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

type TreeCommand *cobra.Command

func New(store *config.Store) TreeCommand {
	var depth int
	var visibility []string
	filterFlags := &filter.NodeFlags{}

	command := &cobra.Command{
		Use:   "tree",
		Short: "Show nodes as a tree",
		Long: `# Display Node Tree

Visualize your entire content hierarchy as a tree.

The tree shows all nodes with ` + "`└──`" + ` and ` + "`├──`" + ` branch indicators, displaying each node's name and slug.

## Examples

Show full tree:
~~~bash
sd node tree
~~~

Limit depth:
~~~bash
sd node tree --depth 2
~~~

Filter by visibility (comma-separated):
~~~bash
sd node tree --visibility published,unlisted
~~~

Show only subtrees that match a name:
~~~bash
sd node tree --name-contains design
~~~

Generate site map:
~~~bash
sd node tree > sitemap.txt
~~~

Use ` + "`sd node children`" + ` to list just the direct children of a specific node.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if depth < 0 {
				return fmt.Errorf("--depth must be greater than or equal to zero")
			}
			if err := validateVisibilities(visibility); err != nil {
				return err
			}
			if err := filterFlags.Validate(); err != nil {
				return err
			}

			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			result, err := fetchTree(cmd.Context(), client.OpenAPI, depth, visibility)
			if err != nil {
				return err
			}

			opts := filterFlags.Build()
			nodes := result.Nodes
			if !opts.Empty() {
				nodes = pruneTree(nodes, opts)
			}

			return renderTree(cmd.OutOrStdout(), nodes)
		},
	}

	command.Flags().IntVar(&depth, "depth", 10, "Maximum child depth to request")
	command.Flags().StringSliceVar(&visibility, "visibility", nil, "Filter by visibility (comma-separated: draft, review, published, unlisted)")
	filterFlags.Bind(command)

	help.SetupMarkdownHelp(command)

	return TreeCommand(command)
}

func fetchTree(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	depth int,
	visibility []string,
) (*openapi.NodeListResult, error) {
	depthParam := openapi.TreeDepthParam(strconv.Itoa(depth))
	format := openapi.NodeListParamsFormatTree
	params := &openapi.NodeListParams{
		Depth:  &depthParam,
		Format: &format,
	}

	if len(visibility) > 0 {
		vp := make(openapi.VisibilityParam, 0, len(visibility))
		for _, v := range visibility {
			vp = append(vp, openapi.Visibility(v))
		}
		params.Visibility = &vp
	}

	response, err := client.NodeListWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, nodeTreeError(response)
	}

	return response.JSON200, nil
}

func nodeTreeError(response *openapi.NodeListResponse) error {
	if response.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("node tree request was not authorised; run sd auth login again")
	}

	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("node tree request failed: %s: %s", response.Status(), body)
	}

	return fmt.Errorf("node tree request failed: %s", response.Status())
}

func validateVisibilities(values []string) error {
	for _, v := range values {
		if v == "" {
			continue
		}
		if _, err := domainvisibility.NewVisibility(v); err != nil {
			return fmt.Errorf("invalid --visibility: %s", v)
		}
	}
	return nil
}

// pruneTree returns the subset of nodes (recursively) that match opts or have
// any descendant that matches. Pruning preserves the original parent/child
// linkage so the rendered tree shows context, not just the matches.
func pruneTree(nodes []openapi.NodeWithChildren, opts filter.NodeOptions) []openapi.NodeWithChildren {
	out := []openapi.NodeWithChildren{}
	for _, n := range nodes {
		kept, ok := prunedNode(n, opts)
		if ok {
			out = append(out, kept)
		}
	}
	return out
}

func prunedNode(n openapi.NodeWithChildren, opts filter.NodeOptions) (openapi.NodeWithChildren, bool) {
	prunedChildren := pruneTree(n.Children, opts)

	if !filter.MatchNode(n, opts) && len(prunedChildren) == 0 {
		return n, false
	}

	n.Children = prunedChildren
	return n, true
}

func renderTree(out io.Writer, nodes []openapi.NodeWithChildren) error {
	fmt.Fprintln(out, ".")

	for i, node := range nodes {
		renderNode(out, node, "", i == len(nodes)-1)
	}

	return nil
}

func renderNode(out io.Writer, node openapi.NodeWithChildren, prefix string, last bool) {
	connector := "├── "
	childPrefix := prefix + "│   "
	if last {
		connector = "└── "
		childPrefix = prefix + "    "
	}

	fmt.Fprintf(out, "%s%s%s\n", prefix, connector, nodeLabel(node))

	for i, child := range node.Children {
		renderNode(out, child, childPrefix, i == len(node.Children)-1)
	}
}

func nodeLabel(node openapi.NodeWithChildren) string {
	name := singleLine(string(node.Name))
	slug := singleLine(string(node.Slug))
	id := singleLine(string(node.Id))

	parts := []string{}
	if slug != "" {
		parts = append(parts, fmt.Sprintf("slug=%s", slug))
	}
	if id != "" {
		parts = append(parts, fmt.Sprintf("id=%s", id))
	}
	if len(parts) == 0 {
		return name
	}

	return fmt.Sprintf("%s [%s]", name, strings.Join(parts, " "))
}

func singleLine(value string) string {
	value = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}

		return r
	}, value)

	return strings.Join(strings.Fields(value), " ")
}
