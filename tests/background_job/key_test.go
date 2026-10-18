package background_job_test

import (
	"context"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/backgroundjob"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestJobKeysDeduplicateWithinKind(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		jobs *background_job.Repository,
		db *ent.Client,
	) {
		lc.Append(fx.StartHook(func() {
			id := xid.New()
			key := id.String()
			err := jobs.Ensure(ctx, id, "test.first", key, payload{Requested: 1})
			require.NoError(t, err)

			err = jobs.Ensure(ctx, xid.New(), "test.first", key, payload{Requested: 2})
			require.NoError(t, err)
			stored, err := jobs.Get[payload](ctx, id)
			require.NoError(t, err)
			require.Equal(t, key, stored.Key)
			require.Equal(t, 1, stored.Payload.Requested)

			err = jobs.Ensure(ctx, xid.New(), "test.second", key, payload{})
			require.NoError(t, err)
			count, err := db.BackgroundJob.Query().Where(backgroundjob.KeyEQ(key)).Count(ctx)
			require.NoError(t, err)
			require.Equal(t, 2, count)

			err = jobs.Ensure(ctx, xid.New(), "test.first", "", payload{})
			require.Error(t, err)
		}))
	}))
}
