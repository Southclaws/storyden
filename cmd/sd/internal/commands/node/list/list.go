package list

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	domainvisibility "github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/listflags"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/filter"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
	outputfmt "github.com/Southclaws/storyden/cmd/sd/internal/output"
	"github.com/Southclaws/storyden/cmd/sd/internal/render"
)

type ListCommand *cobra.Command

// serverQuery groups every server-side query parameter the user can set via
// flags, so the fetch closure stays tidy.
type serverQuery struct {
	author     string // deprecated alias for owner handle, also reused as server-side filter
	visibility []string
	search     string
	nodeID     string
	depth      *int
	nodeFormat string // tree or flat
}

func New(store *config.Store) ListCommand {
	flags := &listflags.Flags{}
	filterFlags := &filter.NodeFlags{}
	var author string
	var visibility []string
	var search string
	var nodeID string
	var parent string
	var depth int
	var depthSet bool
	var nodeFormat string

	command := &cobra.Command{
		Use: "list",
		Long: `# List Nodes

Browse all nodes with plain output, JSON, or JSONL.

The default ` + "`auto`" + ` format prints a plain table that is easy to read in terminals and scripts. Use ` + "`sd tui`" + ` for the interactive explorer.

## Examples

List nodes:
~~~bash
sd node list
~~~

Stream every page as JSONL:
~~~bash
sd node list --all --format jsonl
~~~

Filter to root-level nodes only:
~~~bash
sd node list --root-only
~~~

Filter by linked domain (repeatable):
~~~bash
sd node list --link-domain youtube.com --link-domain youtu.be
~~~

Filter by visibility (comma-separated):
~~~bash
sd node list --visibility review,draft
~~~

Filter by owner handle and name substring:
~~~bash
sd node list --owner-handle southclaws --name-contains design
~~~

List all nodes under a specific parent (by slug):
~~~bash
sd node list --parent web-development
~~~

Full-text search on the server:
~~~bash
sd node list --search "design system"
~~~

Stop after the first 10 matches across pages:
~~~bash
sd node list --all --limit 10 --link-domain youtube.com
~~~

Export to JSON:
~~~bash
sd node list --format json > nodes.json
~~~
`,
		Short: "List Storyden nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := flags.Validate(); err != nil {
				return err
			}
			if err := filterFlags.Validate(); err != nil {
				return err
			}
			if err := validateVisibilities(visibility); err != nil {
				return err
			}

			// --author is the back-compat alias for --owner-handle. Resolve
			// precedence: explicit --owner-handle wins.
			if filterFlags.OwnerHandle == "" && author != "" {
				filterFlags.OwnerHandle = author
			}

			if parent != "" && nodeID != "" {
				return fmt.Errorf("--parent and --node-id are mutually exclusive")
			}

			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			// --parent resolves a human-readable slug to the node ID required
			// by the server-side NodeId filter.
			if parent != "" {
				parentNode, err := nodeapi.Fetch(cmd.Context(), client.OpenAPI, parent)
				if err != nil {
					return fmt.Errorf("could not find parent node %q: %w", parent, err)
				}
				nodeID = string(parentNode.Id)
			}

			if nodeFormat != "" && nodeFormat != "tree" && nodeFormat != "flat" {
				return fmt.Errorf("--node-format must be tree or flat")
			}

			query := serverQuery{
				author:     author,
				visibility: visibility,
				search:     search,
				nodeID:     nodeID,
				nodeFormat: nodeFormat,
			}
			if depthSet {
				d := depth
				query.depth = &d
			}

			fetch := func(page int) (*openapi.NodeListResult, error) {
				return fetchNodes(cmd.Context(), client.OpenAPI, page, query)
			}

			return run(cmd.OutOrStdout(), flags, filterFlags.Build(), fetch)
		},
	}

	flags.Bind(command)
	filterFlags.Bind(command)
	command.Flags().StringVar(&author, "author", "", "Filter by author handle (alias for --owner-handle)")
	command.Flags().StringSliceVar(&visibility, "visibility", nil, "Filter by visibility (comma-separated: draft, review, published, unlisted)")
	command.Flags().StringVarP(&search, "search", "q", "", "Server-side full-text search query")
	command.Flags().StringVar(&nodeID, "node-id", "", "Limit to descendants of this node id")
	command.Flags().StringVar(&parent, "parent", "", "Limit to children of this node (accepts slug; resolves to node id automatically)")
	command.Flags().IntVar(&depth, "depth", 0, "Maximum child depth to return (0 = root nodes only)")
	command.Flags().StringVar(&nodeFormat, "node-format", "flat", "Server response shape: flat (default, every match) or tree (roots with nested children)")

	// Track whether --depth was explicitly set; cobra has no native way to
	// distinguish 0-from-default and 0-from-user, so we read it post-parse.
	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		depthSet = cmd.Flags().Changed("depth")
		return nil
	}

	command.Flags().MarkDeprecated("author", "use --owner-handle")

	help.SetupMarkdownHelp(command)

	return ListCommand(command)
}

func run(out io.Writer, flags *listflags.Flags, opts filter.NodeOptions, fetch func(int) (*openapi.NodeListResult, error)) error {
	format := flags.ResolveFormat(out)

	switch format {
	case listflags.FormatJSON:
		if flags.All {
			return runJSONAll(out, flags, opts, fetch)
		}
		result, err := fetch(flags.Page)
		if err != nil {
			return err
		}
		// Apply client-side filters to the single-page JSON output too, so
		// the shape is consistent regardless of whether --all is set.
		result.Nodes = filter.FilterNodes(result.Nodes, opts)
		if flags.Limit > 0 && len(result.Nodes) > flags.Limit {
			result.Nodes = result.Nodes[:flags.Limit]
		}
		return outputfmt.JSON(out, result)

	case listflags.FormatJSONL:
		return runJSONL(out, flags, opts, fetch)

	case listflags.FormatPlain:
		return runPlain(out, flags, opts, fetch)

	default:
		return fmt.Errorf("unsupported format %q", flags.Format)
	}
}

func runPlain(out io.Writer, flags *listflags.Flags, opts filter.NodeOptions, fetch func(int) (*openapi.NodeListResult, error)) error {
	if flags.All {
		all := []openapi.NodeWithChildren{}
		err := iterPages(flags, fetch, func(page *openapi.NodeListResult) (bool, error) {
			for _, n := range page.Nodes {
				if !filter.MatchNode(n, opts) {
					continue
				}
				all = append(all, n)
				if flags.Limit > 0 && len(all) >= flags.Limit {
					return false, nil
				}
			}
			return true, nil
		})
		if err != nil {
			return err
		}
		return render.Render(out, all, nodeProfile(), flags.Wide(), render.PageInfo{})
	}

	result, err := fetch(flags.Page)
	if err != nil {
		return err
	}
	nodes := filter.FilterNodes(result.Nodes, opts)
	if flags.Limit > 0 && len(nodes) > flags.Limit {
		nodes = nodes[:flags.Limit]
	}
	page := render.PageInfo{
		CurrentPage: result.CurrentPage,
		TotalPages:  result.TotalPages,
		PageSize:    result.PageSize,
		Results:     result.Results,
	}
	return render.Render(out, nodes, nodeProfile(), flags.Wide(), page)
}

func runJSONL(out io.Writer, flags *listflags.Flags, opts filter.NodeOptions, fetch func(int) (*openapi.NodeListResult, error)) error {
	emitted := 0
	encoder := json.NewEncoder(out)
	emit := func(page *openapi.NodeListResult) (bool, error) {
		for _, n := range page.Nodes {
			if !filter.MatchNode(n, opts) {
				continue
			}
			if err := encoder.Encode(n); err != nil {
				return false, err
			}
			emitted++
			if flags.Limit > 0 && emitted >= flags.Limit {
				return false, nil
			}
		}
		return true, nil
	}
	if flags.All {
		return iterPages(flags, fetch, emit)
	}
	result, err := fetch(flags.Page)
	if err != nil {
		return err
	}
	_, err = emit(result)
	return err
}

func runJSONAll(out io.Writer, flags *listflags.Flags, opts filter.NodeOptions, fetch func(int) (*openapi.NodeListResult, error)) error {
	all := []openapi.NodeWithChildren{}
	err := iterPages(flags, fetch, func(page *openapi.NodeListResult) (bool, error) {
		for _, n := range page.Nodes {
			if !filter.MatchNode(n, opts) {
				continue
			}
			all = append(all, n)
			if flags.Limit > 0 && len(all) >= flags.Limit {
				return false, nil
			}
		}
		return true, nil
	})
	if err != nil {
		return err
	}
	return outputfmt.JSON(out, struct {
		Nodes []openapi.NodeWithChildren `json:"nodes"`
	}{Nodes: all})
}

func iterPages(flags *listflags.Flags, fetch func(int) (*openapi.NodeListResult, error), onPage func(*openapi.NodeListResult) (bool, error)) error {
	page := flags.Page
	for {
		result, err := fetch(page)
		if err != nil {
			return err
		}
		keep, err := onPage(result)
		if err != nil {
			return err
		}
		if !keep {
			return nil
		}
		if result.NextPage == nil {
			return nil
		}
		page = *result.NextPage
	}
}

func nodeProfile() render.Profile[openapi.NodeWithChildren] {
	return render.Profile[openapi.NodeWithChildren]{
		Columns: []render.Column[openapi.NodeWithChildren]{
			{Header: "NAME", Render: func(n openapi.NodeWithChildren) string { return string(n.Name) }},
			{Header: "UPDATED", Render: func(n openapi.NodeWithChildren) string { return render.FormatTime(n.UpdatedAt) }},
			{Header: "AUTHOR", Render: func(n openapi.NodeWithChildren) string { return render.AuthorName(n.Owner) }},
			{Header: "CHILDREN", Render: func(n openapi.NodeWithChildren) string { return strconv.Itoa(len(n.Children)) }, Wide: true},
			{Header: "VISIBILITY", Render: func(n openapi.NodeWithChildren) string { return string(n.Visibility) }, Wide: true},
			{Header: "PARENT", Render: func(n openapi.NodeWithChildren) string {
				if n.Parent != nil {
					return string(n.Parent.Slug)
				}
				return ""
			}, Wide: true},
			{Header: "SLUG", Render: func(n openapi.NodeWithChildren) string { return string(n.Slug) }, Wide: true},
		},
	}
}

func fetchNodes(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	page int,
	q serverQuery,
) (*openapi.NodeListResult, error) {
	pageQuery := openapi.PaginationQuery(strconv.Itoa(page))

	params := &openapi.NodeListParams{
		Page: &pageQuery,
	}

	if q.author != "" {
		handle := openapi.AccountHandle(q.author)
		params.Author = &handle
	}

	if len(q.visibility) > 0 {
		vp := make(openapi.VisibilityParam, 0, len(q.visibility))
		for _, v := range q.visibility {
			vp = append(vp, openapi.Visibility(v))
		}
		params.Visibility = &vp
	}

	if q.search != "" {
		s := openapi.SearchQuery(q.search)
		params.Q = &s
	}

	if q.nodeID != "" {
		id := openapi.Identifier(q.nodeID)
		params.NodeId = &id
	}

	if q.depth != nil {
		d := openapi.TreeDepthParam(strconv.Itoa(*q.depth))
		params.Depth = &d
	}

	if q.nodeFormat != "" {
		// "flat" makes nested matches surface in the top-level list so client
		// filters and --limit work intuitively for triage. Tree is still
		// available for callers who want the hierarchy preserved.
		f := openapi.NodeListParamsFormat(q.nodeFormat)
		params.Format = &f

		// Default to a deep traversal when no explicit depth is set.
		// For flat format this surfaces all descendants as a flat list.
		// For tree format scoped via --parent/--node-id this ensures children
		// are returned (without depth the server defaults to depth=0 and
		// returns only the root node itself).
		if q.depth == nil {
			d := openapi.TreeDepthParam("10")
			params.Depth = &d
		}
	}

	response, err := client.NodeListWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, nodeListError(response)
	}

	return response.JSON200, nil
}

func nodeListError(response *openapi.NodeListResponse) error {
	if response.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("node list request was not authorised; run sd auth login again")
	}

	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("node list request failed: %s: %s", response.Status(), body)
	}

	return fmt.Errorf("node list request failed: %s", response.Status())
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
