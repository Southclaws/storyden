package background_job

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/rs/xid"
	"go.uber.org/fx"

	jobs "github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Runner struct {
	mu       sync.Mutex
	started  bool
	stopping bool
	active   sync.WaitGroup
	root     context.Context
	workers  map[string]*worker

	jobs        *jobs.Repository
	bus         *pubsub.Bus
	logger      *slog.Logger
	maxAttempts int
}

func New(
	lc fx.Lifecycle,
	ctx context.Context,
	jobs *jobs.Repository,
	bus *pubsub.Bus,
	logger *slog.Logger,
	cfg config.Config,
) *Runner {
	run, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	r := &Runner{
		root:        run,
		workers:     make(map[string]*worker),
		jobs:        jobs,
		bus:         bus,
		logger:      logger,
		maxAttempts: max(1, cfg.QueueMaxRetries+1),
	}

	lc.Append(fx.StartHook(func() error {
		r.mu.Lock()
		r.started = true
		r.mu.Unlock()

		_, err := pubsub.SubscribeCommand(run, bus, "background_job.attempt", r.handle)
		if err != nil {
			return err
		}

		go func() {
			defer close(done)
			r.dispatch(run)
		}()

		return nil
	}))
	lc.Append(fx.StopHook(func(ctx context.Context) error {
		r.mu.Lock()
		r.stopping = true
		cancel()
		r.mu.Unlock()

		workersDone := make(chan struct{})
		go func() {
			r.active.Wait()
			close(workersDone)
		}()

		select {
		case <-workersDone:
		case <-ctx.Done():
			return ctx.Err()
		}
		select {
		case <-done:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}))

	return r
}

func (r *Runner) register(kind string, workers int, handler handler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.started {
		return errors.New("background job handlers must be registered before startup")
	}
	if _, exists := r.workers[kind]; exists {
		return fmt.Errorf("background job handler already registered for %q", kind)
	}
	if kind == "" || workers < 1 || handler.Run == nil {
		return errors.New("background job handler requires a kind, worker count and run function")
	}

	r.workers[kind] = &worker{
		handler: handler,
		slots:   make(chan struct{}, workers),
	}
	return nil
}

// Dispatch is best effort after commit; the recovery scan retains responsibility for delivery.
func (r *Runner) Dispatch(ctx context.Context, id xid.ID) {
	row, err := r.jobs.Get[json.RawMessage](ctx, id)
	if err != nil {
		r.logger.Error("background job dispatch lookup failed", "job_id", id, "error", err)
		return
	}

	r.dispatchJob(ctx, row)
}

func (r *Runner) dispatch(ctx context.Context) {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	for {
		r.dispatchBatch(ctx)
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
		}
	}
}

func (r *Runner) dispatchBatch(ctx context.Context) {
	for kind, worker := range r.workers {
		limit := cap(worker.slots) - len(worker.slots)
		if limit <= 0 {
			continue
		}

		rows, err := r.jobs.Candidates(ctx, kind, limit)
		if err != nil {
			if ctx.Err() == nil {
				r.logger.Error("background job dispatch query failed", "kind", kind, "error", err)
			}
			continue
		}

		for _, row := range rows {
			r.dispatchJob(ctx, row)
		}
	}
}

func (r *Runner) dispatchJob(ctx context.Context, row *jobs.Job[json.RawMessage]) {
	worker, ok := r.workers[row.Kind]
	if !ok || len(worker.slots) == cap(worker.slots) {
		return
	}

	dispatched, err := r.jobs.Dispatch(ctx, row.ID)
	if err != nil {
		r.logger.Error("background job dispatch update failed", "job_id", row.ID, "error", err)
		return
	}
	if !dispatched {
		return
	}

	err = r.bus.SendCommand(ctx, &message.CommandAttemptBackgroundJob{ID: row.ID})
	if err != nil {
		r.logger.Error("background job command dispatch failed", "job_id", row.ID, "error", err)
	}
}

func (r *Runner) handle(ctx context.Context, cmd *message.CommandAttemptBackgroundJob) error {
	r.mu.Lock()
	if r.stopping {
		r.mu.Unlock()
		return nil
	}
	r.active.Add(1)
	r.mu.Unlock()
	defer r.active.Done()

	row, err := r.jobs.Get[json.RawMessage](ctx, cmd.ID)
	if err != nil {
		return err
	}
	worker, ok := r.workers[row.Kind]
	if !ok {
		return fmt.Errorf("no background job handler registered for %q", row.Kind)
	}

	select {
	case worker.slots <- struct{}{}:
		defer func() {
			<-worker.slots
		}()
	default:
		return r.jobs.Defer(ctx, cmd.ID)
	}

	job, err := r.jobs.Claim(ctx, cmd.ID, r.maxAttempts)
	if err != nil || job == nil {
		if err == nil && row.Attempts >= r.maxAttempts {
			r.notifyStateChange(ctx, row)
		}
		return err
	}
	r.publishStateChange(ctx, row.State, job.State)

	run, cancel := context.WithCancel(ctx)
	stop := context.AfterFunc(r.root, cancel)
	defer stop()
	defer cancel()

	heartbeatDone := make(chan struct{})
	go func() {
		defer close(heartbeatDone)
		r.heartbeat(run, cancel, job)
	}()

	started := time.Now()
	result, err := worker.handler.Run(run, job)
	r.logger.Debug("background job finished",
		"job_id", job.ID,
		"kind", job.Kind,
		"duration", time.Since(started),
		"error", err,
	)
	cancel()
	<-heartbeatDone

	finishCtx := context.WithoutCancel(ctx)
	err = r.finish(finishCtx, job, worker.handler, result, err)
	r.notifyStateChange(finishCtx, job)
	return err
}

func (r *Runner) notifyStateChange(ctx context.Context, previous *jobs.Job[json.RawMessage]) {
	current, err := r.jobs.Get[json.RawMessage](ctx, previous.ID)
	if err != nil {
		r.logger.Error("background job state lookup failed", "job_id", previous.ID, "error", err)
		return
	}
	if previous.LeaseToken != "" && current.LeaseToken != previous.LeaseToken {
		return
	}
	r.publishStateChange(ctx, previous.State, current.State)
}

func (r *Runner) publishStateChange(ctx context.Context, previous, current jobs.State) {
	if previous.Status == current.Status {
		return
	}

	r.bus.Publish(ctx, &message.EventBackgroundJobStateChanged{
		ID:             current.ID,
		Kind:           current.Kind,
		Revision:       current.Revision,
		PreviousStatus: previous.Status,
		Status:         current.Status,
		Attempts:       current.Attempts,
		Reason:         current.Reason,
	})

	if current.Status == jobs.StatusComplete {
		r.bus.Publish(ctx, &message.EventBackgroundJobCompleted{
			ID:       current.ID,
			Kind:     current.Kind,
			Revision: current.Revision,
		})
	}
}

func (r *Runner) finish(
	ctx context.Context,
	job *jobs.Job[json.RawMessage],
	handler handler,
	result Result,
	runErr error,
) error {
	reason := ""
	if runErr != nil && !result.Deferred {
		reason = runErr.Error()
		r.logger.Warn("background job failed",
			"job_id", job.ID,
			"kind", job.Kind,
			"attempt", job.Attempts,
			"error", runErr,
		)
	}

	complete := handler.Complete
	if result.Deferred {
		complete = func(json.RawMessage) (bool, error) {
			return false, nil
		}
	}

	err := r.jobs.Finish(ctx, job.State, reason, result.Terminal, r.maxAttempts, complete)
	if err != nil {
		return err
	}
	if result.Terminal || result.Deferred {
		return nil
	}
	return runErr
}

func (r *Runner) heartbeat(ctx context.Context, cancel context.CancelFunc, job *jobs.Job[json.RawMessage]) {
	tick := time.NewTicker(20 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			if err := r.jobs.Heartbeat(ctx, job.State); err != nil {
				cancel()
				return
			}
		}
	}
}
