package account

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/rbac"
)

type Service interface {
	Get(ctx context.Context, id account.AccountID) (*account.Account, error)
	Update(ctx context.Context, id account.AccountID, params Partial) (*account.Account, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l    *zap.Logger
	rbac rbac.AccessManager

	account_repo account.Repository
}

func New(
	l *zap.Logger,
	rbac rbac.AccessManager,

	account_repo account.Repository,
) Service {
	return &service{
		l:            l.With(zap.String("service", "account")),
		rbac:         rbac,
		account_repo: account_repo,
	}
}
