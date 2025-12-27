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

	"ariga.io/atlas/sql/migrate"
	atlas_schema "ariga.io/atlas/sql/schema"
	entgo "entgo.io/ent"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	_ "github.com/glebarez/go-sqlite"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
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

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			schemaLock.Lock()
			defer schemaLock.Unlock()

			// Run migrations with hooks and index cleanup.
			if err := client.Schema.Create(
				ctx,
				schema.WithDropIndex(true),
				schema.WithDropColumn(true),
				schema.WithApplyHook(populateLastReplyAt()),
				schema.WithApplyHook(migrateReplyVisibility()),
			); err != nil {
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

	case "libsql":
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

		// Try to write to the directory. This provides a better error message
		// compared to SQLite which will give you nonsense if it can't write.
		testwrite := filepath.Join(filepath.Dir(path), ".perm_check")
		if err := os.WriteFile(testwrite, []byte("ok"), 0o644); err != nil {
			return "", "", fault.Wrap(err, fmsg.With(fmt.Sprintf("cannot write to directory for sqlite database: %s", u)))
		}

		return "sqlite", path, nil

	case "libsql":
		// NOTE: Only remote Turso, local file-based libSQL is not supported.
		return "libsql", databaseURL, nil

	default:
		return "", "", fault.Newf("unsupported scheme: %s", u.Scheme)
	}
}

// populateLastReplyAt is a data migration hook that fills NULL last_reply_at values
// with created_at for threads. This only runs when the last_reply_at column is being
// modified (e.g., changing from nullable to non-nullable).
//
// This is a bit of a hack because there's no versioned migrations set up now.
// It shouldn't run again after first run though, and if we change the column
// again at some point in the future, this hook will just be removed.
func populateLastReplyAt() schema.ApplyHook {
	return func(next schema.Applier) schema.Applier {
		return schema.ApplyFunc(func(ctx context.Context, conn dialect.ExecQuerier, plan *migrate.Plan) error {
			// Check if the last_reply_at column is being modified in this migration
			hasChange := func() bool {
				for _, c := range plan.Changes {
					m, ok := c.Source.(*atlas_schema.ModifyTable)
					if ok && m.T.Name == ent_post.Table {
						// Check if last_reply_at column is being modified
						if atlas_schema.Changes(m.Changes).IndexModifyColumn(ent_post.FieldLastReplyAt) != -1 {
							return true
						}
					}
				}
				return false
			}()

			if hasChange {
				err := conn.Exec(ctx, `
					UPDATE posts
					SET last_reply_at = created_at
					WHERE last_reply_at IS NULL
				`, []any{}, nil)
				if err != nil {
					return fault.Wrap(err, fmsg.With("failed to populate last_reply_at"))
				}
			}

			return next.Apply(ctx, conn, plan)
		})
	}
}

// migrateReplyVisibility is a data migration hook that updates all replies in
// draft visibility to published visibility. This is needed for upgrades from
// versions ≤ v1.25.12 where replies defaulted to 'draft' visibility (even
// though draft replies were not a functional feature at that time).
//
// Starting in v1.25.12, replies are created with 'published' visibility, and
// v1.25.14+ adds content moderation with 'review' visibility. This migration
// ensures old draft replies don't disappear when visibility filtering is done.
//
// Safe to run unconditionally because:
// - Draft replies were never a functional feature before v1.25.14
// - Only affects replies (posts with root_post_id set)
// - Idempotent (can run multiple times safely)
//
// TODO: Remove this hook after v1.26.0 once version tracking is implemented.
func migrateReplyVisibility() schema.ApplyHook {
	return func(next schema.Applier) schema.Applier {
		return schema.ApplyFunc(func(ctx context.Context, conn dialect.ExecQuerier, plan *migrate.Plan) error {
			// Always run this migration on schema creation/update.
			// It's idempotent and safe since draft replies weren't functional.
			err := conn.Exec(ctx, `
				UPDATE posts
				SET visibility = 'published'
				WHERE root_post_id IS NOT NULL AND visibility = 'draft'
			`, []any{}, nil)
			if err != nil {
				return fault.Wrap(err, fmsg.With("failed to migrate reply visibility from draft to published"))
			}

			return next.Apply(ctx, conn, plan)
		})
	}
}
