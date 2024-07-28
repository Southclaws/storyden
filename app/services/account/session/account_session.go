package session

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

type SessionProvider struct {
	fx.In
	Repo account_querier.Querier
}

func (s *SessionProvider) Account(ctx context.Context) (*account.Account, error) {
	id, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, err
	}

	acc, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}

func (s *SessionProvider) AccountOpt(ctx context.Context) opt.Optional[account.Account] {
	id, err := session.GetAccountID(ctx)
	if err != nil {
		return opt.NewEmpty[account.Account]()
	}

	acc, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		return opt.NewEmpty[account.Account]()
	}

	return opt.New(*acc)
}
