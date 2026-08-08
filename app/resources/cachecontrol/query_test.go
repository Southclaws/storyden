package cachecontrol

import (
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/assert"
)

func TestStaleETagReturnsColdResponse(t *testing.T) {
	updated := time.Date(2026, time.August, 8, 12, 0, 0, 500, time.UTC)
	oldETag := NewETag(updated.Add(-time.Nanosecond)).String()
	query := NewQuery(opt.New(oldETag))

	_, notModified := query.Check(func() *time.Time { return &updated })

	assert.False(t, notModified)
}

func TestFutureETagReturnsColdResponse(t *testing.T) {
	updated := time.Date(2026, time.August, 8, 12, 0, 0, 500, time.UTC)
	futureETag := NewETag(updated.Add(time.Hour)).String()
	query := NewQuery(opt.New(futureETag))

	_, notModified := query.Check(func() *time.Time { return &updated })

	assert.False(t, notModified, "clock skew or a future validator must not produce 304")
}

func TestExactETagReturnsNotModified(t *testing.T) {
	updated := time.Date(2026, time.August, 8, 12, 0, 0, 500, time.UTC)
	query := NewQuery(opt.New(NewETag(updated).String()))

	_, notModified := query.Check(func() *time.Time { return &updated })

	assert.True(t, notModified)
}
