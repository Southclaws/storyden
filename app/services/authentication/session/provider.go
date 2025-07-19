package session

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
)

type Provider struct {
	querier *account_querier.Querier
}

func New(querier *account_querier.Querier) *Provider {
	return &Provider{querier: querier}
}

func (s *Provider) Account(ctx context.Context) (*account.AccountWithEdges, error) {
	id, err := GetAccountID(ctx)
	if err != nil {
		return nil, err
	}

	acc, err := s.querier.GetByID(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}

func (s *Provider) AccountMaybe(ctx context.Context) opt.Optional[account.AccountWithEdges] {
	id, err := GetAccountID(ctx)
	if err != nil {
		return opt.NewEmpty[account.AccountWithEdges]()
	}

	acc, err := s.querier.GetByID(ctx, id)
	if err != nil {
		return opt.NewEmpty[account.AccountWithEdges]()
	}

	return opt.New(*acc)
}
