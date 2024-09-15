package account_suspension

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	authentication_repo "github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/services/authentication"
)

type Service interface {
	Suspend(ctx context.Context, id account.AccountID) (*account.Account, error)
	Reinstate(ctx context.Context, id account.AccountID) (*account.Account, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l *zap.Logger

	auth_repo      authentication_repo.Repository
	account_writer *account_writer.Writer

	auth_svc *authentication.Manager
}

func New(
	l *zap.Logger,

	auth_repo authentication_repo.Repository,
	account_writer *account_writer.Writer,

	auth_svc *authentication.Manager,
) Service {
	return &service{
		l:              l.With(zap.String("service", "account")),
		auth_repo:      auth_repo,
		account_writer: account_writer,
		auth_svc:       auth_svc,
	}
}

func (s *service) Suspend(ctx context.Context, id account.AccountID) (*account.Account, error) {
	acc, err := s.account_writer.Update(ctx, id, account_writer.SetDeleted(opt.New(time.Now())))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}

func (s *service) Reinstate(ctx context.Context, id account.AccountID) (*account.Account, error) {
	acc, err := s.account_writer.Update(ctx, id, account_writer.SetDeleted(opt.NewEmpty[time.Time]()))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}
