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
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/generative"
	"github.com/Southclaws/storyden/app/services/library/node_mutate"
	"github.com/Southclaws/storyden/app/services/library/node_property_schema"
	"github.com/Southclaws/storyden/app/services/library/node_read"
	"github.com/Southclaws/storyden/app/services/library/node_visibility"
	"github.com/Southclaws/storyden/app/services/library/nodetree"
	"github.com/Southclaws/storyden/app/services/tag/autotagger"
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
	}

	handler.tools = []server.ServerTool{
		{Tool: nodeTreeTool, Handler: handler.nodeTree},
		{Tool: nodeGetTool, Handler: handler.nodeGet},
		{Tool: nodeCreateTool, Handler: handler.nodeCreate},
		{Tool: nodeSearchTool, Handler: handler.nodeSearch},
	}

	return handler
}

var nodeTreeTool = mcp.NewTool("getNodeTree",
	mcp.WithDescription("Get the full tree of nodes in the library"),
)

func (t *nodeTools) nodeTree(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	opts := []node_traversal.Filter{}

	depth := request.GetInt("depth", -1)
	if depth != -1 {
		opts = append(opts, node_traversal.WithDepth(uint(depth)))
	}

	tree, err := t.ntr.Subtree(ctx, opt.NewEmpty[library.NodeID](), true, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	obj := mapNodes(tree)
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mcp.NewToolResultText(string(b)), nil
}

var nodeGetTool = mcp.NewTool("getNode",
	mcp.WithDescription("Get a specific node from the library"),
	mcp.WithString("node_slug", mcp.Required()),
)

func (t *nodeTools) nodeGet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slug, err := request.RequireString("node_slug")
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
		"id":          n.Mark.ID(),
		"slug":        n.Mark.Slug(),
		"name":        n.Name,
		"description": n.Description,
		"tags":        dt.Map(n.Tags, mapTag),
	}
}

func mapNodes(nodes []*library.Node) []map[string]any {
	return dt.Map(nodes, mapNode)
}

func mapTag(t *tag_ref.Tag) string {
	return t.Name.String()
}

var nodeCreateTool = mcp.NewTool("createNode",
	mcp.WithDescription("Create a new node in the library"),
	mcp.WithString("name", mcp.Required()),
	mcp.WithString("content"),
	mcp.WithString("url"),
	mcp.WithString("slug"),
	mcp.WithString("parent"),
)

func (t *nodeTools) nodeCreate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	content := request.GetString("content", "")
	urlStr := request.GetString("url", "")
	slugStr := request.GetString("slug", "")
	parentStr := request.GetString("parent", "")

	var richContent opt.Optional[datagraph.Content]
	if content != "" {
		rc, err := datagraph.NewRichText(content)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		richContent = opt.New(rc)
	}

	var urlParsed opt.Optional[url.URL]
	if urlStr != "" {
		u, err := url.Parse(urlStr)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		urlParsed = opt.New(*u)
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
		parent = opt.New(library.QueryKey{mark.NewQueryKey(parentStr)})
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	node, err := t.nodeMutator.Create(ctx,
		accountID,
		name,
		node_mutate.Partial{
			Slug:    slug,
			Content: richContent,
			URL:     urlParsed,
			Parent:  parent,
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

var nodeSearchTool = mcp.NewTool("searchNodes",
	mcp.WithDescription("Search for nodes in the library"),
	mcp.WithString("query"),
	mcp.WithString("author"),
	mcp.WithNumber("depth"),
)

func (t *nodeTools) nodeSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	author := request.GetString("author", "")
	depth := request.GetInt("depth", -1)

	opts := []node_traversal.Filter{}

	if author != "" {
		opts = append(opts, node_traversal.WithRootOwner(author))
	}

	if depth != -1 {
		opts = append(opts, node_traversal.WithDepth(uint(depth)))
	}

	nodes, err := t.ntr.Subtree(ctx, opt.NewEmpty[library.NodeID](), true, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	obj := mapNodes(nodes)
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mcp.NewToolResultText(string(b)), nil
}
