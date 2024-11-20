package middleware

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/transports/http/middleware/limiter"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
)

func Build() fx.Option {
	return fx.Provide(
		session_cookie.New,
		limiter.New,
	)
}
