package config

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"
)

// Config represents environment variable configuration parameters
type Config struct {
	Production bool          `envconfig:"PRODUCTION" default:"false"`
	LogLevel   zapcore.Level `envconfig:"LOG_LEVEL"  default:"info"`

	DatabaseURL      string `envconfig:"DATABASE_URL"           default:"sqlite://data.db?_pragma=foreign_keys(1)"`
	ListenAddr       string `envconfig:"LISTEN_ADDR"            default:"0.0.0.0:8000"`
	CookieDomain     string `envconfig:"COOKIE_DOMAIN"          default:"localhost"`
	SessionKey       string `envconfig:"SESSION_KEY"            default:"0000000000000000"`
	PublicWebAddress string `envconfig:"PUBLIC_WEB_ADDRESS"     default:"http://localhost:3000"`

	AssetStorageType      string `envconfig:"ASSET_STORAGE_TYPE"`
	AssetStorageLocalPath string `envconfig:"ASSET_STORAGE_LOCAL_PATH"`
	S3Endpoint            string `envconfig:"S3_ENDPOINT"`
	S3Bucket              string `envconfig:"S3_BUCKET"`
	S3Region              string `envconfig:"S3_REGION"`
	S3AccessKey           string `envconfig:"S3_ACCESS_KEY"`
	S3SecretKey           string `envconfig:"S3_SECRET_KEY"`
}

func Build() fx.Option {
	return fx.Provide(func() (c Config, err error) {
		if err = envconfig.Process("", &c); err != nil {
			return c, fault.Wrap(err, fmsg.With("failed to parse configuration from environment variables"))
		}

		return
	})
}
