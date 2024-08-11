package config

import (
	"net/url"

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

	DatabaseURL      string  `envconfig:"DATABASE_URL"           default:"sqlite://data/data.db?_pragma=foreign_keys(1)"`
	ListenAddr       string  `envconfig:"LISTEN_ADDR"            default:"0.0.0.0:8000"`
	SessionKey       string  `envconfig:"SESSION_KEY"            default:"0000000000000000"`
	PublicWebAddress url.URL `envconfig:"PUBLIC_WEB_ADDRESS"     default:"http://localhost:3000"`
	PublicAPIAddress url.URL `envconfig:"PUBLIC_API_ADDRESS"     default:"http://localhost:8000"`

	EmailProvider string `envconfig:"EMAIL_PROVIDER"         default:""`

	AssetStorageType      string `envconfig:"ASSET_STORAGE_TYPE"`
	AssetStorageLocalPath string `envconfig:"ASSET_STORAGE_LOCAL_PATH"`
	S3Endpoint            string `envconfig:"S3_ENDPOINT"`
	S3Bucket              string `envconfig:"S3_BUCKET"`
	S3Region              string `envconfig:"S3_REGION"`
	S3AccessKey           string `envconfig:"S3_ACCESS_KEY"`
	S3SecretKey           string `envconfig:"S3_SECRET_KEY"`

	QueueType string `envconfig:"QUEUE_TYPE" default:"internal"`
	AmqpURL   string `envconfig:"AMQP_URL"   default:"amqp://guest:guest@localhost:5672/"`
}

func Build() fx.Option {
	return fx.Provide(func() (c Config, err error) {
		if err = envconfig.Process("", &c); err != nil {
			return c, fault.Wrap(err, fmsg.With("failed to parse configuration from environment variables"))
		}

		return
	})
}
