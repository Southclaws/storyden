package glue

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/transports/openapi/openapi"
	"github.com/Southclaws/storyden/internal/errctx"
	"github.com/Southclaws/storyden/internal/errtag"
)

// HTTPErrorHandler provides an error handler function for use with the Echo
// router. The purpose of this implementation is to map application level errors
// to HTTP status codes. This is achieved (currently) with the use of a library
// called errtag which enables the decoration of error chains with a basic kind
// of category which helps organise the kind of errors that occur within an app.
func HTTPErrorHandler(l *zap.Logger) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		status := 500

		switch errtag.Tag(err).(type) {
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

		meta := lo.MapValues(ec, func(v, k string) any {
			return v
		})

		c.JSON(status, openapi.APIError{
			Error:    err.Error(),
			Metadata: &meta,
		})
	}
}
