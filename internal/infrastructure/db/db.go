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
	ent_account "github.com/Southclaws/storyden/internal/ent/account"
	ent_email "github.com/Southclaws/storyden/internal/ent/email"
	ent_oauth_device_authorisation "github.com/Southclaws/storyden/internal/ent/oauthdeviceauthorisation"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
	ent_robot_session "github.com/Southclaws/storyden/internal/ent/robotsession"
	ent_robot_session_message "github.com/Southclaws/storyden/internal/ent/robotsessionmessage"
	ent_session "github.com/Southclaws/storyden/internal/ent/session"
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
		return nil, nil, fault.Wrap(err, fmsg.With("failed to open driver"))
	}

	if err := d.Ping(); err != nil {
		_ = d.Close()
		return nil, nil, fault.Wrap(err, fmsg.With("failed to ping database"))
	}

	x := sqlx.NewDb(d, driver)
	err = x.Ping()
	if err != nil {
		_ = d.Close()
		return nil, nil, fault.Wrap(err, fmsg.With("failed to connect to database"))
	}

	return d, x, nil
}

// This is only used in tests to allow simple concurrent tests without needing
// to write too much test-specific code for DB stuff. We should use enttest tbh.
var schemaLock = sync.Mutex{}

func newEntClient(lc fx.Lifecycle, tf tracing.Factory, cfg config.Config, db *sql.DB) (*ent.Client, error) {
	wctx, cancel := context.WithCancel(context.Background())

	client, driver, err := connect(cfg, db)
	if err != nil {
		cancel()
		return nil, fault.Wrap(err, fctx.With(wctx), fmsg.With("failed to connect ent client"))
	}

	tr := tf.Build(lc)

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

			// June 2026: Turso seems to have introduced an undocumented change
			// where foreign key constraints are no longer enabled by default,
			// and they don't seem to be able to be enabled via connection
			// string pragmas. So we have to enable them manually...
			if driver == "libsql" {
				_, err = db.ExecContext(ctx, "PRAGMA foreign_keys = ON")
				if err != nil {
					return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to enable foreign key constraints for libsql (remote)"))
				}
			}

			// Run migrations with hooks and index cleanup.
			if err := client.Schema.Create(
				ctx,
				schema.WithDropIndex(true),
				schema.WithDropColumn(true),
				schema.WithHooks(migrateRobotSessionMessageSequences(db, driver)),
				schema.WithApplyHook(populateLastReplyAt()),
				schema.WithApplyHook(migrateReplyVisibility()),
				schema.WithApplyHook(migrateAccountVerifiedStatus()),
				schema.WithApplyHook(migrateOAuthDeviceUserCodeUniqueness()),
				schema.WithApplyHook(migrateSessionTokenHash()),
			); err != nil {
				return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to run schema migration"))
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

func migrateRobotSessionMessageSequences(db *sql.DB, driver string) schema.Hook {
	return func(next schema.Creator) schema.Creator {
		return schema.CreateFunc(func(ctx context.Context, tables ...*schema.Table) error {
			if err := prepareRobotSessionMessageSequences(ctx, db, driver); err != nil {
				return fault.Wrap(err, fmsg.With("failed to backfill robot session message sequences"))
			}
			if err := next.Create(ctx, tables...); err != nil {
				return err
			}
			if err := reconcileRobotSessionMessageSequences(ctx, db); err != nil {
				return fault.Wrap(err, fmsg.With("failed to reconcile robot session message sequences"))
			}
			return nil
		})
	}
}

func prepareRobotSessionMessageSequences(ctx context.Context, db *sql.DB, driver string) error {
	migrationRequired, err := robotSessionMessageSequenceMigrationRequired(ctx, db, driver)
	if err != nil {
		return err
	}
	if !migrationRequired {
		return nil
	}

	columnType := "INTEGER"
	if driver == "pgx" {
		columnType = "BIGINT"
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, fmt.Sprintf(
		"ALTER TABLE %s ADD COLUMN %s %s NOT NULL DEFAULT 0",
		ent_robot_session_message.Table,
		ent_robot_session_message.FieldSequence,
		columnType,
	)); err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, fmt.Sprintf(
		`
		WITH ranked AS (
			SELECT %s, ROW_NUMBER() OVER (
				PARTITION BY %s
				ORDER BY %s, %s
			) AS sequence
			FROM %s
		)
		UPDATE %s
		SET %s = (
			SELECT ranked.sequence
			FROM ranked
			WHERE ranked.%s = %s.%s
		)
	`,
		ent_robot_session_message.FieldID,
		ent_robot_session_message.FieldSessionID,
		ent_robot_session_message.FieldCreatedAt,
		ent_robot_session_message.FieldID,
		ent_robot_session_message.Table,
		ent_robot_session_message.Table,
		ent_robot_session_message.FieldSequence,
		ent_robot_session_message.FieldID,
		ent_robot_session_message.Table,
		ent_robot_session_message.FieldID,
	))
	if err != nil {
		return err
	}

	return tx.Commit()
}

func robotSessionMessageSequenceMigrationRequired(ctx context.Context, db *sql.DB, driver string) (bool, error) {
	if driver == "pgx" {
		var required bool
		err := db.QueryRowContext(ctx, `
			SELECT EXISTS (
					SELECT 1 FROM information_schema.tables
					WHERE table_schema = current_schema() AND table_name = $1
				)
				AND NOT EXISTS (
					SELECT 1 FROM information_schema.columns
					WHERE table_schema = current_schema() AND table_name = $1 AND column_name = $2
				)
		`, ent_robot_session_message.Table, ent_robot_session_message.FieldSequence).Scan(&required)
		return required, err
	}

	rows, err := db.QueryContext(ctx, fmt.Sprintf("PRAGMA table_info(%s)", ent_robot_session_message.Table))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	tableExists := false
	columnExists := false
	for rows.Next() {
		var cid, notNull, primaryKey int
		var name, columnType string
		var defaultValue sql.NullString
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			return false, err
		}
		tableExists = true
		if name == ent_robot_session_message.FieldSequence {
			columnExists = true
		}
	}
	if err := rows.Err(); err != nil {
		return false, err
	}

	return tableExists && !columnExists, nil
}

func reconcileRobotSessionMessageSequences(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf(
		`
		UPDATE %s
		SET %s = (
			SELECT COALESCE(MAX(%s.%s), 0)
			FROM %s
			WHERE %s.%s = %s.%s
		)
		WHERE %s < (
			SELECT COALESCE(MAX(%s.%s), 0)
			FROM %s
			WHERE %s.%s = %s.%s
		)
	`,
		ent_robot_session.Table,
		ent_robot_session.FieldNextEventSequence,
		ent_robot_session_message.Table,
		ent_robot_session_message.FieldSequence,
		ent_robot_session_message.Table,
		ent_robot_session_message.Table,
		ent_robot_session_message.FieldSessionID,
		ent_robot_session.Table,
		ent_robot_session.FieldID,
		ent_robot_session.FieldNextEventSequence,
		ent_robot_session_message.Table,
		ent_robot_session_message.FieldSequence,
		ent_robot_session_message.Table,
		ent_robot_session_message.Table,
		ent_robot_session_message.FieldSessionID,
		ent_robot_session.Table,
		ent_robot_session.FieldID,
	))
	return err
}

// migrateSessionTokenHash invalidates sessions created before bearer secrets
// were stored as hashes. The original secrets cannot be reconstructed, so a
// placeholder hash would leave broken sessions behind. This only runs while
// adding token_hash; later startup migrations preserve new sessions.
func migrateSessionTokenHash() schema.ApplyHook {
	return func(next schema.Applier) schema.Applier {
		return schema.ApplyFunc(func(ctx context.Context, conn dialect.ExecQuerier, plan *migrate.Plan) error {
			if !addsSessionTokenHash(plan) {
				return next.Apply(ctx, conn, plan)
			}

			query := fmt.Sprintf("DELETE FROM %s", ent_session.Table)
			if err := conn.Exec(ctx, query, []any{}, nil); err != nil {
				return fault.Wrap(err, fmsg.With("failed to invalidate sessions without token hashes"))
			}

			return next.Apply(ctx, conn, plan)
		})
	}
}

func addsSessionTokenHash(plan *migrate.Plan) bool {
	for _, change := range plan.Changes {
		modifyTable, ok := change.Source.(*atlas_schema.ModifyTable)
		if !ok || modifyTable.T.Name != ent_session.Table {
			continue
		}

		if atlas_schema.Changes(modifyTable.Changes).IndexAddColumn(ent_session.FieldTokenHash) != -1 {
			return true
		}
	}

	return false
}

const oauthDeviceUserCodeHashIndex = "oauthdeviceauthorisation_user_code_hash"

// migrateOAuthDeviceUserCodeUniqueness removes every pending device flow before
// the non-unique index is replaced with a unique one. Existing rows were hashed
// with the old normalization rules, so retaining even non-duplicate hashes could
// make an old code containing I/L/O resolve to a different row after aliases are
// introduced. Device flows are short-lived and safe to restart, so fail closed.
func migrateOAuthDeviceUserCodeUniqueness() schema.ApplyHook {
	return func(next schema.Applier) schema.Applier {
		return schema.ApplyFunc(func(ctx context.Context, conn dialect.ExecQuerier, plan *migrate.Plan) error {
			if !addsUniqueOAuthDeviceUserCodeIndex(plan) {
				return next.Apply(ctx, conn, plan)
			}

			if err := deleteOAuthDeviceAuthorisations(ctx, conn); err != nil {
				return fault.Wrap(err, fmsg.With("failed to invalidate oauth device authorisations created with legacy user-code normalization"))
			}

			return next.Apply(ctx, conn, plan)
		})
	}
}

func addsUniqueOAuthDeviceUserCodeIndex(plan *migrate.Plan) bool {
	for _, change := range plan.Changes {
		modifyTable, ok := change.Source.(*atlas_schema.ModifyTable)
		if !ok || modifyTable.T.Name != ent_oauth_device_authorisation.Table {
			continue
		}

		for _, change := range modifyTable.Changes {
			switch indexChange := change.(type) {
			case *atlas_schema.AddIndex:
				if indexChange.I.Name == oauthDeviceUserCodeHashIndex && indexChange.I.Unique {
					return true
				}
			case *atlas_schema.ModifyIndex:
				if indexChange.To.Name == oauthDeviceUserCodeHashIndex && indexChange.To.Unique {
					return true
				}
			}
		}
	}

	return false
}

func deleteOAuthDeviceAuthorisations(ctx context.Context, conn dialect.ExecQuerier) error {
	query := fmt.Sprintf("DELETE FROM %s", ent_oauth_device_authorisation.Table)

	return conn.Exec(ctx, query, []any{}, nil)
}

func connect(cfg config.Config, driver *sql.DB) (*ent.Client, string, error) {
	d, _, err := getDriver(cfg.DatabaseURL)
	if err != nil {
		return nil, "", fault.Wrap(err)
	}

	opts := []ent.Option{}

	switch d {
	case "pgx":
		opts = append(opts, ent.Driver(entsql.OpenDB(dialect.Postgres, driver)))

	case "sqlite":
		opts = append(
			opts,
			ent.Driver(entsql.OpenDB(dialect.SQLite, driver)),
		)

	case "libsql":
		opts = append(
			opts,
			ent.Driver(entsql.OpenDB(dialect.SQLite, driver)),
		)

	default:
		panic(fmt.Sprintf("unsupported driver '%s' in ent connect", d))
	}

	return ent.NewClient(opts...), d, nil
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
		if err := ensureDatabasePathWritable(path, "sqlite", u); err != nil {
			return "", "", err
		}

		return "sqlite", path, nil

	case "libsql":
		// NOTE: We consider URLs of the form:
		// libsql://./path
		// or
		// libsql:///path
		// to be local disk databases and normalise to an absolute path.
		if (u.Host == "" || u.Host == ".") && u.Path != "" {
			path := u.Path
			if u.Host == "." {
				path = "." + u.Path
			}

			if !filepath.IsAbs(path) {
				absPath, err := filepath.Abs(path)
				if err != nil {
					return "", "", fault.Wrap(err, fmsg.With(fmt.Sprintf("failed to resolve libsql relative path: %s", u)))
				}
				path = absPath
			}

			if err := ensureDatabasePathWritable(path, "libsql", u); err != nil {
				return "", "", err
			}

			q := u.Query()
			pragmas := q["_pragma"]
			hasForeignKeysPragma := false
			for _, pragma := range pragmas {
				if pragma == "foreign_keys(1)" {
					hasForeignKeysPragma = true
					break
				}
			}
			if !hasForeignKeysPragma {
				q.Add("_pragma", "foreign_keys(1)")
			}

			return "libsql", (&url.URL{
				Scheme:   "file",
				Path:     path,
				RawQuery: q.Encode(),
			}).String(), nil
		}

		return "libsql", databaseURL, nil

	default:
		return "", "", fault.Newf("unsupported scheme: %s", u.Scheme)
	}
}

func ensureDatabasePathWritable(path string, driver string, u *url.URL) error {
	// NOTE: SQLite-backed drivers can return misleading "out of memory" errors
	// when the target directory does not exist or is not writable.
	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return fault.Wrap(err, fmsg.With(fmt.Sprintf("could not create directory for %s database: %s", driver, u)))
			}
		} else {
			return fault.Wrap(err, fmsg.With(fmt.Sprintf("could not read directory: %s", u)))
		}
	}

	testwrite := filepath.Join(filepath.Dir(path), ".perm_check")
	if err := os.WriteFile(testwrite, []byte("ok"), 0o644); err != nil {
		return fault.Wrap(err, fmsg.With(fmt.Sprintf("cannot write to directory for %s database: %s", driver, u)))
	}

	return nil
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
			if err := next.Apply(ctx, conn, plan); err != nil {
				return err
			}

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

			return nil
		})
	}
}

// migrateAccountVerifiedStatus is a data migration hook for backfilling the
// accounts.verified_status field from existing verified emails when the column
// is added or modified.
func migrateAccountVerifiedStatus() schema.ApplyHook {
	return func(next schema.Applier) schema.Applier {
		return schema.ApplyFunc(func(ctx context.Context, conn dialect.ExecQuerier, plan *migrate.Plan) error {
			hasChange := false
			for _, c := range plan.Changes {
				m, ok := c.Source.(*atlas_schema.ModifyTable)
				if !ok || m.T.Name != ent_account.Table {
					continue
				}

				changes := atlas_schema.Changes(m.Changes)
				if changes.IndexAddColumn(ent_account.FieldVerifiedStatus) != -1 ||
					changes.IndexModifyColumn(ent_account.FieldVerifiedStatus) != -1 {
					hasChange = true
					break
				}
			}

			if err := next.Apply(ctx, conn, plan); err != nil {
				return err
			}

			if !hasChange {
				return nil
			}

			err := conn.Exec(ctx, fmt.Sprintf(
				`
				UPDATE %s
				SET %s = 'email'
				WHERE EXISTS (
					SELECT 1
					FROM %s
					WHERE %s.%s = %s.%s
					  AND %s.%s = true
				)
			`,
				ent_account.Table,
				ent_account.FieldVerifiedStatus,
				ent_email.Table,
				ent_email.Table, ent_email.FieldAccountID,
				ent_account.Table, ent_account.FieldID,
				ent_email.Table, ent_email.FieldVerified,
			), []any{}, nil)
			if err != nil {
				return fault.Wrap(err, fmsg.With("failed to backfill account verified_status"))
			}

			return nil
		})
	}
}
