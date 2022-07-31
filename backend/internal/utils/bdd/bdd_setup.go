package bdd

import (
	"context"
	"fmt"
	"testing"

	"github.com/imdario/mergo"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/internal/config"
	"github.com/Southclaws/storyden/backend/internal/infrastructure"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/backend/pkg/resources"
	"github.com/Southclaws/storyden/backend/pkg/resources/seed"
	"github.com/Southclaws/storyden/backend/pkg/services"
)

// Test provides a BDD style setup for testing service behaviour. It returns a
// context cancellation function for immediate shutdown once all test functions
// have finished. Usage is a simple call and defer:
//
// func TestMyThing(t *testing.T) {
//     ctx := context.Background()
//     defer bdd.Test(t, nil, fx.Invoke(func(test dependencies...) {
//         r := require.New(t)
//         a := assert.New(t)
//
//         your behavioural test code...
//
//     }))
// }
//
func Test(t *testing.T, cfg *config.Config, o ...fx.Option) func() {
	defaultConfig := config.Config{}
	defaultConfig.DatabaseURL = "postgresql://default:default@localhost:5432/postgres"

	o = append(o,
		// main application dependencies
		application(),

		// seeded database
		seed.Create(),

		// database client
		fx.Invoke(func() *model.Client { return db.TestDB(t) }),
	)

	// if this test has a custom config, merge+overwrite with the defaults.
	if cfg != nil {
		mergo.MergeWithOverwrite(&defaultConfig, cfg)
	}

	o = append(o, fx.Provide(func() config.Config { return defaultConfig }))

	ctx, cf := context.WithCancel(context.Background())

	err := fx.New(o...).Start(ctx)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	return cf
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
