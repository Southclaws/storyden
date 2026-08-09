package openapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestValidatorErrorHandlerHandlesSchemaErrorWithoutSchema(t *testing.T) {
	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())

	err := ValidatorErrorHandler()(c, echo.NewHTTPError(http.StatusBadRequest).SetInternal(&openapi3.SchemaError{
		Reason: "value is not one of the allowed values",
	}))
	require.Error(t, err)
	require.Equal(t, ftag.InvalidArgument, ftag.Get(err))
	require.Equal(t, "value is not one of the allowed values", fctx.Unwrap(err)["schema_error"])
}
