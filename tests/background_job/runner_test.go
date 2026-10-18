package background_job_test

import (
	"context"
	"encoding/json"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/app/resources/message"
	job_runner "github.com/Southclaws/storyden/app/services/background_job"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	ej "github.com/Southclaws/storyden/internal/ent/backgroundjob"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/integration"
)

type payload struct {
	Requested int `json:"requested"`
	Completed int `json:"completed"`
}

func TestRunnerKeepsWorkerLimitsIndependent(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		runner *job_runner.Runner,
		jobs *background_job.Repository,
		bus *pubsub.Bus,
	) {
		started := make(chan xid.ID, 4)
		release := make(chan struct{})
		err := runner.Register("test.blocking", 1, job_runner.Handler[payload]{
			Run: func(ctx context.Context, job *background_job.Job[payload]) (job_runner.Result, error) {
				started <- job.ID
				select {
				case <-release:
					return job_runner.Result{}, nil
				case <-ctx.Done():
					return job_runner.Result{}, ctx.Err()
				}
			},
		})
		require.NoError(t, err)

		otherStarted := make(chan xid.ID, 1)
		err = runner.Register("test.other", 1, job_runner.Handler[payload]{
			Run: func(ctx context.Context, job *background_job.Job[payload]) (job_runner.Result, error) {
				otherStarted <- job.ID
				return job_runner.Result{}, nil
			},
		})
		require.NoError(t, err)

		lc.Append(fx.StartHook(func() {
			firstID := xid.New()
			err := jobs.Ensure(ctx, firstID, "test.blocking", firstID.String(), payload{})
			require.NoError(t, err)
			runner.Dispatch(ctx, firstID)

			select {
			case id := <-started:
				require.Equal(t, firstID, id)
			case <-time.After(5 * time.Second):
				t.Fatal("first job did not start")
			}

			secondID := xid.New()
			err = jobs.Ensure(ctx, secondID, "test.blocking", secondID.String(), payload{})
			require.NoError(t, err)
			before := time.Now()
			require.NoError(t, bus.SendCommand(ctx, &message.CommandAttemptBackgroundJob{ID: secondID}))
			require.NoError(t, bus.SendCommand(ctx, &message.CommandAttemptBackgroundJob{ID: firstID}))

			require.Eventually(t, func() bool {
				job, err := jobs.Get[json.RawMessage](ctx, secondID)
				return err == nil && job.DispatchUntil.After(before)
			}, 5*time.Second, 20*time.Millisecond)
			pending, err := jobs.Get[json.RawMessage](ctx, secondID)
			require.NoError(t, err)
			require.Equal(t, background_job.StatusPending, pending.Status)
			require.Zero(t, pending.Attempts)
			require.Empty(t, started)

			otherID := xid.New()
			err = jobs.Ensure(ctx, otherID, "test.other", otherID.String(), payload{})
			require.NoError(t, err)
			runner.Dispatch(ctx, otherID)
			select {
			case id := <-otherStarted:
				require.Equal(t, otherID, id)
			case <-time.After(5 * time.Second):
				t.Fatal("another job kind was blocked by the first kind")
			}

			close(release)
			select {
			case id := <-started:
				require.Equal(t, secondID, id)
			case <-time.After(5 * time.Second):
				t.Fatal("pending job did not recover after capacity became available")
			}

			require.Eventually(t, func() bool {
				for _, id := range []xid.ID{firstID, secondID, otherID} {
					job, err := jobs.Get[json.RawMessage](ctx, id)
					if err != nil || job.Status != background_job.StatusComplete {
						return false
					}
				}
				return true
			}, 5*time.Second, 20*time.Millisecond)
			require.Empty(t, started)
		}))
	}))
}

func TestRunnerRecoversUndeliveredAndExpiredJobs(t *testing.T) {
	t.Parallel()

	integration.Test(t, &config.Config{QueueMaxRetries: 1}, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		runner *job_runner.Runner,
		jobs *background_job.Repository,
		db *ent.Client,
	) {
		started := make(chan xid.ID, 2)
		err := runner.Register("test.recovery", 1, job_runner.Handler[payload]{
			Run: func(ctx context.Context, job *background_job.Job[payload]) (job_runner.Result, error) {
				started <- job.ID
				return job_runner.Result{}, nil
			},
		})
		require.NoError(t, err)

		lc.Append(fx.StartHook(func() {
			undeliveredID := xid.New()
			err := db.BackgroundJob.Create().
				SetID(undeliveredID).
				SetKey(undeliveredID.String()).
				SetKind("test.recovery").
				SetPayload(json.RawMessage(`{}`)).
				SetDispatchUntil(time.Now().Add(-time.Second)).
				Exec(ctx)
			require.NoError(t, err)

			expiredID := xid.New()
			err = db.BackgroundJob.Create().
				SetID(expiredID).
				SetKey(expiredID.String()).
				SetKind("test.recovery").
				SetPayload(json.RawMessage(`{}`)).
				SetStatus(ej.StatusProcessing).
				SetAttempts(1).
				SetLeaseToken(xid.New().String()).
				SetLeaseUntil(time.Now().Add(-time.Second)).
				Exec(ctx)
			require.NoError(t, err)

			require.Eventually(t, func() bool {
				for _, id := range []xid.ID{undeliveredID, expiredID} {
					job, err := jobs.Get[json.RawMessage](ctx, id)
					if err != nil || job.Status != background_job.StatusComplete {
						return false
					}
				}
				return true
			}, 5*time.Second, 20*time.Millisecond)
			require.ElementsMatch(t, []xid.ID{undeliveredID, expiredID}, []xid.ID{<-started, <-started})
		}))
	}))
}

func TestRunnerPersistsTerminalAndRetryOutcomes(t *testing.T) {
	t.Parallel()

	integration.Test(t, &config.Config{QueueMaxRetries: 1}, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		runner *job_runner.Runner,
		jobs *background_job.Repository,
	) {
		err := runner.Register("test.retry", 1, job_runner.Handler[payload]{
			Run: func(context.Context, *background_job.Job[payload]) (job_runner.Result, error) {
				return job_runner.Result{}, errors.New("temporary storage failure")
			},
		})
		require.NoError(t, err)
		err = runner.Register("test.skip", 1, job_runner.Handler[payload]{
			Run: func(context.Context, *background_job.Job[payload]) (job_runner.Result, error) {
				return job_runner.Result{Terminal: true}, errors.New("unsupported input")
			},
		})
		require.NoError(t, err)

		lc.Append(fx.StartHook(func() {
			retryID := xid.New()
			err := jobs.Ensure(ctx, retryID, "test.retry", retryID.String(), payload{})
			require.NoError(t, err)
			runner.Dispatch(ctx, retryID)

			skipID := xid.New()
			err = jobs.Ensure(ctx, skipID, "test.skip", skipID.String(), payload{})
			require.NoError(t, err)
			runner.Dispatch(ctx, skipID)

			require.Eventually(t, func() bool {
				retry, err := jobs.Get[json.RawMessage](ctx, retryID)
				if err != nil || retry.Status != background_job.StatusFailed {
					return false
				}
				skipped, err := jobs.Get[json.RawMessage](ctx, skipID)
				return err == nil && skipped.Status == background_job.StatusSkipped
			}, 8*time.Second, 20*time.Millisecond)

			retry, err := jobs.Get[json.RawMessage](ctx, retryID)
			require.NoError(t, err)
			require.Equal(t, 2, retry.Attempts)
			require.Equal(t, "temporary storage failure", retry.Reason)
			skipped, err := jobs.Get[json.RawMessage](ctx, skipID)
			require.NoError(t, err)
			require.Equal(t, 1, skipped.Attempts)
			require.Equal(t, "unsupported input", skipped.Reason)

			for _, id := range []xid.ID{retryID, skipID} {
				marked, err := jobs.Dispatch(ctx, id)
				require.NoError(t, err)
				require.False(t, marked)
			}
		}))
	}))
}

func TestRunnerDefersAdmissionWithoutCompletingTheJob(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		runner *job_runner.Runner,
		jobs *background_job.Repository,
	) {
		var attempts atomic.Int32
		err := runner.Register("test.admission", 1, job_runner.Handler[payload]{
			Run: func(context.Context, *background_job.Job[payload]) (job_runner.Result, error) {
				attempts.Add(1)
				return job_runner.Result{Deferred: true}, errors.New("memory busy")
			},
		})
		require.NoError(t, err)

		lc.Append(fx.StartHook(func() {
			id := xid.New()
			err := jobs.Ensure(ctx, id, "test.admission", id.String(), payload{})
			require.NoError(t, err)
			runner.Dispatch(ctx, id)

			require.Eventually(t, func() bool {
				job, err := jobs.Get[json.RawMessage](ctx, id)
				if err != nil || attempts.Load() == 0 {
					return false
				}
				return job.Status == background_job.StatusPending && job.LeaseUntil.Before(time.Now())
			}, 5*time.Second, 20*time.Millisecond)

			job, err := jobs.Get[json.RawMessage](ctx, id)
			require.NoError(t, err)
			require.Zero(t, job.Attempts)
			require.Empty(t, job.Reason)
		}))
	}))
}

func TestRunnerCompletionUsesRequestsAddedDuringTheLease(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		runner *job_runner.Runner,
		jobs *background_job.Repository,
	) {
		started := make(chan int, 2)
		release := make(chan struct{})
		err := runner.Register("test.progress", 1, job_runner.Handler[payload]{
			Run: func(ctx context.Context, job *background_job.Job[payload]) (job_runner.Result, error) {
				requested := job.Payload
				started <- requested.Requested
				select {
				case <-release:
				case <-ctx.Done():
					return job_runner.Result{}, ctx.Err()
				}

				err := jobs.UpdatePayload(ctx, job.ID, job.LeaseToken, func(current *payload) error {
					current.Completed = requested.Requested
					return nil
				})
				return job_runner.Result{}, err
			},
			Complete: func(value payload) bool {
				return value.Requested == value.Completed
			},
		})
		require.NoError(t, err)

		lc.Append(fx.StartHook(func() {
			id := xid.New()
			err := jobs.Ensure(ctx, id, "test.progress", id.String(), payload{Requested: 1})
			require.NoError(t, err)
			runner.Dispatch(ctx, id)
			select {
			case requested := <-started:
				require.Equal(t, 1, requested)
			case <-time.After(5 * time.Second):
				t.Fatal("initial job did not start")
			}

			err = jobs.Request(ctx, id, func(current *payload) error {
				current.Requested = 2
				return nil
			}, func(current payload) bool {
				return current.Requested == current.Completed
			})
			require.NoError(t, err)
			close(release)

			select {
			case requested := <-started:
				require.Equal(t, 2, requested)
			case <-time.After(5 * time.Second):
				t.Fatal("newly requested work did not run")
			}
			require.Eventually(t, func() bool {
				job, err := jobs.Get[json.RawMessage](ctx, id)
				return err == nil && job.Status == background_job.StatusComplete
			}, 5*time.Second, 20*time.Millisecond)
		}))
	}))
}
