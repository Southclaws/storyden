package background_job

import (
	"context"
	"encoding/json"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
	ent_backgroundjob "github.com/Southclaws/storyden/internal/ent/backgroundjob"
)

// UpdatePayload may retry the callback on revision conflicts. The callback must
// only mutate the supplied payload; database writes belong in Checkpoint.
func (r *Repository) UpdatePayload[T any](
	ctx context.Context,
	id xid.ID,
	token string,
	update func(*T) error,
) error {
	return r.updatePayload(ctx, id, token, func(_ *ent.BackgroundJob, payload *T) error {
		return update(payload)
	})
}

// Request merges intent into an existing job, reopening completed jobs when needed.
// Failed and skipped jobs remain terminal. The callback follows UpdatePayload's rules.
func (r *Repository) Request[T any](
	ctx context.Context,
	id xid.ID,
	update func(*T) error,
	complete func(T) bool,
) error {
	return r.updatePayload(ctx, id, "", func(job *ent.BackgroundJob, payload *T) error {
		if err := update(payload); err != nil {
			return err
		}
		if job.Status == ent_backgroundjob.StatusComplete && !complete(*payload) {
			job.Status = ent_backgroundjob.StatusPending
		}
		return nil
	})
}

func (r *Repository) updatePayload[T any](
	ctx context.Context,
	id xid.ID,
	token string,
	mutate func(*ent.BackgroundJob, *T) error,
) error {
	return r.update(ctx, id, token, func(row *ent.BackgroundJob) error {
		payload, err := Decode[T](row.Payload)
		if err != nil {
			return err
		}
		if err := mutate(row, &payload); err != nil {
			return err
		}

		encoded, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		row.Payload = encoded
		return nil
	})
}
