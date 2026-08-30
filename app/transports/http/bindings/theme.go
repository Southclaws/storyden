package bindings

import (
	"context"
	"fmt"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"

	"github.com/Southclaws/storyden/app/services/branding/theme"
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

	etag := quoteETag(manifest.Revision)
	if etagMatches(request.Params.IfNoneMatch, etag) {
		return openapi.ThemeGet304Response{}, nil
	}
	return openapi.ThemeGet200JSONResponse{
		ThemeGetOKJSONResponse: openapi.ThemeGetOKJSONResponse{
			Body: serialiseThemeManifest(manifest),
			Headers: openapi.ThemeGetOKResponseHeaders{
				CacheControl: ptr("public, max-age=60, stale-while-revalidate=300"),
				ETag:         &etag,
			},
		},
	}, nil
}

func (t Theme) ThemeAssetGet(ctx context.Context, request openapi.ThemeAssetGetRequestObject) (openapi.ThemeAssetGetResponseObject, error) {
	a, r, err := t.service.GetAsset(ctx, request.AssetFilename)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	etag := quoteETag(a.Integrity)
	if etagMatches(request.Params.IfNoneMatch, etag) {
		return openapi.ThemeAssetGet304Response{}, nil
	}

	return openapi.ThemeAssetGet200AsteriskResponse{
		ThemeAssetGetOKAsteriskResponse: openapi.ThemeAssetGetOKAsteriskResponse{
			Body:          r,
			ContentType:   a.MIMEType,
			ContentLength: int64(a.Size),
			Headers: openapi.ThemeAssetGetOKResponseHeaders{
				CacheControl: ptr("public, max-age=31536000, immutable"),
				ETag:         &etag,
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
		AdminThemeGetOKJSONResponse: openapi.AdminThemeGetOKJSONResponse(serialiseAdminThemeManifest(manifest)),
	}, nil
}

func (t Theme) AdminThemeUpdate(ctx context.Context, request openapi.AdminThemeUpdateRequestObject) (openapi.AdminThemeUpdateResponseObject, error) {
	if request.Body == nil {
		return nil, fault.Wrap(fault.New("theme manifest is required"), fctx.With(ctx), fmsg.WithDesc("missing theme manifest", "Provide stylesheet and script asset IDs."))
	}
	manifest, err := t.service.Publish(ctx, request.Body.Stylesheets, request.Body.Scripts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return openapi.AdminThemeUpdate200JSONResponse{
		AdminThemeUpdateOKJSONResponse: openapi.AdminThemeUpdateOKJSONResponse(serialiseAdminThemeManifest(manifest)),
	}, nil
}

func (t Theme) AdminThemeDelete(ctx context.Context, request openapi.AdminThemeDeleteRequestObject) (openapi.AdminThemeDeleteResponseObject, error) {
	if err := t.service.Disable(ctx); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return openapi.AdminThemeDelete204Response{}, nil
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

func serialiseThemeManifest(in theme.Manifest) openapi.ThemeManifest {
	return openapi.ThemeManifest{
		ApiVersion:  openapi.ThemeManifestApiVersion(in.APIVersion),
		Enabled:     in.Enabled,
		Revision:    in.Revision,
		Stylesheets: serialiseThemeAssets(in.Stylesheets),
		Scripts:     serialiseThemeAssets(in.Scripts),
	}
}

func serialiseAdminThemeManifest(in theme.Manifest) openapi.AdminThemeManifest {
	return openapi.AdminThemeManifest{
		ApiVersion:      openapi.AdminThemeManifestApiVersion(in.APIVersion),
		Enabled:         in.Enabled,
		Revision:        in.Revision,
		RuntimeDisabled: in.RuntimeDisabled,
		Stylesheets:     serialiseThemeAssets(in.Stylesheets),
		Scripts:         serialiseThemeAssets(in.Scripts),
	}
}

func serialiseThemeAssets(in []theme.Asset) []openapi.ThemeAsset {
	out := make([]openapi.ThemeAsset, 0, len(in))
	for _, item := range in {
		out = append(out, serialiseThemeAsset(item))
	}
	return out
}

func serialiseThemeAsset(in theme.Asset) openapi.ThemeAsset {
	return openapi.ThemeAsset{
		Id:        openapi.Identifier(in.ID),
		Filename:  in.Filename,
		Path:      in.Path,
		MimeType:  openapi.ThemeAssetMimeType(in.MIMEType),
		Size:      in.Size,
		Integrity: in.Integrity,
	}
}

func quoteETag(value string) string {
	if value == "" {
		value = "empty"
	}
	return fmt.Sprintf("%q", value)
}

func etagMatches(header *openapi.IfNoneMatch, etag string) bool {
	if header == nil {
		return false
	}
	for _, candidate := range strings.Split(string(*header), ",") {
		candidate = strings.TrimSpace(candidate)
		if candidate == "*" || candidate == etag || strings.TrimPrefix(candidate, "W/") == etag {
			return true
		}
	}
	return false
}
