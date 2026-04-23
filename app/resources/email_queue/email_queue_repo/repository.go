package email_queue_repo

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/email_queue"
	"github.com/Southclaws/storyden/internal/ent"
	ent_emailqueue "github.com/Southclaws/storyden/internal/ent/emailqueue"
	entschema "github.com/Southclaws/storyden/internal/ent/schema"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
)

var (
	ErrNotFound        = fault.New("email queue record not found", ftag.With(ftag.NotFound))
	ErrRetryNotAllowed = fault.New("only failed emails can be retried", ftag.With(ftag.InvalidArgument))
)

type Repository struct {
	db *ent.Client
}

func New(db *ent.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateOrGetByID(ctx context.Context, id xid.ID, msg mailer.Message) (*email_queue.Email, error) {
	err := r.db.EmailQueue.Create().
		SetID(id).
		SetRecipientAddress(msg.Address.Address).
		SetRecipientName(msg.Name).
		SetSubject(msg.Subject).
		SetContentPlain(msg.Content.Plain).
		SetContentHTML(msg.Content.HTML).
		SetStatus(ent_emailqueue.StatusPending).
		SetAttempts([]entschema.EmailAttempt{}).
		OnConflictColumns(ent_emailqueue.FieldID).
		Ignore().
		Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	row, err := r.db.EmailQueue.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return email_queue.Map(row)
}

func (r *Repository) Get(ctx context.Context, id email_queue.ID) (*email_queue.Email, error) {
	row, err := r.db.EmailQueue.Get(ctx, xid.ID(id))
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(ErrNotFound, fctx.With(ctx))
		}
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return email_queue.Map(row)
}

func (r *Repository) RetryNow(ctx context.Context, id email_queue.ID, now time.Time) (*email_queue.Email, error) {
	row, err := r.db.EmailQueue.Get(ctx, xid.ID(id))
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(ErrNotFound, fctx.With(ctx))
		}
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	if row.Status != ent_emailqueue.StatusFailed {
		return nil, fault.Wrap(ErrRetryNotAllowed, fctx.With(ctx))
	}

	_, err = r.db.EmailQueue.UpdateOneID(row.ID).
		SetAvailableAt(now).
		SetStatus(ent_emailqueue.StatusPending).
		Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return r.Get(ctx, email_queue.ID(row.ID))
}

func (r *Repository) ClaimNext(ctx context.Context, now time.Time) (*email_queue.Email, bool, error) {
	row, err := r.db.EmailQueue.Query().
		Where(
			ent_emailqueue.AvailableAtLTE(now),
			ent_emailqueue.StatusIn(ent_emailqueue.StatusPending, ent_emailqueue.StatusFailed),
		).
		Order(ent_emailqueue.ByCreatedAt(sql.OrderAsc())).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, false, nil
		}
		return nil, false, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	updated, err := r.db.EmailQueue.Update().
		Where(
			ent_emailqueue.IDEQ(row.ID),
			ent_emailqueue.StatusIn(ent_emailqueue.StatusPending, ent_emailqueue.StatusFailed),
		).
		SetStatus(ent_emailqueue.StatusProcessing).
		Save(ctx)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}
	if updated == 0 {
		return nil, false, nil
	}

	item, err := email_queue.Map(row)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}
	item.Status = email_queue.StatusProcessing

	return item, true, nil
}

func (r *Repository) Claim(ctx context.Context, id email_queue.ID, statuses ...email_queue.Status) (bool, error) {
	allowed := dt.Map(statuses, func(status email_queue.Status) ent_emailqueue.Status {
		return status.Ent()
	})

	updated, err := r.db.EmailQueue.Update().
		Where(
			ent_emailqueue.IDEQ(xid.ID(id)),
			ent_emailqueue.StatusIn(allowed...),
		).
		SetStatus(ent_emailqueue.StatusProcessing).
		Save(ctx)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return updated > 0, nil
}

func (r *Repository) MarkFailed(ctx context.Context, id email_queue.ID, attempts []*email_queue.Attempt, availableAt time.Time) error {
	_, err := r.db.EmailQueue.UpdateOneID(xid.ID(id)).
		SetStatus(ent_emailqueue.StatusFailed).
		SetAttempts(toEntAttempts(attempts)).
		SetAvailableAt(availableAt).
		ClearProcessedAt().
		Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return nil
}

func (r *Repository) MarkSent(ctx context.Context, id email_queue.ID, attempts []*email_queue.Attempt, processedAt time.Time) error {
	_, err := r.db.EmailQueue.UpdateOneID(xid.ID(id)).
		SetStatus(ent_emailqueue.StatusSent).
		SetAttempts(toEntAttempts(attempts)).
		SetProcessedAt(processedAt).
		Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return nil
}

func toEntAttempts(in []*email_queue.Attempt) []entschema.EmailAttempt {
	return dt.Map(in, func(attempt *email_queue.Attempt) entschema.EmailAttempt {
		return entschema.EmailAttempt{
			Timestamp: attempt.Timestamp,
			Status:    attempt.Status.String(),
			Error:     attempt.Error.Ptr(),
		}
	})
}
