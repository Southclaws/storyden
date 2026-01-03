package tools

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/mcp"
)

type searchTools struct {
	logger         *slog.Logger
	searcher       searcher.Searcher
	accountQuerier *account_querier.Querier
	categoryRepo   *category.Repository
}

func newSearchTools(
	logger *slog.Logger,
	registry *Registry,
	searcher searcher.Searcher,
	accountQuerier *account_querier.Querier,
	categoryRepo *category.Repository,
) *searchTools {
	t := &searchTools{
		logger:         logger,
		searcher:       searcher,
		accountQuerier: accountQuerier,
		categoryRepo:   categoryRepo,
	}

	registry.Register(t.newSearchTool())

	return t
}

func (st *searchTools) newSearchTool() *Tool {
	toolDef := mcp.GetSearchTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				st.ExecuteSearch,
			)
		},
	}
}

func (st *searchTools) ExecuteSearch(ctx tool.Context, args mcp.ToolSearchInput) (*mcp.ToolSearchOutput, error) {
	maxResults := 10
	if args.MaxResults != nil {
		maxResults = *args.MaxResults
	}

	var validationErrors []string

	var kindFilter opt.Optional[[]datagraph.Kind]
	if len(args.Kind) > 0 {
		kinds := make([]datagraph.Kind, 0, len(args.Kind))
		var invalidKinds []string
		for _, k := range args.Kind {
			kind, err := datagraph.NewKind(string(k))
			if err != nil {
				invalidKinds = append(invalidKinds, string(k))
				continue
			}
			kinds = append(kinds, kind)
		}
		if len(invalidKinds) > 0 {
			validationErrors = append(validationErrors, fmt.Sprintf("invalid kinds: %v", invalidKinds))
		}
		if len(kinds) > 0 {
			kindFilter = opt.New(kinds)
		}
	}

	var authorFilter opt.Optional[[]account.AccountID]
	if len(args.Authors) > 0 {
		accounts, err := st.accountQuerier.ProbeMany(ctx, args.Authors...)
		if err != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("failed to look up authors: %v", err))
		} else if len(accounts) < len(args.Authors) {
			foundHandles := dt.Map(accounts, func(acc *account.Account) string {
				return acc.Handle
			})
			var notFound []string
			for _, handle := range args.Authors {
				found := false
				for _, fh := range foundHandles {
					if fh == handle {
						found = true
						break
					}
				}
				if !found {
					notFound = append(notFound, handle)
				}
			}
			if len(notFound) > 0 {
				validationErrors = append(validationErrors, fmt.Sprintf("authors not found: %v", notFound))
			}
		}
		if len(accounts) > 0 {
			authorIDs := dt.Map(accounts, func(acc *account.Account) account.AccountID {
				return acc.ID
			})
			authorFilter = opt.New(authorIDs)
		}
	}

	var categoryFilter opt.Optional[[]category.CategoryID]
	if len(args.Categories) > 0 {
		categories, err := st.categoryRepo.GetCategories(ctx, false)
		if err != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("failed to look up categories: %v", err))
		} else {
			matchedIDs := []category.CategoryID{}
			var notMatched []string
			for _, searchName := range args.Categories {
				searchNameLower := strings.ToLower(searchName)
				found := false
				for _, cat := range categories {
					if strings.ToLower(cat.Name) == searchNameLower {
						matchedIDs = append(matchedIDs, cat.ID)
						found = true
						break
					}
				}
				if !found {
					notMatched = append(notMatched, searchName)
				}
			}
			if len(notMatched) > 0 {
				validationErrors = append(validationErrors, fmt.Sprintf("categories not found: %v", notMatched))
			}
			if len(matchedIDs) > 0 {
				categoryFilter = opt.New(matchedIDs)
			}
		}
	}

	var tagFilter opt.Optional[[]tag_ref.Name]
	if len(args.Tags) > 0 {
		tagNames := dt.Map(args.Tags, func(t string) tag_ref.Name {
			return tag_ref.NewName(t)
		})
		tagFilter = opt.New(tagNames)
	}

	opts := searcher.Options{
		Kinds:      kindFilter,
		Authors:    authorFilter,
		Categories: categoryFilter,
		Tags:       tagFilter,
	}

	pageParams := pagination.NewPageParams(1, uint(maxResults))

	result, err := st.searcher.Search(ctx, args.Query, pageParams, opts)
	if err != nil {
		return nil, err
	}

	items := make([]mcp.SearchedItem, 0, len(result.Items))
	for _, item := range result.Items {
		desc := item.GetDesc()
		items = append(items, mcp.SearchedItem{
			Id:          item.GetID().String(),
			Kind:        item.GetKind().String(),
			Slug:        item.GetSlug(),
			Name:        item.GetName(),
			Description: &desc,
		})
	}

	output := mcp.ToolSearchOutput{
		Results: result.Results,
		Items:   items,
	}

	if len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation errors: %s", strings.Join(validationErrors, "; "))
	}

	return &(output), nil
}
