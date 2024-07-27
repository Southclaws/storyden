package email

import (
	"context"
	"net/mail"

	"go.uber.org/fx"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	account_ent "github.com/Southclaws/storyden/internal/ent/account"
	email_ent "github.com/Southclaws/storyden/internal/ent/email"
)

type EmailRepo struct {
	fx.In

	Ent *ent.Client
}

func (r *EmailRepo) Add(ctx context.Context,
	accountID account.AccountID,
	email mail.Address,
	code string,
	authRecordID opt.Optional[xid.ID],
) (*account.EmailAddress, error) {
	create := r.Ent.Email.Create().
		SetAccountID(xid.ID(accountID)).
		SetEmailAddress(email.Address).
		SetVerificationCode(code)

	authRecordID.Call(func(id xid.ID) {
		create.SetAuthenticationRecordID(id)
	})

	result, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return account.MapEmail(result), nil
}

func (r *EmailRepo) GetCode(ctx context.Context, emailAddress mail.Address) (string, error) {
	q := r.Ent.Email.Query().
		Where(email_ent.EmailAddress(emailAddress.Address))

	result, err := q.Only(ctx)
	if err != nil {
		return "", fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return result.VerificationCode, nil
}

func (r *EmailRepo) LookupCode(ctx context.Context, emailAddress mail.Address, code string) (*account.Account, bool, error) {
	q := r.Ent.Account.
		Query().
		Where(
			account_ent.HasEmailsWith(
				email_ent.EmailAddress(emailAddress.Address),
				email_ent.VerificationCode(code),
			),
		).
		WithTags().
		WithEmails().
		WithAuthentication()

	result, err := q.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, false, nil
		}

		return nil, false, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	acc, err := account.FromModel(result)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, true, nil
}

func (r *EmailRepo) Verify(ctx context.Context, accountID account.AccountID, email mail.Address) error {
	_, err := r.Ent.Email.Update().
		Where(email_ent.EmailAddress(email.Address)).
		SetVerified(true).
		Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (r *EmailRepo) LookupAccount(ctx context.Context, emailAddress mail.Address) (*account.Account, bool, error) {
	q := r.Ent.Account.
		Query().
		Where(account_ent.HasEmailsWith(email_ent.EmailAddress(emailAddress.Address))).
		WithTags().
		WithEmails().
		WithAuthentication()

	result, err := q.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, false, nil
		}

		return nil, false, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	acc, err := account.FromModel(result)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, true, nil
}
