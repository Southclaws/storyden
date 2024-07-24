package email

import (
	"context"
	"net/mail"

	"go.uber.org/fx"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	email_ent "github.com/Southclaws/storyden/internal/ent/email"
)

type EmailRepo struct {
	fx.In

	Ent *ent.Client
}

func (r *EmailRepo) Add(ctx context.Context, accountID account.AccountID, email mail.Address, isAuth bool) (*account.EmailAddress, error) {
	result, err := r.Ent.Email.Create().
		SetAccountID(xid.ID(accountID)).
		SetEmailAddress(email.String()).
		SetIsAuth(isAuth).
		Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return account.MapEmail(result), nil
}

func (r *EmailRepo) Verify(ctx context.Context, accountID account.AccountID, email string) error {
	_, err := r.Ent.Email.Update().
		Where(email_ent.EmailAddress(email)).
		SetVerified(true).
		Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
