package icon

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"

	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/asset/asset_download"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/mime"
)

//go:embed embed/default_icon.png
var defaultIcon []byte

type Service interface {
	Upload(context.Context, io.Reader) error
	Get(context.Context, string) (*asset.Asset, io.Reader, error)
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
		asset.NewFilename("icon-original"),
		asset_upload.Options{ImagePresets: asset.IconImagePresets},
	)
	if err != nil {
		return err
	}

	_, err = s.settings.Set(ctx, settings.Settings{IconAssetID: opt.New(source.ID)})
	return err
}

func (s *service) Get(ctx context.Context, size string) (*asset.Asset, io.Reader, error) {
	var pixels int
	if _, err := fmt.Sscanf(size, "%dx", &pixels); err != nil {
		return nil, nil, err
	}

	variant := asset.Variant(fmt.Sprintf("icon_%d", pixels))
	if preset, ok := asset.Preset(variant); !ok || !preset.PNG {
		return nil, nil, fmt.Errorf("invalid icon size %q", size)
	}

	configured, err := s.settings.Get(ctx)
	if err != nil {
		return nil, nil, err
	}
	if id, ok := configured.IconAssetID.Get(); ok {
		icon, stream, pending, err := s.download.VariantByID(ctx, id, variant)
		if err != nil {
			return nil, nil, err
		}
		if !pending && icon.MIME.String() == "image/png" {
			return icon, stream, nil
		}
		if closer, ok := stream.(io.Closer); ok {
			closer.Close()
		}
	}

	icon := &asset.Asset{
		Name: asset.NewFilepathFilename("default-icon.png"),
		Size: len(defaultIcon),
		MIME: mime.New("image/png"),
	}

	return icon, bytes.NewReader(defaultIcon), nil
}
