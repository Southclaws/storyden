package openapi

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/labstack/echo/v5"
)

func TestHTTPErrorHandlerWritesProblemDetails(t *testing.T) {
	ctx := fctx.WithMeta(context.Background(), "code", "validation_error")
	err := fault.New("invalid account payload",
		fctx.With(ctx),
		fmsg.WithDesc("invalid account", "The account payload is invalid."),
		ftag.With(ftag.InvalidArgument),
	)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := HTTPErrorHandler(slog.New(slog.NewTextHandler(io.Discard, nil)))
	handler(c, err)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	contentType := rec.Header().Get(echo.HeaderContentType)
	if !strings.HasPrefix(contentType, ProblemJSONMediaType) {
		t.Fatalf("expected content type %q, got %q", ProblemJSONMediaType, contentType)
	}

	var problem APIError
	if err := json.Unmarshal(rec.Body.Bytes(), &problem); err != nil {
		t.Fatalf("failed to unmarshal problem details: %v", err)
	}

	if problem.Title == nil || *problem.Title != "The account payload is invalid." {
		t.Fatalf("expected fmsg description as title, got %#v", problem.Title)
	}

	if problem.Detail == nil || *problem.Detail != err.Error() {
		t.Fatalf("expected err.Error as detail, got %#v", problem.Detail)
	}

	if problem.Type == nil || *problem.Type != "urn:storyden:problem:invalid-argument" {
		t.Fatalf("expected problem type from ftag kind, got %#v", problem.Type)
	}

	if problem.TraceId == "" {
		t.Fatal("expected trace_id")
	}
}

func TestHTTPErrorHandlerPreservesEchoStatus(t *testing.T) {
	router := echo.New()
	router.HTTPErrorHandler = HTTPErrorHandler(slog.New(slog.NewTextHandler(io.Discard, nil)))
	router.GET("/exists", func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	router.GET("/forbidden", func(c *echo.Context) error {
		return fault.Wrap(echo.ErrForbidden)
	})

	cases := []struct {
		name   string
		method string
		path   string
		status int
	}{
		{name: "missing route", method: http.MethodGet, path: "/missing", status: http.StatusNotFound},
		{name: "unsupported method", method: http.MethodPost, path: "/exists", status: http.StatusMethodNotAllowed},
		{name: "wrapped status", method: http.MethodGet, path: "/forbidden", status: http.StatusForbidden},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(tc.method, tc.path, nil))

			if response.Code != tc.status {
				t.Fatalf("expected status %d, got %d: %s", tc.status, response.Code, response.Body.String())
			}

			if contentType := response.Header().Get(echo.HeaderContentType); !strings.HasPrefix(contentType, ProblemJSONMediaType) {
				t.Fatalf("expected problem details content type, got %q", contentType)
			}

			var problem APIError
			if err := json.Unmarshal(response.Body.Bytes(), &problem); err != nil {
				t.Fatalf("failed to unmarshal problem details: %v", err)
			}

			if problem.TraceId == "" {
				t.Fatal("expected trace_id")
			}
		})
	}
}
