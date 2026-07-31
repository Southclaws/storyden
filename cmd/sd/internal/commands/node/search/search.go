// Package search implements `sd node search <query>` — a thin wrapper around
// the datagraph search endpoint, filtered to node items. It uses the same
// shared list rendering as `sd node list` so column profiles and format flags
// are consistent.
package search

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Southclaws/opt"
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligenconv"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/listflags"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/render"
)

func New(store *config.Store) cligen.NodeSearchHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.NodeSearchParams) (cligen.SearchNodeListResult, error) {
		if strings.TrimSpace(p.Query) == "" {
			return cligen.SearchNodeListResult{}, fmt.Errorf("search query must not be empty")
		}

		flags := &listflags.Flags{
			Page:   p.Page,
			Limit:  p.Limit,
			All:    p.All,
			Format: string(p.Format),
			Output: string(p.Output),
		}
		if err := flags.Validate(); err != nil {
			return cligen.SearchNodeListResult{}, err
		}

		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return cligen.SearchNodeListResult{}, err
		}

		fetch := func(page int) (*openapi.DatagraphSearchResult, error) {
			return fetchSearch(ctx, client.OpenAPI, p.Query, page)
		}

		return run(cio.Out, flags, fetch)
	}
}

// run dispatches on the resolved format. json's result is returned rather
// than written directly — codegen's generated RunE encodes it according to
// the OpenCLI-declared SearchNodeListResult schema. Every other format still
// writes directly to out, same as before.
func run(out io.Writer, flags *listflags.Flags, fetch func(int) (*openapi.DatagraphSearchResult, error)) (cligen.SearchNodeListResult, error) {
	format := flags.ResolveFormat(out)

	switch format {
	case listflags.FormatJSON:
		return runJSON(flags, fetch)

	case listflags.FormatJSONL:
		return cligen.SearchNodeListResult{}, runJSONL(out, flags, fetch)

	case listflags.FormatPlain:
		return cligen.SearchNodeListResult{}, runPlain(out, flags, fetch)

	default:
		return cligen.SearchNodeListResult{}, fmt.Errorf("unsupported format %q", flags.Format)
	}
}

func runJSON(flags *listflags.Flags, fetch func(int) (*openapi.DatagraphSearchResult, error)) (cligen.SearchNodeListResult, error) {
	if flags.All {
		all := []openapi.Node{}
		err := iterPages(flags, fetch, func(page *openapi.DatagraphSearchResult) (bool, error) {
			for _, n := range extractNodes(page) {
				all = append(all, n)
				if flags.Limit > 0 && len(all) >= flags.Limit {
					return false, nil
				}
			}
			return true, nil
		})
		if err != nil {
			return cligen.SearchNodeListResult{}, err
		}
		return cligenconv.Convert[cligen.SearchNodeListResult](struct {
			Nodes []openapi.Node `json:"nodes"`
		}{Nodes: all})
	}

	result, err := fetch(flags.Page)
	if err != nil {
		return cligen.SearchNodeListResult{}, err
	}
	nodes := extractNodes(result)
	if flags.Limit > 0 && len(nodes) > flags.Limit {
		nodes = nodes[:flags.Limit]
	}
	// Single page: include pagination metadata, matching the spec's
	// SearchNodeListResult schema — consistent with every other list-shaped
	// command's single-page JSON, unlike the bare array this used to return
	// regardless of --all.
	return cligenconv.Convert[cligen.SearchNodeListResult](struct {
		PageSize    int            `json:"page_size"`
		Results     int            `json:"results"`
		TotalPages  int            `json:"total_pages"`
		CurrentPage int            `json:"current_page"`
		NextPage    *int           `json:"next_page,omitempty"`
		Nodes       []openapi.Node `json:"nodes"`
	}{
		PageSize:    result.PageSize,
		Results:     result.Results,
		TotalPages:  result.TotalPages,
		CurrentPage: result.CurrentPage,
		NextPage:    result.NextPage,
		Nodes:       nodes,
	})
}

func runPlain(out io.Writer, flags *listflags.Flags, fetch func(int) (*openapi.DatagraphSearchResult, error)) error {
	if flags.All {
		all := []openapi.Node{}
		err := iterPages(flags, fetch, func(page *openapi.DatagraphSearchResult) (bool, error) {
			for _, n := range extractNodes(page) {
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
		return render.Render(out, all, searchProfile(), flags.Wide(), render.PageInfo{})
	}

	result, err := fetch(flags.Page)
	if err != nil {
		return err
	}
	nodes := extractNodes(result)
	if flags.Limit > 0 && len(nodes) > flags.Limit {
		nodes = nodes[:flags.Limit]
	}
	return render.Render(out, nodes, searchProfile(), flags.Wide(), render.PageInfo{
		CurrentPage: result.CurrentPage,
		TotalPages:  result.TotalPages,
		PageSize:    result.PageSize,
		Results:     result.Results,
	})
}

func runJSONL(out io.Writer, flags *listflags.Flags, fetch func(int) (*openapi.DatagraphSearchResult, error)) error {
	emitted := 0
	encoder := json.NewEncoder(out)
	emit := func(page *openapi.DatagraphSearchResult) (bool, error) {
		for _, n := range extractNodes(page) {
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

func iterPages(flags *listflags.Flags, fetch func(int) (*openapi.DatagraphSearchResult, error), onPage func(*openapi.DatagraphSearchResult) (bool, error)) error {
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

// extractNodes pulls out the node-kind items from a mixed-kind search page.
// Non-node items are silently skipped so the output stays consistent with the
// rest of the node list commands.
func extractNodes(result *openapi.DatagraphSearchResult) []openapi.Node {
	nodes := make([]openapi.Node, 0, len(result.Items))
	for _, item := range result.Items {
		kind, err := item.Discriminator()
		if err != nil || kind != "node" {
			continue
		}
		ni, err := item.AsDatagraphItemNode()
		if err != nil {
			continue
		}
		nodes = append(nodes, ni.Ref)
	}
	return nodes
}

func searchProfile() render.Profile[openapi.Node] {
	return render.Profile[openapi.Node]{
		Columns: []render.Column[openapi.Node]{
			{Header: "NAME", Render: func(n openapi.Node) string { return string(n.Name) }},
			{Header: "UPDATED", Render: func(n openapi.Node) string { return render.FormatTime(n.UpdatedAt) }},
			{Header: "AUTHOR", Render: func(n openapi.Node) string { return render.AuthorName(n.Owner) }},
			{Header: "VISIBILITY", Render: func(n openapi.Node) string { return string(n.Visibility) }, Wide: true},
			{Header: "SLUG", Render: func(n openapi.Node) string { return string(n.Slug) }, Wide: true},
		},
	}
}

func fetchSearch(ctx context.Context, client *openapi.ClientWithResponses, query string, page int) (*openapi.DatagraphSearchResult, error) {
	pageQuery := openapi.PaginationQuery(strconv.Itoa(page))
	kind := openapi.DatagraphKindQuery{openapi.DatagraphItemKind("node")}
	params := &openapi.DatagraphSearchParams{
		Q:    opt.New(openapi.RequiredSearchQuery(query)).Ptr(),
		Kind: &kind,
		Page: &pageQuery,
	}

	response, err := client.DatagraphSearchWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, searchError(response)
	}

	return response.JSON200, nil
}

func searchError(response *openapi.DatagraphSearchResponse) error {
	if response.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("search request was not authorised; run sd auth login again")
	}
	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("search request failed: %s: %s", response.Status(), body)
	}
	return fmt.Errorf("search request failed: %s", response.Status())
}
