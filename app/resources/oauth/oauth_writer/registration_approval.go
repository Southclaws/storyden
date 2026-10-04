package oauth_writer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/oauthregistrationapproval"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

var ErrRegistrationApprovalUnavailable = fault.New("registration approval is unavailable", ftag.With(ftag.PermissionDenied))

type RegistrationApprovalCreate struct {
	RegistrationCodeHash string
	VerificationCodeHash string
	VerificationCode     string
	Metadata             json.RawMessage
	ExpiresAt            time.Time
	NextPollAt           time.Time
	PollIntervalSeconds  int
}

func (w *Writer) CreateRegistrationApproval(ctx context.Context, input RegistrationApprovalCreate) (*oauth.RegistrationApproval, error) {
	row, err := w.db.OAuthRegistrationApproval.Create().
		SetRegistrationCodeHash(input.RegistrationCodeHash).
		SetVerificationCodeHash(input.VerificationCodeHash).
		SetVerificationCodeDisplay(input.VerificationCode).
		SetMetadata(input.Metadata).SetExpiresAt(input.ExpiresAt).
		SetNextPollAt(input.NextPollAt).SetPollIntervalSeconds(input.PollIntervalSeconds).Save(ctx)
	if err != nil {
		return nil, wrapWriteError(ctx, err)
	}

	return oauth.MapRegistrationApproval(row), nil
}

func pendingRegistrationApprovals(now time.Time) []predicate.OAuthRegistrationApproval {
	return []predicate.OAuthRegistrationApproval{
		oauthregistrationapproval.ExpiresAtGT(now),
		oauthregistrationapproval.ApprovedAtIsNil(), oauthregistrationapproval.DeniedAtIsNil(),
		oauthregistrationapproval.CancelledAtIsNil(), oauthregistrationapproval.ConsumedAtIsNil(),
	}
}

func (w *Writer) DecideRegistrationApproval(ctx context.Context, id oauth.RegistrationApprovalID, approver account.AccountID, approved bool, now time.Time) (bool, error) {
	update := w.db.OAuthRegistrationApproval.Update().Where(oauthregistrationapproval.ID(id.XID())).Where(pendingRegistrationApprovals(now)...)
	if approved {
		update.SetApprovedByAccountID(xid.ID(approver)).SetApprovedAt(now)
	} else {
		update.SetDeniedAt(now)
	}

	n, err := update.Save(ctx)
	if err != nil {
		return false, wrapWriteError(ctx, err)
	}

	return n == 1, nil
}

func (w *Writer) CancelRegistrationApproval(ctx context.Context, id oauth.RegistrationApprovalID, now time.Time) (bool, error) {
	n, err := w.db.OAuthRegistrationApproval.Update().Where(oauthregistrationapproval.ID(id.XID())).Where(pendingRegistrationApprovals(now)...).SetCancelledAt(now).Save(ctx)
	if err != nil {
		return false, wrapWriteError(ctx, err)
	}

	return n == 1, nil
}

func (w *Writer) RecordRegistrationApprovalPoll(ctx context.Context, rec *oauth.RegistrationApproval, nextPoll time.Time, interval int) (bool, error) {
	n, err := w.db.OAuthRegistrationApproval.Update().Where(
		oauthregistrationapproval.ID(rec.ID.XID()), oauthregistrationapproval.NextPollAt(rec.NextPollAt),
		oauthregistrationapproval.PollIntervalSeconds(rec.PollIntervalSeconds),
		oauthregistrationapproval.ConsumedAtIsNil(), oauthregistrationapproval.CancelledAtIsNil(),
		oauthregistrationapproval.DeniedAtIsNil(), oauthregistrationapproval.ExpiresAtGT(time.Now()),
	).SetNextPollAt(nextPoll).SetPollIntervalSeconds(interval).Save(ctx)
	if err != nil {
		return false, wrapWriteError(ctx, err)
	}

	return n == 1, nil
}

func consumeRegistrationApproval(ctx context.Context, tx *ent.Tx, rec *oauth.RegistrationApproval) error {
	approver, ok := rec.ApprovedByAccountID.Get()
	if !ok {
		return ErrRegistrationApprovalUnavailable
	}

	n, err := tx.OAuthRegistrationApproval.Update().Where(
		oauthregistrationapproval.ID(rec.ID.XID()), oauthregistrationapproval.ExpiresAtGT(time.Now()),
		oauthregistrationapproval.ApprovedByAccountID(xid.ID(approver)), oauthregistrationapproval.ApprovedAtNotNil(),
		oauthregistrationapproval.ConsumedAtIsNil(), oauthregistrationapproval.CancelledAtIsNil(), oauthregistrationapproval.DeniedAtIsNil(),
	).SetConsumedAt(time.Now()).Save(ctx)
	if err != nil {
		return err
	}

	if n != 1 {
		return ErrRegistrationApprovalUnavailable
	}

	return nil
}

func (w *Writer) DeleteExpiredRegistrationApprovals(ctx context.Context, before time.Time) (int, error) {
	n, err := w.db.OAuthRegistrationApproval.Delete().Where(oauthregistrationapproval.ExpiresAtLT(before)).Exec(ctx)
	if err != nil {
		return 0, wrapWriteError(ctx, err)
	}

	return n, nil
}

func (w *Writer) DecideAllRegistrationApprovals(ctx context.Context, approver account.AccountID, approved bool, before, now time.Time) (int, error) {
	update := w.db.OAuthRegistrationApproval.Update().Where(pendingRegistrationApprovals(now)...).Where(oauthregistrationapproval.CreatedAtLTE(before))
	if approved {
		update.SetApprovedByAccountID(xid.ID(approver)).SetApprovedAt(now)
	} else {
		update.SetDeniedAt(now)
	}

	count, err := update.Save(ctx)
	if err != nil {
		return 0, wrapWriteError(ctx, err)
	}

	return count, nil
}
