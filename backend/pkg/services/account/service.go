package account

import (
	"context"

	"github.com/el-mike/restrict"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/backend/pkg/resources/account"
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

func New() Service {
	return &service{}
}
