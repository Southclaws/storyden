package config

import (
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
	CookieDomain           string `envconfig:"COOKIE_DOMAIN"          required:"true"`
	PublicWebAddress       string `envconfig:"PUBLIC_WEB_ADDRESS"     required:"true"`
	PublicApiAddress       string `envconfig:"PUBLIC_API_ADDRESS"     required:"true"`
	AmqpAddress            string `envconfig:"AMQP_ADDRESS"           default:"amqp://rabbit:5672"`
	HashKey                []byte `envconfig:"HASH_KEY"               required:"true"`
	BlockKey               []byte `envconfig:"BLOCK_KEY"              required:"true"`
}

func Build() fx.Option {
	return fx.Provide(func() (c Config, err error) {
		if err = envconfig.Process("", &c); err != nil {
			return c, err
		}

		return
	})
}
