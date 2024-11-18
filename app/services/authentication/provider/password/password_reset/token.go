package password_reset

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/infrastructure/endec"
	"github.com/rs/xid"
)

var errMalformedToken = fault.New("missing account_id in token")

type TokenProvider struct {
	endec endec.EncrypterDecrypter
}

func NewTokenProvider(endec endec.EncrypterDecrypter) *TokenProvider {
	return &TokenProvider{
		endec: endec,
	}
}

const (
	accountIDKey       = "account_id"
	resetTokenLifespan = time.Hour
)

func (r *TokenProvider) GetResetToken(ctx context.Context, accountID account.AccountID) (string, error) {
	claims := endec.Claims{
		accountIDKey: accountID.String(),
	}

	token, err := r.endec.Encrypt(claims, resetTokenLifespan)
	if err != nil {
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	return token, nil
}

func (r *TokenProvider) Validate(ctx context.Context, tokenString string) (account.AccountID, error) {
	token, err := r.endec.Decrypt(tokenString)
	if err != nil {
		return account.AccountID{}, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to decrypt token"))
	}

	rawID, ok := token[accountIDKey]
	if !ok {
		return account.AccountID{}, fault.Wrap(errMalformedToken, fctx.With(ctx), fmsg.With("failed to find account_id in token"))
	}

	stringID, ok := rawID.(string)
	if !ok {
		return account.AccountID{}, fault.Wrap(errMalformedToken, fctx.With(ctx), fmsg.With("failed to convert account_id in token"))
	}

	accountID, err := xid.FromString(stringID)
	if err != nil {
		return account.AccountID{}, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to parse account_id in token"))
	}

	return account.AccountID(accountID), nil
}
