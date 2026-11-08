package config

import (
	"log/slog"
	"math"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestLoadDefaults(t *testing.T) {
	t.Parallel()

	c, err := load(envLookup(nil))
	require.NoError(t, err)
	require.Equal(t, Config{
		LogLevel:                    slog.LevelInfo,
		DatabaseURL:                 "sqlite://data/data.db?_pragma=foreign_keys(1)",
		ListenAddr:                  "0.0.0.0:8000",
		PublicWebAddress:            url.URL{Scheme: "http", Host: "localhost:3000"},
		PublicAPIAddress:            url.URL{Scheme: "http", Host: "localhost:8000"},
		RateLimit:                   5000,
		RateLimitPeriod:             time.Hour,
		RateLimitBucket:             time.Minute,
		RateLimitGuestCost:          1,
		WebPushMaxSubscriptions:     5,
		SMTPUseTLS:                  true,
		KeycloakDisplayName:         "Keycloak",
		OAuthAccessTokenTTL:         15 * time.Minute,
		OAuthRefreshTokenTTL:        720 * time.Hour,
		OAuthDeviceCodeTTL:          10 * time.Minute,
		OAuthDevicePollEvery:        5 * time.Second,
		S3Secure:                    true,
		ImageWorkers:                1,
		ImageMemoryBudget:           268435456,
		SearchProvider:              "database",
		SearchIndexChunkSize:        1000,
		BlevePath:                   "data/bleve",
		RedisSearchIndexName:        "storyden",
		QueueType:                   "internal",
		AmqpURL:                     "amqp://guest:guest@localhost:5672/",
		QueueMaxRetries:             5,
		QueueRetryInitialInterval:   time.Second,
		QueueRetryMaxInterval:       time.Minute,
		PluginDataPath:              "./data/plugins",
		RobotWorkspaceDataPath:      "./data/robot-workspaces",
		PluginRuntimeProvider:       "local",
		PluginMaxRestartAttempts:    3,
		PluginMaxBackoffDuration:    time.Minute,
		PluginRuntimeCrashThreshold: 30 * time.Second,
		PluginRuntimeCrashBackoff:   5 * time.Second,
		SemdexLocalPath:             "data/semdex",
	}, c)
}

func TestLoadTypedEnvironment(t *testing.T) {
	t.Parallel()

	c, err := load(envLookup(map[string]string{
		"LOG_LEVEL":            "DEBUG+2",
		"LOG_FORMAT":           "json",
		"RATE_LIMIT":           "0x100",
		"IMAGE_MEMORY_BUDGET":  "4294967296",
		"PINECONE_DIMENSIONS":  "1536",
		"DEV_CHAOS_FAIL_RATE":  "0.25",
		"RATE_LIMIT_PERIOD":    "2h30m",
		"SMTP_USE_TLS":         "FALSE",
		"JWT_SECRET":           "raw UTF-8 secret 🔑",
		"CORS_ALLOWED_ORIGINS": "https://one.example, https://two.example",
		"PUBLIC_WEB_ADDRESS":   "https://community.example/forum?q=1",
		"REDIS_URL":            "redis://user:password@cache.example:6379/1",
	}))
	require.NoError(t, err)
	require.Equal(t, slog.LevelDebug+2, c.LogLevel)
	require.Equal(t, "json", c.LogFormat)
	require.Equal(t, 256, c.RateLimit)
	require.Equal(t, int64(4294967296), c.ImageMemoryBudget)
	require.Equal(t, int32(1536), c.PineconeDimensions)
	require.Equal(t, 0.25, c.DevChaosFailRate)
	require.Equal(t, 150*time.Minute, c.RateLimitPeriod)
	require.False(t, c.SMTPUseTLS)
	require.Equal(t, []byte("raw UTF-8 secret 🔑"), c.JWTSecret)
	require.Equal(t, []string{"https://one.example", " https://two.example"}, c.CORSAllowedOrigins)
	require.Equal(t, "https://community.example/forum?q=1", c.PublicWebAddress.String())
	require.Equal(t, "user", c.RedisURL.User.Username())
	require.Equal(t, "cache.example:6379", c.RedisURL.Host)
}

func TestLoadExplicitEmptyAndZero(t *testing.T) {
	t.Parallel()

	c, err := load(envLookup(map[string]string{
		"DATABASE_URL":         "",
		"PUBLIC_WEB_ADDRESS":   "",
		"JWT_SECRET":           "",
		"CORS_ALLOWED_ORIGINS": " \t ",
		"RATE_LIMIT":           "0",
		"IMAGE_WORKERS":        "0",
		"IMAGE_MEMORY_BUDGET":  "0",
	}))
	require.NoError(t, err)
	require.Empty(t, c.DatabaseURL)
	require.Equal(t, url.URL{}, c.PublicWebAddress)
	require.Equal(t, []byte{}, c.JWTSecret)
	require.Equal(t, []string{}, c.CORSAllowedOrigins)
	require.Zero(t, c.RateLimit)
	require.Equal(t, 1, c.ImageWorkers)
	require.Equal(t, int64(268435456), c.ImageMemoryBudget)
}

func TestLoadRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	cases := []struct {
		key   string
		value string
	}{
		{"LOG_LEVEL", "verbose"},
		{"LOG_LEVEL", ""},
		{"RATE_LIMIT", "not-a-number"},
		{"RATE_LIMIT", ""},
		{"PINECONE_DIMENSIONS", "2147483648"},
		{"IMAGE_MEMORY_BUDGET", "9223372036854775808"},
		{"IMAGE_MEMORY_BUDGET", "-1"},
		{"IMAGE_WORKERS", "-1"},
		{"SMTP_USE_TLS", "yes"},
		{"SMTP_USE_TLS", ""},
		{"RATE_LIMIT_PERIOD", "tomorrow"},
		{"RATE_LIMIT_PERIOD", ""},
		{"DEV_CHAOS_FAIL_RATE", "NaN"},
		{"DEV_CHAOS_FAIL_RATE", "+Inf"},
		{"DEV_CHAOS_FAIL_RATE", "-0.01"},
		{"DEV_CHAOS_FAIL_RATE", "1.01"},
		{"REDIS_URL", "redis://user:private-password@host/%zz"},
	}

	for _, tc := range cases {
		t.Run(tc.key+"/"+tc.value, func(t *testing.T) {
			t.Parallel()

			_, err := load(envLookup(map[string]string{tc.key: tc.value}))
			require.ErrorIs(t, err, errInvalidConfiguration)
			require.ErrorContains(t, err, tc.key)
			require.NotContains(t, err.Error(), "private-password")
		})
	}
}

func TestApplyDefaultsAndValidate(t *testing.T) {
	t.Parallel()

	c, err := ApplyDefaults(Config{})
	require.NoError(t, err)
	require.Equal(t, 1, c.ImageWorkers)
	require.Equal(t, int64(268435456), c.ImageMemoryBudget)
	require.Empty(t, c.DatabaseURL)

	_, err = ApplyDefaults(Config{ImageWorkers: -1})
	require.ErrorContains(t, err, "IMAGE_WORKERS")

	_, err = ApplyDefaults(Config{ImageMemoryBudget: -1})
	require.ErrorContains(t, err, "IMAGE_MEMORY_BUDGET")

	require.ErrorContains(t, (Config{DevChaosFailRate: math.NaN()}).Validate(), "DEV_CHAOS_FAIL_RATE")
	require.NoError(t, (Config{DevChaosFailRate: 1}).Validate())
}

func TestBuildLoadsEnvironment(t *testing.T) {
	t.Setenv("RATE_LIMIT", "1234")
	t.Setenv("IMAGE_WORKERS", "2")

	var c Config
	app := fx.New(fx.NopLogger, Build(), fx.Populate(&c))
	require.NoError(t, app.Err())
	require.Equal(t, 1234, c.RateLimit)
	require.Equal(t, 2, c.ImageWorkers)
}

func envLookup(values map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		value, present := values[key]
		return value, present
	}
}
