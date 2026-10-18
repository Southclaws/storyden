package background_job_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/internal/ent"
	ej "github.com/Southclaws/storyden/internal/ent/backgroundjob"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestRequestPreservesJobLifecycle(t *testing.T) {
	t.Parallel()

	type refreshPayload struct {
		Desired string `json:"desired"`
		Applied string `json:"applied"`
	}

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		jobs *background_job.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			for _, test := range []struct {
				name   string
				status ej.Status
				desire string
				want   ej.Status
			}{
				{name: "pending", status: ej.StatusPending, desire: "b", want: ej.StatusPending},
				{name: "processing", status: ej.StatusProcessing, desire: "b", want: ej.StatusProcessing},
				{name: "already complete", status: ej.StatusComplete, desire: "a", want: ej.StatusComplete},
				{name: "more work", status: ej.StatusComplete, desire: "b", want: ej.StatusPending},
				{name: "failed", status: ej.StatusFailed, desire: "b", want: ej.StatusFailed},
				{name: "skipped", status: ej.StatusSkipped, desire: "b", want: ej.StatusSkipped},
			} {
				t.Run(test.name, func(t *testing.T) {
					id := xid.New()
					initial := refreshPayload{Desired: "a", Applied: "a"}
					require.NoError(t, jobs.Ensure(ctx, id, "test.refresh", id.String(), initial))
					before, err := db.BackgroundJob.UpdateOneID(id).
						SetStatus(test.status).
						SetAttempts(2).
						SetReason("previous outcome").
						SetLeaseToken("current lease").
						SetLeaseUntil(time.Now().Add(time.Minute)).
						SetDispatchUntil(time.Now().Add(time.Minute)).
						SetAvailableAt(time.Now().Add(time.Minute)).
						Save(ctx)
					require.NoError(t, err)

					err = jobs.Request(ctx, id, func(current *refreshPayload) error {
						current.Desired = test.desire
						return nil
					}, func(current refreshPayload) bool {
						return current.Desired == current.Applied
					})
					require.NoError(t, err)
					after, err := jobs.Get[refreshPayload](ctx, id)
					require.NoError(t, err)
					require.Equal(t, refreshPayload{Desired: test.desire, Applied: "a"}, after.Payload)
					require.Equal(t, string(test.want), after.Status.String())
					require.Equal(t, before.Revision+1, after.Revision)
					require.Equal(t, before.Attempts, after.Attempts)
					require.Equal(t, before.Reason, after.Reason)
					require.Equal(t, before.LeaseToken, after.LeaseToken)
					require.True(t, before.LeaseUntil.Equal(after.LeaseUntil))
					require.True(t, before.DispatchUntil.Equal(after.DispatchUntil))
					require.True(t, before.AvailableAt.Equal(after.AvailableAt))
				})
			}
		}))
	}))
}

func TestUpdatePayloadPreservesPrecisionAndRejectsInvalidMutations(t *testing.T) {
	t.Parallel()

	type sequencePayload struct {
		Sequence int64 `json:"sequence"`
	}

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		jobs *background_job.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			id := xid.New()
			initial := sequencePayload{Sequence: 1152921504606846977}
			require.NoError(t, jobs.Ensure(ctx, id, "test.sequence", id.String(), initial))
			job, err := jobs.Claim(ctx, id, 3)
			require.NoError(t, err)
			require.NotNil(t, job)
			err = jobs.UpdatePayload(ctx, id, job.LeaseToken, func(current *sequencePayload) error {
				current.Sequence++
				return nil
			})
			require.NoError(t, err)
			stored, err := jobs.Get[json.RawMessage](ctx, id)
			require.NoError(t, err)
			require.JSONEq(t, `{"sequence":1152921504606846978}`, string(stored.Payload))
			require.Equal(t, background_job.StatusProcessing, stored.Status)

			callbackError := errors.New("invalid sequence")
			err = jobs.UpdatePayload(ctx, id, job.LeaseToken, func(current *sequencePayload) error {
				current.Sequence = 0
				return callbackError
			})
			require.ErrorIs(t, err, callbackError)
			unchanged, err := jobs.Get[json.RawMessage](ctx, id)
			require.NoError(t, err)
			require.Equal(t, stored.Payload, unchanged.Payload)
			require.Equal(t, stored.Revision, unchanged.Revision)

			require.NoError(t, db.BackgroundJob.UpdateOneID(id).SetLeaseUntil(time.Now().Add(-time.Second)).Exec(ctx))
			called := false
			err = jobs.UpdatePayload(ctx, id, job.LeaseToken, func(current *sequencePayload) error {
				called = true
				return nil
			})
			require.ErrorIs(t, err, background_job.ErrLeaseLost)
			require.False(t, called)

			malformedID := xid.New()
			_, err = db.BackgroundJob.Create().
				SetID(malformedID).
				SetKey(malformedID.String()).
				SetKind("test.sequence").
				SetPayload(json.RawMessage(`{"sequence":"invalid"}`)).
				Save(ctx)
			require.NoError(t, err)
			err = jobs.UpdatePayload(ctx, malformedID, "", func(current *sequencePayload) error {
				called = true
				return nil
			})
			require.Error(t, err)
			require.False(t, called)
			malformed, err := jobs.Get[json.RawMessage](ctx, malformedID)
			require.NoError(t, err)
			require.Zero(t, malformed.Revision)
			require.JSONEq(t, `{"sequence":"invalid"}`, string(malformed.Payload))

			invalid, err := jobs.Get[sequencePayload](ctx, malformedID)
			require.Error(t, err)
			require.Nil(t, invalid)
		}))
	}))
}

func TestRepositoryUsesTransactionClient(t *testing.T) {
	t.Parallel()

	type valuePayload struct {
		Value string `json:"value"`
	}

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		jobs *background_job.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			rollback := errors.New("roll back job changes")
			for _, commit := range []bool{false, true} {
				id := xid.New()
				createdID := xid.New()
				require.NoError(t, jobs.Ensure(ctx, id, "test.transaction", id.String(), valuePayload{Value: "original"}))

				err := ent.WithTx(ctx, db, func(tx *ent.Tx) error {
					transactionJobs := background_job.New(tx.Client())
					err := transactionJobs.Ensure(ctx, createdID, "test.transaction", createdID.String(), valuePayload{})
					if err != nil {
						return err
					}
					if err := transactionJobs.Request(ctx, id, func(current *valuePayload) error {
						current.Value = "updated"
						return nil
					}, func(valuePayload) bool {
						return false
					}); err != nil {
						return err
					}

					stored, err := transactionJobs.Get[valuePayload](ctx, id)
					if err != nil {
						return err
					}
					require.Equal(t, "updated", stored.Payload.Value)
					if !commit {
						return rollback
					}
					return nil
				})

				if commit {
					require.NoError(t, err)
				} else {
					require.ErrorIs(t, err, rollback)
				}
				stored, err := jobs.Get[valuePayload](ctx, id)
				require.NoError(t, err)
				if commit {
					require.Equal(t, "updated", stored.Payload.Value)
				} else {
					require.Equal(t, "original", stored.Payload.Value)
				}
				_, err = jobs.Get[valuePayload](ctx, createdID)
				if commit {
					require.NoError(t, err)
				} else {
					require.True(t, ent.IsNotFound(err))
				}
			}
		}))
	}))
}
