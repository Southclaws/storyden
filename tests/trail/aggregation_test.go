package trail_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/trail"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestTrailAggregationDoesNotRegress(t *testing.T) {
	tests.ParallelIfNotSharedPostgres(t)
	type pauseAggregationKey struct{}

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		db *ent.Client,
		accounts *account_writer.Writer,
		repository *trail.Repository,
	) {
		db.TrailRun.Use(func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
				if pause, ok := ctx.Value(pauseAggregationKey{}).(func()); ok {
					pause()
				}
				return next.Mutate(ctx, mutation)
			})
		})
		lc.Append(fx.StartHook(func() {
			_, creator := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
			for _, tc := range []struct {
				name     string
				results  []trail.ActionRunStatus
				expected trail.RunStatus
			}{
				{"running", nil, trail.RunStatusRunning},
				{"completed", []trail.ActionRunStatus{trail.ActionRunStatusCompleted, trail.ActionRunStatusCompleted}, trail.RunStatusCompleted},
				{"attention", []trail.ActionRunStatus{trail.ActionRunStatusCompleted, trail.ActionRunStatusBlocked}, trail.RunStatusAttentionRequired},
				{"cancelled", []trail.ActionRunStatus{trail.ActionRunStatusCancelled, trail.ActionRunStatusCancelled}, trail.RunStatusCancelled},
			} {
				t.Run(tc.name, func(t *testing.T) {
					r := require.New(t)
					ctx, cancel := context.WithTimeout(root, 10*time.Second)
					defer cancel()
					config := json.RawMessage(`{"type":"robot_run","robot_ref":"robot-one","instruction":"Prepare the community prompt"}`)
					definition, err := repository.Create(ctx, xid.ID(creator.ID), "Concurrent aggregation", "", trail.StatusPaused,
						domainEventTrigger(t, "EventThreadPublished"),
						[]trail.ActionSpec{{Kind: trail.ActionKindRobotRun, Config: config}, {Kind: trail.ActionKindRobotRun, Config: config}},
						opt.NewEmpty[time.Time]())
					r.NoError(err)
					run, err := repository.MaterialiseManual(ctx, definition.ID, xid.ID(creator.ID))
					r.NoError(err)
					r.Len(run.ActionRuns, 2)
					claimed := make([]*trail.ActionRun, 2)
					if tc.results != nil {
						for i, action := range run.ActionRuns {
							var ok bool
							claimed[i], ok, err = repository.ClaimActionRunByID(ctx, action.ID, time.Now(), time.Minute)
							r.NoError(err)
							r.True(ok)
						}
					}

					ready, resume := make(chan struct{}), make(chan struct{}, 1)
					defer close(resume)
					staleContext := context.WithValue(ctx, pauseAggregationKey{}, func() {
						close(ready)
						select {
						case <-resume:
						case <-ctx.Done():
						}
					})
					done := make(chan error, 1)
					go func() { done <- repository.AggregateRun(staleContext, run.ActionRuns[0].ID) }()
					select {
					case <-ready:
					case <-ctx.Done():
						t.Fatal(ctx.Err())
					}

					// Advance the run after the old snapshot was read, before it is written.
					for i, action := range run.ActionRuns {
						if tc.results == nil {
							_, ok, err := repository.ClaimActionRunByID(ctx, action.ID, time.Now(), time.Minute)
							r.NoError(err)
							r.True(ok)
							r.NoError(repository.AggregateRun(ctx, action.ID))
						} else {
							r.NoError(repository.FinishActionRun(ctx, action.ID, *claimed[i].LeaseToken, tc.results[i], nil, ""))
						}
					}
					before, err := repository.GetRun(ctx, run.ID)
					r.NoError(err)
					r.Equal(tc.expected, before.Status)
					resume <- struct{}{}
					r.NoError(<-done)
					after, err := repository.GetRun(ctx, run.ID)
					r.NoError(err)
					r.Equal(tc.expected, after.Status)
					r.Equal(before.FinishedAt, after.FinishedAt)
				})
			}
		}))
	}))
}
