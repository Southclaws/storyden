package db

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert --feature sql/modifier --feature sql/upsert --feature sql/versioned-migration ./schema --target ./ent

import (
	"context"
	"database/sql"

	"entgo.io/ent/dialect"             // nolint:gci
	entsql "entgo.io/ent/dialect/sql"  // nolint:gci
	"entgo.io/ent/dialect/sql/schema"  // nolint:gci
	"github.com/Southclaws/fault"      // nolint:gci
	"github.com/Southclaws/fault/fctx" // nolint:gci
	"github.com/Southclaws/fault/fmsg" // nolint:gci
	_ "github.com/jackc/pgx/v4/stdlib" // nolint:gci
	"go.uber.org/fx"                   // nolint:gci

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
)

func Build() fx.Option {
	return fx.Provide(newDB)
}

func newDB(lc fx.Lifecycle, cfg config.Config) (*ent.Client, *sql.DB, error) {
	wctx, cancel := context.WithCancel(context.Background())

	client, db, err := connect(wctx, cfg.DatabaseURL)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			defer cancel()

			err := client.Close()
			if err != nil {
				return fault.Wrap(err, fctx.With(ctx))
			}

			return nil
		},
	})

	return client, db, nil
}

func connect(ctx context.Context, url string) (*ent.Client, *sql.DB, error) {
	driver, err := sql.Open("pgx", url)
	if err != nil {
		return nil, nil, fault.Wrap(err, fmsg.With("failed to connect to database"))
	}

	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, driver)))

	opts := []schema.MigrateOption{
		schema.WithAtlas(true),
	}

	// We don't do versioned migrations currently.
	// opts = append(opts, schema.WithDropColumn(true))
	// opts = append(opts, schema.WithDropIndex(true))

	// Run only additive migrations
	if err := client.Schema.Create(ctx, opts...); err != nil {
		return nil, nil, fault.Wrap(err)
	}

	return client, driver, nil
}
