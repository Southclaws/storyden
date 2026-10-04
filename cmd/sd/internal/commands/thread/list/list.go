package list

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	outputfmt "github.com/Southclaws/storyden/cmd/sd/internal/output"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/listflags"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/render"
)

func New(store *config.Store) cligen.ThreadListHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.ThreadListParams) error {
		flags := &listflags.Flags{
			Page:   p.Page,
			Limit:  p.Limit,
			All:    p.All,
			Format: string(p.Output),
			Output: string(p.Columns),
		}
		if err := flags.Validate(); err != nil {
			return err
		}

		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		fetch := func(page int) (*openapi.ThreadListResult, error) {
			return fetchThreads(ctx, client.OpenAPI, page)
		}

		return run(cio.Out, flags, fetch)
	}
}

func run(out io.Writer, flags *listflags.Flags, fetch func(int) (*openapi.ThreadListResult, error)) error {
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
		if flags.Limit > 0 && len(result.Threads) > flags.Limit {
			result.Threads = result.Threads[:flags.Limit]
		}
		return outputfmt.JSON(out, result)

	case listflags.FormatJSONL:
		return runJSONL(out, flags, fetch)

	case listflags.FormatPlain:
		return runPlain(out, flags, fetch)

	default:
		return fmt.Errorf("unsupported format %q", flags.Format)
	}
}

func runPlain(out io.Writer, flags *listflags.Flags, fetch func(int) (*openapi.ThreadListResult, error)) error {
	if flags.All {
		all := []openapi.ThreadReference{}
		err := iterPages(flags, fetch, func(page *openapi.ThreadListResult) (bool, error) {
			all = append(all, page.Threads...)
			if flags.Limit > 0 && len(all) >= flags.Limit {
				all = all[:flags.Limit]
				return false, nil
			}
			return true, nil
		})
		if err != nil {
			return err
		}
		return render.Render(out, all, threadProfile(), flags.Wide(), render.PageInfo{})
	}

	result, err := fetch(flags.Page)
	if err != nil {
		return err
	}
	threads := result.Threads
	if flags.Limit > 0 && len(threads) > flags.Limit {
		threads = threads[:flags.Limit]
	}
	page := render.PageInfo{
		CurrentPage: result.CurrentPage,
		TotalPages:  result.TotalPages,
		PageSize:    result.PageSize,
		Results:     result.Results,
	}
	return render.Render(out, threads, threadProfile(), flags.Wide(), page)
}

func runJSONL(out io.Writer, flags *listflags.Flags, fetch func(int) (*openapi.ThreadListResult, error)) error {
	emitted := 0
	encoder := json.NewEncoder(out)
	emit := func(page *openapi.ThreadListResult) (bool, error) {
		for _, t := range page.Threads {
			if err := encoder.Encode(t); err != nil {
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

func runJSONAll(out io.Writer, flags *listflags.Flags, fetch func(int) (*openapi.ThreadListResult, error)) error {
	all := []openapi.ThreadReference{}
	err := iterPages(flags, fetch, func(page *openapi.ThreadListResult) (bool, error) {
		all = append(all, page.Threads...)
		if flags.Limit > 0 && len(all) >= flags.Limit {
			all = all[:flags.Limit]
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return err
	}
	return outputfmt.JSON(out, struct {
		Threads []openapi.ThreadReference `json:"threads"`
	}{Threads: all})
}

func iterPages(flags *listflags.Flags, fetch func(int) (*openapi.ThreadListResult, error), onPage func(*openapi.ThreadListResult) (bool, error)) error {
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

func threadProfile() render.Profile[openapi.ThreadReference] {
	return render.Profile[openapi.ThreadReference]{
		Columns: []render.Column[openapi.ThreadReference]{
			{Header: "TITLE", Render: func(t openapi.ThreadReference) string { return string(t.Title) }},
			{Header: "UPDATED", Render: func(t openapi.ThreadReference) string { return render.FormatTime(t.UpdatedAt) }},
			{Header: "AUTHOR", Render: func(t openapi.ThreadReference) string { return render.AuthorName(t.Author) }},
			{Header: "REPLIES", Render: func(t openapi.ThreadReference) string { return strconv.Itoa(t.ReplyStatus.Replies) }, Wide: true},
			{Header: "VISIBILITY", Render: func(t openapi.ThreadReference) string { return string(t.Visibility) }, Wide: true},
			{Header: "CATEGORY", Render: func(t openapi.ThreadReference) string { return categoryName(t.Category) }, Wide: true},
			{Header: "SLUG", Render: func(t openapi.ThreadReference) string { return string(t.Slug) }, Wide: true},
		},
	}
}

func fetchThreads(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	page int,
) (*openapi.ThreadListResult, error) {
	pageQuery := openapi.PaginationQuery(strconv.Itoa(page))

	response, err := client.ThreadListWithResponse(ctx, &openapi.ThreadListParams{
		Page: &pageQuery,
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, threadListError(response)
	}

	return response.JSON200, nil
}

func threadListError(response *openapi.ThreadListResponse) error {
	if response.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("thread list request was not authorised; run sd auth login again")
	}

	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("thread list request failed: %s: %s", response.Status(), body)
	}

	return fmt.Errorf("thread list request failed: %s", response.Status())
}

func categoryName(category *openapi.CategoryReference) string {
	if category == nil {
		return ""
	}

	return category.Name
}
