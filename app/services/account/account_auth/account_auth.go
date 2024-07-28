package account_auth

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/samber/lo"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	authentication_repo "github.com/Southclaws/storyden/app/resources/account/authentication"
	authentication_service "github.com/Southclaws/storyden/app/services/authentication"
)

type Manager struct {
	fx.In

	AuthService *authentication_service.Manager
	AuthRepo    authentication_repo.Repository
}

type AuthMethod struct {
	Instance authentication_repo.Authentication
	Provider authentication_service.Provider
}

func (m *Manager) GetAuthMethods(ctx context.Context, id account.AccountID) ([]*AuthMethod, error) {
	ps := m.AuthService.Providers()

	mapping := lo.FromEntries(dt.Map(ps, func(p authentication_service.Provider) lo.Entry[string, authentication_service.Provider] {
		return lo.Entry[string, authentication_service.Provider]{
			Key:   p.ID(),
			Value: p,
		}
	}))

	active, err := m.AuthRepo.GetAuthMethods(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.Map(active, func(a *authentication.Authentication) *AuthMethod {
		p := mapping[string(a.Service)]

		return &AuthMethod{
			Instance: *a,
			Provider: p,
		}
	}), nil
}

func (m *Manager) DeleteAuthMethod(ctx context.Context, id account.AccountID, aid authentication_repo.ID) error {
	_, err := m.AuthRepo.DeleteByID(ctx, id, aid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
