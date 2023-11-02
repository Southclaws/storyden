package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	item_repo "github.com/Southclaws/storyden/app/resources/item"
	"github.com/Southclaws/storyden/app/resources/item_search"
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

	opts := []item_repo.Option{}

	if v := request.Body.Content; v != nil {
		opts = append(opts, item_repo.WithContent(*v))
	}
	if v := request.Body.Properties; v != nil {
		opts = append(opts, item_repo.WithProperties(*v))
	}
	if v := request.Body.ImageUrl; v != nil {
		opts = append(opts, item_repo.WithImageURL(*v))
	}
	if v := request.Body.Url; v != nil {
		opts = append(opts, item_repo.WithURL(*v))
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
	opts := []item_search.Option{}

	if v := request.Params.Q; v != nil {
		opts = append(opts, item_search.WithNameContains(*v))
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
		ImageURL:    opt.NewPtr(request.Body.ImageUrl),
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

func serialiseItem(in *datagraph.Item) openapi.Item {
	return openapi.Item{
		Id:          in.ID.String(),
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		Name:        in.Name,
		Slug:        in.Slug,
		ImageUrl:    in.ImageURL.Ptr(),
		Link:        opt.Map(in.Link, serialiseLink).Ptr(),
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
		ImageUrl:    in.ImageURL.Ptr(),
		Link:        opt.Map(in.Link, serialiseLink).Ptr(),
		Description: in.Description,
		Content:     in.Content.Ptr(),
		Owner:       serialiseProfileReference(in.Owner),
		Properties:  in.Properties,
		Clusters:    clusters,
	}
}
