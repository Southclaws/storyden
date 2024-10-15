package bindings

import (
	"context"
	"fmt"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/services/asset/asset_download"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
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
			ContentType:   a.Metadata.GetMIMEType(),
			ContentLength: int64(a.Size),
		},
	}, nil
}

func (i *Assets) AssetUpload(ctx context.Context, request openapi.AssetUploadRequestObject) (openapi.AssetUploadResponseObject, error) {
	name := opt.NewPtr(request.Params.Filename)

	// NOTE: Should we enforce a filename for upload? If none is available, the
	// client can decide on a suitable placeholder or generate a nonsense slug.
	filename := asset.NewFilename(name.Or("untitled"))

	parentID := opt.NewPtrMap(request.Params.ParentAssetId, deserialiseAssetID)

	contentFillCmd, err := getContentFillRuleCommand(request.Params.ContentFillRule, request.Params.NodeContentFillTarget)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	opts := asset_upload.Options{
		ParentID:    parentID,
		ContentFill: contentFillCmd,
	}

	a, err := i.uploader.Upload(ctx, request.Body, int64(request.Params.ContentLength), filename, opts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AssetUpload200JSONResponse{
		AssetUploadOKJSONResponse: openapi.AssetUploadOKJSONResponse(serialiseAssetPtr(a)),
	}, nil
}

func getContentFillRuleCommand(contentFillRuleParam *openapi.ContentFillRule, contentFillTargetParam *string) (opt.Optional[asset.ContentFillCommand], error) {
	if contentFillRuleParam != nil {
		if contentFillTargetParam == nil {
			return nil, fault.New("node_content_fill_target is required when content_fill_rule is specified")
		}

		rule, err := asset.NewContentFillRule((string)(*contentFillRuleParam))
		if err != nil {
			return nil, fault.Wrap(err)
		}

		nodeID, err := xid.FromString(*contentFillTargetParam)
		if err != nil {
			return nil, fault.Wrap(err)
		}

		return opt.New(asset.ContentFillCommand{
			TargetNodeID: nodeID,
			FillRule:     rule,
		}), nil
	}

	return opt.NewEmpty[asset.ContentFillCommand](), nil
}

func serialiseAsset(a asset.Asset) openapi.Asset {
	path := fmt.Sprintf(`/api/assets/%s`, a.Name.String())

	return openapi.Asset{
		Id:       a.ID.String(),
		Filename: a.Name.String(),
		Path:     path,
		Parent:   opt.Map(a.Parent, serialiseAsset).Ptr(),
		MimeType: a.Metadata.GetMIMEType(),
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
