package db

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert --feature sql/modifier --feature sql/upsert --feature sql/versioned-migration ./schema --target ./model

import (
	"context"
	"database/sql"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/fx"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
)

func Build() fx.Option {
	return fx.Provide(newDB)
}

func newDB(lc fx.Lifecycle, cfg config.Config) (*model.Client, *sql.DB, error) {
	wctx, cancel := context.WithCancel(context.Background())

	client, db, err := connect(wctx, cfg.DatabaseURL)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			defer cancel()
			return client.Close()
		},
	})

	return client, db, nil
}

func connect(ctx context.Context, url string) (*model.Client, *sql.DB, error) {
	driver, err := sql.Open("pgx", url)
	if err != nil {
		return nil, nil, fault.Wrap(err, fmsg.With("failed to connect to database"))
	}

	client := model.NewClient(model.Driver(entsql.OpenDB(dialect.Postgres, driver)))

	opts := []schema.MigrateOption{
		schema.WithAtlas(true),
	}

	// We don't do versioned migrations currently.
	opts = append(opts, schema.WithDropColumn(true))
	opts = append(opts, schema.WithDropIndex(true))

	// Run only additive migrations
	if err := client.Schema.Create(ctx, opts...); err != nil {
		return nil, nil, err
	}

	return client, driver, nil
}
