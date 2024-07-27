package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type WithStatusCode interface {
	StatusCode() int
}

func Ok(t *testing.T, err error, resp WithStatusCode) {
	t.Helper()

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode())
}

func Status(t *testing.T, err error, resp WithStatusCode, status int) {
	t.Helper()

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, status, resp.StatusCode())
}

func AssertRequest[T interface {
	StatusCode() int
}](v T, err error) func(t *testing.T, want int) T {
	return func(t *testing.T, want int) T {
		require.Equal(t, want, v.StatusCode())

		return v
	}
}
