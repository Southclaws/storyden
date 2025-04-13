package avatar

import (
	"context"
	"io"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/services/avatar_gen"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
)

type Service interface {
	Exists(ctx context.Context, accountID account.AccountID) bool
	Set(ctx context.Context, accountID account.AccountID, stream io.Reader, size int64) error
	Get(ctx context.Context, accountID account.AccountID) (io.Reader, int64, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	accountQuery *account_querier.Querier
	generator    avatar_gen.AvatarGenerator
	storage      object.Storer
}

func New(
	accountQuery *account_querier.Querier,
	generator avatar_gen.AvatarGenerator,
	storage object.Storer,
) Service {
	return &service{
		accountQuery: accountQuery,
		generator:    generator,
		storage:      storage,
	}
}
