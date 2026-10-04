package settings_test

import (
	"context"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestSettingsRepository(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(lc fx.Lifecycle, sr *settings.SettingsRepository) {
		lc.Append(fx.StartHook(func(ctx context.Context) {
			t.Run("partial_update", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				content, err := datagraph.NewRichText("<body><p>Hello, Makeroom!</p></body>")
				r.NoError(err)

				set, err := sr.Set(ctx, settings.Settings{
					Title:   opt.New("Makeroom"),
					Content: opt.New(content),
				})
				r.NoError(err)
				r.NotNil(set)

				got, err := sr.Get(ctx)
				r.NoError(err)
				r.NotNil(got)

				a.Equal("Makeroom", got.Title.OrZero())
				a.Equal(content.HTML(), got.Content.OrZero().HTML())
				a.Equal(settings.DefaultDescription, got.Description.OrZero())
			})
		}))
	}))
}

func TestOAuthSettingsSnapshot(t *testing.T) {
	t.Parallel()
	integration.Test(t, nil, fx.Invoke(func(lc fx.Lifecycle, sr *settings.SettingsRepository) {
		lc.Append(fx.StartHook(func(ctx context.Context) {
			snapshot, err := sr.Get(ctx)
			require.NoError(t, err)

			// Independent concurrent patches must both survive, without changing a
			// settings snapshot already held by an in-flight request.
			start := make(chan struct{})
			results := make(chan error, 2)
			for _, patch := range []settings.OAuthServiceSettings{
				{AutonomousRegistrationMode: opt.New(settings.OAuthAutonomousRegistrationModeApproval)},
				{RegistrationApprovalTTL: opt.New(7 * time.Minute)},
			} {
				go func() {
					<-start
					_, err := sr.Set(ctx, settings.Settings{Services: opt.New(settings.ServiceSettings{OAuth: opt.New(patch)})})
					results <- err
				}()
			}
			close(start)
			for range 2 {
				require.NoError(t, <-results)
			}

			old := snapshot.Services.OrZero().OAuth.OrZero()
			require.Equal(t, "disabled", old.AutonomousRegistrationMode.OrZero().String())
			require.Equal(t, 10*time.Minute, old.RegistrationApprovalTTL.OrZero())
			current, err := sr.Get(ctx)
			require.NoError(t, err)
			updated := current.Services.OrZero().OAuth.OrZero()
			require.Equal(t, "approval", updated.AutonomousRegistrationMode.OrZero().String())
			require.Equal(t, 7*time.Minute, updated.RegistrationApprovalTTL.OrZero())
		}))
	}))
}
