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

	"github.com/spf13/cobra"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/listflags"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	outputfmt "github.com/Southclaws/storyden/cmd/sd/internal/output"
	"github.com/Southclaws/storyden/cmd/sd/internal/render"
)

type SearchCommand *cobra.Command

func New(store *config.Store) SearchCommand {
	flags := &listflags.Flags{}

	command := &cobra.Command{
		Use:   "search <query>",
		Short: "Full-text search for nodes",
		Long: `# Search Nodes

Run a full-text query against the datagraph and return matching nodes. Posts, threads, and other kinds are filtered out so the output shape matches the rest of the node list commands.

## Examples

Search for nodes containing a phrase:
~~~bash
sd node search "design system"
~~~

Stream every match across all pages as JSONL:
~~~bash
sd node search "agents" --all --format jsonl
~~~

Stop after the first 10 matches:
~~~bash
sd node search "go" --limit 10
~~~
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]
			if strings.TrimSpace(query) == "" {
				return fmt.Errorf("search query must not be empty")
			}
			if err := flags.Validate(); err != nil {
				return err
			}

			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			fetch := func(page int) (*openapi.DatagraphSearchResult, error) {
				return fetchSearch(cmd.Context(), client.OpenAPI, query, page)
			}

			return run(cmd.OutOrStdout(), flags, fetch)
		},
	}

	flags.Bind(command)
	help.SetupMarkdownHelp(command)

	return SearchCommand(command)
}

func run(out io.Writer, flags *listflags.Flags, fetch func(int) (*openapi.DatagraphSearchResult, error)) error {
	format := flags.ResolveFormat(out)

	switch format {
	case listflags.FormatJSON:
		if flags.All {
			return runJSONAll(out, flags, fetch)
		}
		result, err := fetch(flags.Page)
		if err != nil {
			return err
		}
		nodes := extractNodes(result)
		if flags.Limit > 0 && len(nodes) > flags.Limit {
			nodes = nodes[:flags.Limit]
		}
		return outputfmt.JSON(out, struct {
			Nodes []openapi.Node `json:"nodes"`
		}{Nodes: nodes})

	case listflags.FormatJSONL:
		return runJSONL(out, flags, fetch)

	case listflags.FormatPlain:
		return runPlain(out, flags, fetch)

	default:
		return fmt.Errorf("unsupported format %q", flags.Format)
	}
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

func runJSONAll(out io.Writer, flags *listflags.Flags, fetch func(int) (*openapi.DatagraphSearchResult, error)) error {
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
	return outputfmt.JSON(out, struct {
		Nodes []openapi.Node `json:"nodes"`
	}{Nodes: all})
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
