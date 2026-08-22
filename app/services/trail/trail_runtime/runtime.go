package trail_runtime

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/recurrence"
	"github.com/Southclaws/storyden/app/resources/trail"
	"github.com/Southclaws/storyden/app/services/trail/trail_action_registry"
)

const (
	schedulerLeaseDuration = 15 * time.Second
	actionLeaseDuration    = time.Minute
)

type Runtime struct {
	ctx        context.Context
	logger     *slog.Logger
	repository *trail.Repository
	registry   *trail_action_registry.Registry

	mu             sync.Mutex
	schedulerToken string
	schedulerReady bool
}

func New(
	ctx context.Context,
	lifecycle fx.Lifecycle,
	logger *slog.Logger,
	repository *trail.Repository,
	registry *trail_action_registry.Registry,
) *Runtime {
	runtime := &Runtime{ctx: ctx, logger: logger, repository: repository, registry: registry}
	lifecycle.Append(fx.StartHook(func(context.Context) {
		go runtime.schedulerLoop()
		go runtime.dispatchLoop()
	}))
	return runtime
}

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
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.schedulerToken == "" {
		lease, acquired, err := r.repository.AcquireSchedulerLease(ctx, now, schedulerLeaseDuration)
		if err != nil || !acquired {
			return err
		}
		r.schedulerToken = lease.Token

		// the first valid lease deliberately advances over the offline gap.
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

func (r *Runtime) dispatchLoop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.ctx.Done():
			return

		case now := <-ticker.C:
			for claimed := 0; claimed < 16; claimed++ {
				run, ok, err := r.repository.ClaimActionRun(r.ctx, now.UTC(), actionLeaseDuration)
				if err != nil {
					r.logger.Error("failed to claim Trail action", slog.String("error", err.Error()))
					break
				}
				if !ok {
					break
				}

				_ = r.repository.AggregateRun(r.ctx, run.ID)
				go r.executeAction(run)
			}
		}
	}
}

func (r *Runtime) executeAction(run *trail.ActionRun) {
	ctx := context.WithoutCancel(r.ctx)

	if err := r.registry.Start(ctx, run); err != nil {
		r.finishAction(ctx, run, &trail_action_registry.ActionResult{Status: trail.ActionRunStatusFailed, ErrorText: err.Error()})
		return
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		result, err := r.registry.Reconcile(ctx, run)
		if err != nil {
			r.finishAction(ctx, run, &trail_action_registry.ActionResult{Status: trail.ActionRunStatusFailed, ErrorText: err.Error()})
			return
		}

		if result != nil {
			r.finishAction(ctx, run, result)
			return
		}

		if _, err := r.repository.RenewActionLease(ctx, run.ID, *run.LeaseToken, time.Now().UTC().Add(actionLeaseDuration)); err != nil {
			r.logger.Error("failed to renew Trail action lease", slog.String("action_run_id", run.ID.String()), slog.String("error", err.Error()))
			return
		}

		select {
		case <-r.ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (r *Runtime) finishAction(ctx context.Context, run *trail.ActionRun, result *trail_action_registry.ActionResult) {
	if result.Status == (trail.ActionRunStatus{}) {
		result.Status = trail.ActionRunStatusFailed
	}

	if err := r.repository.FinishActionRun(ctx, run.ID, *run.LeaseToken, result.Status, json.RawMessage(result.Output), result.ErrorText); err != nil {
		r.logger.Error("failed to finish Trail action", slog.String("action_run_id", run.ID.String()), slog.String("error", err.Error()))
		return
	}

	if (result.Status == trail.ActionRunStatusBlocked || result.Status == trail.ActionRunStatusFailed) && run.Trail != nil {
		target := run.RunID.String()
		if _, err := r.repository.RecordAttentionNotification(ctx, run.ID, xid.ID(run.Trail.Creator.ID), target); err != nil {
			r.logger.Error("failed to notify Trail creator", slog.String("action_run_id", run.ID.String()), slog.String("error", err.Error()))
		}
	}
}
