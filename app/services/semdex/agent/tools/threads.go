package tools

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	reply_service "github.com/Southclaws/storyden/app/services/reply"
	thread_service "github.com/Southclaws/storyden/app/services/thread"
	"github.com/Southclaws/storyden/app/services/thread_mark"
)

type threadTools struct {
	tools []server.ServerTool

	thread_svc      thread_service.Service
	reply_svc       reply_service.Service
	thread_mark_svc thread_mark.Service
	accountQuery    *account_querier.Querier
	category_repo   *category.Repository
}

func newThreadTools(
	thread_svc thread_service.Service,
	reply_svc reply_service.Service,
	thread_mark_svc thread_mark.Service,
	accountQuery *account_querier.Querier,
	category_repo *category.Repository,
) *threadTools {
	handler := &threadTools{
		thread_svc:      thread_svc,
		reply_svc:       reply_svc,
		thread_mark_svc: thread_mark_svc,
		accountQuery:    accountQuery,
		category_repo:   category_repo,
	}

	handler.tools = []server.ServerTool{
		{Tool: threadCreateTool, Handler: handler.threadCreate},
		{Tool: threadListTool, Handler: handler.threadList},
		{Tool: threadGetTool, Handler: handler.threadGet},
		{Tool: threadUpdateTool, Handler: handler.threadUpdate},
		{Tool: threadReplyTool, Handler: handler.threadReply},
		{Tool: listCategoresTool, Handler: handler.listCategories},
	}

	return handler
}

var threadCreateTool = mcp.NewTool("createThread",
	mcp.WithDescription("Create a new discussion thread in the forum"),
	mcp.WithString("title", mcp.Required(), mcp.Description("The title of the thread")),
	mcp.WithString("body", mcp.Required(), mcp.Description("The content of the thread in HTML format")),
	mcp.WithString("category", mcp.Required(), mcp.Description("The category slug for the thread")),
	mcp.WithString("visibility", mcp.Description("Thread visibility: published or draft")),
	mcp.WithString("url", mcp.Description("Optional URL if this thread is about a specific link, similar to a link aggregator such as Reddit")),
	mcp.WithString("tags", mcp.Description("Optional comma-separated tags for the thread")),
)

func (t *threadTools) threadCreate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, err := request.RequireString("title")
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	body, err := request.RequireString("body")
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	categorySlug, err := request.RequireString("category")
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	visibilityStr := request.GetString("visibility", "published")
	urlStr := request.GetString("url", "")
	tagsStr := request.GetString("tags", "")
	var tags []string
	if tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	vis, err := visibility.NewVisibility(visibilityStr)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	richContent, err := datagraph.NewRichText(body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	var parsedURL opt.Optional[url.URL]
	if urlStr != "" {
		u, err := url.Parse(urlStr)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		parsedURL = opt.New(*u)
	}

	var tagNames opt.Optional[tag_ref.Names]
	if len(tags) > 0 {
		names := dt.Map(tags, func(t string) tag_ref.Name {
			name := tag_ref.NewName(t)
			return name
		})
		tagNames = opt.New(tag_ref.Names(names))
	}

	cats, err := t.category_repo.GetCategories(ctx, false)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	cat, found := lo.Find(cats, func(c *category.Category) bool { return c.Slug == categorySlug })
	if !found {
		return nil, fault.Wrap(fault.New("category not found"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	thread, err := t.thread_svc.Create(ctx,
		title,
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

	obj := mapThread(thread)
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mcp.NewToolResultText(string(b)), nil
}

var threadListTool = mcp.NewTool("listThreads",
	mcp.WithDescription("List and search discussion threads"),
	mcp.WithString("query", mcp.Description("Search query to filter threads")),
	mcp.WithString("visibility", mcp.Description("Filter by visibility: draft or published")),
	mcp.WithString("page", mcp.Description("Page number (defaults to 1)")),
)

func (t *threadTools) threadList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := request.GetString("query", "")
	visibilityStr := request.GetString("visibility", "")
	pageStr := request.GetString("page", "1")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Smaller page size for MCP to avoid overwhelming LLMs
	pageSize := 10

	var visibilities opt.Optional[[]visibility.Visibility]
	if visibilityStr != "" {
		vis, err := visibility.NewVisibility(visibilityStr)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		visibilities = opt.New([]visibility.Visibility{vis})
	}

	var queryOpt opt.Optional[string]
	if query != "" {
		queryOpt = opt.New(query)
	}

	result, err := t.thread_svc.List(ctx, max(0, page-1), pageSize, thread_service.Params{
		Query:      queryOpt,
		Visibility: visibilities,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	response := map[string]any{
		"threads":      dt.Map(result.Threads, mapThreadSummary),
		"currentPage":  page,
		"totalPages":   result.TotalPages,
		"totalResults": result.Results,
	}

	if nextPage, ok := result.NextPage.Get(); ok {
		response["nextPage"] = nextPage + 1
	}

	b, err := json.Marshal(response)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mcp.NewToolResultText(string(b)), nil
}

var threadGetTool = mcp.NewTool("getThread",
	mcp.WithDescription("Get a specific thread with its posts and replies"),
	mcp.WithString("slug", mcp.Required(), mcp.Description("The thread URL slug")),
	mcp.WithString("page", mcp.Description("Page number for replies (defaults to 1)")),
)

func (t *threadTools) threadGet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	threadMark, err := request.RequireString("slug")
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pageStr := request.GetString("page", "1")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	postID, err := t.thread_mark_svc.Lookup(ctx, threadMark)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pageParams := pagination.NewPageParams(uint(max(1, page)), 20)

	thread, err := t.thread_svc.Get(ctx, postID, pageParams)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	obj := mapThread(thread)
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mcp.NewToolResultText(string(b)), nil
}

var threadUpdateTool = mcp.NewTool("updateThread",
	mcp.WithDescription("Update an existing thread"),
	mcp.WithString("slug", mcp.Required(), mcp.Description("The thread slug to update")),
	mcp.WithString("title", mcp.Description("New title for the thread")),
	mcp.WithString("body", mcp.Description("New content for the thread in HTML format")),
	mcp.WithString("visibility", mcp.Description("New visibility: published or draft")),
	mcp.WithString("tags", mcp.Description("New comma-separated tags for the thread")),
)

func (t *threadTools) threadUpdate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	threadMark, err := request.RequireString("slug")
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	title := request.GetString("title", "")
	body := request.GetString("body", "")
	visibilityStr := request.GetString("visibility", "")
	tagsStr := request.GetString("tags", "")
	var tags []string
	if tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	postID, err := t.thread_mark_svc.Lookup(ctx, threadMark)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	partial := thread_service.Partial{}

	if title != "" {
		partial.Title = opt.New(title)
	}

	if body != "" {
		richContent, err := datagraph.NewRichText(body)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		partial.Content = opt.New(richContent)
	}

	if visibilityStr != "" {
		vis, err := visibility.NewVisibility(visibilityStr)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		partial.Visibility = opt.New(vis)
	}

	if len(tags) > 0 {
		names := dt.Map(tags, func(t string) tag_ref.Name {
			name := tag_ref.NewName(t)
			return name
		})
		partial.Tags = opt.New(tag_ref.Names(names))
	}

	thread, err := t.thread_svc.Update(ctx, postID, partial)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	obj := mapThread(thread)
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mcp.NewToolResultText(string(b)), nil
}

var threadReplyTool = mcp.NewTool("replyToThread",
	mcp.WithDescription("Add a reply to an existing thread"),
	mcp.WithString("slug", mcp.Required(), mcp.Description("The thread slug to reply to")),
	mcp.WithString("body", mcp.Required(), mcp.Description("The reply content in HTML format")),
)

func (t *threadTools) threadReply(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	threadMark, err := request.RequireString("slug")
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	body, err := request.RequireString("body")
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	postID, err := t.thread_mark_svc.Lookup(ctx, threadMark)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	richContent, err := datagraph.NewRichText(body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	reply, err := t.reply_svc.Create(ctx, accountID, postID, reply_service.Partial{
		Content: opt.New(richContent),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	obj := mapReply(reply)
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mcp.NewToolResultText(string(b)), nil
}

var listCategoresTool = mcp.NewTool("listCategories",
	mcp.WithDescription("List all thread categories with their names and descriptions"),
)

func (t *threadTools) listCategories(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	categories, err := t.category_repo.GetCategories(ctx, false)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	obj := dt.Map(categories, mapCategory)
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mcp.NewToolResultText(string(b)), nil
}

func mapCategory(c *category.Category) map[string]any {
	return map[string]any{
		"slug":        c.Slug,
		"name":        c.Name,
		"description": c.Description,
	}
}

func mapThread(t *thread.Thread) map[string]any {
	result := map[string]any{
		"slug":       t.Slug,
		"created_at": t.CreatedAt,
		"title":      t.Title,
		"content":    t.Content.Plaintext(),
		"visibility": t.Visibility.String(),
		"author":     t.Author.Handle,
		"category":   t.Category.OrZero().Name,
		"tags":       dt.Map(t.Tags, func(tag *tag_ref.Tag) string { return tag.Name.String() }),
	}

	if webLink, ok := t.WebLink.Get(); ok {
		result["url"] = webLink.URL
	}

	return result
}

func mapThreadSummary(t *thread.Thread) map[string]any {
	result := map[string]any{
		"slug":     t.Slug,
		"title":    t.Title,
		"excerpt":  t.Short,
		"author":   t.Author.Handle,
		"category": t.Category.OrZero().Name,
	}

	return result
}

func mapReply(p *reply.Reply) map[string]any {
	result := map[string]any{
		"updated_at": p.UpdatedAt,
		"author":     p.Author.Handle,
		"created_at": p.CreatedAt,
		"content":    p.Content.Plaintext(),
	}

	return result
}
