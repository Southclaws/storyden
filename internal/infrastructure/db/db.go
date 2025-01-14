package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	entgo "entgo.io/ent"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	_ "github.com/glebarez/go-sqlite"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/tracing"
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
var schemaLock = sync.Mutex{}

func newEntClient(lc fx.Lifecycle, tf tracing.Factory, cfg config.Config, db *sql.DB) (*ent.Client, error) {
	wctx, cancel := context.WithCancel(context.Background())

	client, err := connect(wctx, cfg, db)
	if err != nil {
		cancel()
		return nil, err
	}

	tr := tf.Build(lc, "ent")

	client.Intercept(ent.InterceptFunc(func(next ent.Querier) ent.Querier {
		return ent.QuerierFunc(func(ctx context.Context, query ent.Query) (ent.Value, error) {
			qc := entgo.QueryFromContext(ctx)
			spanName := fmt.Sprintf("ent/%s/%s", qc.Op, qc.Type)

			ctx, span := tr.Start(ctx, spanName, trace.WithAttributes(
				attribute.String("type", qc.Type),
				attribute.String("op", qc.Op),
				attribute.Bool("unique", opt.NewPtr(qc.Unique).OrZero()),
				attribute.Int("limit", opt.NewPtr(qc.Limit).OrZero()),
				attribute.Int("offset", opt.NewPtr(qc.Offset).OrZero()),
				attribute.StringSlice("fields", qc.Fields),
			))
			defer span.End()

			return next.Query(ctx, query)
		})
	}))

	client.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			spanName := fmt.Sprintf("ent/%s/%s", m.Op(), m.Type())

			ctx, span := tr.Start(ctx, spanName, trace.WithAttributes(
				attribute.String("type", m.Type()),
				attribute.String("op", m.Op().String()),
				attribute.StringSlice("fields", m.Fields()),
				attribute.StringSlice("added_edges", m.AddedEdges()),
				attribute.StringSlice("added_fields", m.AddedFields()),
				attribute.StringSlice("removed_edges", m.RemovedEdges()),
			))
			defer span.End()

			return next.Mutate(ctx, m)
		})
	})

	// client.Intercept(ent.InterceptFunc(func(next ent.Querier) ent.Querier {
	// 	return ent.QuerierFunc(func(ctx context.Context, query ent.Query) (ent.Value, error) {
	// 		qc := entgo.QueryFromContext(ctx)
	// 		spanName := fmt.Sprintf("ent %s %s", qc.Op, qc.Type)

	// 		start := time.Now()
	// 		id := xid.New()
	// 		t := reflect.ValueOf(query)

	// 		defer func() {
	// 			logger.Debug("END   "+spanName,
	// 				zap.String("id", id.String()),
	// 				zap.String("type", t.Elem().String()),
	// 				zap.Duration("duration", time.Since(start)),
	// 			)
	// 		}()

	// 		logger.Debug("BEGIN "+spanName,
	// 			zap.String("id", id.String()),
	// 			zap.String("type", t.Elem().String()),
	// 		)

	// 		v, err := next.Query(ctx, query)
	// 		if err != nil {
	// 			return nil, fault.Wrap(err, fctx.With(ctx))
	// 		}

	// 		return v, nil
	// 	})
	// }))

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			schemaLock.Lock()
			defer schemaLock.Unlock()

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

		// NOTE: SQLite has a bug where if the path does not exist, it provides
		// an incorrect and confusing error message about memory allocation. So
		// we need to perform the checks against the path with a proper error.
		if _, err := os.Stat(filepath.Dir(path)); err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					return "", "", fault.Wrap(err, fmsg.With(fmt.Sprintf("could not create directory for sqlite database: %s", u)))
				}
			} else {
				return "", "", fault.Wrap(err, fmsg.With(fmt.Sprintf("could not read directory: %s", u)))
			}
		}

		return "sqlite", path, nil

	default:
		return "", "", fault.Newf("unsupported scheme: %s", u.Scheme)
	}
}
