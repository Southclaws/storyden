package tests

import (
	"os"
	"strings"
)

// IsSharedPostgresDatabase reports whether the current integration run uses a
// shared PostgreSQL database URL, which means global settings writes can leak
// across concurrently running test packages.
//
// We do this on github actions as one of the 3 databases. The test suite is
// actually intentionally meant to be run on a shared database to ensure that
// many API calls in a non-deterministic order doesn't cause any issues. The
// only minor issue is that a small handful of features will affect global state
// and causes that shared db test to fail in odd ways. To get around this, we
// just skip a couple of tests that touch global state as part of their process
// such as global authentication mode (handle vs email) all other tests will
// work perfectly fine in any order and mostly with t.Parallel() which is great!
func IsSharedPostgresDatabase() bool {
	databaseURL := strings.ToLower(os.Getenv("DATABASE_URL"))

	return strings.HasPrefix(databaseURL, "postgresql://")
}
