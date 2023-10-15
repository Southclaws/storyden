package integration

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/imdario/mergo"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources"
	"github.com/Southclaws/storyden/app/services"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure"
)

// Test provides a full app setup for testing service behaviour. It returns a
// context cancellation function for immediate shutdown once all test functions
// have finished. Usage is a simple call and defer:
//
//	func TestMyThing(t *testing.T) {
//	    defer bdd.Test(t, nil, fx.Invoke(func(test dependencies...) {
//	        r := require.New(t)
//	        a := assert.New(t)
//
//	        your behavioural test code...
//
//	    }))
//	}
func Test(t *testing.T, cfg *config.Config, o ...fx.Option) func() {
	defaultConfig := config.Config{
		CookieDomain:     "localhost",
		PublicWebAddress: "http://localhost",
	}

	if url := os.Getenv("DATABASE_URL"); url != "" {
		if isMaybeProdDB(url) {
			panic("maybe accidental prod DATABASE_URL in integration tests!")
		}
		defaultConfig.DatabaseURL = url
	} else {
		defaultConfig.DatabaseURL = "postgresql://default:default@localhost:5432/postgres"
	}

	ctx, cf := context.WithCancel(context.Background())

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

	return cf
}

// application gives you some basics needed by most components.
func application() fx.Option {
	return fx.Options(
		fx.NopLogger,

		infrastructure.Build(false),
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
