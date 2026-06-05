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
	"github.com/labstack/echo/v4"
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
	handler(err, c)

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
