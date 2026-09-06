package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"

	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/branding/theme"
	"github.com/Southclaws/storyden/app/services/reqinfo"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Theme struct {
	service *theme.Service
}

func NewTheme(service *theme.Service) Theme {
	return Theme{service: service}
}

func (t Theme) ThemeGet(ctx context.Context, request openapi.ThemeGetRequestObject) (openapi.ThemeGetResponseObject, error) {
	manifest, err := t.service.GetPublic(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ThemeGet200JSONResponse{
		ThemeGetOKJSONResponse: openapi.ThemeGetOKJSONResponse{
			Body: serialiseThemeManifest(manifest),
			Headers: openapi.ThemeGetOKResponseHeaders{
				CacheControl: ptr("public, max-age=60, stale-while-revalidate=300"),
			},
		},
	}, nil
}

func (t Theme) ThemeAssetGet(ctx context.Context, request openapi.ThemeAssetGetRequestObject) (openapi.ThemeAssetGetResponseObject, error) {
	a, err := t.service.GetAsset(ctx, request.AssetFilename)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	etag := cachecontrol.NewETagValue(a.Integrity)
	if reqinfo.GetCacheQuery(ctx).MatchesETag(etag.String()) {
		return openapi.ThemeAssetGet304Response{
			Headers: openapi.NotModifiedResponseHeaders{
				CacheControl: ptr("public, max-age=31536000, immutable"),
				ETag:         ptr(etag.String()),
			},
		}, nil
	}

	r, size, err := t.service.OpenAsset(ctx, a.Filename)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ThemeAssetGet200AsteriskResponse{
		ThemeAssetGetOKAsteriskResponse: openapi.ThemeAssetGetOKAsteriskResponse{
			Body:          r,
			ContentType:   a.MIMEType,
			ContentLength: size,
			Headers: openapi.ThemeAssetGetOKResponseHeaders{
				CacheControl: ptr("public, max-age=31536000, immutable"),
				ETag:         ptr(etag.String()),
			},
		},
	}, nil
}

func (t Theme) AdminThemeGet(ctx context.Context, request openapi.AdminThemeGetRequestObject) (openapi.AdminThemeGetResponseObject, error) {
	manifest, err := t.service.GetAdmin(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return openapi.AdminThemeGet200JSONResponse{
		AdminThemeGetOKJSONResponse: openapi.AdminThemeGetOKJSONResponse(serialiseThemeManifest(manifest)),
	}, nil
}

func (t Theme) AdminThemeUpdate(ctx context.Context, request openapi.AdminThemeUpdateRequestObject) (openapi.AdminThemeUpdateResponseObject, error) {
	manifest, err := t.service.Publish(ctx, request.Body.Stylesheets, request.Body.Scripts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return openapi.AdminThemeUpdate200JSONResponse{
		AdminThemeUpdateOKJSONResponse: openapi.AdminThemeUpdateOKJSONResponse(serialiseThemeManifest(manifest)),
	}, nil
}

func (t Theme) AdminThemeAssetUpload(ctx context.Context, request openapi.AdminThemeAssetUploadRequestObject) (openapi.AdminThemeAssetUploadResponseObject, error) {
	if request.Params.Filename == nil {
		return nil, fault.Wrap(fault.New("filename is required"), fctx.With(ctx), fmsg.WithDesc("missing filename", "Provide a .css or .js filename."))
	}
	a, err := t.service.Upload(ctx, request.Body, request.Params.ContentLength, *request.Params.Filename)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return openapi.AdminThemeAssetUpload200JSONResponse{
		AdminThemeAssetUploadOKJSONResponse: openapi.AdminThemeAssetUploadOKJSONResponse(serialiseThemeAsset(a)),
	}, nil
}

func (t Theme) AdminThemeAssetDelete(ctx context.Context, request openapi.AdminThemeAssetDeleteRequestObject) (openapi.AdminThemeAssetDeleteResponseObject, error) {
	if err := t.service.DeleteAsset(ctx, request.AssetFilename); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return openapi.AdminThemeAssetDelete204Response{}, nil
}

func serialiseThemeManifest(in settings.ThemeSettings) openapi.ThemeManifest {
	return openapi.ThemeManifest{
		Stylesheets: serialiseThemeAssets(in.Stylesheets),
		Scripts:     serialiseThemeAssets(in.Scripts),
	}
}

func serialiseThemeAssets(in []settings.ThemeAsset) []openapi.ThemeAsset {
	out := make([]openapi.ThemeAsset, 0, len(in))
	for _, item := range in {
		out = append(out, serialiseThemeAsset(item))
	}
	return out
}

func serialiseThemeAsset(in settings.ThemeAsset) openapi.ThemeAsset {
	return openapi.ThemeAsset{
		Id:        openapi.Identifier(in.ID),
		Filename:  in.Filename,
		Path:      "/api/info/theme/assets/" + in.Filename,
		MimeType:  openapi.ThemeAssetMimeType(in.MIMEType),
		Size:      in.Size,
		Integrity: in.Integrity,
	}
}
