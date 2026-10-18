package asset_test

import (
	"bytes"
	"context"
	"image/png"
	"io"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_repo"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_image"
	"github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestAvatarReadsUseCachedAccount(t *testing.T) {
	t.Parallel()

	type queryCountKey struct{}
	type queryCounts struct {
		accounts atomic.Int64
		assets   atomic.Int64
		roles    atomic.Int64
	}

	integration.Test(t, nil, e2e.Setup(),
		fx.Decorate(func(db *ent.Client) *ent.Client {
			db.Intercept(ent.InterceptFunc(func(next ent.Querier) ent.Querier {
				return ent.QuerierFunc(func(ctx context.Context, query ent.Query) (ent.Value, error) {
					if counts, ok := ctx.Value(queryCountKey{}).(*queryCounts); ok {
						switch query.(type) {
						case *ent.AccountQuery:
							counts.accounts.Add(1)
						case *ent.AssetQuery:
							counts.assets.Add(1)
						case *ent.RoleQuery, *ent.AccountRolesQuery:
							counts.roles.Add(1)
						}
					}
					return next.Query(ctx, query)
				})
			}))
			return db
		}),
		fx.Invoke(func(
			lc fx.Lifecycle,
			ctx context.Context,
			accounts *account_repo.Repository,
			writer *account_writer.Writer,
			jobs *asset_image.Repository,
			avatars avatar.Service,
			store cache.Store,
		) {
			lc.Append(fx.StartHook(func() {
				ctx, owner := e2e.WithAccount(ctx, writer, seed.Account_001_Odin)
				payload := imagePayload(t)
				require.NoError(t, avatars.Set(ctx, owner.ID, bytes.NewReader(payload), int64(len(payload))))
				current, err := accounts.GetRefByID(ctx, owner.ID)
				require.NoError(t, err)
				require.True(t, current.AvatarAssetID.Ok())

				require.Eventually(t, func() bool {
					job, err := jobs.Get(ctx, current.AvatarAssetID.OrZero())
					return err == nil && job.Status == background_job.StatusComplete
				}, 20*time.Second, 50*time.Millisecond)

				require.NoError(t, store.Delete(ctx, "account:data:"+owner.ID.String()))
				var coldQueries queryCounts
				coldContext := context.WithValue(ctx, queryCountKey{}, &coldQueries)
				reloaded, err := accounts.GetRefByID(coldContext, owner.ID)
				require.NoError(t, err)
				require.Equal(t, current.AvatarAssetID, reloaded.AvatarAssetID)
				require.EqualValues(t, 1, coldQueries.accounts.Load())

				for _, test := range []struct {
					variant asset.Variant
					queries int64
				}{
					{variant: asset.Original, queries: 1},
					{variant: asset.ThumbnailMedium, queries: 2},
				} {
					t.Run(string(test.variant), func(t *testing.T) {
						var queries queryCounts
						readContext := context.WithValue(ctx, queryCountKey{}, &queries)
						require.True(t, avatars.Exists(readContext, owner.ID))
						_, stream, pending, err := avatars.GetVariant(readContext, owner.ID, test.variant)
						require.NoError(t, err)
						require.False(t, pending)
						_, err = io.Copy(io.Discard, stream)
						require.NoError(t, err)
						if closer, ok := stream.(io.Closer); ok {
							require.NoError(t, closer.Close())
						}
						require.Zero(t, queries.accounts.Load())
						require.Equal(t, test.queries, queries.assets.Load())
						t.Logf("role hydration queries: %d", queries.roles.Load())
					})
				}

				require.NoError(t, avatars.Delete(ctx, owner.ID))
				var defaultQueries queryCounts
				defaultContext := context.WithValue(ctx, queryCountKey{}, &defaultQueries)
				require.False(t, avatars.Exists(defaultContext, owner.ID))
				_, stream, pending, err := avatars.GetVariant(defaultContext, owner.ID, asset.ThumbnailMedium)
				require.NoError(t, err)
				require.False(t, pending)
				_, err = png.Decode(stream)
				require.NoError(t, err)
				require.NoError(t, stream.(io.Closer).Close())
				require.Zero(t, defaultQueries.accounts.Load())
				require.Zero(t, defaultQueries.assets.Load())
				t.Logf("default avatar role hydration queries: %d", defaultQueries.roles.Load())
			}))
		}),
	)
}
