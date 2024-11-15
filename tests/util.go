package tests

import (
	"net/http"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type WithStatusCode interface {
	StatusCode() int
}

func Ok(t *testing.T, err error, resp WithStatusCode) {
	t.Helper()

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode())
	logAPIError(t, resp)
}

func Status(t *testing.T, err error, resp WithStatusCode, status int) {
	t.Helper()

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, status, resp.StatusCode())
	logAPIError(t, resp)
}

func AssertRequest[T interface {
	StatusCode() int
}](v T, err error) func(t *testing.T, want int) T {
	return func(t *testing.T, want int) T {
		require.Equal(t, want, v.StatusCode())

		return v
	}
}

type ResponseShape struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      any
	JSONDefault  *openapi.InternalServerError
}

func logAPIError(t *testing.T, resp WithStatusCode) {
	if resp.StatusCode() != http.StatusOK {
		if ae := getAPIError(resp); ae != nil {
			t.Logf(`%s message: "%v"`,
				ae.Error,
				opt.NewPtr(ae.Message).OrZero(),
			)
		}
	}
}

func getAPIError(resp WithStatusCode) *openapi.APIError {
	var out ResponseShape
	err := mapstructure.Decode(resp, &out)
	if err != nil {
		return nil
	}

	return out.JSONDefault
}
