package httpserver

import "go.uber.org/fx"

func Build() fx.Option {
	return fx.Options(
		fx.Provide(NewRouter),
		fx.Invoke(NewServer),
	)
}
