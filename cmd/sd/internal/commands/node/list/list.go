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
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligenconv"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/listflags"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/filter"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/render"
)

// serverQuery groups every server-side query parameter the user can set via
// flags, so the fetch closure stays tidy.
type serverQuery struct {
	visibility []string
	search     string
	nodeID     string
	depth      *int
	nodeFormat string // tree or flat
}

func New(store *config.Store) cligen.NodeListHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.NodeListParams) (cligen.NodeListResult, error) {
		flags := &listflags.Flags{
			Page:   p.Page,
			Limit:  p.Limit,
			All:    p.All,
			Format: string(p.Format),
			Output: string(p.Output),
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

		if err := flags.Validate(); err != nil {
			return cligen.NodeListResult{}, err
		}
		if err := validateVisibilities(p.Visibility); err != nil {
			return cligen.NodeListResult{}, err
		}

		nodeID := p.NodeId
		if p.Parent != "" && nodeID != "" {
			return cligen.NodeListResult{}, fmt.Errorf("--parent and --node-id are mutually exclusive")
		}

		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return cligen.NodeListResult{}, err
		}

		// --parent resolves a human-readable slug to the node ID required by
		// the server-side NodeId filter.
		if p.Parent != "" {
			parentNode, err := nodeapi.Fetch(ctx, client.OpenAPI, p.Parent)
			if err != nil {
				return cligen.NodeListResult{}, fmt.Errorf("could not find parent node %q: %w", p.Parent, err)
			}
			nodeID = string(parentNode.Id)
		}

		query := serverQuery{
			visibility: p.Visibility,
			search:     p.Search,
			nodeID:     nodeID,
			nodeFormat: string(p.NodeFormat),
		}
		depth := p.Depth
		query.depth = &depth

		fetch := func(page int) (*openapi.NodeListResult, error) {
			return fetchNodes(ctx, client.OpenAPI, page, query)
		}

		return run(cio.Out, flags, opts, fetch)
	}
}

// run dispatches on the resolved format. json's result is returned rather
// than written directly — codegen's generated RunE encodes it according to
// the OpenCLI-declared NodeListResult schema. Every other format still
// writes directly to out, same as before.
func run(out io.Writer, flags *listflags.Flags, opts filter.NodeOptions, fetch func(int) (*openapi.NodeListResult, error)) (cligen.NodeListResult, error) {
	format := flags.ResolveFormat(out)

	switch format {
	case listflags.FormatJSON:
		return runJSON(flags, opts, fetch)

	case listflags.FormatJSONL:
		return cligen.NodeListResult{}, runJSONL(out, flags, opts, fetch)

	case listflags.FormatPlain:
		return cligen.NodeListResult{}, runPlain(out, flags, opts, fetch)

	default:
		return cligen.NodeListResult{}, fmt.Errorf("unsupported format %q", flags.Format)
	}
}

func runJSON(flags *listflags.Flags, opts filter.NodeOptions, fetch func(int) (*openapi.NodeListResult, error)) (cligen.NodeListResult, error) {
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
			return cligen.NodeListResult{}, err
		}
		return cligenconv.Convert[cligen.NodeListResult](struct {
			Nodes []openapi.NodeWithChildren `json:"nodes"`
		}{Nodes: all})
	}

	page, err := fetch(flags.Page)
	if err != nil {
		return cligen.NodeListResult{}, err
	}
	// Apply client-side filters to the single-page JSON output too, so the
	// shape is consistent regardless of whether --all is set.
	page.Nodes = filter.FilterNodes(page.Nodes, opts)
	if flags.Limit > 0 && len(page.Nodes) > flags.Limit {
		page.Nodes = page.Nodes[:flags.Limit]
	}
	return cligenconv.Convert[cligen.NodeListResult](page)
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
