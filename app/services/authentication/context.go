package authentication

import (
	"context"
	"errors"

	"github.com/Southclaws/fault/errtag"

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

	return account.AccountID{}, errtag.Wrap(ErrNoAccountInContext, errtag.Unauthenticated{})
}
