package http

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/internal/config"
)

func newRouter(l *zap.Logger, cfg config.Config) *echo.Echo {
	router := echo.New()

	// TODO: Check errtags or fault context and react accordingly.
	// With: ctx.Response().WriteHeader( derived... )
	router.HTTPErrorHandler = func(err error, ctx echo.Context) {
		l.Info("request error", zap.Error(err))
	}

	// Router must add all middleware before mounting routes. To add middleware,
	// simply depend on the router in a provider or invoker and do `router.Use`.
	// To mount routes use the lifecycle `OnStart` hook and mount them normally.

	return router
}
