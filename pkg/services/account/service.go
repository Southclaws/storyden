package account

import (
	"context"

	"github.com/el-mike/restrict"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/pkg/resources/account"
)

type Service interface {
	Get(ctx context.Context, id account.AccountID) (*account.Account, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l    *zap.Logger
	rbac *restrict.AccessManager

	account_repo account.Repository
}

func New(
	l *zap.Logger,
	rbac *restrict.AccessManager,

	account_repo account.Repository,
) Service {
	return &service{
		l:            l.With(zap.String("service", "account")),
		rbac:         rbac,
		account_repo: account_repo,
	}
}
