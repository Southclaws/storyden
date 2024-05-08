package bindings

import (
	"context"
	"fmt"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/services/asset_manager"
	"github.com/Southclaws/storyden/app/transports/openapi"
	"github.com/Southclaws/storyden/internal/config"
)

type Assets struct {
	a       asset_manager.Service
	address string
}

func NewAssets(cfg config.Config, a asset_manager.Service) Assets {
	return Assets{a, cfg.PublicWebAddress}
}

func (i *Assets) AssetGet(ctx context.Context, request openapi.AssetGetRequestObject) (openapi.AssetGetResponseObject, error) {
	a, r, err := i.a.Get(ctx, asset.NewFilepathFilename(request.AssetFilename))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AssetGet200AsteriskResponse{
		AssetGetOKAsteriskResponse: openapi.AssetGetOKAsteriskResponse{
			Body:          r,
			ContentType:   a.Metadata.GetMIMEType(),
			ContentLength: int64(a.Size),
		},
	}, nil
}

func (i *Assets) AssetUpload(ctx context.Context, request openapi.AssetUploadRequestObject) (openapi.AssetUploadResponseObject, error) {
	name := "" // TODO: get filename from API call

	filename := asset.NewFilename(name)
	url := fmt.Sprintf("%s/api/v1/assets/%s", i.address, filename.String())

	a, err := i.a.Upload(ctx, request.Body, int64(request.Params.ContentLength), filename, url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AssetUpload200JSONResponse{
		AssetUploadOKJSONResponse: openapi.AssetUploadOKJSONResponse(serialiseAssetReference(a)),
	}, nil
}
