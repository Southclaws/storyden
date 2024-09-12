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
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_traversal"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/library/node_mutate"
	"github.com/Southclaws/storyden/app/services/library/node_read"
	"github.com/Southclaws/storyden/app/services/library/node_visibility"
	"github.com/Southclaws/storyden/app/services/library/nodetree"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Nodes struct {
	accountQuery account_querier.Querier
	nodeMutator  node_mutate.Manager
	nodeReader   *node_read.HydratedQuerier
	nv           *node_visibility.Controller
	ntree        nodetree.Graph
	ntr          node_traversal.Repository
}

func NewNodes(
	accountQuery account_querier.Querier,
	nodeMutator node_mutate.Manager,
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

	node, err := c.nodeMutator.Create(ctx,
		session,
		request.Body.Name,
		node_mutate.Partial{
			Slug:         opt.NewPtr(request.Body.Slug),
			Content:      richContent,
			Metadata:     opt.NewPtr((*map[string]any)(request.Body.Meta)),
			URL:          url,
			AssetsAdd:    opt.NewPtrMap(request.Body.AssetIds, deserialiseAssetIDs),
			AssetSources: opt.NewPtrMap(request.Body.AssetSources, deserialiseAssetSources),
			Parent:       opt.NewPtrMap(request.Body.Parent, deserialiseNodeSlug),
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

	acc, err := opt.MapErr(session.GetOptAccountID(ctx), func(aid account.AccountID) (*account.Account, error) {
		return c.accountQuery.GetByID(ctx, aid)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var cs []*library.Node

	opts := []node_traversal.Filter{}

	author := opt.NewPtr(request.Params.Author)
	if v, ok := author.Get(); ok {
		opts = append(opts, node_traversal.WithOwner(v))
	}

	if d, ok := depth.Get(); ok {
		opts = append(opts, node_traversal.WithDepth(uint(d)))
	}

	visibilities, err := opt.MapErr(opt.NewPtr(request.Params.Visibility), deserialiseVisibilityList)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if a, ok := acc.Get(); ok {
		// NOTE: We do not want to allow anyone to request ANY node that is
		// not published, but we also want to allow admins to request nodes
		// that are in review. So we need to check if the requesting account is
		// an admin and if they are not, automatically add a WithOwner filter.

		if v, ok := visibilities.Get(); ok {
			opts = append(opts, node_traversal.WithVisibility(v...))

			authorFilter, filteringByAuthor := author.Get()

			if !filteringByAuthor || authorFilter == a.Handle {
				if lo.Contains(v, visibility.VisibilityDraft) {
					// If the result is to contain drafts, only show the account's.
					opts = append(opts, node_traversal.WithOwner(a.Handle))
				} else if lo.Contains(v, visibility.VisibilityReview) {
					// If the result is to contain nodes that are in-review, then
					// we need to check if the requesting account is an admin first.
					if !a.Admin {
						opts = append(opts, node_traversal.WithOwner(a.Handle))
					}
				}
			}
		}
	} else {
		// When the request is not made by an authenticated account, we do not
		// permit any visibility other than "published".

		opts = append(opts, node_traversal.WithVisibility(visibility.VisibilityPublished))
	}

	nid, err := opt.MapErr(opt.NewPtr(request.Params.NodeId), library.NodeIDFromString)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	cs, err = c.ntr.Subtree(ctx, nid, opts...)
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
	node, err := c.nodeReader.GetBySlug(ctx, library.NodeSlug(request.NodeSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeGet200JSONResponse{
		NodeGetOKJSONResponse: openapi.NodeGetOKJSONResponse(serialiseNodeWithItems(node)),
	}, nil
}

func (c *Nodes) NodeUpdate(ctx context.Context, request openapi.NodeUpdateRequestObject) (openapi.NodeUpdateResponseObject, error) {
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

	node, err := c.nodeMutator.Update(ctx, library.NodeSlug(request.NodeSlug), node_mutate.Partial{
		Name:         opt.NewPtr(request.Body.Name),
		Slug:         opt.NewPtr(request.Body.Slug),
		AssetsAdd:    opt.NewPtrMap(request.Body.AssetIds, deserialiseAssetIDs),
		AssetSources: opt.NewPtrMap(request.Body.AssetSources, deserialiseAssetSources),
		URL:          url,
		Content:      richContent,
		Parent:       opt.NewPtrMap(request.Body.Parent, deserialiseNodeSlug),
		Metadata:     opt.NewPtr((*map[string]any)(request.Body.Meta)),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeUpdate200JSONResponse{
		NodeUpdateOKJSONResponse: openapi.NodeUpdateOKJSONResponse(serialiseNode(node)),
	}, nil
}

func (c *Nodes) NodeUpdateVisibility(ctx context.Context, request openapi.NodeUpdateVisibilityRequestObject) (openapi.NodeUpdateVisibilityResponseObject, error) {
	v, err := visibility.NewVisibility(string(request.Body.Visibility))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	node, err := c.nv.ChangeVisibility(ctx, library.NodeSlug(request.NodeSlug), v)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeUpdateVisibility200JSONResponse{
		NodeUpdateOKJSONResponse: openapi.NodeUpdateOKJSONResponse(serialiseNode(node)),
	}, nil
}

func (c *Nodes) NodeDelete(ctx context.Context, request openapi.NodeDeleteRequestObject) (openapi.NodeDeleteResponseObject, error) {
	destinationNode, err := c.nodeMutator.Delete(ctx, library.NodeSlug(request.NodeSlug), node_mutate.DeleteOptions{
		NewParent: opt.NewPtr((*library.NodeSlug)(request.Params.TargetNode)),
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

func (c *Nodes) NodeAddAsset(ctx context.Context, request openapi.NodeAddAssetRequestObject) (openapi.NodeAddAssetResponseObject, error) {
	id := openapi.ParseID(request.AssetId)

	contentFillCmd, err := getContentFillRuleCommand(request.Params.ContentFillRule, request.Params.NodeContentFillTarget)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	node, err := c.nodeMutator.Update(ctx, library.NodeSlug(request.NodeSlug), node_mutate.Partial{
		AssetsAdd:   opt.New([]asset.AssetID{id}),
		ContentFill: contentFillCmd,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeAddAsset200JSONResponse{
		NodeUpdateOKJSONResponse: openapi.NodeUpdateOKJSONResponse(serialiseNode(node)),
	}, nil
}

func (c *Nodes) NodeRemoveAsset(ctx context.Context, request openapi.NodeRemoveAssetRequestObject) (openapi.NodeRemoveAssetResponseObject, error) {
	id := openapi.ParseID(request.AssetId)

	node, err := c.nodeMutator.Update(ctx, library.NodeSlug(request.NodeSlug), node_mutate.Partial{
		AssetsRemove: opt.New([]asset.AssetID{id}),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeRemoveAsset200JSONResponse{
		NodeUpdateOKJSONResponse: openapi.NodeUpdateOKJSONResponse(serialiseNode(node)),
	}, nil
}

func (c *Nodes) NodeAddNode(ctx context.Context, request openapi.NodeAddNodeRequestObject) (openapi.NodeAddNodeResponseObject, error) {
	node, err := c.ntree.Move(ctx, library.NodeSlug(request.NodeSlugChild), library.NodeSlug(request.NodeSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeAddNode200JSONResponse{
		NodeAddChildOKJSONResponse: openapi.NodeAddChildOKJSONResponse(serialiseNode(node)),
	}, nil
}

func (c *Nodes) NodeRemoveNode(ctx context.Context, request openapi.NodeRemoveNodeRequestObject) (openapi.NodeRemoveNodeResponseObject, error) {
	node, err := c.ntree.Sever(ctx, library.NodeSlug(request.NodeSlugChild), library.NodeSlug(request.NodeSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeRemoveNode200JSONResponse{
		NodeRemoveChildOKJSONResponse: openapi.NodeRemoveChildOKJSONResponse(serialiseNode(node)),
	}, nil
}

func serialiseNode(in *library.Node) openapi.Node {
	return openapi.Node{
		Id:          in.ID.String(),
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		Name:        in.Name,
		Slug:        in.Slug,
		Assets:      dt.Map(in.Assets, serialiseAssetPtr),
		Link:        opt.Map(in.WebLink, serialiseLinkRef).Ptr(),
		Description: in.GetDesc(),
		Content:     opt.Map(in.Content, serialiseContentHTML).Ptr(),
		Owner:       serialiseProfileReference(in.Owner),
		Parent: opt.PtrMap(in.Parent, func(in library.Node) openapi.Node {
			return serialiseNode(&in)
		}),
		Visibility: serialiseVisibility(in.Visibility),
		Meta:       in.Metadata,
	}
}

func serialiseNodeWithItems(in *library.Node) openapi.NodeWithChildren {
	rs := opt.Map(in.RelevanceScore, func(v float64) float32 { return float32(v) })
	return openapi.NodeWithChildren{
		Id:          in.ID.String(),
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		Name:        in.Name,
		Slug:        in.Slug,
		Assets:      dt.Map(in.Assets, serialiseAssetPtr),
		Link:        opt.Map(in.WebLink, serialiseLinkRef).Ptr(),
		Description: in.GetDesc(),
		Content:     opt.Map(in.Content, serialiseContentHTML).Ptr(),
		Owner:       serialiseProfileReference(in.Owner),
		Parent: opt.PtrMap(in.Parent, func(in library.Node) openapi.Node {
			return serialiseNode(&in)
		}),
		Visibility:     serialiseVisibility(in.Visibility),
		RelevanceScore: rs.Ptr(),
		Meta:           in.Metadata,
		Children:       dt.Map(in.Nodes, serialiseNodeWithItems),
	}
}

func deserialiseNodeSlug(in string) library.NodeSlug {
	return library.NodeSlug(in)
}

func deserialiseAssetSources(in openapi.AssetSourceList) []string {
	return dt.Map(in, deserialiseAssetSourceURL)
}

func deserialiseAssetSourceURL(in openapi.AssetSourceURL) string {
	return string(in)
}
