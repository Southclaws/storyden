package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/errctx"
	"github.com/Southclaws/storyden/internal/errtag"
	"github.com/Southclaws/storyden/pkg/transports/http/openapi"
)

func newRouter(l *zap.Logger, cfg config.Config) *echo.Echo {
	router := echo.New()

	// TODO: Check errtags or fault context and react accordingly.
	// With: ctx.Response().WriteHeader( derived... )
	router.HTTPErrorHandler = func(err error, c echo.Context) {
		status := 500

		switch errtag.Tag(err).(type) {
		case errtag.Cancelled:
			// Do nothing
		case errtag.InvalidArgument:
			status = http.StatusBadRequest
		case errtag.NotFound:
			status = http.StatusNotFound
		case errtag.AlreadyExists:
			status = http.StatusConflict
		case errtag.PermissionDenied:
			status = http.StatusForbidden
		case errtag.Unauthenticated:
			status = http.StatusForbidden
		}

		ec := errctx.Unwrap(err)

		l.Info("request error",
			zap.String("error", err.Error()),
			zap.Any("metadata", ec),
		)

		// TODO: Settle on a nice way to do this.
		c.JSON(status, openapi.APIError{
			Error: err.Error(),
			// Message:              utils.Ref("An unhandled error occurred."),
			// Suggested:            utils.Ref("Please try again later or contact the site team/administrator."),
			AdditionalProperties: lo.MapValues(ec, func(k, v string) any { return v }),
		})
	}

	// Router must add all middleware before mounting routes. To add middleware,
	// simply depend on the router in a provider or invoker and do `router.Use`.
	// To mount routes use the lifecycle `OnStart` hook and mount them normally.

	return router
}
