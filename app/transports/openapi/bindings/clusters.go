package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	cluster_repo "github.com/Southclaws/storyden/app/resources/cluster"
	"github.com/Southclaws/storyden/app/resources/cluster_traversal"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/authentication"
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
	session, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := []cluster_repo.Option{}

	if v := request.Body.Content; v != nil {
		opts = append(opts, cluster_repo.WithContent(*v))
	}
	if v := request.Body.Properties; v != nil {
		opts = append(opts, cluster_repo.WithProperties(*v))
	}
	if v := request.Body.ImageUrl; v != nil {
		opts = append(opts, cluster_repo.WithImageURL(*v))
	}
	if v := request.Body.Url; v != nil {
		opts = append(opts, cluster_repo.WithURL(*v))
	}

	clus, err := c.cs.Create(ctx,
		session,
		request.Body.Name,
		request.Body.Slug,
		request.Body.Description,
		opts...,
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
		ImageURL:    opt.NewPtr(request.Body.ImageUrl),
		URL:         opt.NewPtr(request.Body.Url),
		Description: opt.NewPtr(request.Body.Description),
		Content:     opt.NewPtr(request.Body.Content),
		Properties:  opt.NewPtr(request.Body.Properties),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ClusterUpdate200JSONResponse{
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
	item, err := c.itree.Link(ctx, datagraph.ItemSlug(request.ItemSlug), datagraph.ClusterSlug(request.ClusterSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ClusterAddItem200JSONResponse{
		ClusterAddItemOKJSONResponse: openapi.ClusterAddItemOKJSONResponse(serialiseItem(item)),
	}, nil
}

func (c *Clusters) ClusterRemoveItem(ctx context.Context, request openapi.ClusterRemoveItemRequestObject) (openapi.ClusterRemoveItemResponseObject, error) {
	item, err := c.itree.Sever(ctx, datagraph.ItemSlug(request.ItemSlug), datagraph.ClusterSlug(request.ClusterSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ClusterRemoveItem200JSONResponse{
		ClusterRemoveItemOKJSONResponse: openapi.ClusterRemoveItemOKJSONResponse(serialiseItem(item)),
	}, nil
}

func serialiseCluster(in *datagraph.Cluster) openapi.Cluster {
	return openapi.Cluster{
		Id:          in.ID.String(),
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		Name:        in.Name,
		Slug:        in.Slug,
		ImageUrl:    in.ImageURL.Ptr(),
		Url:         in.URL.Ptr(),
		Description: in.Description,
		Content:     in.Content.Ptr(),
		Owner:       serialiseProfileReference(in.Owner),
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
		ImageUrl:    in.ImageURL.Ptr(),
		Url:         in.URL.Ptr(),
		Description: in.Description,
		Content:     in.Content.Ptr(),
		Owner:       serialiseProfileReference(in.Owner),
		Properties:  in.Properties,
		Clusters:    dt.Map(in.Clusters, serialiseCluster),
		Items:       dt.Map(in.Items, serialiseItem),
	}
}
