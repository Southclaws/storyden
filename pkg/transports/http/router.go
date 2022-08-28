package http

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/fault"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/pkg/transports/http/openapi"
)

func newRouter(l *zap.Logger, cfg config.Config) *echo.Echo {
	router := echo.New()

	// TODO: Check errtags or fault context and react accordingly.
	// With: ctx.Response().WriteHeader( derived... )
	router.HTTPErrorHandler = func(err error, c echo.Context) {
		errdata := fault.ErrorData(err)
		l.Info("request error", zap.Any("meta", errdata))

		// TODO: Settle on a nice way to do this.
		// TODO: Handle error categories mapping to HTTP statuses too.
		c.JSON(500, map[string]any{
			"details": openapi.APIError{
				Error:     err.Error(),
				Message:   utils.Ref("An unhandled error occurred."),
				Suggested: utils.Ref("Please try again later or contact the site team/administrator."),
			},
			"context": errdata,
		})
	}

	// Router must add all middleware before mounting routes. To add middleware,
	// simply depend on the router in a provider or invoker and do `router.Use`.
	// To mount routes use the lifecycle `OnStart` hook and mount them normally.

	return router
}
