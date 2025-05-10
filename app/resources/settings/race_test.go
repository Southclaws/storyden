package settings_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestConcurrentAccess(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		sr *settings.SettingsRepository,
	) {
		lc.Append(fx.StartHook(func(ctx context.Context) {
			t.Run("concurrent_access", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				content, err := datagraph.NewRichText("<body><p>Test content</p></body>")
				r.NoError(err)

				_, err = sr.Set(ctx, settings.Settings{
					Title:   opt.New("Test Title"),
					Content: opt.New(content),
				})
				r.NoError(err)

				const numGoroutines = 10
				const numRequests = 20

				var wg sync.WaitGroup
				wg.Add(numGoroutines)

				done := make(chan struct{})
				defer close(done)

				titles := make([]settings.Settings, numGoroutines)
				for i := 0; i < numGoroutines; i++ {
					titles[i] = settings.Settings{
						Title: opt.New(fmt.Sprintf("Title %d", i)),
					}
				}

				for i := 0; i < numGoroutines; i++ {
					go func(idx int) {
						defer wg.Done()

						for j := 0; j < numRequests; j++ {
							select {
							case <-done:
								return
							default:
								if j%2 == 0 {
									_, err := sr.Set(ctx, titles[idx])
									r.NoError(err)
								} else {
									_, err := sr.Get(ctx)
									r.NoError(err)
								}
							}
						}
					}(i)
				}

				wg.Wait()

				settings, err := sr.Get(ctx)
				r.NoError(err)
				r.NotNil(settings)
				a.NotEmpty(settings.Title.OrZero())
			})
		}))
	}))
}
