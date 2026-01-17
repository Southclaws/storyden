package integration

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"dario.cat/mergo"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources"
	"github.com/Southclaws/storyden/app/services"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure"
	"github.com/Southclaws/storyden/internal/utils"
)

// Test provides a full app setup for testing end to end behaviour. Example:
//
//	func TestMyThing(t *testing.T) {
//	    integration.Test(t, nil, fx.Invoke(func(test dependencies...) {
//	        r := require.New(t)
//	        a := assert.New(t)
//
//	        // your e2e test code...
//
//	    }))
//	}
func Test(t *testing.T, cfg *config.Config, o ...fx.Option) {
	defaultConfig := config.Config{
		PublicAPIAddress: *utils.Must(url.Parse("http://localhost")),
		PublicWebAddress: *utils.Must(url.Parse("http://localhost")),
		JWTSecret:        []byte("00000000000000000000000000000000"),
		RateLimit:        5000,
		RateLimitPeriod:  time.Hour,
		RateLimitBucket:  time.Minute,
		EmailProvider:    "mock",
	}

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		if isMaybeProdDB(dbURL) {
			panic("maybe accidental prod DATABASE_URL in integration tests!")
		}
		defaultConfig.DatabaseURL = dbURL
	} else {
		// Generate a unique database per test, avoids SQLite write contention.
		testDatabaseName := time.Now().Format(time.RFC3339) + t.Name()

		opts := url.Values{"_pragma": []string{
			"foreign_keys(1)",
			"busy_timeout(10000)",
			"journal_mode(WAL)",
			"synchronous(NORMAL)",
			"cache_size(1000000000)",
			"temp_store(MEMORY)",
		}}.Encode()

		defaultConfig.DatabaseURL = fmt.Sprintf(
			"sqlite://data/%s.db?%s",
			testDatabaseName,
			opts,
		)
	}

	ctx, cf := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cf()
	})

	o = append(o,
		// main application dependencies
		application(),

		// provide a global context
		fx.Provide(func() context.Context { return ctx }),
	)

	// if this test has a custom config, merge+overwrite with the defaults.
	if cfg != nil {
		mergo.MergeWithOverwrite(&defaultConfig, cfg)
	}

	o = append(o, fx.Provide(func() config.Config { return defaultConfig }))

	app := fx.New(o...)
	err := app.Start(ctx)
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}
}

// application gives you some basics needed by most components.
func application() fx.Option {
	return fx.Options(
		fx.NopLogger,

		infrastructure.Build(),
		resources.Build(),
		services.Build(),
	)
}

func isMaybeProdDB(url string) bool {
	dangerous := []string{
		"free-tier",
		".aws-eu-central",
		"cockroachlabs",
		"cloud",
		"verify-full",
	}

	for _, v := range dangerous {
		if strings.Contains(url, v) {
			return true
		}
	}

	return false
}
