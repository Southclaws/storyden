// Package search implements `sd search`, a general wrapper around the
// datagraph search endpoint.
package search

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/carapace-sh/carapace"
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

type options struct {
	Query      string
	Kinds      []string
	Authors    []string
	Categories []string
	Tags       []string
}

func New(store *config.Store) SearchCommand {
	flags := &listflags.Flags{}
	opts := &options{}

	command := &cobra.Command{
		Use:   "search <query>",
		Short: "Search all Storyden content",
		Long: `# Search

Search the Storyden datagraph across nodes, threads, replies, posts, and profiles.

## Examples

Search all content:
~~~bash
sd search "design system"
~~~

Search only nodes and threads:
~~~bash
sd search "release notes" --kind node --kind thread
~~~

Filter by author, category, and tag:
~~~bash
sd search "triage" --authors southclaws --categories docs --tags review
~~~

Stream all matches as JSONL:
~~~bash
sd search "agents" --all --format jsonl
~~~
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Query = args[0]
			if err := opts.validate(); err != nil {
				return err
			}
			if err := flags.Validate(); err != nil {
				return err
			}

			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			fetch := func(page int) (*openapi.DatagraphSearchResult, error) {
				return fetchSearch(cmd.Context(), client.OpenAPI, opts, page)
			}

			return run(cmd.OutOrStdout(), flags, fetch)
		},
	}

	command.Flags().StringSliceVar(&opts.Kinds, "kind", nil, "Filter by datagraph item kind: "+strings.Join(searchKinds, ", "))
	command.Flags().StringSliceVar(&opts.Authors, "authors", nil, "Filter by author account IDs or handles; repeat or comma-separate")
	command.Flags().StringSliceVar(&opts.Categories, "categories", nil, "Filter by category slugs; repeat or comma-separate")
	command.Flags().StringSliceVar(&opts.Tags, "tags", nil, "Filter by tag names; repeat or comma-separate")
	flags.Bind(command)
	carapace.Gen(command).FlagCompletion(carapace.ActionMap{
		"kind": carapace.ActionValues(searchKinds...),
	})
	help.SetupMarkdownHelp(command)

	return SearchCommand(command)
}

func (o *options) validate() error {
	if strings.TrimSpace(o.Query) == "" {
		return fmt.Errorf("search query must not be empty")
	}

	for _, kind := range o.Kinds {
		if _, ok := validKinds[strings.ToLower(strings.TrimSpace(kind))]; !ok {
			return fmt.Errorf("invalid --kind %q; must be one of: %s", kind, strings.Join(searchKinds, ", "))
		}
	}

	return nil
}

var searchKinds = []string{"post", "thread", "reply", "node", "collection", "profile", "event"}

var validKinds = mapFromSlice(searchKinds)

func mapFromSlice(values []string) map[string]struct{} {
	mapped := make(map[string]struct{}, len(values))
	for _, value := range values {
		mapped[value] = struct{}{}
	}
	return mapped
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
		result.Items = limitItems(result.Items, flags.Limit)
		return outputfmt.JSON(out, result)

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
		all := []searchRow{}
		err := iterPages(flags, fetch, func(page *openapi.DatagraphSearchResult) (bool, error) {
			for _, item := range page.Items {
				all = append(all, itemRow(item))
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
	items := limitItems(result.Items, flags.Limit)
	rows := make([]searchRow, len(items))
	for i, item := range items {
		rows[i] = itemRow(item)
	}
	return render.Render(out, rows, searchProfile(), flags.Wide(), render.PageInfo{
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
		for _, item := range page.Items {
			if err := encoder.Encode(item); err != nil {
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
	_, err = emit(&openapi.DatagraphSearchResult{Items: limitItems(result.Items, flags.Limit)})
	return err
}

func runJSONAll(out io.Writer, flags *listflags.Flags, fetch func(int) (*openapi.DatagraphSearchResult, error)) error {
	all := []openapi.DatagraphItem{}
	err := iterPages(flags, fetch, func(page *openapi.DatagraphSearchResult) (bool, error) {
		for _, item := range page.Items {
			all = append(all, item)
			if flags.Limit > 0 && len(all) >= flags.Limit {
				return false, nil
			}
		}
		return true, nil
	})
	if err != nil {
		return err
	}
	return outputfmt.JSON(out, openapi.DatagraphSearchResult{Items: all})
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

func limitItems(items []openapi.DatagraphItem, limit int) []openapi.DatagraphItem {
	if limit > 0 && len(items) > limit {
		return items[:limit]
	}
	return items
}

func fetchSearch(ctx context.Context, client *openapi.ClientWithResponses, opts *options, page int) (*openapi.DatagraphSearchResult, error) {
	pageQuery := openapi.PaginationQuery(strconv.Itoa(page))
	params := &openapi.DatagraphSearchParams{
		Q:    opt.New(openapi.RequiredSearchQuery(opts.Query)).Ptr(),
		Page: &pageQuery,
	}

	if len(opts.Kinds) > 0 {
		kinds := make(openapi.DatagraphKindQuery, len(opts.Kinds))
		for i, kind := range opts.Kinds {
			kinds[i] = openapi.DatagraphItemKind(strings.ToLower(strings.TrimSpace(kind)))
		}
		params.Kind = &kinds
	}
	if len(opts.Authors) > 0 {
		authors := make(openapi.DatagraphAuthorQuery, len(opts.Authors))
		for i, author := range opts.Authors {
			authors[i] = openapi.Identifier(strings.TrimSpace(author))
		}
		params.Authors = &authors
	}
	if len(opts.Categories) > 0 {
		categories := make(openapi.DatagraphCategoryQuery, len(opts.Categories))
		for i, category := range opts.Categories {
			categories[i] = openapi.Identifier(strings.TrimSpace(category))
		}
		params.Categories = &categories
	}
	if len(opts.Tags) > 0 {
		tags := make(openapi.TagNameListQueryParam, len(opts.Tags))
		for i, tag := range opts.Tags {
			tags[i] = openapi.TagName(strings.TrimSpace(tag))
		}
		params.Tags = &tags
	}

	response, err := client.DatagraphSearchWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, outputfmt.RequestErrorWithMessages("search request", response, response.Body, outputfmt.UnauthorizedMessage("search request"))
	}

	return response.JSON200, nil
}

type searchRow struct {
	Kind    string
	Title   string
	ID      string
	Slug    string
	Author  string
	Updated string
	Summary string
}

func searchProfile() render.Profile[searchRow] {
	return render.Profile[searchRow]{
		Columns: []render.Column[searchRow]{
			{Header: "KIND", Render: func(r searchRow) string { return r.Kind }},
			{Header: "TITLE", Render: func(r searchRow) string { return r.Title }},
			{Header: "ID", Render: func(r searchRow) string { return r.ID }},
			{Header: "SUMMARY", Render: func(r searchRow) string { return r.Summary }},
			{Header: "SLUG", Render: func(r searchRow) string { return r.Slug }, Wide: true},
			{Header: "AUTHOR", Render: func(r searchRow) string { return r.Author }, Wide: true},
			{Header: "UPDATED", Render: func(r searchRow) string { return r.Updated }, Wide: true},
		},
	}
}

func itemRow(item openapi.DatagraphItem) searchRow {
	kind, err := item.Discriminator()
	if err != nil {
		return searchRow{Kind: "unknown", Summary: err.Error()}
	}

	switch kind {
	case "node":
		node, err := item.AsDatagraphItemNode()
		if err != nil {
			return decodeErrorRow(kind, err)
		}
		return searchRow{
			Kind:    kind,
			Title:   string(node.Ref.Name),
			ID:      string(node.Ref.Id),
			Slug:    string(node.Ref.Slug),
			Author:  render.AuthorName(node.Ref.Owner),
			Updated: render.FormatTime(node.Ref.UpdatedAt),
			Summary: string(node.Ref.Description),
		}
	case "thread":
		thread, err := item.AsDatagraphItemThread()
		if err != nil {
			return decodeErrorRow(kind, err)
		}
		return searchRow{
			Kind:    kind,
			Title:   string(thread.Ref.Title),
			ID:      string(thread.Ref.Id),
			Slug:    string(thread.Ref.Slug),
			Author:  render.AuthorName(thread.Ref.Author),
			Updated: render.FormatTime(thread.Ref.UpdatedAt),
			Summary: stringValue(thread.Ref.Description),
		}
	case "post":
		post, err := item.AsDatagraphItemPost()
		if err != nil {
			return decodeErrorRow(kind, err)
		}
		return searchRow{
			Kind:    kind,
			Title:   string(post.Ref.Title),
			ID:      string(post.Ref.Id),
			Slug:    string(post.Ref.Slug),
			Author:  render.AuthorName(post.Ref.Author),
			Updated: render.FormatTime(post.Ref.UpdatedAt),
			Summary: stringValue(post.Ref.Description),
		}
	case "reply":
		reply, err := item.AsDatagraphItemReply()
		if err != nil {
			return decodeErrorRow(kind, err)
		}
		return searchRow{
			Kind:    kind,
			ID:      string(reply.Ref.Id),
			Slug:    string(reply.Ref.RootSlug),
			Author:  render.AuthorName(reply.Ref.Author),
			Updated: render.FormatTime(reply.Ref.UpdatedAt),
			Summary: stringValue(reply.Ref.Description),
		}
	case "profile":
		profile, err := item.AsDatagraphItemProfile()
		if err != nil {
			return decodeErrorRow(kind, err)
		}
		return searchRow{
			Kind:    kind,
			Title:   profile.Ref.Name,
			ID:      string(profile.Ref.Id),
			Slug:    "@" + profile.Ref.Handle,
			Summary: string(profile.Ref.Bio),
		}
	default:
		return searchRow{Kind: kind, Summary: "unsupported result shape"}
	}
}

func decodeErrorRow(kind string, err error) searchRow {
	return searchRow{Kind: kind, Summary: err.Error()}
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
