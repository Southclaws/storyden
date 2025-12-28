package bindings

import (
	"context"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/post/category_cache"
	category_svc "github.com/Southclaws/storyden/app/services/category"
	"github.com/Southclaws/storyden/app/services/reqinfo"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/deletable"
)

type Categories struct {
	category_repo  *category.Repository
	category_svc   category_svc.Service
	category_cache *category_cache.Cache
}

func NewCategories(
	category_repo *category.Repository,
	category_svc category_svc.Service,
	category_cache *category_cache.Cache,
) Categories {
	return Categories{category_repo, category_svc, category_cache}
}

func (c Categories) CategoryCreate(ctx context.Context, request openapi.CategoryCreateRequestObject) (openapi.CategoryCreateResponseObject, error) {
	parentID := opt.NewEmpty[category.CategoryID]()
	if request.Body.Parent != nil {
		parentID = opt.New(category.CategoryID(deserialiseID(*request.Body.Parent)))
	}

	coverImageAssetID := deletable.Value[*xid.ID]{}
	if request.Body.CoverImageAssetId != nil {
		xidValue := openapi.ParseID(*request.Body.CoverImageAssetId)
		coverImageAssetID = deletable.Skip(opt.New(&xidValue))
	}

	cat, err := c.category_svc.Create(ctx, category_svc.Partial{
		Name:              opt.New(request.Body.Name),
		Slug:              opt.NewPtr(request.Body.Slug),
		Description:       opt.New(request.Body.Description),
		Colour:            opt.New(request.Body.Colour),
		Parent:            parentID,
		CoverImageAssetID: coverImageAssetID,
		Meta:              opt.NewPtr((*map[string]any)(request.Body.Meta)),
	})
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

const categoryGetCacheControl = "public, no-cache"

func (c Categories) CategoryGet(ctx context.Context, request openapi.CategoryGetRequestObject) (openapi.CategoryGetResponseObject, error) {
	slug := string(request.CategorySlug)

	etag, notModified := c.category_cache.Check(ctx, reqinfo.GetCacheQuery(ctx), slug)
	if notModified {
		return openapi.CategoryGet304Response{
			Headers: openapi.NotModifiedResponseHeaders{
				CacheControl: categoryGetCacheControl,
				LastModified: etag.Time.Format(time.RFC1123),
				ETag:         etag.String(),
			},
		}, nil
	}

	cat, err := c.category_repo.Get(ctx, request.CategorySlug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if etag == nil {
		c.category_cache.Store(ctx, slug, cat.UpdatedAt)
		etag = cachecontrol.NewETag(cat.UpdatedAt)
	}

	return openapi.CategoryGet200JSONResponse{
		CategoryGetOKJSONResponse: openapi.CategoryGetOKJSONResponse{
			Body: serialiseCategory(cat),
			Headers: openapi.CategoryGetOKResponseHeaders{
				CacheControl: categoryGetCacheControl,
				LastModified: etag.Time.Format(time.RFC1123),
				ETag:         etag.String(),
			},
		},
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

func (c Categories) CategoryDelete(ctx context.Context, request openapi.CategoryDeleteRequestObject) (openapi.CategoryDeleteResponseObject, error) {
	moveToID := category.CategoryID(deserialiseID(request.Body.MoveTo))

	cat, err := c.category_svc.Delete(ctx, request.CategorySlug, moveToID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CategoryDelete200JSONResponse{
		CategoryDeleteOKJSONResponse: openapi.CategoryDeleteOKJSONResponse(serialiseCategory(cat)),
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
