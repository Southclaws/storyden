package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	_ "github.com/glebarez/go-sqlite"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
)

func Build() fx.Option {
	return fx.Options(
		// provide the underlying *sql.DB and sqlx to the system
		fx.Provide(newSQL),

		// construct a new ent client using the *sql.DB provided above
		fx.Provide(newEntClient),
	)
}

func newSQL(cfg config.Config) (*sql.DB, *sqlx.DB, error) {
	driver, path, err := getDriver(cfg.DatabaseURL)
	if err != nil {
		return nil, nil, fault.Wrap(err)
	}

	d, err := sql.Open(driver, path)
	if err != nil {
		return nil, nil, fault.Wrap(err, fmsg.With("failed to connect to database"))
	}

	x, err := sqlx.Connect(driver, path)
	if err != nil {
		return nil, nil, fault.Wrap(err, fmsg.With("failed to connect to database"))
	}

	return d, x, nil
}

// This is only used in tests to allow simple concurrent tests without needing
// to write too much test-specific code for DB stuff. We should use enttest tbh.
var schema = sync.Mutex{}

func newEntClient(lc fx.Lifecycle, cfg config.Config, db *sql.DB) (*ent.Client, error) {
	wctx, cancel := context.WithCancel(context.Background())

	client, err := connect(wctx, cfg, db)
	if err != nil {
		cancel()
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			schema.Lock()
			defer schema.Unlock()

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

func connect(ctx context.Context, cfg config.Config, driver *sql.DB) (*ent.Client, error) {
	d, _, err := getDriver(cfg.DatabaseURL)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	opts := []ent.Option{}

	switch d {
	case "pgx":
		opts = append(opts, ent.Driver(entsql.OpenDB(dialect.Postgres, driver)))

	case "sqlite":
		opts = append(opts,
			ent.Driver(entsql.OpenDB(dialect.SQLite, driver)),
		)

	default:
		panic(fmt.Sprintf("unsupported driver '%s' in ent connect", d))
	}

	return ent.NewClient(opts...), nil
}

func getDriver(databaseURL string) (string, string, error) {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return "", "", fault.Wrap(err, fmsg.With("failed to parse DATABASE_URL"))
	}

	switch u.Scheme {
	case "postgres", "postgresql":
		return "pgx", databaseURL, nil

	case "sqlite", "sqlite3":
		path, _ := strings.CutPrefix(databaseURL, u.Scheme+"://")
		return "sqlite", path, nil

	default:
		return "", "", fault.Newf("unsupported scheme: %s", u.Scheme)
	}
}
