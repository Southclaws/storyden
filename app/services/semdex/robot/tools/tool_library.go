package tools

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

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
	"github.com/Southclaws/storyden/mcp"
)

type libraryTools struct {
	logger *slog.Logger

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

func newLibraryTools(
	logger *slog.Logger,
	registry *Registry,

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
) *libraryTools {
	t := &libraryTools{
		logger:        logger,
		accountQuery:  accountQuery,
		ntr:           ntr,
		nodeReader:    nodeReader,
		nodeMutator:   nodeMutator,
		searcher:      searcher,
		tagger:        tagger,
		summariser:    summariser,
		titler:        titler,
		nv:            nv,
		ntree:         ntree,
		npos:          npos,
		schemaUpdater: schemaUpdater,
	}

	registry.Register(t.newLibraryPageTreeTool())
	registry.Register(t.newLibraryPageGetTool())
	registry.Register(t.newLibraryPageCreateTool())
	registry.Register(t.newLibraryPageUpdateTool())
	registry.Register(t.newLibraryPageSearchTool())

	return t
}

func (lt *libraryTools) newLibraryPageTreeTool() *Tool {
	toolDef := mcp.GetLibraryPageTreeTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				lt.ExecuteLibraryPageTree,
			)
		},
	}
}

func (lt *libraryTools) ExecuteLibraryPageTree(ctx tool.Context, args mcp.ToolLibraryPageTreeInput) ToolResult[mcp.ToolLibraryPageTreeOutput] {
	account, err := session.GetAccount(ctx)
	if err != nil {
		return NewError[mcp.ToolLibraryPageTreeOutput](err)
	}

	acc, err := lt.accountQuery.GetByID(ctx, account.ID)
	if err != nil {
		return NewError[mcp.ToolLibraryPageTreeOutput](err)
	}

	opts := []node_traversal.Filter{
		node_traversal.WithVisibility(opt.New(*acc), visibility.VisibilityDraft, visibility.VisibilityPublished),
	}

	if args.Depth != nil {
		depth := *args.Depth
		if depth != -1 {
			opts = append(opts, node_traversal.WithDepth(uint(depth)))
		}
	}

	tree, err := lt.ntr.Subtree(ctx, opt.NewEmpty[library.NodeID](), true, opts...)
	if err != nil {
		return NewError[mcp.ToolLibraryPageTreeOutput](err)
	}

	pages := dt.Map(tree, mapNodeToTreeNode)

	output := mcp.ToolLibraryPageTreeOutput{
		Pages: pages,
	}

	return NewSuccess(output)
}

func mapNodeToTreeNode(n *library.Node) mcp.LibraryPageTreeNode {
	node := mcp.LibraryPageTreeNode{
		Slug:        n.Mark.Slug(),
		Name:        n.Name,
		Description: n.Description.OrZero(),
		Tags:        dt.Map(n.Tags, func(t *tag_ref.Tag) string { return t.Name.String() }),
	}

	if parent, ok := n.Parent.Get(); ok {
		parentSlug := parent.Mark.Slug()
		node.Parent = &parentSlug
	}

	return node
}

func (lt *libraryTools) newLibraryPageGetTool() *Tool {
	toolDef := mcp.GetLibraryPageGetTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				lt.ExecuteLibraryPageGet,
			)
		},
	}
}

func (lt *libraryTools) ExecuteLibraryPageGet(ctx tool.Context, args mcp.ToolLibraryPageGetInput) ToolResult[mcp.ToolLibraryPageGetOutput] {
	node, err := lt.nodeReader.GetBySlug(ctx, library.NewKey(args.Slug), nil)
	if err != nil {
		return NewError[mcp.ToolLibraryPageGetOutput](err)
	}

	output := mcp.ToolLibraryPageGetOutput{
		Slug:        node.Mark.Slug(),
		Name:        node.Name,
		Description: node.Description.Ptr(),
		Tags:        dt.Map(node.Tags, func(t *tag_ref.Tag) string { return t.Name.String() }),
		ChildPages:  dt.Map(node.Nodes, func(n *library.Node) string { return n.Mark.Slug() }),
	}

	return NewSuccess(output)
}

func (lt *libraryTools) newLibraryPageCreateTool() *Tool {
	toolDef := mcp.GetLibraryPageCreateTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				lt.ExecuteLibraryPageCreate,
			)
		},
	}
}

func (lt *libraryTools) ExecuteLibraryPageCreate(ctx tool.Context, args mcp.ToolLibraryPageCreateInput) ToolResult[mcp.ToolLibraryPageCreateOutput] {
	var richContent opt.Optional[datagraph.Content]
	if args.Content != nil {
		rc, err := datagraph.NewRichText(*args.Content)
		if err != nil {
			return NewError[mcp.ToolLibraryPageCreateOutput](err)
		}
		richContent = opt.New(rc)
	}

	var urlParsed deletable.Value[url.URL]
	if args.Url != nil {
		u, err := url.Parse(*args.Url)
		if err != nil {
			return NewError[mcp.ToolLibraryPageCreateOutput](err)
		}
		urlParsed = deletable.Skip[url.URL](opt.New(*u))
	}

	var slug opt.Optional[mark.Slug]
	if args.Slug != nil {
		s, err := mark.NewSlug(*args.Slug)
		if err != nil {
			return NewError[mcp.ToolLibraryPageCreateOutput](err)
		}
		slug = opt.New(*s)
	}

	var parent opt.Optional[library.QueryKey]
	if args.Parent != nil {
		parent = opt.New(library.NewKey(*args.Parent))
	}

	var vis opt.Optional[visibility.Visibility]
	if args.Visibility != nil {
		v, err := visibility.NewVisibility(string(*args.Visibility))
		if err != nil {
			return NewError[mcp.ToolLibraryPageCreateOutput](err)
		}
		vis = opt.New(v)
	} else {
		vis = opt.New(visibility.VisibilityPublished)
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return NewError[mcp.ToolLibraryPageCreateOutput](err)
	}

	node, err := lt.nodeMutator.Create(ctx,
		accountID,
		args.Name,
		node_mutate.Partial{
			Slug:       slug,
			Content:    richContent,
			URL:        urlParsed,
			Parent:     parent,
			Visibility: vis,
		},
	)
	if err != nil {
		return NewError[mcp.ToolLibraryPageCreateOutput](err)
	}

	output := mcp.ToolLibraryPageCreateOutput{
		Slug: node.Mark.Slug(),
		Name: node.Name,
	}

	return NewSuccess(output)
}

func (lt *libraryTools) newLibraryPageUpdateTool() *Tool {
	toolDef := mcp.GetLibraryPageUpdateTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				lt.ExecuteLibraryPageUpdate,
			)
		},
	}
}

func (lt *libraryTools) ExecuteLibraryPageUpdate(ctx tool.Context, args mcp.ToolLibraryPageUpdateInput) ToolResult[mcp.ToolLibraryPageUpdateOutput] {
	partial := node_mutate.Partial{}

	if args.Name != nil {
		partial.Name = opt.New(*args.Name)
	}

	if args.Content != nil {
		richContent, err := datagraph.NewRichText(*args.Content)
		if err != nil {
			return NewError[mcp.ToolLibraryPageUpdateOutput](err)
		}
		partial.Content = opt.New(richContent)
	}

	if args.Visibility != nil {
		vis, err := visibility.NewVisibility(string(*args.Visibility))
		if err != nil {
			return NewError[mcp.ToolLibraryPageUpdateOutput](err)
		}
		partial.Visibility = opt.New(vis)
	}

	if args.Url != nil {
		u, err := url.Parse(*args.Url)
		if err != nil {
			return NewError[mcp.ToolLibraryPageUpdateOutput](err)
		}
		partial.URL = deletable.Skip(opt.New(*u))
	}

	if args.Parent != nil {
		partial.Parent = opt.New(library.NewKey(*args.Parent))
	}

	node, err := lt.nodeMutator.Update(ctx, library.NewKey(args.Slug), partial)
	if err != nil {
		return NewError[mcp.ToolLibraryPageUpdateOutput](err)
	}

	output := mcp.ToolLibraryPageUpdateOutput{
		Slug: node.Mark.Slug(),
		Name: node.Name,
	}

	return NewSuccess(output)
}

func (lt *libraryTools) newLibraryPageSearchTool() *Tool {
	toolDef := mcp.GetLibraryPageSearchTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				lt.ExecuteLibraryPageSearch,
			)
		},
	}
}

func (lt *libraryTools) ExecuteLibraryPageSearch(ctx tool.Context, args mcp.ToolLibraryPageSearchInput) ToolResult[mcp.ToolLibraryPageSearchOutput] {
	pp := pagination.NewPageParams(1, 50)

	opts := searcher.Options{
		Kinds: opt.New([]datagraph.Kind{datagraph.KindNode}),
	}

	result, err := lt.searcher.Search(ctx, args.Query, pp, opts)
	if err != nil {
		return NewError[mcp.ToolLibraryPageSearchOutput](err)
	}

	items := dt.Map(result.Items, func(item datagraph.Item) mcp.LibraryPageSearchItem {
		desc := item.GetDesc()
		content := item.GetContent().Plaintext()
		return mcp.LibraryPageSearchItem{
			Slug:        item.GetSlug(),
			Name:        item.GetName(),
			Description: &desc,
			Content:     &content,
		}
	})

	output := mcp.ToolLibraryPageSearchOutput{
		Results: result.Results,
		Items:   items,
	}

	return NewSuccess(output)
}
