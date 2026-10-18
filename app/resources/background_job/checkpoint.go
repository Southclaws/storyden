package background_job

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/storyden/internal/ent"
	ent_backgroundjob "github.com/Southclaws/storyden/internal/ent/backgroundjob"
)

// Checkpoint commits database writes and progress together without finishing
// the job. The callback runs at most once and must use the supplied transaction
// for database work only.
func (r *Repository) Checkpoint[T any](
	ctx context.Context,
	state State,
	write func(*ent.Tx, *T) error,
) error {
	return ent.WithTx(ctx, r.db, func(tx *ent.Tx) error {
		// Lock before reading progress.
		// So requests and lease takeovers cannot interleave.
		updated, err := tx.BackgroundJob.Update().
			Where(leasePredicates(state)).
			SetLeaseToken(state.LeaseToken).
			Save(ctx)
		if err != nil {
			return err
		}
		if updated != 1 {
			return ErrLeaseLost
		}

		current, err := tx.BackgroundJob.Get(ctx, state.ID)
		if err != nil {
			return err
		}

		payload, err := Decode[T](current.Payload)
		if err != nil {
			return err
		}
		if err := write(tx, &payload); err != nil {
			return err
		}

		encoded, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		updated, err = tx.BackgroundJob.Update().
			Where(leasePredicates(state), ent_backgroundjob.RevisionEQ(current.Revision)).
			SetPayload(encoded).
			AddRevision(1).
			Save(ctx)
		if err != nil {
			return err
		}
		if updated != 1 {
			return ErrLeaseLost
		}
		return nil
	})
}
