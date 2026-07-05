package cachecontrol

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPDateIsSpecCompliant(t *testing.T) {
	t.Parallel()

	// a non-utc input must still render as a gmt http-date
	loc, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)

	in := time.Date(2026, 7, 5, 8, 30, 0, 0, loc)

	got := HTTPDate(in)

	assert.True(t, strings.HasSuffix(got, " GMT"), "http dates must end in GMT, got %q", got)
	assert.NotContains(t, got, "UTC", "http dates must not use the UTC zone abbreviation")

	// go's own http date parser must accept it and round-trip the same instant
	parsed, err := http.ParseTime(got)
	require.NoError(t, err, "http.ParseTime must accept the emitted Last-Modified/Retry-After value")
	assert.True(t, parsed.Equal(in), "round-tripped time must equal the input instant")
}
