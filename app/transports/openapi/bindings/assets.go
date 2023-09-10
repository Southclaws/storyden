package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/services/asset"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Assets struct {
	a asset.Service
}

func NewAssets(a asset.Service) Assets {
	return Assets{a}
}

func (i *Assets) AssetGet(ctx context.Context, request openapi.AssetGetRequestObject) (openapi.AssetGetResponseObject, error) {
	a, r, err := i.a.Get(ctx, request.Id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AssetGet200AsteriskResponse{
		AssetGetOKAsteriskResponse: openapi.AssetGetOKAsteriskResponse{
			Body:        r,
			ContentType: a.MIMEType,
			// ContentLength: ,
		},
	}, nil
}

func (i *Assets) AssetUpload(ctx context.Context, request openapi.AssetUploadRequestObject) (openapi.AssetUploadResponseObject, error) {
	a, err := i.a.Upload(ctx, request.Body, int64(request.Params.ContentLength))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AssetUpload200JSONResponse{
		AssetUploadOKJSONResponse: openapi.AssetUploadOKJSONResponse(serialiseAssetReference(a)),
	}, nil
}
