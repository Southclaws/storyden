package trail_runtime

import (
	"context"
	"log/slog"
	"time"

	"github.com/Southclaws/storyden/app/resources/recurrence"
)

const schedulerLeaseDuration = 15 * time.Second

func (r *Runtime) schedulerLoop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-r.ctx.Done():
			return
		case now := <-ticker.C:
			if err := r.schedulerTick(r.ctx, now.UTC()); err != nil {
				r.logger.Error("Trail scheduler tick failed", slog.String("error", err.Error()))
			}
		}
	}
}

func (r *Runtime) schedulerTick(ctx context.Context, now time.Time) error {
	r.schedulerMu.Lock()
	defer r.schedulerMu.Unlock()

	if r.schedulerToken == "" {
		lease, acquired, err := r.repository.AcquireSchedulerLease(ctx, now, schedulerLeaseDuration)
		if err != nil || !acquired {
			return err
		}
		r.schedulerToken = lease.Token

		// The first valid lease deliberately advances over the offline gap.
		if err := r.skipOfflineGap(ctx, now); err != nil {
			return err
		}
		r.schedulerReady = true
		return nil
	}

	renewed, err := r.repository.RenewSchedulerLease(ctx, r.schedulerToken, now.Add(schedulerLeaseDuration))
	if err != nil {
		return err
	}

	if !renewed {
		r.schedulerToken = ""
		r.schedulerReady = false
		return nil
	}

	if !r.schedulerReady {
		if err := r.skipOfflineGap(ctx, now); err != nil {
			return err
		}

		r.schedulerReady = true
		return nil
	}

	return r.materialiseDue(ctx, now)
}

func (r *Runtime) skipOfflineGap(ctx context.Context, now time.Time) error {
	for {
		due, err := r.repository.Due(ctx, now, 100)
		if err != nil {
			return err
		}

		if len(due) == 0 {
			return nil
		}

		for _, definition := range due {
			schedule, err := recurrence.Parse(definition.TriggerConfig)
			if err != nil {
				r.logger.Error("invalid persisted Trail schedule", slog.String("trail_id", definition.ID.String()), slog.String("error", err.Error()))
				continue
			}

			first := *definition.NextOccurrenceAt
			latest, following, ok := recurrence.AdvancePast(schedule, first, now)
			if !ok {
				continue
			}

			var next *time.Time
			finished := following.IsZero()
			if !finished {
				next = &following
			}

			if _, _, err := r.repository.MaterialiseScheduled(ctx, definition.ID, first, latest, next, finished, true); err != nil {
				return err
			}
		}
	}
}

func (r *Runtime) materialiseDue(ctx context.Context, now time.Time) error {
	due, err := r.repository.Due(ctx, now, 100)
	if err != nil {
		return err
	}
	for _, definition := range due {
		schedule, err := recurrence.Parse(definition.TriggerConfig)
		if err != nil {
			r.logger.Error("invalid persisted Trail schedule", slog.String("trail_id", definition.ID.String()), slog.String("error", err.Error()))
			continue
		}

		current := *definition.NextOccurrenceAt
		following, hasNext := schedule.NextAfter(current)
		var next *time.Time
		if hasNext {
			following = following.UTC()
			next = &following
		}

		if _, _, err := r.repository.MaterialiseScheduled(ctx, definition.ID, current, current, next, !hasNext, false); err != nil {
			return err
		}
	}
	return nil
}
