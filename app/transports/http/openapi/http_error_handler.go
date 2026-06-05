package openapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"syscall"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/trace"

	"github.com/Southclaws/storyden/internal/ent"
)

const ProblemJSONMediaType = "application/problem+json"
const UnknownTraceID = "unknown"

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
		title := fmsg.GetIssue(err)
		chain := fault.Flatten(err)
		meta := lo.MapValues(errctx, func(v, k string) any { return v })
		problem := NewAPIError(c.Request().Context(), status, errtag, title, errmsg, meta)

		if status == http.StatusInternalServerError {
			logger.Error(errmsg,
				slog.String("package", "http"),
				slog.String("title", title),
				slog.String("path", c.Path()),
				slog.String("tag", string(errtag)),
				slog.String("trace_id", problem.TraceId),
				slog.Any("metadata", errctx),
				slog.Any("trace", chain),
			)
		}

		//nolint:errcheck
		c.Blob(status, ProblemJSONMediaType, mustMarshalProblem(problem))
	}
}

func NewAPIError(ctx context.Context, status int, kind ftag.Kind, title string, detail string, metadata map[string]any) APIError {
	if title == "" {
		title = "Error"
		if statusText := http.StatusText(status); statusText != "" {
			title = statusText
		}
	}

	problem := APIError{
		TraceId: requestTraceID(ctx),
		Title:   &title,
	}

	if problemType := problemTypeFromKind(kind); problemType != nil {
		problem.Type = problemType
	}

	if detail != "" {
		problem.Detail = &detail
	}

	if len(metadata) > 0 {
		problem.Metadata = &metadata
	}

	return problem
}

func problemTypeFromKind(kind ftag.Kind) *string {
	if kind == ftag.None {
		return nil
	}

	problemType := "urn:storyden:problem:" + slugProblemKind(kind)

	return &problemType
}

func slugProblemKind(kind ftag.Kind) string {
	replacer := strings.NewReplacer("_", "-", " ", "-")

	return strings.ToLower(replacer.Replace(string(kind)))
}

func mustMarshalProblem(problem APIError) []byte {
	body, err := json.Marshal(problem)
	if err != nil {
		return []byte(`{"title":"Internal Server Error","trace_id":"` + problem.TraceId + `"}`)
	}

	return body
}

func requestTraceID(ctx context.Context) string {
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.HasTraceID() {
		return spanContext.TraceID().String()
	}

	return UnknownTraceID
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
		return ftag.PermissionDenied
	case http.StatusUnauthorized:
		return ftag.Unauthenticated
	case http.StatusBadGateway:
		return ftag.Cancelled
	default:
		return ftag.Internal
	}
}
