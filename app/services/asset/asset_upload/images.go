package asset_upload

import (
	"context"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"golang.org/x/image/webp"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/settings"
)

var errInvalidUpload = fault.New("invalid asset upload")

func (s *Uploader) validateImage(
	ctx context.Context,
	name asset.Filename,
	mimeType string,
	size int64,
	limits settings.ImageLimits,
) error {
	if size > limits.MaxSize {
		return fault.Wrap(errInvalidUpload,
			fmsg.WithDesc("stored image size exceeds upload limit", fmt.Sprintf(
				"This image exceeds the maximum upload size of %d bytes.", limits.MaxSize,
			)),
			ftag.With(ftag.InvalidArgument))
	}

	var decodeConfig func(io.Reader) (image.Config, error)
	switch mimeType {
	case "image/jpeg":
		decodeConfig = jpeg.DecodeConfig
	case "image/png":
		decodeConfig = png.DecodeConfig
	case "image/webp":
		decodeConfig = webp.DecodeConfig
	case "image/gif":
		decodeConfig = gif.DecodeConfig
	default:
		return nil
	}

	reader, _, err := s.objects.Read(ctx, asset.BuildAssetPath(name))
	if err != nil {
		return err
	}
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	config, err := decodeConfig(reader)
	if err != nil {
		// Malformed originals remain retrievable; the processing job skips them.
		return nil
	}

	width, height := int64(config.Width), int64(config.Height)
	if height > 0 && width > limits.MaxPixels/height {
		return fault.Wrap(errInvalidUpload,
			fmsg.WithDesc("image pixel count exceeds upload limit", fmt.Sprintf(
				"This image exceeds the maximum resolution of %d pixels.", limits.MaxPixels,
			)),
			ftag.With(ftag.InvalidArgument))
	}

	return nil
}

func (s *Uploader) imageUploadLimits(ctx context.Context, size int64) (settings.ImageLimits, error) {
	configured, err := s.settings.Get(ctx)
	if err != nil {
		return settings.ImageLimits{}, err
	}

	limits, err := configured.ImageLimits()
	if err != nil {
		return settings.ImageLimits{}, err
	}
	if size > limits.MaxSize {
		return settings.ImageLimits{}, fault.Wrap(errInvalidUpload,
			fmsg.WithDesc("declared image size exceeds upload limit", fmt.Sprintf(
				"This image exceeds the maximum upload size of %d bytes.", limits.MaxSize,
			)),
			ftag.With(ftag.InvalidArgument))
	}

	return limits, nil
}
