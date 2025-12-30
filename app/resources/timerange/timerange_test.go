package timerange_test

import (
	"testing"
	"time"

	"github.com/Southclaws/storyden/app/resources/timerange"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDateOnly(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	tr, err := timerange.Parse("2025-11-01/2025-11-30")
	r.NoError(err)

	start, ok := tr.Start.Get()
	a.True(ok)
	a.Equal(2025, start.Year())
	a.Equal(time.November, start.Month())
	a.Equal(1, start.Day())

	end, ok := tr.End.Get()
	a.True(ok)
	a.Equal(2025, end.Year())
	a.Equal(time.November, end.Month())
	a.Equal(30, end.Day())
}

func TestParseRFC3339(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	tr, err := timerange.Parse("2025-11-01T00:00:00Z/2025-11-30T23:59:59Z")
	r.NoError(err)

	start, ok := tr.Start.Get()
	r.True(ok)
	a.Equal(2025, start.Year())
	a.Equal(time.November, start.Month())
	a.Equal(1, start.Day())
	a.Equal(0, start.Hour())
	a.Equal(0, start.Minute())
	a.Equal(0, start.Second())

	end, ok := tr.End.Get()
	r.True(ok)
	a.Equal(2025, end.Year())
	a.Equal(time.November, end.Month())
	a.Equal(30, end.Day())
	a.Equal(23, end.Hour())
	a.Equal(59, end.Minute())
	a.Equal(59, end.Second())
}

func TestParseSingleDate(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	tr, err := timerange.Parse("2025-11-01")
	r.NoError(err)

	start, ok := tr.Start.Get()
	a.True(ok)
	a.Equal(2025, start.Year())

	_, ok = tr.End.Get()
	a.False(ok)
}

func TestParseEndOnly(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	tr, err := timerange.Parse("/2025-11-30")
	r.NoError(err)

	_, ok := tr.Start.Get()
	a.False(ok)

	end, ok := tr.End.Get()
	a.True(ok)
	a.Equal(2025, end.Year())
}

func TestParseInvalidFormat(t *testing.T) {
	r := require.New(t)

	_, err := timerange.Parse("2025-11-01/2025-11-15/2025-11-30")
	r.Error(err)
	r.Contains(err.Error(), "invalid time range format")
	r.Contains(err.Error(), "2025-11-01/2025-11-15/2025-11-30")
}
