package tools

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_traversal"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/generative"
	"github.com/Southclaws/storyden/app/services/library/node_mutate"
	"github.com/Southclaws/storyden/app/services/library/node_property_schema"
	"github.com/Southclaws/storyden/app/services/library/node_read"
	"github.com/Southclaws/storyden/app/services/library/node_visibility"
	"github.com/Southclaws/storyden/app/services/library/nodetree"
	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/app/services/tag/autotagger"
	"github.com/Southclaws/storyden/internal/deletable"
)

type nodeTools struct {
	tools []server.ServerTool

	accountQuery  *account_querier.Querier
	nodeMutator   *node_mutate.Manager
	tagger        *autotagger.Tagger
	summariser    generative.Summariser
	titler        generative.Titler
	nodeReader    *node_read.HydratedQuerier
	nv            *node_visibility.Controller
	ntree         nodetree.Graph
	npos          *nodetree.Position
	ntr           node_traversal.Repository
	schemaUpdater *node_property_schema.Updater
	searcher      searcher.Searcher
}

func newNodeTools(
	accountQuery *account_querier.Querier,
	nodeMutator *node_mutate.Manager,
	tagger *autotagger.Tagger,
	summariser generative.Summariser,
	titler generative.Titler,
	nodeReader *node_read.HydratedQuerier,
	nv *node_visibility.Controller,
	ntree nodetree.Graph,
	npos *nodetree.Position,
	ntr node_traversal.Repository,
	schemaUpdater *node_property_schema.Updater,
	searcher searcher.Searcher,
) *nodeTools {
	handler := &nodeTools{
		accountQuery:  accountQuery,
		nodeMutator:   nodeMutator,
		tagger:        tagger,
		summariser:    summariser,
		titler:        titler,
		nodeReader:    nodeReader,
		nv:            nv,
		ntree:         ntree,
		npos:          npos,
		ntr:           ntr,
		schemaUpdater: schemaUpdater,
		searcher:      searcher,
	}

	handler.tools = []server.ServerTool{
		{Tool: libraryPageTreeTool, Handler: handler.libraryPageTree},
		{Tool: libraryPageGetTool, Handler: handler.libraryPageGet},
		{Tool: libraryPageCreateTool, Handler: handler.libraryPageCreate},
		{Tool: libraryPageUpdateTool, Handler: handler.libraryPageUpdate},
		{Tool: libraryPageSearchTool, Handler: handler.libraryPageSearch},
	}

	return handler
}

var libraryPageTreeTool = mcp.NewTool("getLibraryPageTree",
	mcp.WithDescription("Get the full tree of pages in the library"),
)

func (t *nodeTools) libraryPageTree(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	account, err := session.GetAccount(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := t.accountQuery.GetByID(ctx, account.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := []node_traversal.Filter{
		node_traversal.WithVisibility(opt.New(*acc), visibility.VisibilityDraft, visibility.VisibilityPublished),
	}

	depth := request.GetInt("depth", -1)
	if depth != -1 {
		opts = append(opts, node_traversal.WithDepth(uint(depth)))
	}

	tree, err := t.ntr.Subtree(ctx, opt.NewEmpty[library.NodeID](), true, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	obj := mapNodeTree(tree)
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mcp.NewToolResultText(string(b)), nil
}

var libraryPageGetTool = mcp.NewTool("getLibraryPage",
	mcp.WithDescription("Get a specific page from the library"),
	mcp.WithString("slug", mcp.Required()),
)

func (t *nodeTools) libraryPageGet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slug, err := request.RequireString("slug")
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	node, err := t.nodeReader.GetBySlug(ctx, library.NewKey(slug), nil)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	obj := mapNode(node)
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mcp.NewToolResultText(string(b)), nil
}

func mapNode(n *library.Node) map[string]any {
	return map[string]any{
		"slug":        n.Mark.Slug(),
		"name":        n.Name,
		"description": n.Description,
		"tags":        dt.Map(n.Tags, mapTag),
		"child_pages": mapNodes(n.Nodes),
	}
}

func mapNodes(nodes []*library.Node) []map[string]any {
	return dt.Map(nodes, mapNode)
}

func mapNodeTreeItem(n *library.Node) map[string]any {
	i := map[string]any{
		"slug":        n.Mark.Slug(),
		"name":        n.Name,
		"description": n.Description,
		"tags":        dt.Map(n.Tags, mapTag),
	}

	if v, ok := n.Parent.Get(); ok {
		i["parent"] = v.Mark.Slug()
	}

	return i
}

func mapNodeTree(nodes []*library.Node) []map[string]any {
	return dt.Map(nodes, mapNodeTreeItem)
}

func mapTag(t *tag_ref.Tag) string {
	return t.Name.String()
}

var libraryPageCreateTool = mcp.NewTool("createLibraryPage",
	mcp.WithDescription("Create a new page in the library"),
	mcp.WithString("name", mcp.Required(), mcp.Description("The name of the page.")),
	mcp.WithString("slug", mcp.Description("The unique slug within Storyden for this page. If you leave this empty, a slug will be generated for you.")),
	mcp.WithString("content", mcp.Description("The content of the page in HTML format.")),
	mcp.WithString("parent", mcp.Description("Only include the parent if you already have a parent slug available from a page search. If not, this field must be left empty, otherwise the createNode tool will fail catastrophically and everyone will be very sad.")),
	mcp.WithString("visibility", mcp.Description("Visibility of the page. published or draft, defaults to published")),
	mcp.WithString("url", mcp.Description("If this page is about a topic referred to on an external website, use this to reference that website.")),
)

func (t *nodeTools) libraryPageCreate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	content := request.GetString("content", "")
	urlStr := request.GetString("url", "")
	slugStr := request.GetString("slug", "")
	parentStr := request.GetString("parent", "")
	visibilityStr := request.GetString("visibility", "published")

	var richContent opt.Optional[datagraph.Content]
	if content != "" {
		rc, err := datagraph.NewRichText(content)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		richContent = opt.New(rc)
	}

	var urlParsed deletable.Value[url.URL]
	if urlStr != "" {
		u, err := url.Parse(urlStr)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		urlParsed = deletable.Skip[url.URL](opt.New(*u))
	}

	var slug opt.Optional[mark.Slug]
	if slugStr != "" {
		s, err := mark.NewSlug(slugStr)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		slug = opt.New(*s)
	}

	var parent opt.Optional[library.QueryKey]
	if parentStr != "" {
		parent = opt.New(library.NewKey(parentStr))
	}

	var vis opt.Optional[visibility.Visibility]
	if visibilityStr != "" {
		v, err := visibility.NewVisibility(visibilityStr)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		vis = opt.New(v)
	} else {
		vis = opt.New(visibility.VisibilityPublished)
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	node, err := t.nodeMutator.Create(ctx,
		accountID,
		name,
		node_mutate.Partial{
			Slug:       slug,
			Content:    richContent,
			URL:        urlParsed,
			Parent:     parent,
			Visibility: vis,
		},
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	obj := mapNode(node)
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mcp.NewToolResultText(string(b)), nil
}

var libraryPageUpdateTool = mcp.NewTool("updateLibraryPage",
	mcp.WithDescription("Update an existing page in the library"),
	mcp.WithString("slug", mcp.Required(), mcp.Description("The slug of the page to update")),
	mcp.WithString("name", mcp.Description("The new name of the page")),
	mcp.WithString("content", mcp.Description("The new content of the page in HTML format")),
	mcp.WithString("visibility", mcp.Description("New visibility of the page: published or draft")),
	mcp.WithString("url", mcp.Description("If this page is about a topic referred to on an external website, use this to reference that website")),
	mcp.WithString("parent", mcp.Description("New parent page slug. Leave empty to move to root level")),
)

func (t *nodeTools) libraryPageUpdate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slug, err := request.RequireString("slug")
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	name := request.GetString("name", "")
	content := request.GetString("content", "")
	visibilityStr := request.GetString("visibility", "")
	urlStr := request.GetString("url", "")
	parentStr := request.GetString("parent", "")

	partial := node_mutate.Partial{}

	if name != "" {
		partial.Name = opt.New(name)
	}

	if content != "" {
		richContent, err := datagraph.NewRichText(content)
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

	if urlStr != "" {
		u, err := url.Parse(urlStr)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		partial.URL = deletable.Skip(opt.New(*u))
	}

	if parentStr != "" {
		partial.Parent = opt.New(library.NewKey(parentStr))
	}

	node, err := t.nodeMutator.Update(ctx, library.NewKey(slug), partial)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	obj := mapNode(node)
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mcp.NewToolResultText(string(b)), nil
}

var libraryPageSearchTool = mcp.NewTool("searchLibraryPages",
	mcp.WithDescription("Search for pages in the library."),
	mcp.WithString("query", mcp.Required()),
)

func (t *nodeTools) libraryPageSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := request.RequireString("query")
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pp := pagination.NewPageParams(1, 50)

	opts := searcher.Options{
		Kinds: opt.New([]datagraph.Kind{datagraph.KindNode}),
	}

	result, err := t.searcher.Search(ctx, query, pp, opts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	obj := mapDatagraphItems(result.Items)
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mcp.NewToolResultText(string(b)), nil
}

func mapDatagraphItems(items datagraph.ItemList) []map[string]any {
	return dt.Map(items, mapDatagraphItem)
}

func mapDatagraphItem(item datagraph.Item) map[string]any {
	base := map[string]any{
		"slug":        item.GetSlug(),
		"name":        item.GetName(),
		"description": item.GetDesc(),
	}

	if content := item.GetContent(); content.Plaintext() != "" {
		base["content"] = content.Plaintext()
	}

	return base
}
