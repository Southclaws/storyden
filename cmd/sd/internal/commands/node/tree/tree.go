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
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligenconv"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/filter"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
)

func New(store *config.Store) cligen.NodeTreeHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.NodeTreeParams) (cligen.NodeListResult, error) {
		// codegen validates repeatable+choices flags as plain string lists
		// (each item isn't individually enum-checked yet), so this manual
		// check is still load-bearing, not just defensive.
		if err := validateVisibilities(p.Visibility); err != nil {
			return cligen.NodeListResult{}, err
		}
		if p.Depth < 0 {
			return cligen.NodeListResult{}, fmt.Errorf("--depth must be greater than or equal to zero")
		}

		opts := filter.NodeOptions{
			LinkDomains:     p.LinkDomain,
			LinkURLContains: p.LinkUrlContains,
			LinkScheme:      string(p.LinkScheme),
			NoLink:          p.NoLink,
			HasLink:         p.HasLink,
			RootOnly:        p.RootOnly,
			OwnerHandle:     p.OwnerHandle,
			NameContains:    p.NameContains,
		}

		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return cligen.NodeListResult{}, err
		}

		var rootLabel string
		var nodeID string
		if p.Slug != "" {
			root, err := nodeapi.Fetch(ctx, client.OpenAPI, p.Slug)
			if err != nil {
				return cligen.NodeListResult{}, fmt.Errorf("could not find node %q: %w", p.Slug, err)
			}
			rootLabel = string(root.Name) + " [slug=" + p.Slug + "]"
			nodeID = string(root.Id)
		}

		result, err := fetchTree(ctx, client.OpenAPI, p.Depth, p.Visibility, nodeID)
		if err != nil {
			return cligen.NodeListResult{}, err
		}

		nodes := result.Nodes
		if !opts.Empty() {
			nodes = pruneTree(nodes, opts)
		}

		switch p.Format {
		case cligen.NodeTreeFormatJson:
			return cligenconv.Convert[cligen.NodeListResult](struct {
				PageSize    int                        `json:"page_size"`
				Results     int                        `json:"results"`
				TotalPages  int                        `json:"total_pages"`
				CurrentPage int                        `json:"current_page"`
				Nodes       []openapi.NodeWithChildren `json:"nodes"`
			}{
				PageSize:    result.PageSize,
				Results:     result.Results,
				TotalPages:  result.TotalPages,
				CurrentPage: result.CurrentPage,
				Nodes:       nodes,
			})
		default:
			return cligen.NodeListResult{}, renderTree(cio.Out, nodes, rootLabel)
		}
	}
}

func fetchTree(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	depth int,
	visibility []string,
	nodeID string,
) (*openapi.NodeListResult, error) {
	depthParam := openapi.TreeDepthParam(strconv.Itoa(depth))
	format := openapi.NodeListParamsFormatTree
	params := &openapi.NodeListParams{
		Depth:  &depthParam,
		Format: &format,
	}

	if nodeID != "" {
		id := openapi.Identifier(nodeID)
		params.NodeId = &id
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

func renderTree(out io.Writer, nodes []openapi.NodeWithChildren, rootLabel string) error {
	if rootLabel != "" && len(nodes) == 1 {
		// When the API is scoped to a node by ID it returns the root node
		// itself as the single top-level item. Render the node's own label as
		// the tree root and recurse into its children, so the root line is not
		// repeated twice.
		fmt.Fprintln(out, nodeLabel(nodes[0]))
		for i, child := range nodes[0].Children {
			renderNode(out, child, "", i == len(nodes[0].Children)-1)
		}
		return nil
	}

	if rootLabel != "" {
		fmt.Fprintln(out, rootLabel)
	} else {
		fmt.Fprintln(out, ".")
	}

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
