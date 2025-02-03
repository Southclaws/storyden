package bindings

import (
	"context"
	"net/url"
	"strconv"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_traversal"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/library/node_mutate"
	"github.com/Southclaws/storyden/app/services/library/node_read"
	"github.com/Southclaws/storyden/app/services/library/node_visibility"
	"github.com/Southclaws/storyden/app/services/library/nodetree"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/deletable"
)

type Nodes struct {
	accountQuery *account_querier.Querier
	nodeMutator  *node_mutate.Manager
	nodeReader   *node_read.HydratedQuerier
	nv           *node_visibility.Controller
	ntree        nodetree.Graph
	ntr          node_traversal.Repository
}

func NewNodes(
	accountQuery *account_querier.Querier,
	nodeMutator *node_mutate.Manager,
	nodeReader *node_read.HydratedQuerier,
	nv *node_visibility.Controller,
	ntree nodetree.Graph,
	ntr node_traversal.Repository,
) Nodes {
	return Nodes{
		accountQuery: accountQuery,
		nodeMutator:  nodeMutator,
		nodeReader:   nodeReader,
		nv:           nv,
		ntree:        ntree,
		ntr:          ntr,
	}
}

func (c *Nodes) NodeCreate(ctx context.Context, request openapi.NodeCreateRequestObject) (openapi.NodeCreateResponseObject, error) {
	session, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	vis, err := opt.MapErr(opt.NewPtr(request.Body.Visibility), deserialiseVisibility)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	richContent, err := opt.MapErr(opt.NewPtr(request.Body.Content), datagraph.NewRichText)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	url, err := opt.MapErr(opt.NewPtr(request.Body.Url), func(s string) (url.URL, error) {
		u, err := url.Parse(s)
		if err != nil {
			return url.URL{}, err
		}
		return *u, nil
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	tags := opt.Map(opt.NewPtr(request.Body.Tags), func(tags []string) tag_ref.Names {
		return dt.Map(tags, deserialiseTagName)
	})

	slug, err := deserialiseInputSlug(request.Body.Slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	primaryImage := opt.Map(opt.NewPtr(request.Body.PrimaryImageAssetId), deserialiseAssetID)

	node, err := c.nodeMutator.Create(ctx,
		session,
		request.Body.Name,
		node_mutate.Partial{
			Slug:         slug,
			PrimaryImage: deletable.Skip(primaryImage),
			Content:      richContent,
			Metadata:     opt.NewPtr((*map[string]any)(request.Body.Meta)),
			URL:          url,
			AssetsAdd:    opt.NewPtrMap(request.Body.AssetIds, deserialiseAssetIDs),
			AssetSources: opt.NewPtrMap(request.Body.AssetSources, deserialiseAssetSources),
			Parent:       opt.NewPtrMap(request.Body.Parent, deserialiseNodeMark),
			Tags:         tags,
			Visibility:   vis,
		},
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeCreate200JSONResponse{
		NodeCreateOKJSONResponse: openapi.NodeCreateOKJSONResponse(serialiseNode(node)),
	}, nil
}

func (c *Nodes) NodeList(ctx context.Context, request openapi.NodeListRequestObject) (openapi.NodeListResponseObject, error) {
	depth, err := opt.MapErr(opt.NewPtr(request.Params.Depth), func(s string) (int, error) {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return 0, err
		}

		return max(0, int(v)), nil
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// TODO: Clean this mess up.
	acc, err := opt.MapErr(session.GetOptAccountID(ctx), func(aid account.AccountID) (account.Account, error) {
		a, err := c.accountQuery.GetByID(ctx, aid)
		if err != nil {
			return account.Account{}, err
		}

		return *a, nil
	})
	if err != nil {
		if ftag.Get(err) == ftag.NotFound {
			acc = opt.NewEmpty[account.Account]()
		} else {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	var cs []*library.Node

	opts := []node_traversal.Filter{}

	author := opt.NewPtr(request.Params.Author)
	if v, ok := author.Get(); ok {
		opts = append(opts, node_traversal.WithRootOwner(v))
	}

	if d, ok := depth.Get(); ok {
		opts = append(opts, node_traversal.WithDepth(uint(d)))
	}

	visibilities, err := opt.MapErr(opt.NewPtr(request.Params.Visibility), deserialiseVisibilityList)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if v, ok := visibilities.Get(); ok {
		opts = append(opts, node_traversal.WithVisibility(acc, v...))
	}

	nid, err := opt.MapErr(opt.NewPtr(request.Params.NodeId), library.NodeIDFromString)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	flatten := opt.NewPtr(request.Params.Format).Or(openapi.NodeListParamsFormatTree) == openapi.NodeListParamsFormatFlat

	cs, err = c.ntr.Subtree(ctx, nid, flatten, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeList200JSONResponse{
		NodeListOKJSONResponse: openapi.NodeListOKJSONResponse{
			Nodes: dt.Map(cs, serialiseNodeWithItems),
		},
	}, nil
}

func (c *Nodes) NodeGet(ctx context.Context, request openapi.NodeGetRequestObject) (openapi.NodeGetResponseObject, error) {
	node, err := c.nodeReader.GetBySlug(ctx, deserialiseNodeMark(request.NodeSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeGet200JSONResponse{
		NodeGetOKJSONResponse: openapi.NodeGetOKJSONResponse(serialiseNodeWithItems(node)),
	}, nil
}

func (c *Nodes) NodeUpdate(ctx context.Context, request openapi.NodeUpdateRequestObject) (openapi.NodeUpdateResponseObject, error) {
	content, err := opt.MapErr(opt.NewPtr(request.Body.Content), datagraph.NewRichText)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	url, err := opt.MapErr(opt.NewPtr(request.Body.Url), func(s string) (url.URL, error) {
		u, err := url.Parse(s)
		if err != nil {
			return url.URL{}, err
		}
		return *u, nil
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	slug, err := deserialiseInputSlug(request.Body.Slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	titleFillRuleParam, err := opt.MapErr(opt.NewPtr(request.Params.TitleFillRule), deserialiseTitleFillRule)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tagFillRuleParam, err := opt.MapErr(opt.NewPtr(request.Params.TagFillRule), deserialiseTagFillRule)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	fillSource, err := opt.MapErr(opt.NewPtr(request.Params.FillSource), deserialiseFillSource)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	contentFillCmd, err := getContentFillRuleSourceCommand(request.Params.ContentFillRule, request.Params.FillSource)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	tags := opt.Map(opt.NewPtr(request.Body.Tags), func(tags []string) tag_ref.Names {
		return dt.Map(tags, deserialiseTagName)
	})

	primaryImage := deletable.NewMap(request.Body.PrimaryImageAssetId, deserialiseAssetID)

	partial := node_mutate.Partial{
		Name:         opt.NewPtr(request.Body.Name),
		Slug:         slug,
		AssetsAdd:    opt.NewPtrMap(request.Body.AssetIds, deserialiseAssetIDs),
		AssetSources: opt.NewPtrMap(request.Body.AssetSources, deserialiseAssetSources),
		URL:          url,
		Content:      content,
		PrimaryImage: primaryImage,
		Parent:       opt.NewPtrMap(request.Body.Parent, deserialiseNodeMark),
		Tags:         tags,
		Metadata:     opt.NewPtr((*map[string]any)(request.Body.Meta)),
		FillSource:   fillSource,
		ContentFill:  contentFillCmd,
	}

	if tfr, ok := titleFillRuleParam.Get(); ok {
		partial.TitleFill = opt.New(datagraph.TitleFillCommand{FillRule: tfr})
	}

	if tfr, ok := tagFillRuleParam.Get(); ok {
		partial.TagFill = opt.New(tag.TagFillCommand{FillRule: tfr})
	}

	node, err := c.nodeMutator.Update(ctx, deserialiseNodeMark(request.NodeSlug), partial)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeUpdate200JSONResponse{
		NodeUpdateOKJSONResponse: openapi.NodeUpdateOKJSONResponse(serialiseUpdatedNode(node)),
	}, nil
}

func (c *Nodes) NodeUpdateVisibility(ctx context.Context, request openapi.NodeUpdateVisibilityRequestObject) (openapi.NodeUpdateVisibilityResponseObject, error) {
	v, err := visibility.NewVisibility(string(request.Body.Visibility))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	node, err := c.nv.ChangeVisibility(ctx, deserialiseNodeMark(request.NodeSlug), v)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeUpdateVisibility200JSONResponse{
		NodeUpdateOKJSONResponse: openapi.NodeUpdateOKJSONResponse(serialiseNode(node)),
	}, nil
}

func (c *Nodes) NodeDelete(ctx context.Context, request openapi.NodeDeleteRequestObject) (openapi.NodeDeleteResponseObject, error) {
	destinationNode, err := c.nodeMutator.Delete(ctx, deserialiseNodeMark(request.NodeSlug), node_mutate.DeleteOptions{
		NewParent: opt.NewPtrMap(request.Params.TargetNode, deserialiseNodeMark),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeDelete200JSONResponse{
		NodeDeleteOKJSONResponse: openapi.NodeDeleteOKJSONResponse{
			Destination: opt.Map(opt.NewPtr(destinationNode), func(in library.Node) openapi.Node {
				return serialiseNode(&in)
			}).Ptr(),
		},
	}, nil
}

func (c *Nodes) NodeUpdateChildrenPropertySchema(ctx context.Context, request openapi.NodeUpdateChildrenPropertySchemaRequestObject) (openapi.NodeUpdateChildrenPropertySchemaResponseObject, error) {
	return nil, nil
}

func (c *Nodes) NodeUpdateProperties(ctx context.Context, request openapi.NodeUpdatePropertiesRequestObject) (openapi.NodeUpdatePropertiesResponseObject, error) {
	return nil, nil
}

func (c *Nodes) NodeAddAsset(ctx context.Context, request openapi.NodeAddAssetRequestObject) (openapi.NodeAddAssetResponseObject, error) {
	id := openapi.ParseID(request.AssetId)

	contentFillCmd, err := getContentFillRuleCommand(request.Params.ContentFillRule, request.Params.NodeContentFillTarget)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	node, err := c.nodeMutator.Update(ctx, deserialiseNodeMark(request.NodeSlug), node_mutate.Partial{
		AssetsAdd:   opt.New([]asset.AssetID{id}),
		ContentFill: contentFillCmd,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeAddAsset200JSONResponse{
		NodeUpdateOKJSONResponse: openapi.NodeUpdateOKJSONResponse(serialiseUpdatedNode(node)),
	}, nil
}

func (c *Nodes) NodeRemoveAsset(ctx context.Context, request openapi.NodeRemoveAssetRequestObject) (openapi.NodeRemoveAssetResponseObject, error) {
	id := openapi.ParseID(request.AssetId)

	node, err := c.nodeMutator.Update(ctx, deserialiseNodeMark(request.NodeSlug), node_mutate.Partial{
		AssetsRemove: opt.New([]asset.AssetID{id}),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeRemoveAsset200JSONResponse{
		NodeUpdateOKJSONResponse: openapi.NodeUpdateOKJSONResponse(serialiseUpdatedNode(node)),
	}, nil
}

func (c *Nodes) NodeAddNode(ctx context.Context, request openapi.NodeAddNodeRequestObject) (openapi.NodeAddNodeResponseObject, error) {
	node, err := c.ntree.Move(ctx, deserialiseNodeMark(request.NodeSlugChild), deserialiseNodeMark(request.NodeSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeAddNode200JSONResponse{
		NodeAddChildOKJSONResponse: openapi.NodeAddChildOKJSONResponse(serialiseNode(node)),
	}, nil
}

func (c *Nodes) NodeRemoveNode(ctx context.Context, request openapi.NodeRemoveNodeRequestObject) (openapi.NodeRemoveNodeResponseObject, error) {
	node, err := c.ntree.Sever(ctx, deserialiseNodeMark(request.NodeSlugChild), deserialiseNodeMark(request.NodeSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeRemoveNode200JSONResponse{
		NodeRemoveChildOKJSONResponse: openapi.NodeRemoveChildOKJSONResponse(serialiseNode(node)),
	}, nil
}

func serialiseUpdatedNode(in *node_mutate.Updated) openapi.Node {
	n := serialiseNode(&in.Node)

	if ts, ok := in.TitleSuggestion.Get(); ok {
		n.TitleSuggestion = &ts
	}

	if ts, ok := in.TagSuggestions.Get(); ok {
		s := ts.Strings()
		n.TagSuggestions = &s
	}

	if cs, ok := in.ContentSuggestion.Get(); ok {
		html := cs.HTML()
		n.ContentSuggestion = &html
	}

	return n
}

func serialiseNode(in *library.Node) openapi.Node {
	return openapi.Node{
		Id:           in.Mark.ID().String(),
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
		Name:         in.Name,
		Slug:         in.Mark.Slug(),
		Assets:       dt.Map(in.Assets, serialiseAssetPtr),
		Link:         opt.Map(in.WebLink, serialiseLinkRef).Ptr(),
		Description:  in.GetDesc(),
		PrimaryImage: opt.Map(in.PrimaryImage, serialiseAsset).Ptr(),
		Content:      opt.Map(in.Content, serialiseContentHTML).Ptr(),
		Owner:        serialiseProfileReference(in.Owner),
		Parent: opt.PtrMap(in.Parent, func(in library.Node) openapi.Node {
			return serialiseNode(&in)
		}),
		Tags:       serialiseTagReferenceList(in.Tags),
		Visibility: serialiseVisibility(in.Visibility),
		Meta:       in.Metadata,
	}
}

func serialiseNodeWithItems(in *library.Node) openapi.NodeWithChildren {
	rs := opt.Map(in.RelevanceScore, func(v float64) float32 { return float32(v) })
	return openapi.NodeWithChildren{
		Id:           in.Mark.ID().String(),
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
		Name:         in.Name,
		Slug:         in.Mark.Slug(),
		Assets:       dt.Map(in.Assets, serialiseAssetPtr),
		Link:         opt.Map(in.WebLink, serialiseLinkRef).Ptr(),
		Description:  in.GetDesc(),
		PrimaryImage: opt.Map(in.PrimaryImage, serialiseAsset).Ptr(),
		Content:      opt.Map(in.Content, serialiseContentHTML).Ptr(),
		Owner:        serialiseProfileReference(in.Owner),
		Parent: opt.PtrMap(in.Parent, func(in library.Node) openapi.Node {
			return serialiseNode(&in)
		}),
		Properties:          serialisePropertyList(in.Properties),
		ChildPropertySchema: serialisePropertySchemaList(in.ChildProperties),
		Tags:                serialiseTagReferenceList(in.Tags),
		Visibility:          serialiseVisibility(in.Visibility),
		RelevanceScore:      rs.Ptr(),
		Meta:                in.Metadata,
		Children:            dt.Map(in.Nodes, serialiseNodeWithItems),
	}
}

func deserialiseNodeMark(in string) library.QueryKey {
	return library.QueryKey{deserialiseMark(in)}
}

func deserialiseAssetSources(in openapi.AssetSourceList) []string {
	return dt.Map(in, deserialiseAssetSourceURL)
}

func deserialiseAssetSourceURL(in openapi.AssetSourceURL) string {
	return string(in)
}

func deserialiseInputSlug(in *string) (opt.Optional[mark.Slug], error) {
	if in == nil {
		return opt.NewEmpty[mark.Slug](), nil
	}

	if *in == "" {
		return opt.NewEmpty[mark.Slug](), nil
	}

	slug, err := mark.NewSlug(*in)
	if err != nil {
		return nil, err
	}

	return opt.New(*slug), nil
}

func deserialiseTitleFillRule(in openapi.TitleFillRule) (datagraph.TitleFillRule, error) {
	return datagraph.NewTitleFillRule(string(in))
}

func deserialiseContentFillRule(in openapi.ContentFillRule) (asset.ContentFillRule, error) {
	return asset.NewContentFillRule(string(in))
}

func deserialiseFillSource(in openapi.FillSource) (asset.FillSource, error) {
	return asset.NewFillSource(string(in))
}

func serialiseProperty(in library.Property) openapi.Property {
	return openapi.Property{
		Name:  in.Name,
		Type:  in.Type,
		Value: opt.Map(in.Value, func(v string) string { return v }).Ptr(),
	}
}

func serialisePropertyList(in library.PropertyTable) openapi.PropertyList {
	return dt.Map(in, serialiseProperty)
}

func serialisePropertySchema(in library.PropertySchema) openapi.PropertySchema {
	return openapi.PropertySchema{
		Name: in.Name,
		Type: in.Type,
	}
}

func serialisePropertySchemaList(in library.PropertySchemas) *openapi.PropertySchemaList {
	if len(in) == 0 {
		return nil
	}
	schemas := dt.Map(in, serialisePropertySchema)
	return &schemas
}

func getContentFillRuleSourceCommand(contentFillRuleParam *openapi.ContentFillRule, contentFillSourceParam *openapi.FillSourceQuery) (opt.Optional[asset.ContentFillCommand], error) {
	if contentFillRuleParam != nil {
		if contentFillSourceParam == nil {
			return nil, fault.New("node_content_fill_target is required when content_fill_rule is specified")
		}

		rule, err := asset.NewContentFillRule((string)(*contentFillRuleParam))
		if err != nil {
			return nil, fault.Wrap(err)
		}

		sourceType, err := asset.NewFillSource((string)(*contentFillSourceParam))
		if err != nil {
			return nil, fault.Wrap(err)
		}

		return opt.New(asset.ContentFillCommand{
			SourceType: opt.New(sourceType),
			FillRule:   rule,
		}), nil
	}

	return opt.NewEmpty[asset.ContentFillCommand](), nil
}
