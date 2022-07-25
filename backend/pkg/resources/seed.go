package resources

import (
	"fmt"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/backend/pkg/resources/user"
)

// Seeded is a type you can depend on during integration tests which, when used
// will ensure the database is seeded with data before your tests run.
// Usage is simple, use `bdd.Test` and add this value as a parameter to the
// test function invoke call:
//
// bdd.Test(t, nil, fx.Invoke(
//     func(
//         _ resources.Seeded,
//         repo user.Repository,
//     ) {
//         ... your test code
//
type Seeded struct{}

// Seed provides a type to the application which, when present in a component's
// dependency tree, will seed the database with all resource seed data.
func Seed() fx.Option {
	return fx.Options(
		fx.Provide(func(
			lc fx.Lifecycle,
			client *model.Client,
			user_repo user.Repository,
		) Seeded {
			defer func() {
				// recover panics so that test cleanups can run.
				if r := recover(); r != nil {
					fmt.Println(r)
					return
				}
			}()

			user.Seed(user_repo)

			return Seeded{}
		}),
	)
}
