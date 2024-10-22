package settings_test

import (
	"context"
	"testing"

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
