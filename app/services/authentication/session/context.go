package session

import (
	"context"
	"errors"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

var ErrNoAccountInContext = errors.New("no account in context")

var contextKey = struct{}{}

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

func GetSessionFromMessage[T any](ctx context.Context, msg *pubsub.Message[T]) context.Context {
	actorID, ok := msg.ActorID.Get()
	if !ok {
		return ctx
	}

	return WithAccountID(ctx, account.AccountID(actorID))
}
