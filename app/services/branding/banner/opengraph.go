package banner

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/disintegration/imaging"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/services/asset/asset_download"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/mime"
)

const (
	bannerFileTemplate = "banner-%s.png"
)

type Size string

var sizes = []Size{
	"1200x630",
}

var sizeMap = map[Size][2]int{
	"1200x630": {1200, 630},
}

var errBadFormat = fault.Wrap(fault.New("bad format"), ftag.With(ftag.InvalidArgument))

type Service interface {
	Upload(ctx context.Context, r io.Reader) error
	Get(ctx context.Context) (*asset.Asset, io.Reader, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	uploader   *asset_upload.Uploader
	downloader *asset_download.Downloader
}

func New(
	uploader *asset_upload.Uploader,
	downloader *asset_download.Downloader,

	cfg config.Config,
) Service {
	return &service{
		uploader:   uploader,
		downloader: downloader,
	}
}

func (s *service) Upload(ctx context.Context, r io.Reader) error {
	return s.uploadSizes(ctx, r, sizes)
}

func (s *service) uploadSizes(ctx context.Context, or io.Reader, sizes []Size) error {
	mt, r, err := mime.Detect(or)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	mime := mt.String()
	ctx = fctx.WithMeta(ctx, "mimetype", mime)

	if !strings.HasPrefix(mime, "image") {
		return fault.Wrap(errBadFormat, fctx.With(ctx))
	}

	source, t, err := image.Decode(r)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	ctx = fctx.WithMeta(ctx, "type", t)

	// re-used across each size
	resizeBuffer := bytes.NewBuffer(nil)

	for _, size := range sizes {
		px, py := sizeMap[size][0], sizeMap[size][1]
		assetFilename := getFilepathForSize(size)

		resized := imaging.Resize(source, px, py, imaging.Lanczos)

		resizeBuffer.Reset()

		if err := png.Encode(resizeBuffer, resized); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		_, err = s.uploader.Upload(ctx, resizeBuffer, int64(resizeBuffer.Len()), assetFilename, asset_upload.Options{})
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	return nil
}

func (s *service) Get(ctx context.Context) (*asset.Asset, io.Reader, error) {
	filename := getFilepathForSize(Size("1200x630"))

	a, r, err := s.downloader.Get(ctx, filename)
	if err != nil {
		return &asset.Asset{
			ID:   xid.New(),
			Name: filename,
			Size: len(defaultBanner),
			MIME: mime.New("image/png"),
		}, bytes.NewReader(defaultBanner), nil
	}

	return a, r, nil
}

func getFilepathForSize(size Size) asset.Filename {
	return asset.NewFilepathFilename(fmt.Sprintf(bannerFileTemplate, size))
}
