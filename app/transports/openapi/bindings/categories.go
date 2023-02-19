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

func (c Categories) CategoriesList(ctx context.Context, request openapi.CategoriesListRequestObject) (openapi.CategoriesListResponseObject, error) {
	cats, err := c.category_repo.GetCategories(ctx, false)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CategoriesList200JSONResponse{
		CategoriesListSuccessJSONResponse: openapi.CategoriesListSuccessJSONResponse{
			Categories: dt.Map(cats, serialiseCategory),
		},
	}, nil
}
