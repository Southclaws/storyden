package seed

import (
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/db"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/backend/pkg/resources/account"
	"github.com/Southclaws/storyden/backend/pkg/resources/category"
)

// Ready is a type you can depend on during integration tests which, when used
// will ensure the database is seeded with data before your tests run.
// Usage is simple, use `bdd.Test` and add this value as a parameter to the
// test function invoke call:
//
// bdd.Test(t, nil, fx.Invoke(
//     func(
//         _ seed.Ready,
//         repo user.Repository,
//     ) {
//         ... your test code
//
type Ready struct{}

// Seed provides a type to the application which, when present in a component's
// dependency tree, will seed the database with all resource seed data.
func Create() fx.Option {
	fmt.Println("Seed constructor called\n\n---")

	return fx.Options(
		fx.Provide(New),
		fx.Invoke(func(_ Ready) {}),
	)
}

// New runs the data seeding script, creating all fake data for testing/demos.
func New(
	database *sql.DB,
	client *model.Client,
	account_repo account.Repository,
	category_repo category.Repository,
) (r Ready, err error) {
	defer func() {
		// recover panics so that test cleanups can run.
		if r := recover(); r != nil {
			fmt.Println(r)

			err = errors.New("failed to seed")
		}
	}()

	if err := db.Truncate(database); err != nil {
		panic(err)
	}

	fmt.Println("seeding database")

	users(account_repo)
	categories(category_repo)

	return Ready{}, err
}
