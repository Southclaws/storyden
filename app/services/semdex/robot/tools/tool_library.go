package tools

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_properties"
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
	registry.Register(t.newLibraryPagePropertySchemaGetTool())
	registry.Register(t.newLibraryPagePropertySchemaUpdateTool())
	registry.Register(t.newLibraryPagePropertiesUpdateTool())

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

func (lt *libraryTools) ExecuteLibraryPageTree(ctx tool.Context, args mcp.ToolLibraryPageTreeInput) (*mcp.ToolLibraryPageTreeOutput, error) {
	account, err := session.GetAccount(ctx)
	if err != nil {
		return nil, err
	}

	acc, err := lt.accountQuery.GetByID(ctx, account.ID)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	pages := dt.Map(tree, mapNodeToTreeNode)

	output := mcp.ToolLibraryPageTreeOutput{
		Pages: pages,
	}

	return &output, nil
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

func (lt *libraryTools) ExecuteLibraryPageGet(ctx tool.Context, args mcp.ToolLibraryPageGetInput) (*mcp.ToolLibraryPageGetOutput, error) {
	node, err := lt.nodeReader.GetBySlug(ctx, library.NewKey(args.Id), nil)
	if err != nil {
		return nil, err
	}

	output := mcp.ToolLibraryPageGetOutput{
		Slug:        node.Mark.Slug(),
		Name:        node.Name,
		Description: node.Description.Ptr(),
		Tags:        dt.Map(node.Tags, func(t *tag_ref.Tag) string { return t.Name.String() }),
		ChildPages:  dt.Map(node.Nodes, func(n *library.Node) string { return n.Mark.Slug() }),
	}

	return &output, nil
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

func (lt *libraryTools) ExecuteLibraryPageCreate(ctx tool.Context, args mcp.ToolLibraryPageCreateInput) (*mcp.ToolLibraryPageCreateOutput, error) {
	var richContent opt.Optional[datagraph.Content]
	if args.Content != nil {
		rc, err := datagraph.NewRichText(*args.Content)
		if err != nil {
			return nil, err
		}
		richContent = opt.New(rc)
	}

	var urlParsed deletable.Value[url.URL]
	if args.Url != nil {
		u, err := url.Parse(*args.Url)
		if err != nil {
			return nil, err
		}
		urlParsed = deletable.Skip[url.URL](opt.New(*u))
	}

	var slug opt.Optional[mark.Slug]
	if args.Slug != nil {
		s, err := mark.NewSlug(*args.Slug)
		if err != nil {
			return nil, err
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
			return nil, err
		}
		vis = opt.New(v)
	} else {
		vis = opt.New(visibility.VisibilityPublished)
	}

	var tagNames opt.Optional[tag_ref.Names]
	if len(args.Tags) > 0 {
		names := dt.Map(args.Tags, func(t string) tag_ref.Name {
			return tag_ref.NewName(t)
		})
		tagNames = opt.New(tag_ref.Names(names))
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, err
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
			Tags:       tagNames,
		},
	)
	if err != nil {
		return nil, err
	}

	output := mcp.ToolLibraryPageCreateOutput{
		Slug: node.Mark.Slug(),
		Name: node.Name,
	}

	return &output, nil
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

func (lt *libraryTools) ExecuteLibraryPageUpdate(ctx tool.Context, args mcp.ToolLibraryPageUpdateInput) (*mcp.ToolLibraryPageUpdateOutput, error) {
	partial := node_mutate.Partial{}

	if args.Name != nil {
		partial.Name = opt.New(*args.Name)
	}

	if args.Content != nil {
		richContent, err := datagraph.NewRichText(*args.Content)
		if err != nil {
			return nil, err
		}
		partial.Content = opt.New(richContent)
	}

	if args.Visibility != nil {
		vis, err := visibility.NewVisibility(string(*args.Visibility))
		if err != nil {
			return nil, err
		}
		partial.Visibility = opt.New(vis)
	}

	if args.Url != nil {
		u, err := url.Parse(*args.Url)
		if err != nil {
			return nil, err
		}
		partial.URL = deletable.Skip(opt.New(*u))
	}

	if args.Parent != nil {
		partial.Parent = opt.New(library.NewKey(*args.Parent))
	}

	if len(args.Tags) > 0 {
		names := dt.Map(args.Tags, func(t string) tag_ref.Name {
			return tag_ref.NewName(t)
		})
		partial.Tags = opt.New(tag_ref.Names(names))
	}

	node, err := lt.nodeMutator.Update(ctx, library.NewKey(args.Id), partial)
	if err != nil {
		return nil, err
	}

	output := mcp.ToolLibraryPageUpdateOutput{
		Slug: node.Mark.Slug(),
		Name: node.Name,
	}

	return &output, nil
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

func (lt *libraryTools) ExecuteLibraryPageSearch(ctx tool.Context, args mcp.ToolLibraryPageSearchInput) (*mcp.ToolLibraryPageSearchOutput, error) {
	pp := pagination.NewPageParams(1, 50)

	opts := searcher.Options{
		Kinds: opt.New([]datagraph.Kind{datagraph.KindNode}),
	}

	result, err := lt.searcher.Search(ctx, args.Query, pp, opts)
	if err != nil {
		return nil, err
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

	return &output, nil
}

func (lt *libraryTools) newLibraryPagePropertySchemaGetTool() *Tool {
	toolDef := mcp.GetLibraryPagePropertySchemaGetTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				lt.ExecuteLibraryPagePropertySchemaGet,
			)
		},
	}
}

func (lt *libraryTools) ExecuteLibraryPagePropertySchemaGet(ctx tool.Context, args mcp.ToolLibraryPagePropertySchemaGetInput) (*mcp.ToolLibraryPagePropertySchemaGetOutput, error) {
	node, err := lt.nodeReader.GetBySlug(ctx, library.NewKey(args.Id), nil)
	if err != nil {
		return nil, err
	}

	output := mcp.ToolLibraryPagePropertySchemaGetOutput{
		Fields: []mcp.PropertySchemaField{},
	}

	childSchema, hasChildSchema := node.ChildProperties.Get()
	if !hasChildSchema {
		hasSchema := false
		output.HasSchema = &hasSchema
		return &output, nil
	}

	hasSchema := true
	output.HasSchema = &hasSchema

	output.Fields = dt.Map(childSchema.Fields, func(f *library.PropertySchemaField) mcp.PropertySchemaField {
		return mcp.PropertySchemaField{
			Id:   f.ID.String(),
			Name: f.Name,
			Type: mcp.PropertySchemaFieldType(f.Type.String()),
			Sort: &f.Sort,
		}
	})

	return &output, nil
}

func (lt *libraryTools) newLibraryPagePropertySchemaUpdateTool() *Tool {
	toolDef := mcp.GetLibraryPagePropertySchemaUpdateTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				lt.ExecuteLibraryPagePropertySchemaUpdate,
			)
		},
	}
}

func (lt *libraryTools) ExecuteLibraryPagePropertySchemaUpdate(ctx tool.Context, args mcp.ToolLibraryPagePropertySchemaUpdateInput) (*mcp.ToolLibraryPagePropertySchemaUpdateOutput, error) {
	mutations, err := dt.MapErr(args.Fields, func(f mcp.PropertySchemaFieldMutation) (*node_properties.SchemaFieldMutation, error) {
		propType, err := library.NewPropertyType(string(f.Type))
		if err != nil {
			return nil, err
		}

		var fieldID opt.Optional[xid.ID]
		if f.Id != nil {
			id, err := xid.FromString(*f.Id)
			if err != nil {
				return nil, err
			}
			fieldID = opt.New(id)
		}

		return &node_properties.SchemaFieldMutation{
			ID:   fieldID,
			Name: f.Name,
			Type: propType,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	updated, err := lt.schemaUpdater.UpdateChildren(ctx, library.NewKey(args.Id), mutations)
	if err != nil {
		return nil, err
	}

	output := mcp.ToolLibraryPagePropertySchemaUpdateOutput{
		Fields: dt.Map(updated.Fields, func(f *library.PropertySchemaField) mcp.PropertySchemaFieldResult {
			return mcp.PropertySchemaFieldResult{
				Id:   f.ID.String(),
				Name: f.Name,
				Type: mcp.PropertySchemaFieldResultType(f.Type.String()),
			}
		}),
	}

	return &output, nil
}

func (lt *libraryTools) newLibraryPagePropertiesUpdateTool() *Tool {
	toolDef := mcp.GetLibraryPagePropertiesUpdateTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				lt.ExecuteLibraryPagePropertiesUpdate,
			)
		},
	}
}

func (lt *libraryTools) ExecuteLibraryPagePropertiesUpdate(ctx tool.Context, args mcp.ToolLibraryPagePropertiesUpdateInput) (*mcp.ToolLibraryPagePropertiesUpdateOutput, error) {
	mutations, err := dt.MapErr(args.Properties, func(p mcp.PropertyValueMutation) (*library.PropertyMutation, error) {
		fieldID, err := xid.FromString(p.FieldId)
		if err != nil {
			return nil, err
		}

		return &library.PropertyMutation{
			ID:    opt.New(fieldID),
			Value: p.Value,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	node, err := lt.nodeMutator.Update(ctx, library.NewKey(args.Id), node_mutate.Partial{
		Properties: opt.New(library.PropertyMutationList(mutations)),
	})
	if err != nil {
		return nil, err
	}

	output := mcp.ToolLibraryPagePropertiesUpdateOutput{
		Properties: []mcp.PropertyValueResult{},
	}

	if props, ok := node.Properties.Get(); ok {
		output.Properties = dt.Map(props.Properties, func(p *library.Property) mcp.PropertyValueResult {
			propType := mcp.PropertyValueResultType(p.Field.Type.String())
			return mcp.PropertyValueResult{
				FieldId: p.Field.ID.String(),
				Name:    p.Field.Name,
				Type:    &propType,
				Value:   p.Value.OrZero(),
			}
		})
	}

	return &output, nil
}
