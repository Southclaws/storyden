package db

import (
	"context"
	"database/sql"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"ariga.io/atlas/sql/migrate"
	atlas_schema "ariga.io/atlas/sql/schema"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	entschema "entgo.io/ent/dialect/sql/schema"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	entmigrate "github.com/Southclaws/storyden/internal/ent/migrate"
	ent_oauth_device_authorisation "github.com/Southclaws/storyden/internal/ent/oauthdeviceauthorisation"
	ent_robot_session "github.com/Southclaws/storyden/internal/ent/robotsession"
	ent_robot_session_message "github.com/Southclaws/storyden/internal/ent/robotsessionmessage"
	ent_session "github.com/Southclaws/storyden/internal/ent/session"
)

func TestGetDriverLibsqlRemote(t *testing.T) {
	t.Parallel()

	dsn := "libsql://example-org.turso.io?authToken=abc123"

	driver, path, err := getDriver(dsn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if driver != "libsql" {
		t.Fatalf("expected driver libsql, got %q", driver)
	}

	if path != dsn {
		t.Fatalf("expected path %q, got %q", dsn, path)
	}
}

func TestGetDriverLibsqlLocalFilePath(t *testing.T) {
	t.Parallel()

	dbFile := filepath.Join(t.TempDir(), "storyden-libsql.db")
	dsn := "libsql://" + dbFile + "?_pragma=foreign_keys(1)"

	driver, path, err := getDriver(dsn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if driver != "libsql" {
		t.Fatalf("expected driver libsql, got %q", driver)
	}

	if !strings.HasPrefix(path, "file://") {
		t.Fatalf("expected local libsql path to be rewritten as file:// URL, got %q", path)
	}

	u, err := url.Parse(path)
	if err != nil {
		t.Fatalf("failed to parse rewritten URL: %v", err)
	}

	if !strings.Contains(strings.Join(u.Query()["_pragma"], ","), "foreign_keys(1)") {
		t.Fatalf("expected rewritten path to include _pragma=foreign_keys(1), got %q", path)
	}
}

func TestGetDriverLibsqlLocalAddsForeignKeysPragma(t *testing.T) {
	t.Parallel()

	dbFile := filepath.Join(t.TempDir(), "storyden-libsql-no-pragma.db")
	dsn := "libsql://" + dbFile

	driver, path, err := getDriver(dsn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if driver != "libsql" {
		t.Fatalf("expected driver libsql, got %q", driver)
	}

	u, err := url.Parse(path)
	if err != nil {
		t.Fatalf("failed to parse rewritten URL: %v", err)
	}

	if !strings.Contains(strings.Join(u.Query()["_pragma"], ","), "foreign_keys(1)") {
		t.Fatalf("expected local libsql URL to include _pragma=foreign_keys(1), got %q", path)
	}
}

func TestGetDriverLibsqlLocalFilePathOpen(t *testing.T) {
	t.Parallel()

	dbFile := filepath.Join(t.TempDir(), "storyden-libsql-open.db")
	dsn := "libsql://" + dbFile + "?_pragma=foreign_keys(1)"

	driver, path, err := getDriver(dsn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	db, err := sql.Open(driver, path)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if _, err := db.Exec(`create table if not exists ping (id integer primary key, value text)`); err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	if _, err := db.Exec(`insert into ping (value) values (?)`, "ok"); err != nil {
		t.Fatalf("failed to insert row: %v", err)
	}

	var got string
	if err := db.QueryRow(`select value from ping limit 1`).Scan(&got); err != nil {
		t.Fatalf("failed to query row: %v", err)
	}

	if got != "ok" {
		t.Fatalf("expected value %q, got %q", "ok", got)
	}
}

func TestGetDriverLibsqlRemotePreservesQueryParams(t *testing.T) {
	t.Parallel()

	dsn := "libsql://example-org.turso.io?authToken=abc123&_pragma=foreign_keys(1)"

	driver, path, err := getDriver(dsn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if driver != "libsql" {
		t.Fatalf("expected driver libsql, got %q", driver)
	}

	if path != dsn {
		t.Fatalf("expected remote libsql URL to be preserved as-is, got %q", path)
	}
}

func TestAddsUniqueOAuthDeviceUserCodeIndex(t *testing.T) {
	tests := map[string]struct {
		change atlas_schema.Change
		want   bool
	}{
		"added unique index": {
			change: &atlas_schema.AddIndex{I: &atlas_schema.Index{Name: oauthDeviceUserCodeHashIndex, Unique: true}},
			want:   true,
		},
		"modified unique index": {
			change: &atlas_schema.ModifyIndex{To: &atlas_schema.Index{Name: oauthDeviceUserCodeHashIndex, Unique: true}},
			want:   true,
		},
		"non-unique index": {
			change: &atlas_schema.AddIndex{I: &atlas_schema.Index{Name: oauthDeviceUserCodeHashIndex}},
			want:   false,
		},
		"different index": {
			change: &atlas_schema.AddIndex{I: &atlas_schema.Index{Name: "other", Unique: true}},
			want:   false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			plan := &migrate.Plan{Changes: []*migrate.Change{{
				Source: &atlas_schema.ModifyTable{
					T:       &atlas_schema.Table{Name: ent_oauth_device_authorisation.Table},
					Changes: []atlas_schema.Change{test.change},
				},
			}}}

			assert.Equal(t, test.want, addsUniqueOAuthDeviceUserCodeIndex(plan))
		})
	}
}

func TestMigrateOAuthDeviceUserCodeUniqueness(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open("sqlite", "file:"+t.Name()+"?mode=memory&cache=shared&_pragma=foreign_keys(1)")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	driver := entsql.OpenDB(dialect.SQLite, db)

	legacyTables, err := entschema.CopyTables(entmigrate.Tables)
	require.NoError(t, err)

	foundIndex := false
	for _, table := range legacyTables {
		if table.Name != ent_oauth_device_authorisation.Table {
			continue
		}
		for _, index := range table.Indexes {
			if index.Name == oauthDeviceUserCodeHashIndex {
				index.Unique = false
				foundIndex = true
			}
		}
	}
	require.True(t, foundIndex)

	legacySchema := entmigrate.NewSchema(driver)
	require.NoError(t, entmigrate.Create(ctx, legacySchema, legacyTables))

	clientID := xid.New().String()
	_, err = db.ExecContext(ctx, `
		INSERT INTO oauth_clients
			(id, updated_at, client_id, name, redirect_uris, allowed_scopes, allowed_grants)
		VALUES (?, ?, ?, ?, '[]', '[]', '[]')
	`, clientID, time.Now(), "migration-test-client", "Migration test client")
	require.NoError(t, err)

	insertDevice := func(deviceHash, userHash, display string) error {
		_, err := db.ExecContext(ctx, `
			INSERT INTO oauth_device_authorisations
				(id, device_code_hash, user_code_hash, user_code_display, scope, expires_at, client_id)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, xid.New().String(), deviceHash, userHash, display, "openid", time.Now().Add(10*time.Minute), clientID)
		return err
	}

	// These distinct legacy hashes can become ambiguous once Crockford aliases
	// map O/I/L to 0/1. The migration must invalidate all in-flight rows, not
	// only exact duplicate hashes.
	require.NoError(t, insertDevice("legacy-device-o", "legacy-hash-o", "O234-5678"))
	require.NoError(t, insertDevice("legacy-device-zero", "legacy-hash-zero", "0234-5678"))

	require.NoError(t, entmigrate.Create(
		ctx,
		entmigrate.NewSchema(driver),
		entmigrate.Tables,
		entschema.WithDropIndex(true),
		entschema.WithApplyHook(migrateOAuthDeviceUserCodeUniqueness()),
	))

	var count int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM oauth_device_authorisations`).Scan(&count))
	assert.Zero(t, count, "legacy in-flight authorisations must be invalidated")

	require.NoError(t, insertDevice("new-device", "unique-user-hash", "ABCD-EFGH"))

	// Once the index is unique, ordinary startup migrations must preserve new
	// in-flight flows because there is no longer a normalization transition.
	require.NoError(t, entmigrate.Create(
		ctx,
		entmigrate.NewSchema(driver),
		entmigrate.Tables,
		entschema.WithDropIndex(true),
		entschema.WithApplyHook(migrateOAuthDeviceUserCodeUniqueness()),
	))
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM oauth_device_authorisations`).Scan(&count))
	assert.Equal(t, 1, count)

	err = insertDevice("duplicate-device", "unique-user-hash", "ABCD-EFGH")
	require.Error(t, err, "the migrated index must reject duplicate user-code hashes")
}

func TestMigrateSessionTokenHash(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open("sqlite", "file:"+t.Name()+"?mode=memory&cache=shared&_pragma=foreign_keys(1)")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	driver := entsql.OpenDB(dialect.SQLite, db)

	legacyTables, err := entschema.CopyTables(entmigrate.Tables)
	require.NoError(t, err)

	foundTokenHash := false
	for _, table := range legacyTables {
		if table.Name != ent_session.Table {
			continue
		}

		columns := table.Columns[:0]
		for _, column := range table.Columns {
			if column.Name == ent_session.FieldTokenHash {
				foundTokenHash = true
				continue
			}
			columns = append(columns, column)
		}
		table.Columns = columns
	}
	require.True(t, foundTokenHash)

	require.NoError(t, entmigrate.Create(ctx, entmigrate.NewSchema(driver), legacyTables))

	accountID := xid.New().String()
	_, err = db.ExecContext(ctx, `
		INSERT INTO accounts (id, updated_at, handle, name)
		VALUES (?, ?, ?, ?)
	`, accountID, time.Now(), "migration-test-account", "Migration test account")
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO sessions (id, expires_at, account_id)
		VALUES (?, ?, ?)
	`, xid.New().String(), time.Now().Add(time.Hour), accountID)
	require.NoError(t, err)

	require.NoError(t, entmigrate.Create(
		ctx,
		entmigrate.NewSchema(driver),
		entmigrate.Tables,
		entschema.WithDropColumn(true),
		entschema.WithApplyHook(migrateSessionTokenHash()),
	))

	var count int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions`).Scan(&count))
	assert.Zero(t, count, "sessions without recoverable bearer secrets must be invalidated")

	_, err = db.ExecContext(ctx, `
		INSERT INTO sessions (id, token_hash, expires_at, account_id)
		VALUES (?, ?, ?, ?)
	`, xid.New().String(), "new-session-token-hash", time.Now().Add(time.Hour), accountID)
	require.NoError(t, err)

	// Once token_hash exists, routine migrations must preserve sessions.
	require.NoError(t, entmigrate.Create(
		ctx,
		entmigrate.NewSchema(driver),
		entmigrate.Tables,
		entschema.WithDropColumn(true),
		entschema.WithApplyHook(migrateSessionTokenHash()),
	))
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions`).Scan(&count))
	assert.Equal(t, 1, count)
}

func TestMigrateRobotSessionMessageSequences(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open("sqlite", "file:"+t.Name()+"?mode=memory&cache=shared&_pragma=foreign_keys(1)")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	driver := entsql.OpenDB(dialect.SQLite, db)
	legacyTables, err := entschema.CopyTables(entmigrate.Tables)
	require.NoError(t, err)

	foundSequence := false
	foundAllocator := false
	for _, table := range legacyTables {
		switch table.Name {
		case ent_robot_session_message.Table:
			columns := table.Columns[:0]
			for _, column := range table.Columns {
				if column.Name == ent_robot_session_message.FieldSequence {
					foundSequence = true
					continue
				}
				columns = append(columns, column)
			}
			table.Columns = columns
			table.Indexes = nil

		case ent_robot_session.Table:
			columns := table.Columns[:0]
			for _, column := range table.Columns {
				if column.Name == ent_robot_session.FieldNextEventSequence {
					foundAllocator = true
					continue
				}
				columns = append(columns, column)
			}
			table.Columns = columns
		}
	}
	require.True(t, foundSequence)
	require.True(t, foundAllocator)
	require.NoError(t, entmigrate.Create(ctx, entmigrate.NewSchema(driver), legacyTables))

	accountID := xid.New().String()
	_, err = db.ExecContext(ctx, `
		INSERT INTO accounts (id, updated_at, handle, name)
		VALUES (?, ?, ?, ?)
	`, accountID, time.Now(), "robot-migration-account", "Robot migration account")
	require.NoError(t, err)

	sessionA := xid.New().String()
	sessionB := xid.New().String()
	for _, sessionID := range []string{sessionA, sessionB} {
		_, err = db.ExecContext(ctx, `
			INSERT INTO robot_sessions (id, updated_at, name, account_id, state)
			VALUES (?, ?, ?, ?, '{}')
		`, sessionID, time.Now(), "Migration session", accountID)
		require.NoError(t, err)
	}

	firstID := xid.New().String()
	secondID := xid.New().String()
	thirdID := xid.New().String()
	createdAt := time.Now().Add(-time.Hour)
	insertMessage := func(id, sessionID string, created time.Time) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO robot_session_messages (id, created_at, session_id)
			VALUES (?, ?, ?)
		`, id, created, sessionID)
		require.NoError(t, err)
	}
	insertMessage(secondID, sessionA, createdAt.Add(time.Minute))
	insertMessage(firstID, sessionA, createdAt)
	insertMessage(thirdID, sessionB, createdAt)

	require.NoError(t, entmigrate.Create(
		ctx,
		entmigrate.NewSchema(driver),
		entmigrate.Tables,
		entschema.WithDropIndex(true),
		entschema.WithDropColumn(true),
		entschema.WithHooks(migrateRobotSessionMessageSequences(db, "sqlite")),
	))

	var firstSequence, secondSequence, thirdSequence uint64
	require.NoError(t, db.QueryRowContext(ctx, `SELECT sequence FROM robot_session_messages WHERE id = ?`, firstID).Scan(&firstSequence))
	require.NoError(t, db.QueryRowContext(ctx, `SELECT sequence FROM robot_session_messages WHERE id = ?`, secondID).Scan(&secondSequence))
	require.NoError(t, db.QueryRowContext(ctx, `SELECT sequence FROM robot_session_messages WHERE id = ?`, thirdID).Scan(&thirdSequence))
	assert.Equal(t, uint64(1), firstSequence)
	assert.Equal(t, uint64(2), secondSequence)
	assert.Equal(t, uint64(1), thirdSequence)

	var sessionASequence, sessionBSequence uint64
	require.NoError(t, db.QueryRowContext(ctx, `SELECT next_event_sequence FROM robot_sessions WHERE id = ?`, sessionA).Scan(&sessionASequence))
	require.NoError(t, db.QueryRowContext(ctx, `SELECT next_event_sequence FROM robot_sessions WHERE id = ?`, sessionB).Scan(&sessionBSequence))
	assert.Equal(t, uint64(2), sessionASequence)
	assert.Equal(t, uint64(1), sessionBSequence)

	_, err = db.ExecContext(ctx, `UPDATE robot_sessions SET next_event_sequence = 10 WHERE id = ?`, sessionA)
	require.NoError(t, err)
	require.NoError(t, entmigrate.Create(
		ctx,
		entmigrate.NewSchema(driver),
		entmigrate.Tables,
		entschema.WithDropIndex(true),
		entschema.WithDropColumn(true),
		entschema.WithHooks(migrateRobotSessionMessageSequences(db, "sqlite")),
	))
	require.NoError(t, db.QueryRowContext(ctx, `SELECT next_event_sequence FROM robot_sessions WHERE id = ?`, sessionA).Scan(&sessionASequence))
	assert.Equal(t, uint64(10), sessionASequence)
}
