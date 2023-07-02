package bindings

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/object"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Assets struct {
	os      object.Storer
	address string
}

func NewAssets(cfg config.Config, os object.Storer) Assets {
	return Assets{os, cfg.PublicWebAddress}
}

const assetsSubdirectory = "assets"

func (i *Assets) AssetGetUploadURL(ctx context.Context, request openapi.AssetGetUploadURLRequestObject) (openapi.AssetGetUploadURLResponseObject, error) {
	// TODO: Check if S3 is available and create a pre-signed upload URL if so.

	url := fmt.Sprintf("%s/api/v1/assets", i.address)

	return openapi.AssetGetUploadURL200JSONResponse{
		AssetGetUploadURLOKJSONResponse: openapi.AssetGetUploadURLOKJSONResponse{
			Url: url,
		},
	}, nil
}

func (i *Assets) AssetGet(ctx context.Context, request openapi.AssetGetRequestObject) (openapi.AssetGetResponseObject, error) {
	path := filepath.Join(assetsSubdirectory, request.Id)

	r, err := i.os.Read(ctx, path)
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
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	assetID := fmt.Sprintf("%s-%s", accountID.String(), xid.New().String())
	path := filepath.Join(assetsSubdirectory, assetID)

	if err := i.os.Write(ctx, path, request.Body); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	url := fmt.Sprintf("%s/api/v1/assets/%s", i.address, assetID)

	return openapi.AssetUpload200JSONResponse{
		AssetUploadOKJSONResponse: openapi.AssetUploadOKJSONResponse{
			Url: url,
		},
	}, nil
}
