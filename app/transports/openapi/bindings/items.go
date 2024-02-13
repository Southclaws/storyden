package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/item_search"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/item_crud"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Items struct {
	im item_crud.Manager
	is item_search.Search
}

func NewItems(
	im item_crud.Manager,
	is item_search.Search,
) Items {
	return Items{
		im: im,
		is: is,
	}
}

func (i *Items) ItemCreate(ctx context.Context, request openapi.ItemCreateRequestObject) (openapi.ItemCreateResponseObject, error) {
	session, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	itm, err := i.im.Create(ctx,
		session,
		request.Body.Name,
		request.Body.Slug,
		request.Body.Description,
		item_crud.Partial{
			Content:    opt.NewPtr(request.Body.Content),
			Properties: opt.NewPtr(request.Body.Properties),
			URL:        opt.NewPtr(request.Body.Url),
			AssetsAdd:  opt.NewPtrMap(request.Body.AssetIds, deserialiseAssetIDs),
		},
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ItemCreate200JSONResponse{
		ItemCreateOKJSONResponse: openapi.ItemCreateOKJSONResponse(serialiseItem(itm)),
	}, nil
}

func (i *Items) ItemList(ctx context.Context, request openapi.ItemListRequestObject) (openapi.ItemListResponseObject, error) {
	opts := []item_search.Option{}

	if v := request.Params.Q; v != nil {
		opts = append(opts, item_search.WithNameContains(*v))
	}

	visibility, err := opt.MapErr(opt.NewPtr(request.Params.Visibility), deserialiseVisibilityList)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if v, ok := visibility.Get(); ok {
		opts = append(opts, item_search.WithVisibility(v))
	}

	items, err := i.is.Search(ctx, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ItemList200JSONResponse{
		ItemListOKJSONResponse: openapi.ItemListOKJSONResponse{
			Items: dt.Map(items, serialiseItemWithParents),
		},
	}, nil
}

func (i *Items) ItemGet(ctx context.Context, request openapi.ItemGetRequestObject) (openapi.ItemGetResponseObject, error) {
	item, err := i.im.Get(ctx, datagraph.ItemSlug(request.ItemSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ItemGet200JSONResponse{
		ItemGetOKJSONResponse: openapi.ItemGetOKJSONResponse(serialiseItemWithParents(item)),
	}, nil
}

func (i *Items) ItemUpdate(ctx context.Context, request openapi.ItemUpdateRequestObject) (openapi.ItemUpdateResponseObject, error) {
	item, err := i.im.Update(ctx, datagraph.ItemSlug(request.ItemSlug), item_crud.Partial{
		Name:        opt.NewPtr(request.Body.Name),
		Slug:        opt.NewPtr(request.Body.Slug),
		AssetsAdd:   opt.NewPtrMap(request.Body.AssetIds, deserialiseAssetIDs),
		URL:         opt.NewPtr(request.Body.Url),
		Description: opt.NewPtr(request.Body.Description),
		Content:     opt.NewPtr(request.Body.Content),
		Properties:  opt.NewPtr(request.Body.Properties),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ItemUpdate200JSONResponse{
		ItemUpdateOKJSONResponse: openapi.ItemUpdateOKJSONResponse(serialiseItem(item)),
	}, nil
}

func (i *Items) ItemUpdateVisibility(ctx context.Context, request openapi.ItemUpdateVisibilityRequestObject) (openapi.ItemUpdateVisibilityResponseObject, error) {
	v, err := post.NewVisibility(string(request.Body.Visibility))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	item, err := i.im.Update(ctx, datagraph.ItemSlug(request.ItemSlug), item_crud.Partial{
		Visibility: opt.New(v),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ItemUpdateVisibility200JSONResponse{
		ItemUpdateOKJSONResponse: openapi.ItemUpdateOKJSONResponse(serialiseItem(item)),
	}, nil
}

func (i *Items) ItemDelete(ctx context.Context, request openapi.ItemDeleteRequestObject) (openapi.ItemDeleteResponseObject, error) {
	_, err := i.im.Delete(ctx, datagraph.ItemSlug(request.ItemSlug))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ItemDelete200Response{}, nil
}

func (c *Items) ItemAddAsset(ctx context.Context, request openapi.ItemAddAssetRequestObject) (openapi.ItemAddAssetResponseObject, error) {
	id := openapi.ParseID(request.AssetId)

	item, err := c.im.Update(ctx, datagraph.ItemSlug(request.ItemSlug), item_crud.Partial{
		AssetsAdd: opt.New([]asset.AssetID{id}),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ItemAddAsset200JSONResponse{
		ItemUpdateOKJSONResponse: openapi.ItemUpdateOKJSONResponse(serialiseItem(item)),
	}, nil
}

func (c *Items) ItemRemoveAsset(ctx context.Context, request openapi.ItemRemoveAssetRequestObject) (openapi.ItemRemoveAssetResponseObject, error) {
	id := openapi.ParseID(request.AssetId)

	item, err := c.im.Update(ctx, datagraph.ItemSlug(request.ItemSlug), item_crud.Partial{
		AssetsRemove: opt.New([]asset.AssetID{id}),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ItemRemoveAsset200JSONResponse{
		ItemUpdateOKJSONResponse: openapi.ItemUpdateOKJSONResponse(serialiseItem(item)),
	}, nil
}

func serialiseItem(in *datagraph.Item) openapi.Item {
	return openapi.Item{
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
		Properties:  in.Properties,
	}
}

func serialiseItemWithParents(in *datagraph.Item) openapi.ItemWithParents {
	clusters := dt.Map(in.In, serialiseCluster)
	return openapi.ItemWithParents{
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
		Properties:  in.Properties,
		Clusters:    clusters,
	}
}
