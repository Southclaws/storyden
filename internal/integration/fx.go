package integration

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	"dario.cat/mergo"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources"
	"github.com/Southclaws/storyden/app/services"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure"
	"github.com/Southclaws/storyden/internal/utils"
)

// Test provides a full app setup for testing service behaviour. It returns a
// context cancellation function for immediate shutdown once all test functions
// have finished. Usage is a simple call and defer:
//
//	func TestMyThing(t *testing.T) {
//	    defer integration.Test(t, nil, fx.Invoke(func(test dependencies...) {
//	        r := require.New(t)
//	        a := assert.New(t)
//
//	        your behavioural test code...
//
//	    }))
//	}
func Test(t *testing.T, cfg *config.Config, o ...fx.Option) {
	defaultConfig := config.Config{
		PublicAPIAddress: *utils.Must(url.Parse("http://localhost")),
		PublicWebAddress: *utils.Must(url.Parse("http://localhost")),
	}

	if url := os.Getenv("DATABASE_URL"); url != "" {
		if isMaybeProdDB(url) {
			panic("maybe accidental prod DATABASE_URL in integration tests!")
		}
		defaultConfig.DatabaseURL = url
	} else {
		defaultConfig.DatabaseURL = "sqlite://data.db?_pragma=foreign_keys(1)&_pragma=busy_timeout(1000)"
	}

	ctx, cf := context.WithCancel(context.Background())
	t.Cleanup(func() {
		fmt.Println("integration test cleanup")
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

	err := fx.New(o...).Start(ctx)
	if err != nil {
		fmt.Println(err)
		t.Error()
	}

	return
}

// application gives you some basics needed by most components.
func application() fx.Option {
	return fx.Options(
		// fx.NopLogger,

		infrastructure.Build(),
		resources.Build(),
		services.Build(),

		// Tests can depend on Migrated to trigger migrations pre-test.
		// This is not parallel safe.
		fx.Provide(WithMigrated),
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

type Migrated interface{}

func WithMigrated(ctx context.Context, client *ent.Client) (Migrated, error) {
	if err := client.Schema.Create(ctx); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return 1, nil
}
