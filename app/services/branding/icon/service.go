package icon

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"net/url"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/disintegration/imaging"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/services/asset/asset_download"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/mime"
)

//go:embed embed/default_icon.png
var defaultIcon []byte

const (
	iconFileTemplate = "icon-%s.png"
)

type Size string

var sizes = []Size{
	"512x512",
	"180x180",
	"167x167",
	"152x152",
	"120x120",
	"32x32",
}

var sizeMap = map[Size]int{
	"512x512": 512,
	"180x180": 180,
	"167x167": 167,
	"152x152": 152,
	"120x120": 120,
	"32x32":   32,
}

var errBadFormat = fault.Wrap(fault.New("bad format"), ftag.With(ftag.InvalidArgument))

type Service interface {
	Upload(ctx context.Context, r io.Reader) error
	Get(ctx context.Context, size string) (*asset.Asset, io.Reader, error)
}

type service struct {
	accountQuery *account_querier.Querier
	uploader     *asset_upload.Uploader
	downloader   *asset_download.Downloader

	os object.Storer

	address url.URL
}

func New(
	accountQuery *account_querier.Querier,
	uploader *asset_upload.Uploader,
	downloader *asset_download.Downloader,

	os object.Storer,
	cfg config.Config,
) Service {
	return &service{
		accountQuery: accountQuery,
		uploader:     uploader,
		downloader:   downloader,
		os:           os,
		address:      cfg.PublicWebAddress,
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
		px := sizeMap[size]
		filename := fmt.Sprintf(iconFileTemplate, size)
		assetFilename := asset.NewFilepathFilename(filename)

		resized := imaging.Resize(source, px, px, imaging.Lanczos)

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

func (s *service) Get(ctx context.Context, size string) (*asset.Asset, io.Reader, error) {
	filename := asset.NewFilepathFilename(fmt.Sprintf(iconFileTemplate, size))

	a, r, err := s.downloader.Get(ctx, filename)
	if err != nil {
		a, r = s.getDefaultIcon()
	}

	return a, r, nil
}

func (s *service) getDefaultIcon() (*asset.Asset, io.Reader) {
	defaultIconReader := bytes.NewReader(defaultIcon)
	return &asset.Asset{
		ID:   xid.NilID(),
		Name: asset.NewFilepathFilename("default-icon.png"),
		Size: len(defaultIcon),
		MIME: mime.New("image/png"),
	}, defaultIconReader
}
