package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/cluster_traversal"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	cluster_svc "github.com/Southclaws/storyden/app/services/cluster"
	"github.com/Southclaws/storyden/app/services/clustertree"
	"github.com/Southclaws/storyden/app/services/item_tree"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Clusters struct {
	cs    cluster_svc.Manager
	ctree clustertree.Graph
	ctr   cluster_traversal.Repository
	itree item_tree.Graph
}

func NewClusters(
	cs cluster_svc.Manager,
	ctree clustertree.Graph,
	ctr cluster_traversal.Repository,
	itree item_tree.Graph,
) Clusters {
	return Clusters{
		cs:    cs,
		ctree: ctree,
		ctr:   ctr,
		itree: itree,
	}
}

func (c *Clusters) ClusterCreate(ctx context.Context, request openapi.ClusterCreateRequestObject) (openapi.ClusterCreateResponseObject, error) {
	session, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	clus, err := c.cs.Create(ctx,
		session,
		request.Body.Name,
		request.Body.Slug,
		request.Body.Description,
		cluster_svc.Partial{
			Content:    opt.NewPtr(request.Body.Content),
			Properties: opt.NewPtr(request.Body.Properties),
			URL:        opt.NewPtr(request.Body.Url),
			AssetsAdd:  opt.NewPtrMap(request.Body.AssetIds, deserialiseAssetIDs),
			Parent:     opt.NewPtrMap(request.Body.Parent, deserialiseClusterSlug),
		},
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ClusterCreate200JSONResponse{
		ClusterCreateOKJSONResponse: openapi.ClusterCreateOKJSONResponse(serialiseCluster(clus)),
	}, nil
}

func (c *Clusters) ClusterList(ctx context.Context, request openapi.ClusterListRequestObject) (openapi.ClusterListResponseObject, error) {
	var cs []*datagraph.Cluster
	var err error

	opts := []cluster_traversal.Filter{}

	if v := request.Params.Author; v != nil {
		opts = append(opts, cluster_traversal.WithOwner(*v))
	}

	if id := request.Params.ClusterId; id != nil {
		cid, err := xid.FromString(*id)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}

		cs, err = c.ctr.Subtree(ctx, datagraph.ClusterID(cid), opts...)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	} else {
		cs, err = c.ctr.Root(ctx, opts...)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	return openapi.ClusterList200JSONResponse{
		ClusterListOKJSONResponse: openapi.ClusterListOKJSONResponse{
			Clusters: dt.Map(cs, serialiseCluster),
		},
	}, nil
}

func (c *Clusters) ClusterGet(ctx context.Context, request openapi.ClusterGetRequestObject) (openapi.ClusterGetResponseObject, error) {
	clus, err := c.cs.Get(ctx, datagraph.ClusterSlug(request.ClusterSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ClusterGet200JSONResponse{
		ClusterGetOKJSONResponse: openapi.ClusterGetOKJSONResponse(serialiseClusterWithItems(clus)),
	}, nil
}

func (c *Clusters) ClusterUpdate(ctx context.Context, request openapi.ClusterUpdateRequestObject) (openapi.ClusterUpdateResponseObject, error) {
	clus, err := c.cs.Update(ctx, datagraph.ClusterSlug(request.ClusterSlug), cluster_svc.Partial{
		Name:        opt.NewPtr(request.Body.Name),
		Slug:        opt.NewPtr(request.Body.Slug),
		AssetsAdd:   opt.NewPtrMap(request.Body.AssetIds, deserialiseAssetIDs),
		URL:         opt.NewPtr(request.Body.Url),
		Description: opt.NewPtr(request.Body.Description),
		Content:     opt.NewPtr(request.Body.Content),
		Parent:      opt.NewPtrMap(request.Body.Parent, deserialiseClusterSlug),
		Properties:  opt.NewPtr(request.Body.Properties),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ClusterUpdate200JSONResponse{
		ClusterUpdateOKJSONResponse: openapi.ClusterUpdateOKJSONResponse(serialiseCluster(clus)),
	}, nil
}

func (c *Clusters) ClusterDelete(ctx context.Context, request openapi.ClusterDeleteRequestObject) (openapi.ClusterDeleteResponseObject, error) {
	destinationCluster, err := c.cs.Delete(ctx, datagraph.ClusterSlug(request.ClusterSlug), cluster_svc.DeleteOptions{
		MoveTo:   opt.NewPtr((*datagraph.ClusterSlug)(request.Params.TargetCluster)),
		Clusters: opt.NewPtr(request.Params.MoveChildClusters).OrZero(),
		Items:    opt.NewPtr(request.Params.MoveChildItems).OrZero(),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ClusterDelete200JSONResponse{
		ClusterDeleteOKJSONResponse: openapi.ClusterDeleteOKJSONResponse{
			Destination: opt.Map(opt.NewPtr(destinationCluster), func(in datagraph.Cluster) openapi.Cluster {
				return serialiseCluster(&in)
			}).Ptr(),
		},
	}, nil
}

func (c *Clusters) ClusterAddAsset(ctx context.Context, request openapi.ClusterAddAssetRequestObject) (openapi.ClusterAddAssetResponseObject, error) {
	clus, err := c.cs.Update(ctx, datagraph.ClusterSlug(request.ClusterSlug), cluster_svc.Partial{
		AssetsAdd: opt.New([]asset.AssetID{asset.AssetID(request.Id)}),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ClusterAddAsset200JSONResponse{
		ClusterUpdateOKJSONResponse: openapi.ClusterUpdateOKJSONResponse(serialiseCluster(clus)),
	}, nil
}

func (c *Clusters) ClusterRemoveAsset(ctx context.Context, request openapi.ClusterRemoveAssetRequestObject) (openapi.ClusterRemoveAssetResponseObject, error) {
	clus, err := c.cs.Update(ctx, datagraph.ClusterSlug(request.ClusterSlug), cluster_svc.Partial{
		AssetsRemove: opt.New([]asset.AssetID{asset.AssetID(request.Id)}),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ClusterRemoveAsset200JSONResponse{
		ClusterUpdateOKJSONResponse: openapi.ClusterUpdateOKJSONResponse(serialiseCluster(clus)),
	}, nil
}

func (c *Clusters) ClusterAddCluster(ctx context.Context, request openapi.ClusterAddClusterRequestObject) (openapi.ClusterAddClusterResponseObject, error) {
	clus, err := c.ctree.Move(ctx, datagraph.ClusterSlug(request.ClusterSlugChild), datagraph.ClusterSlug(request.ClusterSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ClusterAddCluster200JSONResponse{
		ClusterAddChildOKJSONResponse: openapi.ClusterAddChildOKJSONResponse(serialiseCluster(clus)),
	}, nil
}

func (c *Clusters) ClusterRemoveCluster(ctx context.Context, request openapi.ClusterRemoveClusterRequestObject) (openapi.ClusterRemoveClusterResponseObject, error) {
	clus, err := c.ctree.Sever(ctx, datagraph.ClusterSlug(request.ClusterSlugChild), datagraph.ClusterSlug(request.ClusterSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ClusterRemoveCluster200JSONResponse{
		ClusterRemoveChildOKJSONResponse: openapi.ClusterRemoveChildOKJSONResponse(serialiseCluster(clus)),
	}, nil
}

func (c *Clusters) ClusterAddItem(ctx context.Context, request openapi.ClusterAddItemRequestObject) (openapi.ClusterAddItemResponseObject, error) {
	_, err := c.itree.Link(ctx, datagraph.ItemSlug(request.ItemSlug), datagraph.ClusterSlug(request.ClusterSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	clus, err := c.cs.Get(ctx, datagraph.ClusterSlug(request.ClusterSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ClusterAddItem200JSONResponse{
		ClusterAddItemOKJSONResponse: openapi.ClusterAddItemOKJSONResponse(serialiseClusterWithItems(clus)),
	}, nil
}

func (c *Clusters) ClusterRemoveItem(ctx context.Context, request openapi.ClusterRemoveItemRequestObject) (openapi.ClusterRemoveItemResponseObject, error) {
	_, err := c.itree.Sever(ctx, datagraph.ItemSlug(request.ItemSlug), datagraph.ClusterSlug(request.ClusterSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	clus, err := c.cs.Get(ctx, datagraph.ClusterSlug(request.ClusterSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ClusterRemoveItem200JSONResponse{
		ClusterRemoveItemOKJSONResponse: openapi.ClusterRemoveItemOKJSONResponse(serialiseClusterWithItems(clus)),
	}, nil
}

func serialiseCluster(in *datagraph.Cluster) openapi.Cluster {
	return openapi.Cluster{
		Id:          in.ID.String(),
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		Name:        in.Name,
		Slug:        in.Slug,
		Assets:      dt.Map(in.Assets, serialiseAssetReference),
		Link:        opt.Map(in.Links.Latest(), serialiseLink).Ptr(),
		Description: in.Description,
		Content:     in.Content.Ptr(),
		Owner:       serialiseProfileReference(in.Owner),
		Parent:      opt.Map(in.Parent, serialiseCluster).Ptr(),
		Properties:  in.Properties,
	}
}

func serialiseClusterWithItems(in *datagraph.Cluster) openapi.ClusterWithItems {
	return openapi.ClusterWithItems{
		Id:          in.ID.String(),
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		Name:        in.Name,
		Slug:        in.Slug,
		Assets:      dt.Map(in.Assets, serialiseAssetReference),
		Link:        opt.Map(in.Links.Latest(), serialiseLink).Ptr(),
		Description: in.Description,
		Content:     in.Content.Ptr(),
		Owner:       serialiseProfileReference(in.Owner),
		Parent:      opt.Map(in.Parent, serialiseCluster).Ptr(),
		Properties:  in.Properties,
		Clusters:    dt.Map(in.Clusters, serialiseCluster),
		Items:       dt.Map(in.Items, serialiseItem),
	}
}

func deserialiseClusterSlug(in string) datagraph.ClusterSlug {
	return datagraph.ClusterSlug(in)
}
