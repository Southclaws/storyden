package cachecontrol

import (
	"time"

	"github.com/Southclaws/opt"
)

// Query represents a HTTP conditional request query.
type Query struct {
	ETag          opt.Optional[string]
	ModifiedSince opt.Optional[time.Time]
}

// NewQuery must be constructed from a HTTP request's conditional headers.
func NewQuery(ifNoneMatch opt.Optional[string], ifModifiedSince opt.Optional[time.Time]) Query {
	return Query{
		ETag:          ifNoneMatch,
		ModifiedSince: ifModifiedSince,
	}
}

// NotModified takes the current updated date of a resource and returns true if
// the cache control query includes a Is-Modified-Since header and the resource
// updated date is not after the header value. True means a 304 response header.
func (q Query) NotModified(fn func() *time.Time) bool {
	if ms, ok := q.ModifiedSince.Get(); ok {
		resourceUpdated := fn()
		if resourceUpdated == nil {
			return false
		}

		// truncate the resourceUpdated to the nearest second because the actual
		// HTTP header is already truncated but the DB time is in nanoseconds.
		// If we didn't do this the resource time will always be slightly ahead.
		truncated := resourceUpdated.Truncate(time.Second)

		// If the resource update time is ahead of the HTTP Last-Modified check,
		// modified = 1, meaning the resource has been modified since the last
		// request and should be returned from the DB, instead of a 304 status.
		modified := truncated.Compare(ms)

		return modified <= 0
	}

	return false
}
