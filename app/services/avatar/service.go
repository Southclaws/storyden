package avatar

import (
	"context"
	"io"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/object"
)

type Service interface {
	Exists(ctx context.Context, accountID account.AccountID) bool
	Set(ctx context.Context, accountID account.AccountID, stream io.Reader) error
	Get(ctx context.Context, accountID account.AccountID) (io.Reader, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l *zap.Logger

	account_repo account.Repository
	storage      object.Storer
}

func New(
	l *zap.Logger,

	account_repo account.Repository,
	storage object.Storer,
) Service {
	return &service{
		l:            l.With(zap.String("service", "avatar")),
		account_repo: account_repo,
		storage:      storage,
	}
}
