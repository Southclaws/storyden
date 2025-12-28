package bindings

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_cache"
	"github.com/Southclaws/storyden/app/resources/library/node_properties"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_traversal"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/generative"
	"github.com/Southclaws/storyden/app/services/library/node_mutate"
	"github.com/Southclaws/storyden/app/services/library/node_property_schema"
	"github.com/Southclaws/storyden/app/services/library/node_read"
	"github.com/Southclaws/storyden/app/services/library/node_visibility"
	"github.com/Southclaws/storyden/app/services/library/nodetree"
	"github.com/Southclaws/storyden/app/services/reqinfo"
	"github.com/Southclaws/storyden/app/services/tag/autotagger"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/deletable"
)

type Nodes struct {
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
	node_cache    *node_cache.Cache
}

func NewNodes(
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
	node_cache *node_cache.Cache,
) Nodes {
	return Nodes{
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
		node_cache:    node_cache,
	}
}

func (c *Nodes) NodeCreate(ctx context.Context, request openapi.NodeCreateRequestObject) (openapi.NodeCreateResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := session.Authorise(ctx, nil, rbac.PermissionManageLibrary, rbac.PermissionSubmitLibraryNode); err != nil {
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

	pml, err := opt.MapErr(opt.NewPtr(request.Body.Properties), deserialisePropertyMutationList)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
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
		accountID,
		request.Body.Name,
		node_mutate.Partial{
			Slug:         slug,
			PrimaryImage: deletable.Skip(primaryImage),
			Content:      richContent,
			Metadata:     opt.NewPtr((*map[string]any)(request.Body.Meta)),
			URL:          deletable.Skip(url),
			Description:  opt.NewPtr(request.Body.Description),
			AssetsAdd:    opt.NewPtrMap(request.Body.AssetIds, deserialiseAssetIDs),
			AssetSources: opt.NewPtrMap(request.Body.AssetSources, deserialiseAssetSources),
			Parent:       opt.NewPtrMap(request.Body.Parent, deserialiseNodeMark),
			HideChildren: opt.NewPtr(request.Body.HideChildTree), Tags: tags,
			Visibility: vis,
			Properties: pml,
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

	acc, err := opt.MapErr(session.GetOptAccountID(ctx), func(aid account.AccountID) (account.AccountWithEdges, error) {
		a, err := c.accountQuery.GetByID(ctx, aid)
		if err != nil {
			return account.AccountWithEdges{}, err
		}

		return *a, nil
	})
	if err != nil {
		if ftag.Get(err) == ftag.NotFound {
			acc = opt.NewEmpty[account.AccountWithEdges]()
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
	pp := deserialisePageParams(request.Params.Page, 100)
	sortChildrenBy := opt.NewPtrMap(request.Params.ChildrenSort, func(cs string) node_querier.ChildSortRule {
		return node_querier.NewChildSortRule(cs, pp)
	})

	qk := deserialiseNodeMark(request.NodeSlug)
	cacheKey := qk.String()

	etag, notModified := c.node_cache.Check(ctx, reqinfo.GetCacheQuery(ctx), cacheKey)
	if notModified {
		return openapi.NodeGet304Response{
			Headers: openapi.NotModifiedResponseHeaders{
				CacheControl: getAuthStateCacheControl(ctx, "no-cache"),
				LastModified: etag.Time.Format(time.RFC1123),
				ETag:         etag.String(),
			},
		}, nil
	}

	node, err := c.nodeReader.GetBySlug(ctx, qk, sortChildrenBy)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if etag == nil {
		c.node_cache.Store(ctx, cacheKey, node.UpdatedAt)
		etag = cachecontrol.NewETag(node.UpdatedAt)
	}

	return openapi.NodeGet200JSONResponse{
		NodeGetOKJSONResponse: openapi.NodeGetOKJSONResponse{
			Body: serialiseNodeWithItems(node),
			Headers: openapi.NodeGetOKResponseHeaders{
				CacheControl: getAuthStateCacheControl(ctx, "no-cache"),
				LastModified: etag.Time.Format(time.RFC1123),
				ETag:         etag.String(),
			},
		},
	}, nil
}

func (c *Nodes) NodeListChildren(ctx context.Context, request openapi.NodeListChildrenRequestObject) (openapi.NodeListChildrenResponseObject, error) {
	pp := deserialisePageParams(request.Params.Page, 100)

	opts := []node_querier.Option{}

	if request.Params.ChildrenSort != nil {
		opts = append(opts, node_querier.WithSortChildrenBy(node_querier.NewChildSortRule(*request.Params.ChildrenSort, pp)))
	}

	if request.Params.Q != nil {
		opts = append(opts, node_querier.WithSearchChildren(*request.Params.Q))
	}

	if request.Params.Tags != nil {
		tags := dt.Map(*request.Params.Tags, deserialiseTagName)
		opts = append(opts, node_querier.WithFilterChildrenByTags(tags...))
	}

	// NOTE: Visibility rules are automatically applied by the node_read.HydratedQuerier
	// service layer. This ensures that non-published nodes are not visible to
	// unauthorized users, addressing issue #450. The rules are the same as those
	// applied by other node endpoints:
	// - Published nodes are visible to everyone
	// - Draft/Unlisted nodes are only visible to their owners
	// - Review nodes are visible to admins/library managers + owners
	// - Unauthenticated users only see published nodes
	r, err := c.nodeReader.ListChildren(ctx, deserialiseNodeMark(request.NodeSlug), pp, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeListChildren200JSONResponse{
		NodeListOKJSONResponse: openapi.NodeListOKJSONResponse{
			CurrentPage: r.CurrentPage,
			NextPage:    r.NextPage.Ptr(),
			PageSize:    r.Size,
			Results:     r.Results,
			Nodes:       dt.Map(r.Items, serialiseNodeWithItems),
			TotalPages:  r.TotalPages,
		},
	}, nil
}

func (c *Nodes) NodeUpdate(ctx context.Context, request openapi.NodeUpdateRequestObject) (openapi.NodeUpdateResponseObject, error) {
	content, err := opt.MapErr(opt.NewPtr(request.Body.Content), datagraph.NewRichText)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	url, err := deletable.NewMapErr(request.Body.Url, func(s string) (url.URL, error) {
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

	pml, err := opt.MapErr(opt.NewPtr(request.Body.Properties), deserialisePropertyMutationList)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
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
		Description:  opt.NewPtr(request.Body.Description),
		Content:      content,
		PrimaryImage: primaryImage,
		Parent:       opt.NewPtrMap(request.Body.Parent, deserialiseNodeMark),
		HideChildren: opt.NewPtr(request.Body.HideChildTree),
		Properties:   pml,
		Tags:         tags,
		Metadata:     opt.NewPtr((*map[string]any)(request.Body.Meta)),
	}

	node, err := c.nodeMutator.Update(ctx, deserialiseNodeMark(request.NodeSlug), partial)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeUpdate200JSONResponse{
		NodeUpdateOKJSONResponse: openapi.NodeUpdateOKJSONResponse{
			Body: serialiseUpdatedNode(node),
			Headers: openapi.NodeUpdateOKResponseHeaders{
				LastModified: node.UpdatedAt.UTC().Format(time.RFC1123),
				CacheControl: "private, no-cache, no-store, must-revalidate",
			},
		},
	}, nil
}

func (c *Nodes) NodeGenerateContent(ctx context.Context, request openapi.NodeGenerateContentRequestObject) (openapi.NodeGenerateContentResponseObject, error) {
	content, err := datagraph.NewRichText(request.Body.Content)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	summary, err := c.summariser.Summarise(ctx, content)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeGenerateContent200JSONResponse{
		NodeGenerateContentOKJSONResponse: openapi.NodeGenerateContentOKJSONResponse{
			Content: summary,
		},
	}, nil
}

func (c *Nodes) NodeGenerateTags(ctx context.Context, request openapi.NodeGenerateTagsRequestObject) (openapi.NodeGenerateTagsResponseObject, error) {
	content, err := datagraph.NewRichText(request.Body.Content)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	tags, err := c.tagger.Gather(ctx, content)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeGenerateTags200JSONResponse{
		NodeGenerateTagsOKJSONResponse: openapi.NodeGenerateTagsOKJSONResponse{
			Tags: dt.Map(tags, func(t tag_ref.Name) string {
				return t.String()
			}),
		},
	}, nil
}

func (c *Nodes) NodeGenerateTitle(ctx context.Context, request openapi.NodeGenerateTitleRequestObject) (openapi.NodeGenerateTitleResponseObject, error) {
	content, err := datagraph.NewRichText(request.Body.Content)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	title, err := c.titler.SuggestTitle(ctx, content)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// TODO: Support multiple suggestions in API.
	if len(title) == 0 {
		return nil, fault.New("no title suggestions returned", fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return openapi.NodeGenerateTitle200JSONResponse{
		NodeGenerateTitleOKJSONResponse: openapi.NodeGenerateTitleOKJSONResponse{
			Title: title[0],
		},
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
		NodeUpdateOKJSONResponse: openapi.NodeUpdateOKJSONResponse{
			Body: serialiseNodeWithItems(node),
			Headers: openapi.NodeUpdateOKResponseHeaders{
				LastModified: node.UpdatedAt.UTC().Format(time.RFC1123),
			},
		},
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
	schemas, err := dt.MapErr(*request.Body, deserialisePropertySchemaMutation)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	updated, err := c.schemaUpdater.UpdateChildren(ctx, deserialiseNodeMark(request.NodeSlug), schemas)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeUpdateChildrenPropertySchema200JSONResponse{
		NodeUpdatePropertySchemaOKJSONResponse: openapi.NodeUpdatePropertySchemaOKJSONResponse{
			Properties: serialisePropertySchemas(*updated),
		},
	}, nil
}

func (c *Nodes) NodeUpdatePropertySchema(ctx context.Context, request openapi.NodeUpdatePropertySchemaRequestObject) (openapi.NodeUpdatePropertySchemaResponseObject, error) {
	schemas, err := dt.MapErr(*request.Body, deserialisePropertySchemaMutation)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	updated, err := c.schemaUpdater.UpdateSiblings(ctx, deserialiseNodeMark(request.NodeSlug), schemas)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeUpdatePropertySchema200JSONResponse{
		NodeUpdatePropertySchemaOKJSONResponse: openapi.NodeUpdatePropertySchemaOKJSONResponse{
			Properties: serialisePropertySchemas(*updated),
		},
	}, nil
}

func (c *Nodes) NodeUpdateProperties(ctx context.Context, request openapi.NodeUpdatePropertiesRequestObject) (openapi.NodeUpdatePropertiesResponseObject, error) {
	pml, err := deserialisePropertyMutationList(request.Body.Properties)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	updated, err := c.nodeMutator.Update(ctx, deserialiseNodeMark(request.NodeSlug), node_mutate.Partial{
		Properties: opt.New(pml),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	properties := serialisePropertyTableOpt(updated.Properties)

	return openapi.NodeUpdateProperties200JSONResponse{
		NodeUpdatePropertiesOKJSONResponse: openapi.NodeUpdatePropertiesOKJSONResponse{
			Properties: properties,
		},
	}, nil
}

func (c *Nodes) NodeAddAsset(ctx context.Context, request openapi.NodeAddAssetRequestObject) (openapi.NodeAddAssetResponseObject, error) {
	id := openapi.ParseID(request.AssetId)

	node, err := c.nodeMutator.Update(ctx, deserialiseNodeMark(request.NodeSlug), node_mutate.Partial{
		AssetsAdd: opt.New([]asset.AssetID{id}),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeAddAsset200JSONResponse{
		NodeUpdateOKJSONResponse: openapi.NodeUpdateOKJSONResponse{
			Body: serialiseNodeWithItems(node),
			Headers: openapi.NodeUpdateOKResponseHeaders{
				LastModified: node.UpdatedAt.UTC().Format(time.RFC1123),
			},
		},
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
		NodeUpdateOKJSONResponse: openapi.NodeUpdateOKJSONResponse{
			Body: serialiseNodeWithItems(node),
			Headers: openapi.NodeUpdateOKResponseHeaders{
				LastModified: node.UpdatedAt.UTC().Format(time.RFC1123),
			},
		},
	}, nil
}

func (c *Nodes) NodeAddNode(ctx context.Context, request openapi.NodeAddNodeRequestObject) (openapi.NodeAddNodeResponseObject, error) {
	_, err := c.ntree.Move(ctx, deserialiseNodeMark(request.NodeSlugChild), deserialiseNodeMark(request.NodeSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	node, err := c.nodeReader.GetBySlug(ctx, deserialiseNodeMark(request.NodeSlug), opt.NewEmpty[node_querier.ChildSortRule]())
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

func (c *Nodes) NodeUpdatePosition(ctx context.Context, request openapi.NodeUpdatePositionRequestObject) (openapi.NodeUpdatePositionResponseObject, error) {
	opts := nodetree.Options{}

	beforeID, err := opt.MapErr(opt.NewPtr(request.Body.Before), func(s string) (library.NodeID, error) {
		id, err := xid.FromString(s)
		return library.NodeID(id), err
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	afterID, err := opt.MapErr(opt.NewPtr(request.Body.After), func(s string) (library.NodeID, error) {
		id, err := xid.FromString(s)
		return library.NodeID(id), err
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	parentID, err := deletable.NewMapErr(request.Body.Parent, func(s string) (library.QueryKey, error) {
		id, err := xid.FromString(s)
		return library.NewID(id), err
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	opts.Parent = parentID
	opts.Before = beforeID
	opts.After = afterID

	n, err := c.npos.Move(ctx, deserialiseNodeMark(request.NodeSlug), opts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeUpdatePosition200JSONResponse{
		NodeUpdateOKJSONResponse: openapi.NodeUpdateOKJSONResponse{
			Body: serialiseNodeWithItems(n),
			Headers: openapi.NodeUpdateOKResponseHeaders{
				LastModified: n.UpdatedAt.UTC().Format(time.RFC1123),
			},
		},
	}, nil
}

func serialiseUpdatedNode(in *library.Node) openapi.NodeWithChildren {
	return serialiseNodeWithItems(in)
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
		HideChildTree: in.HideChildTree,
		Tags:          serialiseTagReferenceList(in.Tags),
		Visibility:    serialiseVisibility(in.Visibility),
		Meta:          in.Metadata,
	}
}

func serialiseNodeWithItems(in *library.Node) openapi.NodeWithChildren {
	rs := opt.Map(in.RelevanceScore, func(v float64) float32 { return float32(v) })
	properties := opt.Map(in.Properties, serialisePropertyTable)
	childPropertySchema := opt.Map(in.ChildProperties, serialisePropertySchemaList)

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
		HideChildTree:       in.HideChildTree,
		Properties:          properties.Or([]openapi.Property{}),
		ChildPropertySchema: childPropertySchema.Or([]openapi.PropertySchema{}),
		Tags:                serialiseTagReferenceList(in.Tags),
		Visibility:          serialiseVisibility(in.Visibility),
		RelevanceScore:      rs.Ptr(),
		Meta:                in.Metadata,
		Children:            dt.Map(in.Nodes, serialiseNodeWithItems),
	}
}

func deserialiseNodeMark(in string) library.QueryKey {
	return library.NewKey(in)
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

func serialiseProperty(in *library.Property) openapi.Property {
	return openapi.Property{
		Fid:   in.Field.ID.String(),
		Name:  in.Field.Name,
		Type:  openapi.PropertyType(in.Field.Type.String()),
		Sort:  in.Field.Sort,
		Value: in.Value.OrZero(),
	}
}

func serialisePropertyTable(in library.PropertyTable) openapi.PropertyList {
	if len(in.Properties) == 0 {
		return openapi.PropertyList{}
	}

	propertyFieldMap := lo.KeyBy(in.Properties, func(f *library.Property) xid.ID {
		return f.Field.ID
	})

	pl := dt.Map(in.Schema.Fields, func(f *library.PropertySchemaField) openapi.Property {
		p, ok := propertyFieldMap[f.ID]
		if !ok {
			return openapi.Property{
				Fid:   f.ID.String(),
				Name:  f.Name,
				Type:  openapi.PropertyType(f.Type.String()),
				Sort:  f.Sort,
				Value: "",
			}
		}

		return serialiseProperty(p)
	})

	return pl
}

func serialisePropertyTableOpt(in opt.Optional[library.PropertyTable]) openapi.PropertyList {
	pt, ok := in.Get()
	if !ok {
		return openapi.PropertyList{}
	}
	return serialisePropertyTable(pt)
}

func serialisePropertySchema(in *library.PropertySchemaField) openapi.PropertySchema {
	return openapi.PropertySchema{
		Fid:  in.ID.String(),
		Name: in.Name,
		Type: openapi.PropertyType(in.Type.String()),
		Sort: in.Sort,
	}
}

func serialisePropertySchemas(in library.PropertySchema) openapi.PropertySchemaList {
	return dt.Map(in.Fields, serialisePropertySchema)
}

func serialisePropertySchemaList(in library.PropertySchema) openapi.PropertySchemaList {
	if len(in.Fields) == 0 {
		return openapi.PropertySchemaList{}
	}

	return dt.Map(in.Fields, serialisePropertySchema)
}

func deserialisePropertySchemaMutation(in openapi.PropertySchemaMutableProps) (*node_properties.SchemaFieldMutation, error) {
	t, err := library.NewPropertyType(string(in.Type))
	if err != nil {
		return nil, fault.Wrap(err, ftag.With(ftag.InvalidArgument))
	}

	return &node_properties.SchemaFieldMutation{
		ID:   opt.Map(opt.NewPtr(in.Fid), deserialiseID),
		Name: in.Name,
		Type: t,
		Sort: in.Sort,
	}, nil
}

func deserialisePropertyMutationList(in openapi.PropertyMutationList) (library.PropertyMutationList, error) {
	return dt.MapErr(in, deserialisePropertyMutation)
}

func deserialisePropertyMutation(in openapi.PropertyMutation) (*library.PropertyMutation, error) {
	t, err := opt.MapErr(opt.NewPtr(in.Type), func(s openapi.PropertyType) (library.PropertyType, error) {
		return library.NewPropertyType(string(s))
	})
	if err != nil {
		return nil, fault.Wrap(err, ftag.With(ftag.InvalidArgument))
	}

	return &library.PropertyMutation{
		ID:    opt.Map(opt.NewPtr(in.Fid), deserialiseID),
		Name:  in.Name,
		Value: in.Value,
		Type:  t,
		Sort:  opt.NewPtr(in.Sort),
	}, nil
}
