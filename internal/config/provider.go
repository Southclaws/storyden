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

		return ApplyDefaults(c)
	})
}

var errInvalidConfiguration = fault.New("invalid configuration")

func ApplyDefaults(c Config) (Config, error) {
	if c.ImageMemoryBudget == 0 {
		c.ImageMemoryBudget = DefaultImageMemoryBudget
	}
	if c.ImageWorkers == 0 {
		c.ImageWorkers = DefaultImageWorkers
	}

	if c.ImageMemoryBudget < 0 {
		return c, fault.Wrap(errInvalidConfiguration, fmsg.With("negative image processing memory budget"))
	}
	if c.ImageWorkers < 0 {
		return c, fault.Wrap(errInvalidConfiguration, fmsg.With("negative image worker count"))
	}

	return c, nil
}
