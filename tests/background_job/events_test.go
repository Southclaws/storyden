package background_job_test

import (
	"context"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/app/resources/message"
	job_runner "github.com/Southclaws/storyden/app/services/background_job"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestRunnerPublishesCommittedStateChanges(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		runner *job_runner.Runner,
		jobs *background_job.Repository,
		bus *pubsub.Bus,
	) {
		id := xid.New()
		changes := make(chan message.EventBackgroundJobStateChanged, 2)
		completed := make(chan message.EventBackgroundJobCompleted, 1)

		err := runner.Register("test.events", 1, job_runner.Handler[payload]{
			Run: func(context.Context, *background_job.Job[payload]) (job_runner.Result, error) {
				return job_runner.Result{}, nil
			},
		})
		require.NoError(t, err)

		lc.Append(fx.StartHook(func() {
			_, err := pubsub.Subscribe(ctx, bus, "test.job_state_changes",
				func(_ context.Context, event *message.EventBackgroundJobStateChanged) error {
					if event.ID == id {
						changes <- *event
					}

					return nil
				})
			require.NoError(t, err)

			_, err = pubsub.Subscribe(ctx, bus, "test.job_completion",
				func(_ context.Context, event *message.EventBackgroundJobCompleted) error {
					if event.ID == id {
						completed <- *event
					}

					return nil
				})
			require.NoError(t, err)

			require.NoError(t, jobs.Ensure(ctx, id, "test.events", id.String(), payload{}))
			runner.Dispatch(ctx, id)

			select {
			case event := <-completed:
				job, err := jobs.Get[payload](ctx, id)
				require.NoError(t, err)
				require.Equal(t, background_job.StatusComplete, job.Status)
				require.Equal(t, job.Kind, event.Kind)
				require.Equal(t, job.Revision, event.Revision)
			case <-time.After(5 * time.Second):
				t.Fatal("job completion event was not published")
			}

			require.Eventually(t, func() bool {
				return len(changes) == 2
			}, 5*time.Second, 20*time.Millisecond)

			first, second := <-changes, <-changes
			if first.Revision > second.Revision {
				first, second = second, first
			}

			require.Equal(t, background_job.StatusPending, first.PreviousStatus)
			require.Equal(t, background_job.StatusProcessing, first.Status)
			require.Equal(t, background_job.StatusProcessing, second.PreviousStatus)
			require.Equal(t, background_job.StatusComplete, second.Status)
			require.Greater(t, second.Revision, first.Revision)
		}))
	}))
}
