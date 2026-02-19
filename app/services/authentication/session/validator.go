package session

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/authentication/access_key"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/account/token"
	"github.com/Southclaws/storyden/app/resources/settings"
)

type Validator struct {
	tokenRepo      token.Repository
	accountQuerier *account_querier.Querier
	roleQuerier    *role_repo.Repository
	akRepo         *access_key.Repository
	settings       *settings.SettingsRepository
}

// NewValidator creates a new session validator with the required dependencies.
func NewValidator(tokenRepo token.Repository, accountQuerier *account_querier.Querier, roleQuerier *role_repo.Repository, akRepo *access_key.Repository, settings *settings.SettingsRepository) *Validator {
	return &Validator{
		tokenRepo:      tokenRepo,
		accountQuerier: accountQuerier,
		roleQuerier:    roleQuerier,
		akRepo:         akRepo,
		settings:       settings,
	}
}

// func (v *Validator) ValidateSession(ctx context.Context, raw string) (context.Context, error) {
// 	var err error
// 	if len(raw) == 20 /* xid.encodedLen */ {
// 		ctx, err = v.ValidateSessionToken(ctx, raw)
// 	} else if len(raw) == access_key.AccessKeyLength {
// 		ctx, err = v.ValidateAccessKeyToken(ctx, raw)
// 	}
// 	if err != nil {
// 		return nil, fault.Wrap(err, fctx.With(ctx))
// 	}

// 	// Hydrate the context with role information

// 	return ctx, nil
// }

func (v *Validator) resolveRolesForAccount(ctx context.Context, acc *account.AccountWithEdges) (role.Roles, error) {
	if acc.Admin {
		return acc.Roles.Roles(), nil
	}

	requiresEmailVerification, err := v.installationRequiresEmailVerification(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if requiresEmailVerification && acc.VerifiedStatus != account.VerifiedStatusVerifiedEmail {
		guestRole, err := v.roleQuerier.GetGuestRole(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		return role.Roles{guestRole}, nil
	}

	return acc.Roles.Roles(), nil
}

func (v *Validator) installationRequiresEmailVerification(ctx context.Context) (bool, error) {
	s, err := v.settings.Get(ctx)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}

	authMode := s.AuthenticationMode.Or(authentication.ModeHandle)

	return authMode == authentication.ModeEmail, nil
}

// ValidateSessionToken validates a session token and returns a context with account info.
func (v *Validator) ValidateSessionToken(ctx context.Context, raw string) (context.Context, error) {
	t, err := token.FromString(raw)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tv, err := v.tokenRepo.Validate(ctx, t)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := v.accountQuerier.GetByID(ctx, tv.AccountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roles, err := v.resolveRolesForAccount(ctx, acc)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return WithAccountAndToken(ctx, acc.Account, roles, raw), nil
}

// ValidateAccessKeyToken validates an access key token and returns a context with account info.
func (v *Validator) ValidateAccessKeyToken(ctx context.Context, raw string) (context.Context, error) {
	ak, err := access_key.ParseAccessKeyToken(raw)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	ar, err := v.akRepo.LookupByToken(ctx, ak)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	ark, err := access_key.AccessKeyRecordFromAuthenticationRecord(*ar)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, err = ak.Validate(*ark)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := v.accountQuerier.GetByID(ctx, ar.Account.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roles, err := v.resolveRolesForAccount(ctx, acc)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return WithAccessKey(ctx, acc.Account, roles), nil
}

// WithUnauthenticatedRoles returns a context with guest role permissions.
func (v *Validator) WithUnauthenticatedRoles(ctx context.Context) (context.Context, error) {
	guestRole, err := v.roleQuerier.GetGuestRole(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return WithGuest(ctx, role.Roles{guestRole}), nil
}
