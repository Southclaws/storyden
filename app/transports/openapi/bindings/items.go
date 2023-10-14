package bindings

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	collection_svc "github.com/Southclaws/storyden/app/services/collection"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Items struct {
	collection_repo collection.Repository
	collection_svc  collection_svc.Service
}

func NewItems(
	collection_repo collection.Repository,
	collection_svc collection_svc.Service,
) Items {
	return Items{
		collection_repo: collection_repo,
		collection_svc:  collection_svc,
	}
}

func (i *Items) ItemList(ctx context.Context, request openapi.ItemListRequestObject) (openapi.ItemListResponseObject, error) {
	return nil, nil
}

func (i *Items) ItemCreate(ctx context.Context, request openapi.ItemCreateRequestObject) (openapi.ItemCreateResponseObject, error) {
	return nil, nil
}

func (i *Items) ItemGet(ctx context.Context, request openapi.ItemGetRequestObject) (openapi.ItemGetResponseObject, error) {
	return nil, nil
}

func (i *Items) ItemUpdate(ctx context.Context, request openapi.ItemUpdateRequestObject) (openapi.ItemUpdateResponseObject, error) {
	return nil, nil
}

func serialiseItem(in *datagraph.Item) openapi.Item {
	return openapi.Item{
		Id:          in.ID.String(),
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		Name:        in.Name,
		Slug:        in.Slug,
		ImageUrl:    in.ImageURL.Ptr(),
		Description: in.Description,
		Owner:       serialiseProfileReference(in.Owner),
		Properties:  in.Properties,
	}
}
