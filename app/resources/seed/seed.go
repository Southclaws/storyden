package seed

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph/node"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/db"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/utils"
)

// NOTE: identifiers in the system use the xid format. This format has a couple
// of checks when reading from a string format. Because of this, the constant ID
// string literals used in the seed data are written to work properly but also
// be super simple and readable for debugging purposes when working with seed
// data. The format is just avoiding setting the final character so the first 2
// characters of the final section are used. In the documentation this is
// referred to as: "3-byte counter, starting with a random value."
var id = func(s string) xid.ID { return utils.Must(xid.FromString(s)) }

type Ready struct{}

type Empty struct{}

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
	client *ent.Client,
	settings settings.Repository,
	account_repo account.Repository,
	auth_repo authentication.Repository,
	category_repo category.Repository,
	thread_repo thread.Repository,
	post_repo reply.Repository,
	react_repo react.Repository,
	asset_repo asset.Repository,
	node_repo node.Repository,
) (r Ready) {
	if err := client.Schema.Create(context.Background()); err != nil {
		panic(err)
	}

	if err := db.Truncate(database); err != nil {
		panic(err)
	}

	fmt.Println("seeding database")

	utils.Must[any](nil, settings.Init(context.Background()))

	accounts(account_repo, auth_repo)
	categories(category_repo)
	threads(thread_repo, post_repo, react_repo, asset_repo)

	return Ready{}
}

func NewEmpty(database *sql.DB) Empty {
	if err := db.Truncate(database); err != nil {
		panic(err)
	}

	return Empty{}
}
