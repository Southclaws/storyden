package config

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Provide(func() (Config, error) {
		c, err := Load()
		if err != nil {
			return c, fault.Wrap(err, fmsg.With("failed to parse configuration from environment variables"))
		}

		return c, nil
	})
}
