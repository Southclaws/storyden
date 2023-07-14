package asset

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/object"
)

const assetsSubdirectory = "assets"

type Service interface {
	Upload(ctx context.Context, pid post.PostID, r io.Reader) (string, error)
	Read(ctx context.Context, path string) (io.Reader, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l    *zap.Logger
	rbac rbac.AccessManager

	account_repo account.Repository
	thread_repo  thread.Repository
	post_repo    post.Repository

	os object.Storer

	address string
}

func New(
	l *zap.Logger,
	rbac rbac.AccessManager,

	account_repo account.Repository,
	thread_repo thread.Repository,
	post_repo post.Repository,

	os object.Storer,
	cfg config.Config,
) Service {
	return &service{
		l:            l.With(zap.String("service", "post")),
		rbac:         rbac,
		account_repo: account_repo,
		thread_repo:  thread_repo,
		post_repo:    post_repo,
		os:           os,
		address:      cfg.PublicWebAddress,
	}
}

func (s *service) Upload(ctx context.Context, pid post.PostID, r io.Reader) (string, error) {
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	assetID := fmt.Sprintf("%s-%s", accountID.String(), xid.New().String())
	path := filepath.Join(assetsSubdirectory, assetID)

	if err := s.os.Write(ctx, path, r); err != nil {
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	url := fmt.Sprintf("%s/api/v1/assets/%s", s.address, assetID)

	_, err = s.post_repo.Update(ctx, pid, post.WithAssets(
	// TODO: Insert asset record and append to post.
	))
	if err != nil {
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	return url, nil
}

func (s *service) Read(ctx context.Context, assetID string) (io.Reader, error) {
	path := filepath.Join(assetsSubdirectory, assetID)
	r, err := s.os.Read(ctx, path)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r, nil
}
