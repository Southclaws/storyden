package session

import (
	"context"
	"errors"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
)

var ErrNoAccountInContext = errors.New("no account in context")

var contextKey = struct{}{}

// WithAccountID stores the ID of the account making the request.
func WithAccountID(ctx context.Context, u account.AccountID) context.Context {
	return context.WithValue(ctx, contextKey, u)
}

// GetAccountID pulls out an account ID associated with the call.
func GetAccountID(ctx context.Context) (account.AccountID, error) {
	if auth, ok := ctx.Value(contextKey).(account.AccountID); ok {
		return auth, nil
	}

	return account.AccountID{}, fault.Wrap(ErrNoAccountInContext, fctx.With(ctx), ftag.With(ftag.Unauthenticated))
}

func GetOptAccountID(ctx context.Context) opt.Optional[account.AccountID] {
	if auth, ok := ctx.Value(contextKey).(account.AccountID); ok {
		return opt.New(auth)
	}

	return opt.NewEmpty[account.AccountID]()
}
