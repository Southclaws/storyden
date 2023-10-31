package account

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/authentication"
	authentication_service "github.com/Southclaws/storyden/app/services/authentication"
)

func (s *service) Get(ctx context.Context, id account.AccountID) (*account.Account, error) {
	acc, err := s.account_repo.GetByID(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get account by ID"))
	}

	return acc, nil
}

func (s *service) GetAuthMethods(ctx context.Context, id account.AccountID) ([]*AuthMethod, error) {
	ps := s.auth_svc.Providers()

	mapping := lo.FromEntries(dt.Map(ps, func(p authentication_service.Provider) lo.Entry[string, authentication_service.Provider] {
		return lo.Entry[string, authentication_service.Provider]{
			Key:   p.ID(),
			Value: p,
		}
	}))

	active, err := s.auth_repo.GetAuthMethods(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.Map(active, func(a *authentication.Authentication) *AuthMethod {
		p := mapping[string(a.Service)]

		return &AuthMethod{
			ID:       a.ID.String(),
			Name:     a.Name.Or(p.Name()),
			Provider: p,
		}
	}), nil
}
