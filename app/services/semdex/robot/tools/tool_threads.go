package tools

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	reply_service "github.com/Southclaws/storyden/app/services/reply"
	thread_service "github.com/Southclaws/storyden/app/services/thread"
	"github.com/Southclaws/storyden/app/services/thread_mark"
	"github.com/Southclaws/storyden/mcp"
)

type threadTools struct {
	logger          *slog.Logger
	thread_svc      thread_service.Service
	replyMutator    *reply_service.Mutator
	thread_mark_svc thread_mark.Service
	category_repo   *category.Repository
}

func newThreadTools(
	logger *slog.Logger,
	registry *Registry,
	thread_svc thread_service.Service,
	replyMutator *reply_service.Mutator,
	thread_mark_svc thread_mark.Service,
	category_repo *category.Repository,
) *threadTools {
	t := &threadTools{
		logger:          logger,
		thread_svc:      thread_svc,
		replyMutator:    replyMutator,
		thread_mark_svc: thread_mark_svc,
		category_repo:   category_repo,
	}

	registry.Register(t.newThreadCreateTool())
	registry.Register(t.newThreadListTool())
	registry.Register(t.newThreadGetTool())
	registry.Register(t.newThreadUpdateTool())
	registry.Register(t.newThreadReplyTool())
	registry.Register(t.newCategoryListTool())

	return t
}

func (tt *threadTools) newThreadCreateTool() *Tool {
	toolDef := mcp.GetThreadCreateTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				tt.ExecuteThreadCreate,
			)
		},
	}
}

func (tt *threadTools) ExecuteThreadCreate(ctx tool.Context, args mcp.ToolThreadCreateInput) (*mcp.ToolThreadCreateOutput, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	vis := visibility.VisibilityPublished
	if args.Visibility != nil {
		switch *args.Visibility {
		case mcp.ToolThreadCreateInputVisibilityDraft:
			vis = visibility.VisibilityDraft
		case mcp.ToolThreadCreateInputVisibilityPublished:
			vis = visibility.VisibilityPublished
		}
	}

	richContent, err := datagraph.NewRichText(args.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	var parsedURL opt.Optional[url.URL]
	if args.Url != nil && *args.Url != "" {
		u, err := url.Parse(*args.Url)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		parsedURL = opt.New(*u)
	}

	var tagNames opt.Optional[tag_ref.Names]
	if len(args.Tags) > 0 {
		names := dt.Map(args.Tags, func(t string) tag_ref.Name {
			name := tag_ref.NewName(t)
			return name
		})
		tagNames = opt.New(tag_ref.Names(names))
	}

	cats, err := tt.category_repo.GetCategories(ctx, false)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	cat, found := lo.Find(cats, func(c *category.Category) bool { return c.Slug == args.Category })
	if !found {
		return nil, fault.Wrap(fault.New("category not found"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	createdThread, err := tt.thread_svc.Create(ctx,
		args.Title,
		accountID,
		map[string]any{},
		thread_service.Partial{
			Content:    opt.New(richContent),
			URL:        parsedURL,
			Tags:       tagNames,
			Visibility: opt.New(vis),
			Category:   opt.New(xid.ID(cat.ID)),
		},
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var urlStr *string
	if webLink, ok := createdThread.WebLink.Get(); ok {
		urlStr = &webLink.URL
	}

	output := mcp.ToolThreadCreateOutput{
		Slug:       createdThread.Slug,
		Title:      createdThread.Title,
		Content:    func() *string { s := createdThread.Content.Plaintext(); return &s }(),
		CreatedAt:  func() *string { s := createdThread.CreatedAt.Format(time.RFC3339); return &s }(),
		Visibility: func() *string { s := createdThread.Visibility.String(); return &s }(),
		Author:     func() *string { s := createdThread.Author.Handle; return &s }(),
		Category:   func() *string { s := createdThread.Category.OrZero().Name; return &s }(),
		Tags:       dt.Map(createdThread.Tags, func(tag *tag_ref.Tag) string { return tag.Name.String() }),
		Url:        urlStr,
	}

	return &output, nil
}

func (tt *threadTools) newThreadListTool() *Tool {
	toolDef := mcp.GetThreadListTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				tt.ExecuteThreadList,
			)
		},
	}
}

func (tt *threadTools) ExecuteThreadList(ctx tool.Context, args mcp.ToolThreadListInput) (*mcp.ToolThreadListOutput, error) {
	page := 1
	if args.Page != nil {
		page = *args.Page
	}

	pageSize := 10

	var visibilities opt.Optional[[]visibility.Visibility]
	if args.Visibility != nil {
		switch *args.Visibility {
		case mcp.ToolThreadListInputVisibilityDraft:
			visibilities = opt.New([]visibility.Visibility{visibility.VisibilityDraft})
		case mcp.ToolThreadListInputVisibilityPublished:
			visibilities = opt.New([]visibility.Visibility{visibility.VisibilityPublished})
		}
	}

	var queryOpt opt.Optional[string]
	if args.Query != nil {
		queryOpt = opt.New(*args.Query)
	}

	result, err := tt.thread_svc.List(ctx, max(0, page-1), pageSize, thread_service.Params{
		Query:      queryOpt,
		Visibility: visibilities,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var nextPage *int
	if np, ok := result.NextPage.Get(); ok {
		nextPage = func() *int { i := int(np) + 1; return &i }()
	}

	output := mcp.ToolThreadListOutput{
		Threads:      dt.Map(result.Threads, mapThreadSummary),
		CurrentPage:  page,
		TotalPages:   int(result.TotalPages),
		TotalResults: int(result.Results),
		NextPage:     nextPage,
	}

	return &output, nil
}

func (tt *threadTools) newThreadGetTool() *Tool {
	toolDef := mcp.GetThreadGetTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				tt.ExecuteThreadGet,
			)
		},
	}
}

func (tt *threadTools) ExecuteThreadGet(ctx tool.Context, args mcp.ToolThreadGetInput) (*mcp.ToolThreadGetOutput, error) {
	postID, err := tt.thread_mark_svc.Lookup(ctx, args.Id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	page := 1
	if args.Page != nil {
		page = *args.Page
	}

	pageParams := pagination.NewPageParams(uint(page), 20)

	thread, err := tt.thread_svc.Get(ctx, postID, pageParams)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var urlStr *string
	if webLink, ok := thread.WebLink.Get(); ok {
		urlStr = &webLink.URL
	}

	output := mcp.ToolThreadGetOutput{
		Slug:       thread.Slug,
		Title:      thread.Title,
		Content:    thread.Content.Plaintext(),
		CreatedAt:  thread.CreatedAt.Format(time.RFC3339),
		Visibility: thread.Visibility.String(),
		Author:     thread.Author.Handle,
		Category:   thread.Category.OrZero().Name,
		Tags:       dt.Map(thread.Tags, func(tag *tag_ref.Tag) string { return tag.Name.String() }),
		Url:        urlStr,
	}

	return &output, nil
}

func mapThreadSummary(t *thread.Thread) mcp.ThreadSummary {
	return mcp.ThreadSummary{
		Slug:     t.Slug,
		Title:    t.Title,
		Excerpt:  t.Short,
		Author:   t.Author.Handle,
		Category: t.Category.OrZero().Name,
	}
}

func (tt *threadTools) newThreadUpdateTool() *Tool {
	toolDef := mcp.GetThreadUpdateTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				tt.ExecuteThreadUpdate,
			)
		},
	}
}

func (tt *threadTools) ExecuteThreadUpdate(ctx tool.Context, args mcp.ToolThreadUpdateInput) (*mcp.ToolThreadUpdateOutput, error) {
	postID, err := tt.thread_mark_svc.Lookup(ctx, args.Id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	partial := thread_service.Partial{}

	if args.Title != nil && *args.Title != "" {
		partial.Title = opt.New(*args.Title)
	}

	if args.Body != nil && *args.Body != "" {
		richContent, err := datagraph.NewRichText(*args.Body)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		partial.Content = opt.New(richContent)
	}

	if args.Visibility != nil {
		switch *args.Visibility {
		case mcp.ToolThreadUpdateInputVisibilityDraft:
			partial.Visibility = opt.New(visibility.VisibilityDraft)
		case mcp.ToolThreadUpdateInputVisibilityPublished:
			partial.Visibility = opt.New(visibility.VisibilityPublished)
		}
	}

	if args.Tags != nil && len(args.Tags) > 0 {
		names := dt.Map(args.Tags, func(t string) tag_ref.Name {
			name := tag_ref.NewName(t)
			return name
		})
		partial.Tags = opt.New(tag_ref.Names(names))
	}

	thread, err := tt.thread_svc.Update(ctx, postID, partial)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var urlStr *string
	if webLink, ok := thread.WebLink.Get(); ok {
		urlStr = &webLink.URL
	}

	output := mcp.ToolThreadUpdateOutput{
		Slug:       thread.Slug,
		Title:      thread.Title,
		Content:    func() *string { s := thread.Content.Plaintext(); return &s }(),
		CreatedAt:  func() *string { s := thread.CreatedAt.Format(time.RFC3339); return &s }(),
		Visibility: func() *string { s := thread.Visibility.String(); return &s }(),
		Author:     func() *string { s := thread.Author.Handle; return &s }(),
		Category:   func() *string { s := thread.Category.OrZero().Name; return &s }(),
		Tags:       dt.Map(thread.Tags, func(tag *tag_ref.Tag) string { return tag.Name.String() }),
		Url:        urlStr,
	}

	return &output, nil
}

func (tt *threadTools) newThreadReplyTool() *Tool {
	toolDef := mcp.GetThreadReplyTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				tt.ExecuteThreadReply,
			)
		},
	}
}

func (tt *threadTools) ExecuteThreadReply(ctx tool.Context, args mcp.ToolThreadReplyInput) (*mcp.ToolThreadReplyOutput, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	postID, err := tt.thread_mark_svc.Lookup(ctx, args.Id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	richContent, err := datagraph.NewRichText(args.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	reply, err := tt.replyMutator.Create(ctx, accountID, postID, reply_service.Partial{
		Content: opt.New(richContent),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	output := mcp.ToolThreadReplyOutput{
		Author:    reply.Author.Handle,
		Content:   reply.Content.Plaintext(),
		CreatedAt: reply.CreatedAt.Format(time.RFC3339),
		UpdatedAt: reply.UpdatedAt.Format(time.RFC3339),
	}

	return &output, nil
}

func (tt *threadTools) newCategoryListTool() *Tool {
	toolDef := mcp.GetCategoryListTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				tt.ExecuteCategoryList,
			)
		},
	}
}

func (tt *threadTools) ExecuteCategoryList(ctx tool.Context, args map[string]any) (*mcp.ToolCategoryListOutput, error) {
	categories, err := tt.category_repo.GetCategories(ctx, false)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	cats := make([]mcp.CategoryItem, 0, len(categories))
	for _, cat := range categories {
		var desc *string
		if cat.Description != "" {
			desc = &cat.Description
		}

		cats = append(cats, mcp.CategoryItem{
			Slug:        cat.Slug,
			Name:        cat.Name,
			Description: desc,
		})
	}

	output := mcp.ToolCategoryListOutput{
		Categories: cats,
	}

	return &output, nil
}
