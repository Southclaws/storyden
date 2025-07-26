package bindings

import (
	"context"
	"fmt"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/asset/asset_download"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Assets struct {
	uploader   *asset_upload.Uploader
	downloader *asset_download.Downloader
}

func NewAssets(uploader *asset_upload.Uploader, downloader *asset_download.Downloader) Assets {
	return Assets{uploader, downloader}
}

func (i *Assets) AssetGet(ctx context.Context, request openapi.AssetGetRequestObject) (openapi.AssetGetResponseObject, error) {
	a, r, err := i.downloader.Get(ctx, asset.NewFilepathFilename(request.AssetFilename))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AssetGet200AsteriskResponse{
		AssetGetOKAsteriskResponse: openapi.AssetGetOKAsteriskResponse{
			Body:          r,
			ContentType:   a.MIME.String(),
			ContentLength: int64(a.Size),
			Headers: openapi.AssetGetOKResponseHeaders{
				CacheControl: "public, max-age=31536000",
			},
		},
	}, nil
}

func (i *Assets) AssetUpload(ctx context.Context, request openapi.AssetUploadRequestObject) (openapi.AssetUploadResponseObject, error) {
	// NOTE: This op doesn't run the authorisation validator for some reason.
	if !session.GetOptAccountID(ctx).Ok() {
		return nil, fault.Wrap(fault.New("session required for upload", fctx.With(ctx)), fctx.With(ctx), ftag.With(ftag.Unauthenticated))
	}

	if err := session.Authorise(ctx, nil, rbac.PermissionUploadAsset); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
	}

	name := opt.NewPtr(request.Params.Filename)

	// NOTE: Should we enforce a filename for upload? If none is available, the
	// client can decide on a suitable placeholder or generate a nonsense slug.
	filename := asset.NewFilename(name.Or("untitled"))

	parentID := opt.NewPtrMap(request.Params.ParentAssetId, deserialiseAssetID)

	opts := asset_upload.Options{
		ParentID: parentID,
	}

	a, err := i.uploader.Upload(ctx, request.Body, request.Params.ContentLength, filename, opts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AssetUpload200JSONResponse{
		AssetUploadOKJSONResponse: openapi.AssetUploadOKJSONResponse(serialiseAssetPtr(a)),
	}, nil
}

func serialiseAsset(a asset.Asset) openapi.Asset {
	path := fmt.Sprintf(`/api/assets/%s`, a.Name.String())

	return openapi.Asset{
		Id:       a.ID.String(),
		Filename: a.Name.String(),
		Path:     path,
		Parent:   opt.Map(a.Parent, serialiseAsset).Ptr(),
		MimeType: a.MIME.String(),
		Width:    float32(a.Metadata.GetWidth()),
		Height:   float32(a.Metadata.GetHeight()),
	}
}

func serialiseAssetPtr(a *asset.Asset) openapi.Asset {
	return serialiseAsset(*a)
}

func deserialiseAssetID(in string) asset.AssetID {
	return asset.AssetID(openapi.ParseID(in))
}

func deserialiseAssetIDs(ids []string) []asset.AssetID {
	return dt.Map(ids, deserialiseAssetID)
}
