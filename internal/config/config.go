package config

import (
	"net/url"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"
)

// Config represents environment variable configuration parameters
type Config struct {
	Production  bool          `envconfig:"PRODUCTION"   default:"false"`
	LogLevel    zapcore.Level `envconfig:"LOG_LEVEL"    default:"info"`
	RunFrontend string        `envconfig:"RUN_FRONTEND" default:""`

	DevChaosSlowMode time.Duration `envconfig:"DEV_CHAOS_SLOW_MODE"` // Simulates slow requests
	DevChaosFailRate float64       `envconfig:"DEV_CHAOS_FAIL_RATE"` // Simulates failed requests

	DatabaseURL      string  `envconfig:"DATABASE_URL"           default:"sqlite://data/data.db?_pragma=foreign_keys(1)"`
	ListenAddr       string  `envconfig:"LISTEN_ADDR"            default:"0.0.0.0:8000"`
	SessionKey       string  `envconfig:"SESSION_KEY"            default:"0000000000000000"`
	PublicWebAddress url.URL `envconfig:"PUBLIC_WEB_ADDRESS"     default:"http://localhost:3000"`
	PublicAPIAddress url.URL `envconfig:"PUBLIC_API_ADDRESS"     default:"http://localhost:8000"`

	RateLimit       int           `envconfig:"RATE_LIMIT"        default:"1000"`
	RateLimitPeriod time.Duration `envconfig:"RATE_LIMIT_PERIOD" default:"1h"`
	RateLimitExpire time.Duration `envconfig:"RATE_LIMIT_EXPIRE" default:"1m"`

	EmailProvider string `envconfig:"EMAIL_PROVIDER"         default:""`

	AssetStorageType      string `envconfig:"ASSET_STORAGE_TYPE"`
	AssetStorageLocalPath string `envconfig:"ASSET_STORAGE_LOCAL_PATH"`
	S3Secure              bool   `envconfig:"S3_SECURE" default:"true"`
	S3Endpoint            string `envconfig:"S3_ENDPOINT"`
	S3Bucket              string `envconfig:"S3_BUCKET"`
	S3Region              string `envconfig:"S3_REGION"`
	S3AccessKey           string `envconfig:"S3_ACCESS_KEY"`
	S3SecretKey           string `envconfig:"S3_SECRET_KEY"`

	CacheProvider string `envconfig:"CACHE_PROVIDER" default:""`
	RedisHost     string `envconfig:"REDIS_HOST"     default:""`

	QueueType string `envconfig:"QUEUE_TYPE" default:"internal"`
	AmqpURL   string `envconfig:"AMQP_URL"   default:"amqp://guest:guest@localhost:5672/"`

	SemdexEnabled     bool   `envconfig:"SEMDEX_ENABLED" default:"false"`
	WeaviateURL       string `envconfig:"WEAVIATE_URL"`
	WeaviateToken     string `envconfig:"WEAVIATE_API_TOKEN"`
	WeaviateClassName string `envconfig:"WEAVIATE_CLASS_NAME"`
	OpenAIKey         string `envconfig:"OPENAI_API_KEY"`
}

func Build() fx.Option {
	return fx.Provide(func() (c Config, err error) {
		if err = envconfig.Process("", &c); err != nil {
			return c, fault.Wrap(err, fmsg.With("failed to parse configuration from environment variables"))
		}

		return
	})
}
