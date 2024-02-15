package testutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func AssertRequest[T interface {
	StatusCode() int
}](v T, err error) func(t *testing.T, want int) T {
	return func(t *testing.T, want int) T {
		require.Equal(t, want, v.StatusCode())

		return v
	}
}
