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
	"github.com/Southclaws/storyden/app/resources/account/role/role_querier"
	"github.com/Southclaws/storyden/app/resources/account/token"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/kv"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/spanner"
)

type Validator struct {
	ins            spanner.Instrumentation
	tokenRepo      token.Repository
	accountQuerier *account_querier.Querier
	roleQuerier    *role_querier.Querier
	akRepo         *access_key.Repository
	settings       *settings.SettingsRepository
}

// NewValidator creates a new session validator with the required dependencies.
func NewValidator(ins spanner.Builder, tokenRepo token.Repository, accountQuerier *account_querier.Querier, roleQuerier *role_querier.Querier, akRepo *access_key.Repository, settings *settings.SettingsRepository) *Validator {
	return &Validator{
		ins:            ins.Build(),
		tokenRepo:      tokenRepo,
		accountQuerier: accountQuerier,
		roleQuerier:    roleQuerier,
		akRepo:         akRepo,
		settings:       settings,
	}
}

func (v *Validator) resolveRolesForAccount(ctx context.Context, acc *account.Account) (role.Roles, error) {
	ctx, span := v.ins.Instrument(ctx, kv.String("account_id", acc.ID.String()))
	defer span.End()

	if acc.Admin {
		span.Event("admin account bypassed email verification role gate")
		return acc.Roles.Roles(), nil
	}

	requiresEmailVerification, err := v.installationRequiresEmailVerification(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if requiresEmailVerification && acc.VerifiedStatus != account.VerifiedStatusVerifiedEmail {
		span.Event("account forced to guest role due to unverified email")
		guestRole, err := v.roleQuerier.GetGuestRole(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		return role.Roles{guestRole}, nil
	}

	return acc.Roles.Roles(), nil
}

func (v *Validator) installationRequiresEmailVerification(ctx context.Context) (bool, error) {
	ctx, span := v.ins.Instrument(ctx)
	defer span.End()

	s, err := v.settings.Get(ctx)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}

	authMode := s.AuthenticationMode.Or(authentication.ModeHandle)

	return authMode == authentication.ModeEmail, nil
}

// ValidateSessionToken validates a session token and returns a context with account info.
func (v *Validator) ValidateSessionToken(ctx context.Context, raw string) (context.Context, error) {
	ctx, span := v.ins.Instrument(ctx)
	defer span.End()

	t, err := token.FromString(raw)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tv, err := v.tokenRepo.Validate(ctx, t)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := v.accountQuerier.GetRefByID(ctx, tv.AccountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roles, err := v.resolveRolesForAccount(ctx, acc)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return WithAccountAndToken(ctx, *acc, roles, raw), nil
}

// ValidateAccessKeyToken validates an access key token and returns a context with account info.
func (v *Validator) ValidateAccessKeyToken(ctx context.Context, raw string) (context.Context, error) {
	ctx, span := v.ins.Instrument(ctx)
	defer span.End()

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

	acc, err := v.accountQuerier.GetRefByID(ctx, ar.Account.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roles, err := v.resolveRolesForAccount(ctx, acc)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return WithAccessKey(ctx, *acc, roles), nil
}

// WithUnauthenticatedRoles returns a context with guest role permissions.
func (v *Validator) WithUnauthenticatedRoles(ctx context.Context) (context.Context, error) {
	ctx, span := v.ins.Instrument(ctx)
	defer span.End()

	guestRole, err := v.roleQuerier.GetGuestRole(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return WithGuest(ctx, role.Roles{guestRole}), nil
}
