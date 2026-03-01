package integration

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"dario.cat/mergo"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources"
	"github.com/Southclaws/storyden/app/services"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure"
	"github.com/Southclaws/storyden/internal/utils"
)

// Test provides a full app setup for testing end to end behaviour. Example:
//
//	func TestMyThing(t *testing.T) {
//	    integration.Test(t, nil, fx.Invoke(func(test dependencies...) {
//	        r := require.New(t)
//	        a := assert.New(t)
//
//	        // your e2e test code...
//
//	    }))
//	}
func Test(t *testing.T, cfg *config.Config, o ...fx.Option) {
	defaultConfig := config.Config{
		PublicAPIAddress: *utils.Must(url.Parse("http://localhost")),
		PublicWebAddress: *utils.Must(url.Parse("http://localhost")),
		JWTSecret:        []byte("00000000000000000000000000000000"),
		RateLimit:        5000,
		RateLimitPeriod:  time.Hour,
		RateLimitBucket:  time.Minute,
		EmailProvider:    "mock",
	}

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		if isMaybeProdDB(dbURL) {
			panic("maybe accidental prod DATABASE_URL in integration tests!")
		}
		defaultConfig.DatabaseURL = makePerTestDatabaseURL(dbURL, t.Name())
		fmt.Println("Using database URL from environment: ", defaultConfig.DatabaseURL)
	} else {
		defaultConfig.DatabaseURL = makePerTestDatabaseURL("sqlite://data/data.db", t.Name())
		fmt.Println("Using database URL default: ", defaultConfig.DatabaseURL)
	}

	ctx, cf := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cf()
	})

	o = append(o,
		// main application dependencies
		application(),

		// provide a global context
		fx.Provide(func() context.Context { return ctx }),
	)

	// if this test has a custom config, merge+overwrite with the defaults.
	if cfg != nil {
		mergo.MergeWithOverwrite(&defaultConfig, cfg)
	}

	o = append(o, fx.Provide(func() config.Config { return defaultConfig }))

	app := fx.New(o...)
	err := app.Start(ctx)
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}
}

// application gives you some basics needed by most components.
func application() fx.Option {
	return fx.Options(
		fx.NopLogger,

		infrastructure.Build(),
		resources.Build(),
		services.Build(),
	)
}

func isMaybeProdDB(url string) bool {
	dangerous := []string{
		"free-tier",
		".aws-eu-central",
		"cockroachlabs",
		"cloud",
		"verify-full",
		".turso.io",
	}

	for _, v := range dangerous {
		if strings.Contains(url, v) {
			return true
		}
	}

	return false
}

func makePerTestDatabaseURL(databaseURL string, testName string) string {
	u, err := url.Parse(databaseURL)
	if err != nil {
		panic(err)
	}

	timePrefix := time.Now().Format(time.RFC3339)

	testSuffix := sanitizeDBTestName(fmt.Sprintf("%s-%s", timePrefix, testName))

	switch u.Scheme {
	case "libsql":
		// Keep remote Turso URLs shared. Only local file-backed libsql should
		// get isolated per-test files.
		if isLibsqlRemote(u) {
			return databaseURL
		}

		rewritten := *u
		rewritten.Path = appendSuffixToFilePath(u.Path, testSuffix)

		q := rewritten.Query()
		desired := []string{
			"foreign_keys(1)",
			"busy_timeout(10000)",
			"journal_mode(WAL)",
			"synchronous(NORMAL)",
			"cache_size(1000000000)",
			"temp_store(MEMORY)",
		}
		present := make(map[string]struct{}, len(q["_pragma"]))
		for _, v := range q["_pragma"] {
			present[v] = struct{}{}
		}
		for _, pragma := range desired {
			if _, ok := present[pragma]; ok {
				continue
			}
			q.Add("_pragma", pragma)
		}

		rewritten.RawQuery = q.Encode()

		return rewritten.String()

	case "sqlite", "sqlite3":
		if u.Path == "" {
			return databaseURL
		}

		rewritten := *u
		rewritten.Path = appendSuffixToFilePath(u.Path, testSuffix)

		desired := []string{
			"foreign_keys(1)",
			"busy_timeout(10000)",
			"journal_mode(WAL)",
			"synchronous(NORMAL)",
			"cache_size(1000000000)",
			"temp_store(MEMORY)",
		}

		q := rewritten.Query()
		present := make(map[string]struct{}, len(q["_pragma"]))
		for _, v := range q["_pragma"] {
			present[v] = struct{}{}
		}

		for _, pragma := range desired {
			if _, ok := present[pragma]; ok {
				continue
			}
			q.Add("_pragma", pragma)
		}

		rewritten.RawQuery = q.Encode()

		return rewritten.String()

	default:
		return databaseURL
	}
}

func isLibsqlRemote(u *url.URL) bool {
	switch u.Host {
	case "", ".":
		return false
	default:
		return true
	}
}

func appendSuffixToFilePath(filePath string, suffix string) string {
	ext := path.Ext(filePath)
	base := strings.TrimSuffix(filePath, ext)
	if ext == "" {
		ext = ".db"
	}
	return base + "-" + suffix + ext
}

func sanitizeDBTestName(name string) string {
	// Keep names filesystem-safe and deterministic across OSes.
	mapped := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= 'A' && r <= 'Z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == '-', r == '_':
			return r
		default:
			return '_'
		}
	}, name)

	mapped = strings.Trim(mapped, "_")
	if mapped == "" {
		return "test"
	}

	// Avoid excessively long filenames in deep package test paths.
	return filepath.Base(mapped)
}
