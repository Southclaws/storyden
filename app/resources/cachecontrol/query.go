package cachecontrol

import (
	"time"

	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

// Query represents a HTTP conditional request query.
type Query struct {
	ETag          opt.Optional[string]
	ModifiedSince opt.Optional[time.Time]
}

// NewQuery must be constructed from a HTTP request's conditional headers.
func NewQuery(
	IfNoneMatch *string,
	IfModifiedSince *string,
) opt.Optional[Query] {
	if IfNoneMatch == nil && IfModifiedSince == nil {
		return opt.NewEmpty[Query]()
	}

	modifiedSince, err := opt.MapErr(opt.NewPtr(IfModifiedSince), parseConditionalRequestTime)
	if err != nil {
		return opt.NewEmpty[Query]()
	}

	return opt.New(Query{
		ETag:          opt.NewPtr((*string)(IfNoneMatch)),
		ModifiedSince: modifiedSince,
	})
}

// NotModified takes the current updated date of a resource and returns true if
// the cache control query includes a Is-Modified-Since header and the resource
// updated date is not after the header value. True means a 304 response header.
func (q Query) NotModified(resourceUpdated time.Time) bool {
	// truncate the resourceUpdated to the nearest second because the actual
	// HTTP header is already truncated but the database time is in nanoseconds.
	// If we didn't do this, the resource updated will always be slightly ahead.
	truncated := resourceUpdated.Truncate(time.Second)

	if ms, ok := q.ModifiedSince.Get(); ok {

		// If the resource update time is ahead of the HTTP Last-Modified check,
		// modified = 1, meaning the resource has been modified since the last
		// request and should be returned from the DB, instead of a 304 status.
		modified := truncated.Compare(ms)

		return modified <= 0
	}

	return false
}

func parseConditionalRequestTime(in openapi.IfModifiedSince) (time.Time, error) {
	return time.Parse(time.RFC1123, in)
}
