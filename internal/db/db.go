package db

import (
	"context"
	"database/sql"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
)

func Build() fx.Option {
	return fx.Options(
		// provide the underlying *sql.DB to the system
		fx.Provide(newSQL),

		// provide sqlx to make hand written queries slightly less painful
		fx.Provide(newSQLX),

		// construct a new ent client using the *sql.DB provided above
		fx.Provide(newEntClient),
	)
}

func newSQL(cfg config.Config) (*sql.DB, error) {
	driver, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to connect to database"))
	}

	return driver, nil
}

func newSQLX(cfg config.Config) (*sqlx.DB, error) {
	driver, err := sqlx.Connect("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to connect to database"))
	}

	return driver, nil
}

func newEntClient(lc fx.Lifecycle, db *sql.DB) (*ent.Client, error) {
	wctx, cancel := context.WithCancel(context.Background())

	client, err := connect(wctx, db)
	if err != nil {
		cancel()
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Run create-only migrations after initialisation.
			// This is done in tests and scripts too.
			if err := client.Schema.Create(ctx); err != nil {
				return fault.Wrap(err, fctx.With(ctx))
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			defer cancel()

			err := client.Close()
			if err != nil {
				return fault.Wrap(err, fctx.With(ctx))
			}

			return nil
		},
	})

	return client, nil
}

func connect(ctx context.Context, driver *sql.DB) (*ent.Client, error) {
	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, driver)))

	return client, nil
}
