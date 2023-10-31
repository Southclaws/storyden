package account

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	authentication_repo "github.com/Southclaws/storyden/app/resources/authentication"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication"
)

type Service interface {
	Get(ctx context.Context, id account.AccountID) (*account.Account, error)
	GetAuthMethods(ctx context.Context, id account.AccountID) ([]*AuthMethod, error)
	DeleteAuthMethod(ctx context.Context, id account.AccountID, aid authentication_repo.ID) error
	Update(ctx context.Context, id account.AccountID, params Partial) (*account.Account, error)
}

type AuthMethod struct {
	Instance authentication_repo.Authentication
	Provider authentication.Provider
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l    *zap.Logger
	rbac rbac.AccessManager

	auth_repo    authentication_repo.Repository
	account_repo account.Repository

	auth_svc *authentication.Manager
}

func New(
	l *zap.Logger,
	rbac rbac.AccessManager,

	auth_repo authentication_repo.Repository,
	account_repo account.Repository,

	auth_svc *authentication.Manager,
) Service {
	return &service{
		l:            l.With(zap.String("service", "account")),
		rbac:         rbac,
		auth_repo:    auth_repo,
		account_repo: account_repo,
		auth_svc:     auth_svc,
	}
}
