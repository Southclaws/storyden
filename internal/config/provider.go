package config

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Provide(func() (c Config, err error) {
		if err = envconfig.Process("", &c); err != nil {
			return c, fault.Wrap(err, fmsg.With("failed to parse configuration from environment variables"))
		}

		return
	})
}
