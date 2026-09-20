package model_media

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_querier"
	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/llm_provider"
	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
	"github.com/Southclaws/storyden/internal/config"
)

var supportedImageMIMETypes = map[string]struct{}{
	"image/gif":  {},
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
}

type Image struct {
	AssetID asset.AssetID
	URL     string
	MIME    string
}

type ImageResolver interface {
	ResolveImage(context.Context, asset.AssetID) (Image, error)
}

type Resolver struct {
	assets    *asset_querier.Querier
	publicAPI url.URL
}

func New(cfg config.Config, assets *asset_querier.Querier) *Resolver {
	return &Resolver{
		assets:    assets,
		publicAPI: cfg.PublicAPIAddress,
	}
}

func (r *Resolver) ValidateContents(ctx context.Context, contents []*genai.Content) error {
	ids, err := robotresource.UniqueImageAssetIDs(contents...)
	if err != nil {
		return err
	}
	for _, id := range ids {
		if _, err := r.resolve(ctx, id, false); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) ResolveImage(ctx context.Context, id asset.AssetID) (Image, error) {
	return r.resolve(ctx, id, true)
}

func ValidateModelInput(ctx context.Context, provider llm_provider.Provider, ref model_ref.ModelRef, contents []*genai.Content) error {
	ids, err := robotresource.ImageAssetIDs(contents...)
	if err != nil || len(ids) == 0 {
		return err
	}

	capabilities, err := provider.ModelCapabilities(ctx, ref)
	if err != nil {
		return err
	}
	if capabilities.ImageInput == llm_provider.CapabilitySupportUnsupported {
		return fmt.Errorf("model %q does not support image input", ref.Model)
	}
	return nil
}

func (r *Resolver) resolve(ctx context.Context, id asset.AssetID, includeURL bool) (Image, error) {
	a, err := r.assets.GetByID(ctx, id)
	if err != nil {
		return Image{}, fault.Wrap(err, fctx.With(ctx), fmsg.Withf("image asset %s is not available", id.String()))
	}

	mimeType := a.MIME.String()
	if _, ok := supportedImageMIMETypes[mimeType]; !ok {
		return Image{}, fault.New(
			fmt.Sprintf("asset %s has unsupported image MIME type %q", id.String(), mimeType),
			fctx.With(ctx),
			fmsg.With("Robot image attachments must be GIF, JPEG, PNG, or WebP assets."),
		)
	}

	image := Image{AssetID: id, MIME: mimeType}
	if includeURL {
		image.URL = r.publicAPI.JoinPath("api", "assets", url.PathEscape(a.Name.String())).String()
	}
	return image, nil
}
