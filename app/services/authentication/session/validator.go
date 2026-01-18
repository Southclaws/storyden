package session

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/authentication/access_key"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/role_querier"
	"github.com/Southclaws/storyden/app/resources/account/token"
)

type Validator struct {
	tokenRepo      token.Repository
	accountQuerier *account_querier.Querier
	roleQuerier    *role_querier.Querier
	akRepo         *access_key.Repository
}

func NewValidator(tokenRepo token.Repository, accountQuerier *account_querier.Querier, roleQuerier *role_querier.Querier, akRepo *access_key.Repository) *Validator {
	return &Validator{
		tokenRepo:      tokenRepo,
		accountQuerier: accountQuerier,
		roleQuerier:    roleQuerier,
		akRepo:         akRepo,
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

func (v *Validator) resolveRolesForAccount(ctx context.Context, acc *account.AccountWithRoles) (role.Roles, error) {
	if acc.VerifiedStatus != account.VerifiedStatusVerifiedEmail {
		guestRole, err := v.roleQuerier.GetGuestRole(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		return role.Roles{guestRole}, nil
	}
	return acc.Roles.Roles(), nil
}

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

func (v *Validator) WithUnauthenticatedRoles(ctx context.Context) (context.Context, error) {
	guestRole, err := v.roleQuerier.GetGuestRole(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return WithGuest(ctx, role.Roles{guestRole}), nil
}
