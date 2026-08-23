package trail_runtime

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/trail"
	"github.com/Southclaws/storyden/app/services/trail/trail_action_registry"
)

const actionLeaseDuration = time.Minute

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
