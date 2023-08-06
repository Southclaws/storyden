package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Categories struct {
	category_repo category.Repository
}

func NewCategories(
	category_repo category.Repository,
) Categories {
	return Categories{category_repo}
}

func (c Categories) CategoryList(ctx context.Context, request openapi.CategoryListRequestObject) (openapi.CategoryListResponseObject, error) {
	cats, err := c.category_repo.GetCategories(ctx, false)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CategoryList200JSONResponse{
		CategoryListOKJSONResponse: openapi.CategoryListOKJSONResponse{
			Categories: dt.Map(cats, serialiseCategory),
		},
	}, nil
}

func (c Categories) CategoryUpdateOrder(ctx context.Context, request openapi.CategoryUpdateOrderRequestObject) (openapi.CategoryUpdateOrderResponseObject, error) {
	ids := dt.Map(*request.Body, func(in openapi.Identifier) category.CategoryID {
		return category.CategoryID(openapi.ParseID(in))
	})

	cats, err := c.category_repo.Reorder(ctx, ids)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CategoryUpdateOrder200JSONResponse{
		CategoryListOKJSONResponse: openapi.CategoryListOKJSONResponse{
			Categories: dt.Map(cats, serialiseCategory),
		},
	}, nil
}
