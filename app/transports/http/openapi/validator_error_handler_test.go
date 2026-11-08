package openapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
)

func TestValidatorErrorHandlerHandlesSchemaErrorWithoutSchema(t *testing.T) {
	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())

	wrapped := echo.NewHTTPError(http.StatusBadRequest, "").Wrap(&openapi3.SchemaError{
		Reason: "value is not one of the allowed values",
	})
	var httpErr *echo.HTTPError
	require.ErrorAs(t, wrapped, &httpErr)

	err := ValidatorErrorHandler()(c, httpErr)
	require.Error(t, err)
	require.Equal(t, ftag.InvalidArgument, ftag.Get(err))
	require.Equal(t, "value is not one of the allowed values", fctx.Unwrap(err)["schema_error"])
}
