package httpserver

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/infrastructure/httpserver/ratelimit"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(NewRouter, ratelimit.New),
		fx.Invoke(NewServer),
	)
}
