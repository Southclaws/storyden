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
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/object"
)

const (
	iconRoute          = "api/v1/info/icon"
	iconStoragePath    = "app"
	iconFileTemplate   = "icon-%s.png"
	assetsSubdirectory = "assets"
)

type Size string

var sizes = []Size{
	"32x32",
	"120x120",
	"152x152",
	"167x167",
	"180x180",
	"512x512",
}

var sizeMap = map[Size]int{
	"32x32":   32,
	"120x120": 120,
	"152x152": 152,
	"167x167": 167,
	"180x180": 180,
	"512x512": 512,
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
	asset_repo   asset.Repository
	thread_repo  thread.Repository

	os object.Storer

	address string
}

func New(
	l *zap.Logger,
	rbac rbac.AccessManager,

	account_repo account.Repository,
	asset_repo asset.Repository,
	thread_repo thread.Repository,

	os object.Storer,
	cfg config.Config,
) Service {
	return &service{
		l:            l.With(zap.String("service", "icon")),
		rbac:         rbac,
		account_repo: account_repo,
		asset_repo:   asset_repo,
		thread_repo:  thread_repo,
		os:           os,
		address:      cfg.PublicWebAddress,
	}
}

func (s *service) Upload(ctx context.Context, r io.Reader) error {
	accountID, err := authentication.GetAccountID(ctx)
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

	mt := mimetype.Detect(buf)
	mime := mt.String()
	ctx = fctx.WithMeta(ctx, "mimetype", mime)

	if !strings.HasPrefix(mime, "image") {
		return fault.Wrap(errBadFormat, fctx.With(ctx))
	}

	for _, size := range sizes {
		if err := s.uploadSize(ctx, buf, size); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	return nil
}

func (s *service) Get(ctx context.Context, size string) (*asset.Asset, io.Reader, error) {
	filename := fmt.Sprintf(iconFileTemplate, size)

	a, err := s.asset_repo.Get(ctx, filename)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	filepath := path.Join(iconStoragePath, filename)
	ctx = fctx.WithMeta(ctx, "path", filepath, "asset_id", string(a.ID))

	r, err := s.os.Read(ctx, filepath)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	return a, r, nil
}

func (s *service) uploadSize(ctx context.Context, buf []byte, size Size) error {
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	px := sizeMap[size]

	filename := fmt.Sprintf(iconFileTemplate, size)
	filepath := path.Join(iconStoragePath, filename)

	source, t, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	resized := imaging.Resize(source, px, px, imaging.Lanczos)

	ctx = fctx.WithMeta(ctx, "type", t)

	out := bytes.NewBuffer(nil)

	if err := png.Encode(out, resized); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.os.Write(ctx, filepath, out); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	apiPath := path.Join(iconRoute)
	url := fmt.Sprintf("%s/%s", s.address, apiPath)

	_, err = s.asset_repo.Add(ctx, accountID, filename, url, "image/png", 0, 0)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
