package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	item_repo "github.com/Southclaws/storyden/app/resources/item"
	"github.com/Southclaws/storyden/app/services/authentication"
	item_svc "github.com/Southclaws/storyden/app/services/item"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Items struct {
	im item_svc.Manager
}

func NewItems(
	im item_svc.Manager,
) Items {
	return Items{
		im: im,
	}
}

func (i *Items) ItemCreate(ctx context.Context, request openapi.ItemCreateRequestObject) (openapi.ItemCreateResponseObject, error) {
	session, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := []item_repo.Option{}

	if v := request.Body.Properties; v != nil {
		opts = append(opts, item_repo.WithProperties(*v))
	}

	itm, err := i.im.Create(ctx,
		session,
		request.Body.Name,
		request.Body.Slug,
		request.Body.Description,
		opts...,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ItemCreate200JSONResponse{
		ItemCreateOKJSONResponse: openapi.ItemCreateOKJSONResponse(serialiseItem(itm)),
	}, nil
}

func (i *Items) ItemList(ctx context.Context, request openapi.ItemListRequestObject) (openapi.ItemListResponseObject, error) {
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
