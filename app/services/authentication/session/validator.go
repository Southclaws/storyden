package session

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/account/token"
)

type Validator struct {
	tokenRepo token.Repository
}

func NewValidator(tokenRepo token.Repository) *Validator {
	return &Validator{
		tokenRepo: tokenRepo,
	}
}

func (v *Validator) Validate(ctx context.Context, raw string) (context.Context, error) {
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
