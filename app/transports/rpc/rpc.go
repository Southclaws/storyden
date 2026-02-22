package rpc

import (
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Options(
		fx.Invoke(MountRPC),
	)
}
