package tools

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_traversal"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
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
