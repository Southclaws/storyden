package middleware

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/transports/http/middleware/limiter"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session"
)

func Build() fx.Option {
	return fx.Provide(
		session.New,
		limiter.New,
	)
}
