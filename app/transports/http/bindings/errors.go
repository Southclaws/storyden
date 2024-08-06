package bindings

import (
	"context"
	"errors"
	"io"
	"net/http"
	"syscall"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type ErrorHandler struct {
	fx.In
}

func (e *ErrorHandler) NewError(ctx context.Context, err error) *openapi.InternalServerErrorStatusCode {
	errmsg := err.Error()
	errtag, status := categorise(err)
	errctx := fctx.Unwrap(err)
	message := fmsg.GetIssue(err)
	chain := fault.Flatten(err)

	if status == http.StatusInternalServerError {
		l.Error(errmsg,
			zap.String("package", "http"),
			zap.String("message", message),
			zap.String("path", c.Path()),
			zap.String("tag", string(errtag)),
			zap.Any("metadata", errctx),
			zap.Any("trace", chain),
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

	return &openapi.InternalServerErrorStatusCode{
		StatusCode: 500,
		Response: openapi.APIError{
			Error:    errmsg,
			Message:  openapi.NewOptString(errormessage),
			Metadata: openapi.NewOptAPIErrorMetadata(errormetadata),
		},
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
		return http.StatusUnauthorized
	case ftag.Unauthenticated:
		return http.StatusForbidden
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
	default:
		return ftag.Internal
	}
}
