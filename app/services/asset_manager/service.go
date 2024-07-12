package asset_manager

import (
	"context"
	"io"
	"path"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/object"
)

const assetsSubdirectory = "assets"

type Service interface {
	Upload(ctx context.Context, r io.Reader, size int64, name asset.Filename, url string) (*asset.Asset, error)
	Get(ctx context.Context, id asset.Filename) (*asset.Asset, io.Reader, error)
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
		l:            l.With(zap.String("service", "asset")),
		rbac:         rbac,
		account_repo: account_repo,
		asset_repo:   asset_repo,
		thread_repo:  thread_repo,
		os:           os,
		address:      cfg.PublicWebAddress,
	}
}

func (s *service) Upload(ctx context.Context, r io.Reader, size int64, name asset.Filename, url string) (*asset.Asset, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	a, err := s.asset_repo.Add(ctx, accountID, name, url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	path := buildPath(a.Name)

	if err := s.os.Write(ctx, path, r, size); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return a, nil
}

func (s *service) Get(ctx context.Context, id asset.Filename) (*asset.Asset, io.Reader, error) {
	a, err := s.asset_repo.Get(ctx, id)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	path := buildPath(a.Name)
	ctx = fctx.WithMeta(ctx, "path", path, "asset_id", id.String())

	r, size, err := s.os.Read(ctx, path)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	a.Size = int(size)

	return a, r, nil
}

func buildPath(name asset.Filename) string {
	return path.Join(assetsSubdirectory, name.String())
}
