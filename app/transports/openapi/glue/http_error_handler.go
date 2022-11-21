package glue

import (
	"net/http"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/internal/openapi"
)

// HTTPErrorHandler provides an error handler function for use with the Echo
// router. The purpose of this implementation is to map application level errors
// to HTTP status codes. This is achieved (currently) with the use of a library
// called "ftag" which enables the decoration of error chains with a basic kind
// of category which helps organise the kind of errors that occur within an app.
func HTTPErrorHandler(l *zap.Logger) func(err error, c echo.Context) {
	getLoggerForTag := func(tag ftag.Kind) func(msg string, fields ...zap.Field) {
		switch tag {
		case ftag.Internal:
			return l.Error

		// We don't need to log these in prod.
		case ftag.PermissionDenied:
			return l.Debug
		case ftag.Unauthenticated:
			return l.Debug

		default:
			return l.Info
		}
	}

	return func(err error, c echo.Context) {
		errmsg := err.Error()
		errtag := ftag.Get(err)
		status := statusFromErrorKind(errtag)
		errctx := fctx.Unwrap(err)
		message := fmsg.GetIssue(err)
		chain := fault.Flatten(err)

		fn := getLoggerForTag(errtag)

		fn(errmsg,
			zap.String("package", "http"),
			zap.String("message", message),
			zap.String("path", c.Path()),
			zap.String("tag", string(errtag)),
			zap.Any("metadata", errctx),
			zap.Any("trace", chain.Errors),
		)

		meta := lo.MapValues(errctx, func(v, k string) any { return v })
		errormessage := opt.NewIf(message, func(s string) bool { return s != "" }).Ptr()
		errormetadata := opt.NewIf(meta, func(m map[string]any) bool { return len(m) > 0 }).Ptr()

		c.JSON(status, openapi.APIError{
			Error:    errmsg,
			Message:  errormessage,
			Metadata: errormetadata,
		})
	}
}

func statusFromErrorKind(k ftag.Kind) int {
	switch k {
	case ftag.InvalidArgument:
		return http.StatusBadRequest
	case ftag.NotFound:
		return http.StatusNotFound
	case ftag.AlreadyExists:
		return http.StatusConflict
	case ftag.PermissionDenied:
		return http.StatusForbidden
	case ftag.Unauthenticated:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
