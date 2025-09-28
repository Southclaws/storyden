package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/post/category"
	category_svc "github.com/Southclaws/storyden/app/services/category"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/deletable"
)

type Categories struct {
	category_repo *category.Repository
	category_svc  category_svc.Service
}

func NewCategories(
	category_repo *category.Repository,
	category_svc category_svc.Service,
) Categories {
	return Categories{category_repo, category_svc}
}

func (c Categories) CategoryCreate(ctx context.Context, request openapi.CategoryCreateRequestObject) (openapi.CategoryCreateResponseObject, error) {
	cat, err := c.category_svc.Create(ctx, request.Body.Name, request.Body.Description, request.Body.Colour, request.Body.Admin)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CategoryCreate200JSONResponse{
		CategoryCreateOKJSONResponse: openapi.CategoryCreateOKJSONResponse(serialiseCategory(cat)),
	}, nil
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

func (c Categories) CategoryGet(ctx context.Context, request openapi.CategoryGetRequestObject) (openapi.CategoryGetResponseObject, error) {
	cat, err := c.category_repo.Get(ctx, request.CategorySlug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CategoryGet200JSONResponse{
		CategoryGetOKJSONResponse: openapi.CategoryGetOKJSONResponse(serialiseCategory(cat)),
	}, nil
}

func (c Categories) CategoryUpdatePosition(ctx context.Context, request openapi.CategoryUpdatePositionRequestObject) (openapi.CategoryUpdatePositionResponseObject, error) {
	parent, err := deletable.NewMapErr(request.Body.Parent, func(in openapi.Identifier) (category.CategoryID, error) {
		return category.CategoryID(deserialiseID(in)), nil
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	before := opt.NewEmpty[category.CategoryID]()
	if request.Body.Before != nil {
		before = opt.New(category.CategoryID(deserialiseID(*request.Body.Before)))
	}

	after := opt.NewEmpty[category.CategoryID]()
	if request.Body.After != nil {
		after = opt.New(category.CategoryID(deserialiseID(*request.Body.After)))
	}

	move := category_svc.Move{
		Parent: parent,
		Before: before,
		After:  after,
	}

	cats, err := c.category_svc.Move(ctx, request.CategorySlug, move)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CategoryUpdatePosition200JSONResponse{
		CategoryListOKJSONResponse: openapi.CategoryListOKJSONResponse{
			Categories: dt.Map(cats, serialiseCategory),
		},
	}, nil
}

func (c Categories) CategoryUpdate(ctx context.Context, request openapi.CategoryUpdateRequestObject) (openapi.CategoryUpdateResponseObject, error) {
	coverImageAssetID := deletable.NewMap(request.Body.CoverImageAssetId, func(id openapi.NullableIdentifier) *xid.ID {
		xidValue := openapi.ParseID(openapi.Identifier(id))
		return &xidValue
	})

	cat, err := c.category_svc.Update(ctx, request.CategorySlug, category_svc.Partial{
		Name:              opt.NewPtr(request.Body.Name),
		Slug:              opt.NewPtr(request.Body.Slug),
		Description:       opt.NewPtr(request.Body.Description),
		Colour:            opt.NewPtr(request.Body.Colour),
		Admin:             opt.NewPtr(request.Body.Admin),
		CoverImageAssetID: coverImageAssetID,
		Meta:              opt.NewPtr((*map[string]any)(request.Body.Meta)),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CategoryUpdate200JSONResponse{
		CategoryUpdateOKJSONResponse: openapi.CategoryUpdateOKJSONResponse(serialiseCategory(cat)),
	}, nil
}

func serialiseCategory(c *category.Category) openapi.Category {
	var parentID *openapi.NullableIdentifier
	if c.ParentID != nil {
		pid := openapi.NullableIdentifier(xid.ID(*c.ParentID).String())
		parentID = &pid
	}

	children := dt.Map(c.Children, serialiseCategory)

	return openapi.Category{
		Id:          *openapi.IdentifierFrom(xid.ID(c.ID)),
		Name:        c.Name,
		Slug:        c.Slug,
		Colour:      c.Colour,
		Description: c.Description,
		PostCount:   c.PostCount,
		Admin:       c.Admin,
		Sort:        c.Sort,
		Parent:      parentID,
		CoverImage:  opt.Map(c.CoverImage, serialiseAsset).Ptr(),
		Children:    children,
		Meta:        (*openapi.Metadata)(&c.Metadata),
	}
}

func serialiseCategoryReference(c category.Category) openapi.CategoryReference {
	var parentID *openapi.NullableIdentifier
	if c.ParentID != nil {
		pid := openapi.NullableIdentifier(xid.ID(*c.ParentID).String())
		parentID = &pid
	}

	children := dt.Map(c.Children, serialiseCategory)

	return openapi.CategoryReference{
		Id:          *openapi.IdentifierFrom(xid.ID(c.ID)),
		Name:        c.Name,
		Slug:        c.Slug,
		Admin:       c.Admin,
		Colour:      c.Colour,
		Description: c.Description,
		Sort:        c.Sort,
		Parent:      parentID,
		CoverImage:  opt.Map(c.CoverImage, serialiseAsset).Ptr(),
		Children:    children,
		Meta:        (*openapi.Metadata)(&c.Metadata),
	}
}

func serialiseCategoryReferencePtr(c *category.Category) openapi.CategoryReference {
	return serialiseCategoryReference(*c)
}
