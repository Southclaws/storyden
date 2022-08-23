package thread

import (
	"context"

	"github.com/el-mike/restrict"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/pkg/resources/account"
	"github.com/Southclaws/storyden/pkg/resources/category"
	"github.com/Southclaws/storyden/pkg/resources/thread"
)

type Service interface {
	Create(
		ctx context.Context,
		title string,
		body string,
		authorID account.AccountID,
		categoryID category.CategoryID,
		tags []string,
	) (*thread.Thread, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l    *zap.Logger
	rbac *restrict.AccessManager

	account_repo account.Repository
	thread_repo  thread.Repository
}

func New(
	l *zap.Logger,
	rbac *restrict.AccessManager,

	account_repo account.Repository,
	thread_repo thread.Repository,
) Service {
	return &service{
		l:            l.With(zap.String("service", "thread")),
		rbac:         rbac,
		account_repo: account_repo,
		thread_repo:  thread_repo,
	}
}
