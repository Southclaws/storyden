package session

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/account/authentication/access_key"
	"github.com/Southclaws/storyden/app/resources/account/token"
)

type Validator struct {
	tokenRepo token.Repository
	akRepo    *access_key.Repository
}

func NewValidator(tokenRepo token.Repository, akRepo *access_key.Repository) *Validator {
	return &Validator{
		tokenRepo: tokenRepo,
		akRepo:    akRepo,
	}
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

	return WithAccountID(ctx, tv.AccountID), nil
}

func (v *Validator) ValidateAccessKeyToken(ctx context.Context, raw string) (context.Context, error) {
	ak, err := access_key.ParseAccessKeyToken(raw)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// Lookup the access key record by the access key token identifier
	// this is the part before the hash, such as "sdpak_12345678abef".
	ar, err := v.akRepo.LookupByToken(ctx, ak)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// Convert the retrieved authentication record into an access key record.
	ark, err := access_key.AccessKeyRecordFromAuthenticationRecord(*ar)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// Validate the access key record by argon2 verifying the secret + hash.
	_, err = ak.Validate(*ark)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return WithAccessKey(ctx, ar.Account.ID), nil
}
