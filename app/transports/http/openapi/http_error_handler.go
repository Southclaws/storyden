package openapi

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"syscall"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/internal/ent"
)

// HTTPErrorHandler provides an error handler function for use with the Echo
// router. The purpose of this implementation is to map application level errors
// to HTTP status codes. This is achieved (currently) with the use of a library
// called "ftag" which enables the decoration of error chains with a basic kind
// of category which helps organise the kind of errors that occur within an app.
func HTTPErrorHandler(logger *slog.Logger) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		errmsg := err.Error()
		errtag, status := categorise(err)
		errctx := fctx.Unwrap(err)
		message := fmsg.GetIssue(err)
		chain := fault.Flatten(err)

		if status == http.StatusInternalServerError {
			logger.Error(errmsg,
				slog.String("package", "http"),
				slog.String("message", message),
				slog.String("path", c.Path()),
				slog.String("tag", string(errtag)),
				slog.Any("metadata", errctx),
				slog.Any("trace", chain),
			)
		}

		meta := lo.MapValues(errctx, func(v, k string) any { return v })
		errormessage := opt.NewIf(message, func(s string) bool { return s != "" }).Ptr()
		errormetadata := opt.NewIf(meta, func(m map[string]any) bool { return len(m) > 0 }).Ptr()

		//nolint:errcheck
		c.JSON(status, APIError{
			Error:    errmsg,
			Message:  errormessage,
			Metadata: errormetadata,
		})
	}
}

func categorise(err error) (ftag.Kind, int) {
	errtag := ftag.Get(err)
	status := statusFromErrorKind(errtag)

	var he *echo.HTTPError
	if errors.As(err, &he) {
		return errorKindFromStatus(he.Code), he.Code
	}

	if errors.Is(err, context.Canceled) {
		return ftag.Cancelled, http.StatusBadRequest
	}

	if errors.Is(err, io.ErrClosedPipe) {
		return ftag.Cancelled, http.StatusBadRequest
	}

	if errors.Is(err, syscall.EPIPE) {
		return ftag.Cancelled, http.StatusBadRequest
	}

	if ent.IsNotFound(err) {
		return ftag.NotFound, http.StatusNotFound
	}

	return errtag, status
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
		return http.StatusUnauthorized
	case ftag.Cancelled:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

func errorKindFromStatus(s int) ftag.Kind {
	switch s {
	case http.StatusBadRequest:
		return ftag.InvalidArgument
	case http.StatusNotFound:
		return ftag.NotFound
	case http.StatusConflict:
		return ftag.AlreadyExists
	case http.StatusForbidden:
		return ftag.Unauthenticated
	case http.StatusUnauthorized:
		return ftag.PermissionDenied
	case http.StatusBadGateway:
		return ftag.Cancelled
	default:
		return ftag.Internal
	}
}
