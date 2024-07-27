package icon

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
	"path"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/asset_manager"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
)

const (
	iconRoute        = "api/v1/info/icon"
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

var (
	errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))
	errBadFormat     = fault.Wrap(fault.New("bad format"), ftag.With(ftag.InvalidArgument))
)

type Service interface {
	Upload(ctx context.Context, r io.Reader) error
	Get(ctx context.Context, size string) (*asset.Asset, io.Reader, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l    *zap.Logger
	rbac rbac.AccessManager

	account_repo account.Repository
	am           asset_manager.Service
	thread_repo  thread.Repository

	os object.Storer

	address string
}

func New(
	l *zap.Logger,
	rbac rbac.AccessManager,

	account_repo account.Repository,
	am asset_manager.Service,
	thread_repo thread.Repository,

	os object.Storer,
	cfg config.Config,
) Service {
	return &service{
		l:            l.With(zap.String("service", "icon")),
		rbac:         rbac,
		account_repo: account_repo,
		am:           am,
		thread_repo:  thread_repo,
		os:           os,
		address:      cfg.PublicWebAddress,
	}
}

func (s *service) Upload(ctx context.Context, r io.Reader) error {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.account_repo.GetByID(ctx, accountID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if !acc.Admin {
		return fault.Wrap(errNotAuthorised, fctx.With(ctx))
	}

	return s.uploadSizes(ctx, r, sizes)
}

func (s *service) uploadSizes(ctx context.Context, r io.Reader, sizes []Size) error {
	// NOTE: We load the whole file into memory in order to compute a hash first
	// which isn't the most optimal route as it means 5 people uploading a 100MB
	// file to a 512MB server would result in a crash but this can be optimised.
	// There are a few alternatives, one is to upload the whole file now by just
	// streaming it to its destination then computing hashes and resizes another
	// time, another way is by using a rolling hash on the stream during upload.

	buf, err := io.ReadAll(r)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	// we read r already, but image.Decode needs a reader, so make a new one
	bufferReader := bytes.NewReader(buf)

	mt := mimetype.Detect(buf)
	mime := mt.String()
	ctx = fctx.WithMeta(ctx, "mimetype", mime)

	if !strings.HasPrefix(mime, "image") {
		return fault.Wrap(errBadFormat, fctx.With(ctx))
	}

	source, t, err := image.Decode(bufferReader)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	ctx = fctx.WithMeta(ctx, "type", t)

	// re-used across each size
	resizeBuffer := bytes.NewBuffer(buf)

	for _, size := range sizes {
		px := sizeMap[size]
		filename := fmt.Sprintf(iconFileTemplate, size)
		assetFilename := asset.NewFilepathFilename(filename)

		resized := imaging.Resize(source, px, px, imaging.Lanczos)

		resizeBuffer.Reset()

		if err := png.Encode(resizeBuffer, resized); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		apiPath := path.Join(iconRoute)
		url := fmt.Sprintf("%s/%s", s.address, apiPath)

		_, err = s.am.Upload(ctx, resizeBuffer, int64(resizeBuffer.Len()), assetFilename, url)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	return nil
}

func (s *service) Get(ctx context.Context, size string) (*asset.Asset, io.Reader, error) {
	filename := asset.NewFilepathFilename(fmt.Sprintf(iconFileTemplate, size))

	a, r, err := s.am.Get(ctx, filename)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	return a, r, nil
}
