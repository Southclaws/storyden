package seed

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/infrastructure/db"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/pkg/resources/account"
	"github.com/Southclaws/storyden/pkg/resources/category"
)

// NOTE: identifiers in the system use the xid format. This format has a couple
// of checks when reading from a string format. Because of this, the constant ID
// string literals used in the seed data are written to work properly but also
// be super simple and readable for debugging purposes when working with seed
// data. The format is just avoiding setting the final character so the first 2
// characters of the final section are used. In the documentation this is
// referred to as: "3-byte counter, starting with a random value."
var id = func(s string) xid.ID { return utils.Must(xid.FromString(s)) }

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

	accounts(account_repo)
	categories(category_repo)

	return Ready{}, err
}
