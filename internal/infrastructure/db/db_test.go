package db

import (
	"database/sql"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
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
