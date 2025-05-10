package settings_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/integration"
)

// TestCacheRace specifically tests for the race condition in the cache() method
// where multiple goroutines were writing to the cacheLastFetch field concurrently.
func TestCacheRace(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(lc fx.Lifecycle, sr *settings.SettingsRepository) {
		lc.Append(fx.StartHook(func(ctx context.Context) {
			t.Run("cache_race", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				// Number of goroutines to test with
				const numGoroutines = 20
				const numRequests = 50

				// Create a wait group to sync all goroutines
				var wg sync.WaitGroup
				wg.Add(numGoroutines)

				// Channel to signal completion
				done := make(chan struct{})
				defer close(done)

				// Start multiple goroutines that all call Get() at the same time
				// This will cause multiple cache() calls to happen concurrently
				for i := 0; i < numGoroutines; i++ {
					go func() {
						defer wg.Done()

						// Make multiple requests from each goroutine
						for j := 0; j < numRequests; j++ {
							select {
							case <-done:
								return
							default:
								// Call Get which may trigger cache() internally
								settings, err := sr.Get(ctx)
								if err != nil {
									t.Errorf("Error getting settings: %v", err)
									return
								}
								a.NotNil(settings)

								// Small pause to increase chances of race conditions
								time.Sleep(time.Millisecond)
							}
						}
					}()
				}

				// Wait for all goroutines to complete
				wg.Wait()

				// Final verification that we can still get settings
				settings, err := sr.Get(ctx)
				r.NoError(err)
				r.NotNil(settings)

				// The test passes if we don't have any race detector failures
			})
		}))
	}))
}