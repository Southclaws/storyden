package children

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/listflags"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/filter"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	outputfmt "github.com/Southclaws/storyden/cmd/sd/internal/output"
	"github.com/Southclaws/storyden/cmd/sd/internal/render"
)

type ChildrenCommand *cobra.Command

func New(store *config.Store) ChildrenCommand {
	flags := &listflags.Flags{}
	filterFlags := &filter.NodeFlags{}
	var sort string

	command := &cobra.Command{
		Use:   "children <slug>",
		Short: "List children of a node",
		Long: `# List Node Children

List all direct children of a node. This shows only immediate children, not grandchildren.

## Examples

List children:
~~~bash
sd node children docs
~~~

Wide columns:
~~~bash
sd node children docs -o wide
~~~

Stream every page as JSONL:
~~~bash
sd node children docs --all --format jsonl
~~~

Get as JSON:
~~~bash
sd node children docs --format json
~~~

Sort children:
~~~bash
sd node children docs --sort name
~~~

Use ` + "`sd node tree`" + ` to see the full hierarchy including grandchildren.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			slug := args[0]

			if err := flags.Validate(); err != nil {
				return err
			}
			if err := filterFlags.Validate(); err != nil {
				return err
			}

			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			fetch := func(page int) (*openapi.NodeListResult, error) {
				return fetchChildren(cmd.Context(), client.OpenAPI, slug, page, sort)
			}

			return run(cmd.OutOrStdout(), flags, filterFlags.Build(), fetch)
		},
	}

	flags.Bind(command)
	filterFlags.Bind(command)
	command.Flags().StringVar(&sort, "sort", "", "Sort order")

	help.SetupMarkdownHelp(command)

	return ChildrenCommand(command)
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
		return render.Render(out, all, childProfile(), flags.Wide(), render.PageInfo{})
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
	return render.Render(out, nodes, childProfile(), flags.Wide(), page)
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

func childProfile() render.Profile[openapi.NodeWithChildren] {
	return render.Profile[openapi.NodeWithChildren]{
		Columns: []render.Column[openapi.NodeWithChildren]{
			{Header: "NAME", Render: func(n openapi.NodeWithChildren) string { return string(n.Name) }},
			{Header: "UPDATED", Render: func(n openapi.NodeWithChildren) string { return render.FormatTime(n.UpdatedAt) }},
			{Header: "VISIBILITY", Render: func(n openapi.NodeWithChildren) string { return string(n.Visibility) }},
			{Header: "SLUG", Render: func(n openapi.NodeWithChildren) string { return string(n.Slug) }},
			{Header: "CHILDREN", Render: func(n openapi.NodeWithChildren) string { return strconv.Itoa(len(n.Children)) }, Wide: true},
			{Header: "AUTHOR", Render: func(n openapi.NodeWithChildren) string { return render.AuthorName(n.Owner) }, Wide: true},
		},
	}
}

func fetchChildren(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	slug string,
	page int,
	sort string,
) (*openapi.NodeListResult, error) {
	pageQuery := openapi.PaginationQuery(strconv.Itoa(page))

	params := &openapi.NodeListChildrenParams{
		Page: &pageQuery,
	}

	if sort != "" {
		sortParam := openapi.NodeChildrenSortParam(sort)
		params.ChildrenSort = &sortParam
	}

	response, err := client.NodeListChildrenWithResponse(ctx, slug, params)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, childrenListError(response)
	}

	return response.JSON200, nil
}

func childrenListError(response *openapi.NodeListChildrenResponse) error {
	if response.StatusCode() == http.StatusNotFound {
		return fmt.Errorf("node not found")
	}

	if response.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("children list request was not authorised; run sd auth login again")
	}

	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("children list request failed: %s: %s", response.Status(), body)
	}

	return fmt.Errorf("children list request failed: %s", response.Status())
}
