package bindings

import (
	"context"
	"fmt"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/services/asset"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Assets struct {
	a asset.Service
}

func NewAssets(a asset.Service) Assets {
	return Assets{a}
}

func (i *Assets) AssetGetUploadURL(ctx context.Context, request openapi.AssetGetUploadURLRequestObject) (openapi.AssetGetUploadURLResponseObject, error) {
	// TODO: Check if S3 is available and create a pre-signed upload URL if so.

	url := fmt.Sprintf("%s/api/v1/assets", "i.address")

	return openapi.AssetGetUploadURL200JSONResponse{
		AssetGetUploadURLOKJSONResponse: openapi.AssetGetUploadURLOKJSONResponse{
			Url: url,
		},
	}, nil
}

func (i *Assets) AssetGet(ctx context.Context, request openapi.AssetGetRequestObject) (openapi.AssetGetResponseObject, error) {
	r, err := i.a.Read(ctx, request.Id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AssetGet200AsteriskResponse{
		AssetGetOKAsteriskResponse: openapi.AssetGetOKAsteriskResponse{
			Body: r,
		},
	}, nil
}

func (i *Assets) AssetUpload(ctx context.Context, request openapi.AssetUploadRequestObject) (openapi.AssetUploadResponseObject, error) {
	postID := openapi.ParseID(request.Params.PostId)

	url, err := i.a.Upload(ctx, post.PostID(postID), request.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AssetUpload200JSONResponse{
		AssetUploadOKJSONResponse: openapi.AssetUploadOKJSONResponse{
			Url: url,
		},
	}, nil
}
