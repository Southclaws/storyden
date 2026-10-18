package background_job

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Southclaws/dt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
	ent_backgroundjob "github.com/Southclaws/storyden/internal/ent/backgroundjob"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

var ErrLeaseLost = errors.New("background job lease lost")

const LeaseDuration = 2 * time.Minute

type Repository struct {
	db *ent.Client
}

func New(db *ent.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Ensure(ctx context.Context, id xid.ID, kind, key string, payload any) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return r.db.BackgroundJob.
		Create().
		SetID(id).
		SetKind(kind).
		SetKey(key).
		SetPayload(encoded).
		OnConflictColumns(ent_backgroundjob.FieldKind, ent_backgroundjob.FieldKey).
		Ignore().
		Exec(ctx)
}

func (r *Repository) Get[T any](ctx context.Context, id xid.ID) (*Job[T], error) {
	row, err := r.db.BackgroundJob.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return Map[T](row)
}

func (r *Repository) Candidates(ctx context.Context, kind string, limit int) ([]*Job[json.RawMessage], error) {
	now := time.Now()

	rows, err := r.db.BackgroundJob.Query().
		Where(
			ent_backgroundjob.KindEQ(kind),
			ent_backgroundjob.StatusIn(ent_backgroundjob.StatusPending, ent_backgroundjob.StatusProcessing),
			ent_backgroundjob.AvailableAtLTE(now),
			ent_backgroundjob.DispatchUntilLTE(now),
			ent_backgroundjob.LeaseUntilLTE(now),
		).
		Order(ent.Asc(ent_backgroundjob.FieldCreatedAt)).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return dt.MapErr(rows, Map[json.RawMessage])
}

func (r *Repository) Dispatch(ctx context.Context, id xid.ID) (bool, error) {
	now := time.Now()
	n, err := r.db.BackgroundJob.Update().
		Where(
			ent_backgroundjob.ID(id),
			ent_backgroundjob.StatusIn(ent_backgroundjob.StatusPending, ent_backgroundjob.StatusProcessing),
			ent_backgroundjob.AvailableAtLTE(now),
			ent_backgroundjob.DispatchUntilLTE(now),
			ent_backgroundjob.LeaseUntilLTE(now),
		).
		SetDispatchUntil(now.Add(LeaseDuration)).
		Save(ctx)

	return n == 1, err
}

func (r *Repository) Claim(ctx context.Context, id xid.ID, maxAttempts int) (*Job[json.RawMessage], error) {
	now := time.Now()
	if err := r.db.BackgroundJob.Update().
		Where(
			ent_backgroundjob.ID(id),
			ent_backgroundjob.StatusIn(ent_backgroundjob.StatusPending, ent_backgroundjob.StatusProcessing),
			ent_backgroundjob.LeaseUntilLTE(now),
			ent_backgroundjob.AttemptsGTE(maxAttempts),
		).
		SetStatus(ent_backgroundjob.StatusFailed).
		SetReason("retry limit reached after lease expiry").
		Exec(ctx); err != nil {
		return nil, err
	}

	token := xid.New().String()
	n, err := r.db.BackgroundJob.Update().
		Where(
			ent_backgroundjob.ID(id),
			ent_backgroundjob.StatusIn(ent_backgroundjob.StatusPending, ent_backgroundjob.StatusProcessing),
			ent_backgroundjob.LeaseUntilLTE(now),
			ent_backgroundjob.AvailableAtLTE(now),
			ent_backgroundjob.AttemptsLT(maxAttempts),
		).
		SetLeaseToken(token).
		SetLeaseUntil(now.Add(LeaseDuration)).
		SetStatus(ent_backgroundjob.StatusProcessing).
		AddAttempts(1).
		AddRevision(1).
		Save(ctx)
	if err != nil || n == 0 {
		return nil, err
	}

	row, err := r.db.BackgroundJob.Query().
		Where(ent_backgroundjob.ID(id), ent_backgroundjob.LeaseToken(token)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return Map[json.RawMessage](row)
}

func leasePredicates(state State) predicate.BackgroundJob {
	return ent_backgroundjob.And(
		ent_backgroundjob.ID(state.ID),
		ent_backgroundjob.LeaseTokenEQ(state.LeaseToken),
		ent_backgroundjob.StatusEQ(ent_backgroundjob.StatusProcessing),
		ent_backgroundjob.LeaseUntilGT(time.Now()),
	)
}

func (r *Repository) Heartbeat(ctx context.Context, state State) error {
	n, err := r.db.BackgroundJob.Update().
		Where(leasePredicates(state)).
		SetLeaseUntil(time.Now().Add(LeaseDuration)).
		Save(ctx)
	if err != nil {
		return err
	}
	if n != 1 {
		return ErrLeaseLost
	}
	return nil
}

func (r *Repository) Defer(ctx context.Context, id xid.ID) error {
	return r.db.BackgroundJob.Update().
		Where(
			ent_backgroundjob.ID(id),
			ent_backgroundjob.StatusIn(ent_backgroundjob.StatusPending, ent_backgroundjob.StatusProcessing),
			ent_backgroundjob.LeaseUntilLTE(time.Now()),
		).
		SetDispatchUntil(time.Now().Add(time.Second)).
		Exec(ctx)
}

// update fences payload changes by revision so requests arriving during a lease are preserved.
func (r *Repository) update(ctx context.Context, id xid.ID, token string, fn func(*ent.BackgroundJob) error) error {
	for range 8 {
		row, err := r.db.BackgroundJob.Get(ctx, id)
		if err != nil {
			return err
		}

		if token != "" {
			if row.LeaseToken != token {
				return ErrLeaseLost
			}
			if row.Status != ent_backgroundjob.StatusProcessing {
				return ErrLeaseLost
			}
			if row.LeaseUntil.Before(time.Now()) {
				return ErrLeaseLost
			}
		}

		if err = fn(row); err != nil {
			return err
		}

		predicates := []predicate.BackgroundJob{
			ent_backgroundjob.ID(id),
			ent_backgroundjob.RevisionEQ(row.Revision),
		}
		if token != "" {
			predicates = append(predicates, leasePredicates(State{ID: id, LeaseToken: token}))
		}

		n, err := r.db.BackgroundJob.Update().
			Where(predicates...).
			SetPayload(row.Payload).
			SetStatus(row.Status).
			SetReason(row.Reason).
			SetAttempts(row.Attempts).
			SetAvailableAt(row.AvailableAt).
			AddRevision(1).
			Save(ctx)
		if err != nil {
			return err
		}
		if n == 1 {
			return nil
		}
	}

	return errors.New("background job changed concurrently")
}

func (r *Repository) Finish(
	ctx context.Context,
	state State,
	reason string,
	terminal bool,
	maxAttempts int,
	complete func(json.RawMessage) (bool, error),
) error {
	err := r.update(ctx, state.ID, state.LeaseToken, func(current *ent.BackgroundJob) error {
		current.Reason = reason
		current.Status = ent_backgroundjob.StatusPending
		if terminal {
			current.Status = ent_backgroundjob.StatusSkipped
		} else if reason != "" && current.Attempts >= maxAttempts {
			current.Status = ent_backgroundjob.StatusFailed
		} else if reason != "" {
			current.AvailableAt = time.Now().Add(time.Second * time.Duration(1<<min(current.Attempts, 6)))
		} else {
			current.Attempts = 0
			done := true
			if complete != nil {
				var err error
				done, err = complete(current.Payload)
				if err != nil {
					return err
				}
			}
			if done {
				current.Status = ent_backgroundjob.StatusComplete
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return r.Release(ctx, state)
}

func (r *Repository) Release(ctx context.Context, state State) error {
	return r.db.BackgroundJob.Update().
		Where(ent_backgroundjob.ID(state.ID), ent_backgroundjob.LeaseTokenEQ(state.LeaseToken)).
		SetLeaseUntil(time.Now()).
		SetDispatchUntil(time.Now()).
		Exec(ctx)
}
