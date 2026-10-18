package background_job_test

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestCheckpointPreservesConcurrentRequests(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		jobs *background_job.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			id := xid.New()
			require.NoError(t, jobs.Ensure(ctx, id, "test.checkpoint", id.String(), payload{Requested: 1}))
			job, err := jobs.Claim(ctx, id, 3)
			require.NoError(t, err)
			require.NotNil(t, job)

			err = jobs.Request(ctx, id, func(current *payload) error {
				current.Requested |= 2
				return nil
			}, func(current payload) bool {
				return current.Requested == current.Completed
			})
			require.NoError(t, err)

			started := make(chan struct{})
			release := make(chan struct{})
			checkpointDone := make(chan error, 1)
			go func() {
				checkpointDone <- jobs.Checkpoint(ctx, job.State, func(tx *ent.Tx, current *payload) error {
					close(started)
					select {
					case <-release:
					case <-ctx.Done():
						return ctx.Err()
					}

					if err := tx.Setting.Create().SetID(id.String()).SetValue("published").Exec(ctx); err != nil {
						return err
					}
					current.Completed |= 1
					return nil
				})
			}()
			defer close(release)
			select {
			case <-started:
			case err := <-checkpointDone:
				t.Fatalf("checkpoint stopped before its callback: %v", err)
			case <-time.After(5 * time.Second):
				t.Fatal("checkpoint did not start")
			}

			requestRead := make(chan struct{})
			requestDone := make(chan error, 1)
			var readOnce sync.Once
			go func() {
				requestDone <- jobs.Request(ctx, id, func(current *payload) error {
					readOnce.Do(func() {
						close(requestRead)
					})
					current.Requested |= 4
					return nil
				}, func(current payload) bool {
					return current.Requested == current.Completed
				})
			}()
			select {
			case <-requestRead:
			case <-time.After(5 * time.Second):
				t.Fatal("request did not read the job")
			}
			release <- struct{}{}
			select {
			case err := <-checkpointDone:
				require.NoError(t, err)
			case <-time.After(5 * time.Second):
				t.Fatal("checkpoint did not finish")
			}
			select {
			case err := <-requestDone:
				require.NoError(t, err)
			case <-time.After(5 * time.Second):
				t.Fatal("request did not finish")
			}

			stored, err := jobs.Get[payload](ctx, id)
			require.NoError(t, err)
			require.Equal(t, payload{Requested: 7, Completed: 1}, stored.Payload)
			require.Equal(t, background_job.StatusProcessing, stored.Status)
			require.Equal(t, job.LeaseToken, stored.LeaseToken)
			require.Equal(t, job.Revision+3, stored.Revision)
			published, err := db.Setting.Get(ctx, id.String())
			require.NoError(t, err)
			require.Equal(t, "published", published.Value)
		}))
	}))
}

func TestCheckpointRollsBackWritesAndProgress(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		jobs *background_job.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			callbackError := errors.New("application write failed")
			for _, test := range []struct {
				name      string
				expire    bool
				wantError error
			}{
				{name: "callback failure", wantError: callbackError},
				{name: "lease expires before publication", expire: true, wantError: background_job.ErrLeaseLost},
			} {
				t.Run(test.name, func(t *testing.T) {
					id := xid.New()
					require.NoError(t, jobs.Ensure(ctx, id, "test.checkpoint", id.String(), payload{Requested: 1}))
					job, err := jobs.Claim(ctx, id, 3)
					require.NoError(t, err)
					require.NotNil(t, job)
					calls := 0
					err = jobs.Checkpoint(ctx, job.State, func(tx *ent.Tx, current *payload) error {
						calls++
						if err := tx.Setting.Create().SetID(id.String()).SetValue("candidate").Exec(ctx); err != nil {
							return err
						}
						if !test.expire {
							return callbackError
						}
						if err := tx.BackgroundJob.UpdateOneID(id).SetLeaseUntil(time.Now().Add(-time.Second)).Exec(ctx); err != nil {
							return err
						}
						current.Completed = 1
						return nil
					})
					require.ErrorIs(t, err, test.wantError)
					require.Equal(t, 1, calls)
					_, err = db.Setting.Get(ctx, id.String())
					require.True(t, ent.IsNotFound(err))
					stored, err := jobs.Get[json.RawMessage](ctx, id)
					require.NoError(t, err)
					require.JSONEq(t, string(job.Payload), string(stored.Payload))
					require.Equal(t, job.Revision, stored.Revision)
					require.Equal(t, job.LeaseToken, stored.LeaseToken)
					require.WithinDuration(t, job.LeaseUntil, stored.LeaseUntil, time.Millisecond)
				})
			}
		}))
	}))
}

func TestCheckpointRejectsStaleLeases(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		jobs *background_job.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			for _, takeover := range []bool{false, true} {
				id := xid.New()
				require.NoError(t, jobs.Ensure(ctx, id, "test.checkpoint", id.String(), payload{Requested: 1}))
				job, err := jobs.Claim(ctx, id, 3)
				require.NoError(t, err)
				require.NotNil(t, job)
				require.NoError(t, db.BackgroundJob.UpdateOneID(id).SetLeaseUntil(time.Now().Add(-time.Second)).Exec(ctx))
				if takeover {
					replacement, err := jobs.Claim(ctx, id, 3)
					require.NoError(t, err)
					require.NotNil(t, replacement)
				}

				called := false
				err = jobs.Checkpoint(ctx, job.State, func(tx *ent.Tx, current *payload) error {
					called = true
					return nil
				})
				require.ErrorIs(t, err, background_job.ErrLeaseLost)
				require.False(t, called)
				stored, err := jobs.Get[json.RawMessage](ctx, id)
				require.NoError(t, err)
				require.JSONEq(t, string(job.Payload), string(stored.Payload))
			}
		}))
	}))
}
