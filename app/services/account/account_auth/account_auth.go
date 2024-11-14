package account_auth

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	authentication_repo "github.com/Southclaws/storyden/app/resources/account/authentication"
	authentication_service "github.com/Southclaws/storyden/app/services/authentication"
)

type Manager struct {
	authService *authentication_service.Manager
	authRepo    authentication_repo.Repository
}

func New(
	authService *authentication_service.Manager,
	authRepo authentication_repo.Repository,
) *Manager {
	return &Manager{
		authService: authService,
		authRepo:    authRepo,
	}
}

type AuthMethod struct {
	Instance authentication_repo.Authentication
	Provider authentication_service.Provider
}

func (m *Manager) GetAuthMethods(ctx context.Context, id account.AccountID) ([]*AuthMethod, error) {
	ps, err := m.authService.GetProviderList(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mapping := lo.FromEntries(dt.Map(ps, func(p authentication_service.Provider) lo.Entry[authentication.Service, authentication_service.Provider] {
		return lo.Entry[authentication.Service, authentication_service.Provider]{
			Key:   p.Provides(),
			Value: p,
		}
	}))

	active, err := m.authRepo.GetAuthMethods(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.Map(active, func(a *authentication.Authentication) *AuthMethod {
		p := mapping[a.Service]

		return &AuthMethod{
			Instance: *a,
			Provider: p,
		}
	}), nil
}

func (m *Manager) DeleteAuthMethod(ctx context.Context, id account.AccountID, aid authentication_repo.ID) error {
	_, err := m.authRepo.DeleteByID(ctx, id, aid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
