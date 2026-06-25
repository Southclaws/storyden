package tools

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
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
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/profile/profile_search"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/lib/mcp"
)

type searchTools struct {
	logger         *slog.Logger
	webAddress     url.URL
	searcher       searcher.Searcher
	accountQuerier *account_querier.Querier
	profileSearch  *profile_search.Querier
	categoryRepo   *category.Repository
}

func newSearchTools(
	cfg config.Config,
	logger *slog.Logger,
	registry *Registry,
	searcher searcher.Searcher,
	accountQuerier *account_querier.Querier,
	profileSearch *profile_search.Querier,
	categoryRepo *category.Repository,
) *searchTools {
	t := &searchTools{
		logger:         logger,
		webAddress:     cfg.PublicWebAddress,
		searcher:       searcher,
		accountQuerier: accountQuerier,
		profileSearch:  profileSearch,
		categoryRepo:   categoryRepo,
	}

	registry.Register(t.newContentSearchTool())
	registry.Register(t.newThreadSearchTool())
	registry.Register(t.newReplySearchTool())
	registry.Register(t.newPostSearchTool())
	registry.Register(t.newMemberSearchTool())

	return t
}

func (st *searchTools) newContentSearchTool() *Tool {
	toolDef := mcp.GetContentSearchTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				func(ctx tool.Context, args mcp.ToolContentSearchInput) (*mcp.ToolContentSearchOutput, error) {
					return st.ExecuteContentSearch(ctx, args)
				},
			)
		},
		Handler: makeHandler(st.ExecuteContentSearch),
	}
}

func (st *searchTools) ExecuteContentSearch(ctx context.Context, args mcp.ToolContentSearchInput) (*mcp.ToolContentSearchOutput, error) {
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

	authorFilter, authorErrors := st.resolveAuthors(ctx, args.Authors)
	validationErrors = append(validationErrors, authorErrors...)

	categoryFilter, categoryErrors := st.resolveCategories(ctx, args.Categories)
	validationErrors = append(validationErrors, categoryErrors...)

	tagFilter := resolveTags(args.Tags)

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
			BrowserUrl:  datagraph.CanonicalResolveURL(st.webAddress, item.GetKind(), fmt.Sprintf("%s-%s", item.GetID().String(), item.GetSlug())).String(),
			Id:          item.GetID().String(),
			Kind:        item.GetKind().String(),
			Slug:        item.GetSlug(),
			Name:        item.GetName(),
			Description: &desc,
		})
	}

	output := mcp.ToolContentSearchOutput{
		Results: result.Results,
		Items:   items,
	}

	if len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation errors: %s", strings.Join(validationErrors, "; "))
	}

	return &output, nil
}

func (st *searchTools) newThreadSearchTool() *Tool {
	toolDef := mcp.GetThreadSearchTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				func(ctx tool.Context, args mcp.ToolThreadSearchInput) (*mcp.ToolThreadSearchOutput, error) {
					return st.ExecuteThreadSearch(ctx, args)
				},
			)
		},
		Handler: makeHandler(st.ExecuteThreadSearch),
	}
}

func (st *searchTools) ExecuteThreadSearch(ctx context.Context, args mcp.ToolThreadSearchInput) (*mcp.ToolThreadSearchOutput, error) {
	maxResults := 10
	if args.MaxResults != nil {
		maxResults = *args.MaxResults
	}

	var validationErrors []string

	authorFilter, authorErrors := st.resolveAuthors(ctx, args.Authors)
	validationErrors = append(validationErrors, authorErrors...)

	categoryFilter, categoryErrors := st.resolveCategories(ctx, args.Categories)
	validationErrors = append(validationErrors, categoryErrors...)

	tagFilter := resolveTags(args.Tags)

	opts := searcher.Options{
		Kinds:      opt.New([]datagraph.Kind{datagraph.KindThread}),
		Authors:    authorFilter,
		Categories: categoryFilter,
		Tags:       tagFilter,
	}

	pageParams := pagination.NewPageParams(1, uint(maxResults))

	result, err := st.searcher.Search(ctx, args.Query, pageParams, opts)
	if err != nil {
		return nil, err
	}

	items := dt.Map(result.Items, func(item datagraph.Item) mcp.ThreadSearchItem {
		desc := item.GetDesc()
		return mcp.ThreadSearchItem{
			BrowserUrl:  datagraph.CanonicalResolveURL(st.webAddress, datagraph.KindThread, fmt.Sprintf("%s-%s", item.GetID().String(), item.GetSlug())).String(),
			Id:          item.GetID().String(),
			Slug:        item.GetSlug(),
			Name:        item.GetName(),
			Description: &desc,
		}
	})

	output := mcp.ToolThreadSearchOutput{
		Results: result.Results,
		Items:   items,
	}

	if len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation errors: %s", strings.Join(validationErrors, "; "))
	}

	return &output, nil
}

func (st *searchTools) newReplySearchTool() *Tool {
	toolDef := mcp.GetReplySearchTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				func(ctx tool.Context, args mcp.ToolReplySearchInput) (*mcp.ToolReplySearchOutput, error) {
					return st.ExecuteReplySearch(ctx, args)
				},
			)
		},
		Handler: makeHandler(st.ExecuteReplySearch),
	}
}

func (st *searchTools) ExecuteReplySearch(ctx context.Context, args mcp.ToolReplySearchInput) (*mcp.ToolReplySearchOutput, error) {
	maxResults := 10
	if args.MaxResults != nil {
		maxResults = *args.MaxResults
	}

	var validationErrors []string

	authorFilter, authorErrors := st.resolveAuthors(ctx, args.Authors)
	validationErrors = append(validationErrors, authorErrors...)

	tagFilter := resolveTags(args.Tags)

	opts := searcher.Options{
		Kinds:   opt.New([]datagraph.Kind{datagraph.KindReply}),
		Authors: authorFilter,
		Tags:    tagFilter,
	}

	pageParams := pagination.NewPageParams(1, uint(maxResults))

	result, err := st.searcher.Search(ctx, args.Query, pageParams, opts)
	if err != nil {
		return nil, err
	}

	items := dt.Map(result.Items, func(item datagraph.Item) mcp.ReplySearchItem {
		desc := item.GetDesc()
		return mcp.ReplySearchItem{
			BrowserUrl:  datagraph.CanonicalResolveURL(st.webAddress, datagraph.KindReply, item.GetID().String()).String(),
			Id:          item.GetID().String(),
			Slug:        item.GetSlug(),
			Name:        item.GetName(),
			Description: &desc,
		}
	})

	output := mcp.ToolReplySearchOutput{
		Results: result.Results,
		Items:   items,
	}

	if len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation errors: %s", strings.Join(validationErrors, "; "))
	}

	return &output, nil
}

func (st *searchTools) newPostSearchTool() *Tool {
	toolDef := mcp.GetPostSearchTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				func(ctx tool.Context, args mcp.ToolPostSearchInput) (*mcp.ToolPostSearchOutput, error) {
					return st.ExecutePostSearch(ctx, args)
				},
			)
		},
		Handler: makeHandler(st.ExecutePostSearch),
	}
}

func (st *searchTools) ExecutePostSearch(ctx context.Context, args mcp.ToolPostSearchInput) (*mcp.ToolPostSearchOutput, error) {
	maxResults := 10
	if args.MaxResults != nil {
		maxResults = *args.MaxResults
	}

	var validationErrors []string

	authorFilter, authorErrors := st.resolveAuthors(ctx, args.Authors)
	validationErrors = append(validationErrors, authorErrors...)

	categoryFilter, categoryErrors := st.resolveCategories(ctx, args.Categories)
	validationErrors = append(validationErrors, categoryErrors...)

	tagFilter := resolveTags(args.Tags)

	opts := searcher.Options{
		Kinds:      opt.New([]datagraph.Kind{datagraph.KindThread, datagraph.KindReply}),
		Authors:    authorFilter,
		Categories: categoryFilter,
		Tags:       tagFilter,
	}

	pageParams := pagination.NewPageParams(1, uint(maxResults))

	result, err := st.searcher.Search(ctx, args.Query, pageParams, opts)
	if err != nil {
		return nil, err
	}

	items := dt.Map(result.Items, func(item datagraph.Item) mcp.PostSearchItem {
		desc := item.GetDesc()
		return mcp.PostSearchItem{
			BrowserUrl:  datagraph.CanonicalResolveURL(st.webAddress, item.GetKind(), fmt.Sprintf("%s-%s", item.GetID().String(), item.GetSlug())).String(),
			Id:          item.GetID().String(),
			Kind:        item.GetKind().String(),
			Slug:        item.GetSlug(),
			Name:        item.GetName(),
			Description: &desc,
		}
	})

	output := mcp.ToolPostSearchOutput{
		Results: result.Results,
		Items:   items,
	}

	if len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation errors: %s", strings.Join(validationErrors, "; "))
	}

	return &output, nil
}

func (st *searchTools) newMemberSearchTool() *Tool {
	toolDef := mcp.GetMemberSearchTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				func(ctx tool.Context, args mcp.ToolMemberSearchInput) (*mcp.ToolMemberSearchOutput, error) {
					return st.ExecuteMemberSearch(ctx, args)
				},
			)
		},
		Handler: makeHandler(st.ExecuteMemberSearch),
	}
}

func (st *searchTools) ExecuteMemberSearch(ctx context.Context, args mcp.ToolMemberSearchInput) (*mcp.ToolMemberSearchOutput, error) {
	maxResults := 10
	if args.MaxResults != nil {
		maxResults = *args.MaxResults
	}

	pageParams := pagination.NewPageParams(1, uint(maxResults))

	result, err := st.profileSearch.Search(ctx, pageParams, profile_search.WithNamesLike(args.Query))
	if err != nil {
		return nil, err
	}

	items := dt.Map(result.Items, func(p *profile.Public) mcp.MemberSearchItem {
		bio := p.Bio.Short()
		return mcp.MemberSearchItem{
			BrowserUrl: datagraph.CanonicalResolveURL(st.webAddress, datagraph.KindProfile, fmt.Sprintf("%s-%s", p.ID.String(), p.Handle)).String(),
			Id:         p.ID.String(),
			Handle:     p.Handle,
			Name:       p.Name,
			Bio:        &bio,
		}
	})

	output := mcp.ToolMemberSearchOutput{
		Results: result.Results,
		Items:   items,
	}

	return &output, nil
}

func (st *searchTools) resolveAuthors(ctx context.Context, handles []string) (opt.Optional[[]account.AccountID], []string) {
	if len(handles) == 0 {
		return opt.NewEmpty[[]account.AccountID](), nil
	}

	var validationErrors []string

	accounts, err := st.accountQuerier.ProbeMany(ctx, handles...)
	if err != nil {
		return opt.NewEmpty[[]account.AccountID](), []string{fmt.Sprintf("failed to look up authors: %v", err)}
	}

	if len(accounts) < len(handles) {
		foundHandles := dt.Map(accounts, func(acc *account.Account) string {
			return acc.Handle
		})
		var notFound []string
		for _, handle := range handles {
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

	if len(accounts) == 0 {
		return opt.NewEmpty[[]account.AccountID](), validationErrors
	}

	authorIDs := dt.Map(accounts, func(acc *account.Account) account.AccountID {
		return acc.ID
	})

	return opt.New(authorIDs), validationErrors
}

func (st *searchTools) resolveCategories(ctx context.Context, names []string) (opt.Optional[[]category.CategoryID], []string) {
	if len(names) == 0 {
		return opt.NewEmpty[[]category.CategoryID](), nil
	}

	var validationErrors []string

	categories, err := st.categoryRepo.GetCategories(ctx, false)
	if err != nil {
		return opt.NewEmpty[[]category.CategoryID](), []string{fmt.Sprintf("failed to look up categories: %v", err)}
	}

	matchedIDs := []category.CategoryID{}
	var notMatched []string
	for _, searchName := range names {
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

	if len(matchedIDs) == 0 {
		return opt.NewEmpty[[]category.CategoryID](), validationErrors
	}

	return opt.New(matchedIDs), validationErrors
}

func resolveTags(tags []string) opt.Optional[[]tag_ref.Name] {
	if len(tags) == 0 {
		return opt.NewEmpty[[]tag_ref.Name]()
	}

	tagNames := dt.Map(tags, func(t string) tag_ref.Name {
		return tag_ref.NewName(t)
	})

	return opt.New(tagNames)
}
