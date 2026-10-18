package banner

import (
	"bytes"
	"context"
	"io"

	"github.com/Southclaws/opt"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/asset/asset_download"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/mime"
)

type Service interface {
	Upload(context.Context, io.Reader) error
	Get(context.Context) (*asset.Asset, io.Reader, error)
	GetVariant(context.Context, asset.Variant) (*asset.Asset, io.Reader, bool, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	settings *settings.SettingsRepository
	upload   *asset_upload.Uploader
	download *asset_download.Downloader
}

func New(
	settings *settings.SettingsRepository,
	upload *asset_upload.Uploader,
	download *asset_download.Downloader,
) Service {
	return &service{
		settings: settings,
		upload:   upload,
		download: download,
	}
}

func (s *service) Upload(ctx context.Context, stream io.Reader) error {
	owner, err := session.GetAccountID(ctx)
	if err != nil {
		return err
	}

	source, err := s.upload.UploadImage(
		ctx,
		owner,
		stream,
		-1,
		asset.NewFilename("banner-original"),
		asset_upload.Options{ImagePresets: asset.BannerImagePresets},
	)
	if err != nil {
		return err
	}

	_, err = s.settings.Set(ctx, settings.Settings{BannerAssetID: opt.New(source.ID)})
	return err
}

func (s *service) Get(ctx context.Context) (*asset.Asset, io.Reader, error) {
	banner, stream, _, err := s.GetVariant(ctx, "banner")
	return banner, stream, err
}

func (s *service) GetVariant(ctx context.Context, variant asset.Variant) (*asset.Asset, io.Reader, bool, error) {
	configured, err := s.settings.Get(ctx)
	if err != nil {
		return nil, nil, false, err
	}
	if id, ok := configured.BannerAssetID.Get(); ok {
		banner, stream, pending, err := s.download.VariantByID(ctx, id, variant)
		if err != nil {
			return nil, nil, false, err
		}
		preset, _ := asset.Preset(variant)
		if !preset.PNG || (!pending && banner.MIME.String() == "image/png") {
			return banner, stream, pending, nil
		}
		if closer, ok := stream.(io.Closer); ok {
			closer.Close()
		}
	}

	banner := &asset.Asset{
		Name: asset.NewFilepathFilename("default-banner.png"),
		Size: len(defaultBanner),
		MIME: mime.New("image/png"),
	}

	return banner, bytes.NewReader(defaultBanner), true, nil
}
