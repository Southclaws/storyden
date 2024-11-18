package email

import (
	"context"
	"net/mail"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	account_ent "github.com/Southclaws/storyden/internal/ent/account"
	email_ent "github.com/Southclaws/storyden/internal/ent/email"
)

type Repository struct {
	db *ent.Client
}

func New(db *ent.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Add(ctx context.Context,
	accountID account.AccountID,
	email mail.Address,
	code string,
) (*account.EmailAddress, error) {
	// Check for unclaimed but existing email addresses. Email addresses may be
	// added by admins or integrations for newsletters without being associated
	// with an account yet. As long as the email address becomes verified, good.
	existing, exists, err := r.lookupEmail(ctx, email)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if exists && existing.AccountID != nil {
		return nil, fault.New("email address already claimed", fctx.With(ctx), ftag.With(ftag.AlreadyExists))
	}

	// Exists already, but not claimed, update to new owner.
	if exists {
		updated, err := r.db.Email.UpdateOne(existing).
			Where(email_ent.EmailAddress(email.Address)).
			SetAccountID(xid.ID(accountID)).
			SetVerificationCode(code).
			Save(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return account.MapEmail(updated), nil
	}

	// Does not exist, create a new email record, bind to owner.

	create := r.db.Email.Create().
		SetAccountID(xid.ID(accountID)).
		SetEmailAddress(email.Address).
		SetVerificationCode(code)

	result, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return account.MapEmail(result), nil
}

func (r *Repository) GetCode(ctx context.Context, emailAddress mail.Address) (string, error) {
	q := r.db.Email.Query().
		Where(email_ent.EmailAddress(emailAddress.Address))

	result, err := q.Only(ctx)
	if err != nil {
		return "", fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return result.VerificationCode, nil
}

func (r *Repository) LookupCode(ctx context.Context, emailAddress mail.Address, code string) (*account.Account, bool, error) {
	q := r.db.Account.
		Query().
		Where(
			account_ent.HasEmailsWith(
				email_ent.EmailAddress(emailAddress.Address),
				email_ent.VerificationCode(code),
			),
		).
		WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() }).
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

	acc, err := account.MapAccount(result)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, true, nil
}

func (r *Repository) Verify(ctx context.Context, accountID account.AccountID, email mail.Address) error {
	_, err := r.db.Email.Update().
		Where(email_ent.EmailAddress(email.Address)).
		SetVerified(true).
		Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (r *Repository) lookupEmail(ctx context.Context, emailAddress mail.Address) (*ent.Email, bool, error) {
	result, err := r.db.Email.Query().
		Where(email_ent.EmailAddress(emailAddress.Address)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, false, nil
		}

		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	return result, true, nil
}

func (r *Repository) LookupAccount(ctx context.Context, emailAddress mail.Address) (*account.Account, bool, error) {
	q := r.db.Account.
		Query().
		Where(account_ent.HasEmailsWith(email_ent.EmailAddress(emailAddress.Address))).
		WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() }).
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

	acc, err := account.MapAccount(result)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, true, nil
}

func (r *Repository) Remove(ctx context.Context, accountID account.AccountID, emailID xid.ID) error {
	_, err := r.db.Email.Delete().
		Where(
			email_ent.ID(emailID),
			email_ent.HasAccountWith(account_ent.ID(xid.ID(accountID))),
		).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
