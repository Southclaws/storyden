package bindings

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
	"github.com/Southclaws/storyden/app/services/asset_manager"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
)

type Assets struct {
	a        asset_manager.Service
	uploader *asset_upload.Uploader
	address  url.URL
}

func NewAssets(cfg config.Config, a asset_manager.Service, uploader *asset_upload.Uploader) Assets {
	return Assets{a, uploader, cfg.PublicAPIAddress}
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
	name := opt.NewPtr(request.Params.Filename)

	// NOTE: Should we enforce a filename for upload? If none is available, the
	// client can decide on a suitable placeholder or generate a nonsense slug.
	filename := asset.NewFilename(name.Or("untitled"))

	// TODO: This must be specified on the READ path not the write path.
	// It's not the responsibility of the write-path transport layer to figure
	// out the public URL of the asset. This may change if a direct CDN is used.
	url := fmt.Sprintf("%s/api/v1/assets/%s", i.address.String(), filename.String())

	contentFillCmd, err := getContentFillRuleCommand(request.Params.ContentFillRule, request.Params.NodeContentFillTarget)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	opts := asset_upload.Options{
		ContentFill: contentFillCmd,
	}

	a, err := i.uploader.Upload(ctx, request.Body, int64(request.Params.ContentLength), filename, url, opts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AssetUpload200JSONResponse{
		AssetUploadOKJSONResponse: openapi.AssetUploadOKJSONResponse(serialiseAssetReference(a)),
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
