package asset

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
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

const assetsSubdirectory = "assets"

type Service interface {
	Upload(ctx context.Context, r io.Reader) (*asset.Asset, error)
	Read(ctx context.Context, path string) (io.Reader, error)
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
		l:            l.With(zap.String("service", "post")),
		rbac:         rbac,
		account_repo: account_repo,
		asset_repo:   asset_repo,
		thread_repo:  thread_repo,
		os:           os,
		address:      cfg.PublicWebAddress,
	}
}

func (s *service) Upload(ctx context.Context, r io.Reader) (*asset.Asset, error) {
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// NOTE: We load the whole file into memory in order to compute a hash first
	// which isn't the most optimal route as it means 5 people uploading a 100MB
	// file to a 512MB server would result in a crash but this can be optimised.
	// There are a few alternatives, one is to upload the whole file now by just
	// streaming it to its destination then computing hashes and resizes another
	// time, another way is by using a rolling hash on the stream during upload.

	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	r = bytes.NewReader(buf)

	mt := mimetype.Detect(buf)

	hash := sha1.Sum(buf)
	assetID := hex.EncodeToString(hash[:])
	slug := fmt.Sprintf("%s-%s", assetID, accountID.String())
	filePath := filepath.Join(assetsSubdirectory, slug)

	if err := s.os.Write(ctx, filePath, r); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	apiPath := path.Join("api/v1/assets", slug)
	url := fmt.Sprintf("%s/%s", s.address, apiPath)
	mime := mt.String()

	if strings.HasPrefix(mime, "image") {
		// TODO: figure out width and height
		fmt.Println("IS AN IMAGE")
	}

	ast, err := s.asset_repo.Add(ctx, accountID, assetID, url, mime, 0, 0)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return ast, nil
}

func (s *service) Read(ctx context.Context, assetID string) (io.Reader, error) {
	path := filepath.Join(assetsSubdirectory, assetID)
	r, err := s.os.Read(ctx, path)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r, nil
}
