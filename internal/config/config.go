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

	DatabaseURL            string `envconfig:"DATABASE_URL"           required:"true"`
	MailTemplatesDirectory string `envconfig:"MAIL_TEMPLATES_DIRECTORY" default:"email"`
	ListenAddr             string `envconfig:"LISTEN_ADDR"            default:"0.0.0.0:8000"`
	CookieDomain           string `envconfig:"COOKIE_DOMAIN"          default:"localhost"`
	SessionKey             string `envconfig:"SESSION_KEY"            required:"true"`
	PublicWebAddress       string `envconfig:"PUBLIC_WEB_ADDRESS"     default:"http://localhost:3000"`
	PublicApiAddress       string `envconfig:"PUBLIC_API_ADDRESS"     default:"http://localhost:8000"`
	AmqpAddress            string `envconfig:"AMQP_ADDRESS"           default:"amqp://rabbit:5672"`

	S3Endpoint  string `envconfig:"S3_ENDPOINT"   required:"true"`
	S3Bucket    string `envconfig:"S3_BUCKET"     required:"true"`
	S3Region    string `envconfig:"S3_REGION"     required:"true"`
	S3AccessKey string `envconfig:"S3_ACCESS_KEY" required:"true"`
	S3SecretKey string `envconfig:"S3_SECRET_KEY" required:"true"`
}

func Build() fx.Option {
	return fx.Provide(func() (c Config, err error) {
		if err = envconfig.Process("", &c); err != nil {
			return c, fault.Wrap(err, fmsg.With("failed to parse configuration from environment variables"))
		}

		return
	})
}
