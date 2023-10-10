package bindings

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/collection"
	collection_svc "github.com/Southclaws/storyden/app/services/collection"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Clusters struct {
	collection_repo collection.Repository
	collection_svc  collection_svc.Service
}

func NewClusters(
	collection_repo collection.Repository,
	collection_svc collection_svc.Service,
) Clusters {
	return Clusters{
		collection_repo: collection_repo,
		collection_svc:  collection_svc,
	}
}

func (i *Clusters) ClusterList(ctx context.Context, request openapi.ClusterListRequestObject) (openapi.ClusterListResponseObject, error) {
	return nil, nil
}

func (i *Clusters) ClusterCreate(ctx context.Context, request openapi.ClusterCreateRequestObject) (openapi.ClusterCreateResponseObject, error) {
	return nil, nil
}

func (i *Clusters) ClusterGet(ctx context.Context, request openapi.ClusterGetRequestObject) (openapi.ClusterGetResponseObject, error) {
	return nil, nil
}

func (i *Clusters) ClusterUpdate(ctx context.Context, request openapi.ClusterUpdateRequestObject) (openapi.ClusterUpdateResponseObject, error) {
	return nil, nil
}

func (i *Clusters) ClusterRemoveCluster(ctx context.Context, request openapi.ClusterRemoveClusterRequestObject) (openapi.ClusterRemoveClusterResponseObject, error) {
	return nil, nil
}

func (i *Clusters) ClusterAddCluster(ctx context.Context, request openapi.ClusterAddClusterRequestObject) (openapi.ClusterAddClusterResponseObject, error) {
	return nil, nil
}

func (i *Clusters) ClusterRemoveItem(ctx context.Context, request openapi.ClusterRemoveItemRequestObject) (openapi.ClusterRemoveItemResponseObject, error) {
	return nil, nil
}

func (i *Clusters) ClusterAddItem(ctx context.Context, request openapi.ClusterAddItemRequestObject) (openapi.ClusterAddItemResponseObject, error) {
	return nil, nil
}
