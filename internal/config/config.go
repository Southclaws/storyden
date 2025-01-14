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
	Production    bool          `envconfig:"PRODUCTION"   default:"false"`
	LogLevel      zapcore.Level `envconfig:"LOG_LEVEL"    default:"info"`
	RunFrontend   string        `envconfig:"RUN_FRONTEND" default:""`           // Path to server.js for running frontend process
	FrontendProxy url.URL       `envconfig:"PROXY_FRONTEND_ADDRESS" default:""` // Proxy non-/api requests to this address

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

	// logger (default) or sentry
	OTELProvider string  `envconfig:"OTEL_PROVIDER"               default:""`
	OTELEndpoint url.URL `envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT" default:""`
	SentryDSN    string  `envconfig:"SENTRY_DSN"                  default:""` // required when OTLPProvider = sentry

	EmailProvider string `envconfig:"EMAIL_PROVIDER"         default:""`

	AssetStorageType      string `envconfig:"ASSET_STORAGE_TYPE"`
	AssetStorageLocalPath string `envconfig:"ASSET_STORAGE_LOCAL_PATH"`
	S3Secure              bool   `envconfig:"S3_SECURE" default:"true"`
	S3Endpoint            string `envconfig:"S3_ENDPOINT"`
	S3Bucket              string `envconfig:"S3_BUCKET"`
	S3Region              string `envconfig:"S3_REGION"`
	S3AccessKey           string `envconfig:"S3_ACCESS_KEY"`
	S3SecretKey           string `envconfig:"S3_SECRET_KEY"`

	CacheProvider string  `envconfig:"CACHE_PROVIDER" default:""`
	RedisURL      url.URL `envconfig:"REDIS_URL"      default:""`

	QueueType string `envconfig:"QUEUE_TYPE" default:"internal"`
	AmqpURL   string `envconfig:"AMQP_URL"   default:"amqp://guest:guest@localhost:5672/"`

	LanguageModelProvider string `envconfig:"LANGUAGE_MODEL_PROVIDER"`
	OpenAIKey             string `envconfig:"OPENAI_API_KEY"`

	// default (whatever LanguageModelProvider is) or perplexity
	AskerProvider    string `envconfig:"ASKER_PROVIDER" default:""`
	PerplexityAPIKey string `envconfig:"PERPLEXITY_API_KEY"`

	// chromem (local), weaviate, pinecone
	SemdexProvider string `envconfig:"SEMDEX_PROVIDER" default:""`

	// Chromem
	SemdexLocalPath string `envconfig:"SEMDEX_LOCAL_PATH" default:"data/semdex"`

	// Weaviate
	WeaviateURL       string `envconfig:"WEAVIATE_URL"`
	WeaviateToken     string `envconfig:"WEAVIATE_API_TOKEN"`
	WeaviateClassName string `envconfig:"WEAVIATE_CLASS_NAME"`

	// Pinecone
	PineconeAPIKey     string `envconfig:"PINECONE_API_KEY"`
	PineconeIndex      string `envconfig:"PINECONE_INDEX"`
	PineconeDimensions int32  `envconfig:"PINECONE_DIMENSIONS"`
	PineconeCloud      string `envconfig:"PINECONE_CLOUD"`
	PineconeRegion     string `envconfig:"PINECONE_REGION"`
}

func Build() fx.Option {
	return fx.Provide(func() (c Config, err error) {
		if err = envconfig.Process("", &c); err != nil {
			return c, fault.Wrap(err, fmsg.With("failed to parse configuration from environment variables"))
		}

		return
	})
}
