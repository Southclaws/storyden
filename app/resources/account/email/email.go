package email

import (
	"context"
	"crypto/subtle"
	"log/slog"
	"net/mail"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_repo"
	"github.com/Southclaws/storyden/internal/ent"
	account_ent "github.com/Southclaws/storyden/internal/ent/account"
	email_ent "github.com/Southclaws/storyden/internal/ent/email"
)

const (
	VerificationCodeLifespan = time.Hour
	MaxVerificationAttempts  = 5
)

var ErrCodeSuperseded = fault.New("verification code was replaced before it was used")

type Repository struct {
	db          *ent.Client
	accountRepo *account_querier.Querier
}

func New(db *ent.Client, accountRepo *account_querier.Querier) *Repository {
	return &Repository{db: db, accountRepo: accountRepo}
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

	if exists {
		// Already been claimed, by a different account
		if existing.AccountID != nil && *existing.AccountID != xid.ID(accountID) {
			return nil, fault.New("email address already claimed", fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}

		// Already claimed by this account, update the record
		update := r.db.Email.UpdateOne(existing).
			Where(email_ent.EmailAddress(email.Address)).
			SetVerificationCode(code).
			SetVerificationCodeExpiresAt(time.Now().Add(VerificationCodeLifespan)).
			SetVerificationAttempts(0)

		if existing.AccountID == nil {
			update.SetAccountID(xid.ID(accountID))
		}

		updated, err := update.Save(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		r.trySyncVerifiedStatus(ctx, accountID)

		return account.MapEmail(updated), nil
	}

	// Does not exist, create a new email record, bind to owner.

	create := r.db.Email.Create().
		SetAccountID(xid.ID(accountID)).
		SetEmailAddress(email.Address).
		SetVerificationCode(code).
		SetVerificationCodeExpiresAt(time.Now().Add(VerificationCodeLifespan))

	result, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r.trySyncVerifiedStatus(ctx, accountID)

	return account.MapEmail(result), nil
}

func (r *Repository) RotateCode(ctx context.Context, emailAddress mail.Address, code string) error {
	_, err := r.db.Email.Update().
		Where(email_ent.EmailAddress(emailAddress.Address)).
		SetVerificationCode(code).
		SetVerificationCodeExpiresAt(time.Now().Add(VerificationCodeLifespan)).
		SetVerificationAttempts(0).
		Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return nil
}

func (r *Repository) LookupCode(ctx context.Context, emailAddress mail.Address, code string) (*account.Account, bool, error) {
	record, exists, err := r.lookupEmail(ctx, emailAddress)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}
	if !exists {
		return nil, false, nil
	}

	allowed, err := r.recordAttempt(ctx, record.ID)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}
	if !allowed {
		return nil, false, nil
	}

	if subtle.ConstantTimeCompare([]byte(record.VerificationCode), []byte(code)) != 1 {
		return nil, false, nil
	}

	result, err := r.db.Account.
		Query().
		Where(account_ent.HasEmailsWith(email_ent.ID(record.ID))).
		WithEmails().
		WithAuthentication().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, false, nil
		}

		return nil, false, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	acc, err := account.MapRef(result)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, true, nil
}

// the limit is checked inside the update so concurrent guesses cannot all pass a stale read
func (r *Repository) recordAttempt(ctx context.Context, id xid.ID) (bool, error) {
	affected, err := r.db.Email.Update().
		Where(
			email_ent.ID(id),
			email_ent.VerificationAttemptsLT(MaxVerificationAttempts),
			email_ent.VerificationCodeExpiresAtGT(time.Now()),
		).
		AddVerificationAttempts(1).
		Save(ctx)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return affected > 0, nil
}

func (r *Repository) Verify(ctx context.Context, accountID account.AccountID, email mail.Address, code string) error {
	affected, err := r.db.Email.Update().
		Where(
			email_ent.EmailAddress(email.Address),
			email_ent.AccountID(xid.ID(accountID)),
			email_ent.VerificationCode(code),
		).
		SetVerified(true).
		SetVerificationAttempts(0).
		ClearVerificationCodeExpiresAt().
		Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	if affected == 0 {
		return fault.Wrap(ErrCodeSuperseded, fctx.With(ctx), ftag.With(ftag.Unauthenticated))
	}

	r.trySyncVerifiedStatus(ctx, accountID)

	return nil
}

func (r *Repository) SetVerifiedStatus(ctx context.Context, accountID account.AccountID, emailID xid.ID, verified bool) (*account.EmailAddress, error) {
	result, err := r.db.Email.UpdateOneID(emailID).
		Where(email_ent.AccountID(xid.ID(accountID))).
		SetVerified(verified).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := r.syncVerifiedStatus(ctx, accountID); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return account.MapEmail(result), nil
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

func (r *Repository) LookupAccount(ctx context.Context, emailAddress mail.Address) (*account.AccountWithEdges, bool, error) {
	q := r.db.Account.
		Query().
		Where(account_ent.HasEmailsWith(email_ent.EmailAddress(emailAddress.Address))).
		WithEmails().
		WithAuthentication().
		WithTags()

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

	r.trySyncVerifiedStatus(ctx, accountID)

	return nil
}

func (r *Repository) syncVerifiedStatus(ctx context.Context, accountID account.AccountID) error {
	verified, err := r.db.Email.Query().
		Where(
			email_ent.AccountID(xid.ID(accountID)),
			email_ent.Verified(true),
		).
		Exist(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	status := account.VerifiedStatusNone
	if verified {
		status = account.VerifiedStatusVerifiedEmail
	} else {
		acc, err := r.accountRepo.GetByID(ctx, accountID)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
		if acc.VerifiedStatus == account.VerifiedStatusManual {
			return nil
		}
	}

	_, err = r.accountRepo.Update(ctx, accountID, account_repo.SetVerifiedStatus(status))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (r *Repository) trySyncVerifiedStatus(ctx context.Context, accountID account.AccountID) {
	if err := r.syncVerifiedStatus(ctx, accountID); err != nil {
		slog.Error("failed to sync account verified status",
			slog.String("account_id", accountID.String()),
			slog.String("error", err.Error()),
		)
	}
}
