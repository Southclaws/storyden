package account_suspension

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	authentication_repo "github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Service interface {
	Suspend(ctx context.Context, id account.AccountID) (*account.AccountWithEdges, error)
	Reinstate(ctx context.Context, id account.AccountID) (*account.AccountWithEdges, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	auth_repo      authentication_repo.Repository
	account_writer *account_writer.Writer
	bus            *pubsub.Bus

	auth_svc *authentication.Manager
}

func New(
	auth_repo authentication_repo.Repository,
	account_writer *account_writer.Writer,
	bus *pubsub.Bus,

	auth_svc *authentication.Manager,
) Service {
	return &service{
		auth_repo:      auth_repo,
		account_writer: account_writer,
		bus:            bus,
		auth_svc:       auth_svc,
	}
}

func (s *service) Suspend(ctx context.Context, id account.AccountID) (*account.AccountWithEdges, error) {
	acc, err := s.account_writer.Update(ctx, id, account_writer.SetDeleted(opt.New(time.Now())))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &rpc.EventAccountSuspended{
		ID: id,
	})

	return acc, nil
}

func (s *service) Reinstate(ctx context.Context, id account.AccountID) (*account.AccountWithEdges, error) {
	acc, err := s.account_writer.Update(ctx, id, account_writer.SetDeleted(opt.NewEmpty[time.Time]()))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &rpc.EventAccountUnsuspended{
		ID: id,
	})

	return acc, nil
}
