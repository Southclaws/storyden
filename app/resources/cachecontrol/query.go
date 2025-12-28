package cachecontrol

import (
	"time"

	"github.com/Southclaws/opt"
)

type ETag struct {
	Value string
	Time  time.Time
}

func NewETag(t time.Time) *ETag {
	return &ETag{
		Value: "t-" + t.UTC().Format(time.RFC3339Nano),
		Time:  t,
	}
}

func ParseETag(t string) ETag {
	// Strip quotes if present (ETags are quoted in HTTP headers)
	if len(t) >= 2 && t[0] == '"' && t[len(t)-1] == '"' {
		t = t[1 : len(t)-1]
	}

	// expected format: t-<time in RFC3339Nano>
	if len(t) < 3 || t[:2] != "t-" {
		return ETag{}
	}

	parsedTime, err := time.Parse(time.RFC3339Nano, t[2:])
	if err != nil {
		return ETag{}
	}

	return ETag{
		Value: t,
		Time:  parsedTime,
	}
}

func (t ETag) String() string {
	// Return quoted ETag for HTTP header
	return `"` + t.Value + `"`
}

// Query represents a HTTP conditional request query.
type Query struct {
	ETag          opt.Optional[ETag]
	ModifiedSince opt.Optional[time.Time]
}

// NewQuery must be constructed from a HTTP request's conditional headers.
func NewQuery(ifNoneMatch opt.Optional[string], ifModifiedSince opt.Optional[time.Time]) Query {
	return Query{
		ETag:          opt.Map(ifNoneMatch, ParseETag),
		ModifiedSince: ifModifiedSince,
	}
}

// NotModified takes the current updated date of a resource and returns true if
// the cache control query includes a Is-Modified-Since header and the resource
// updated date is not after the header value. True means a 304 response header.
func (q Query) Check(fn func() *time.Time) (*ETag, bool) {
	resourceUpdated := fn()
	if resourceUpdated == nil {
		return nil, false
	}

	if etag, ok := q.ETag.Get(); ok {
		modified := resourceUpdated.Compare(etag.Time)

		return NewETag(*resourceUpdated), modified <= 0
	}

	if ms, ok := q.ModifiedSince.Get(); ok {

		// truncate the resourceUpdated to the nearest second because the actual
		// HTTP header is already truncated but the DB time is in nanoseconds.
		// If we didn't do this the resource time will always be slightly ahead.
		truncated := resourceUpdated.Truncate(time.Second)

		// If the resource update time is ahead of the HTTP Last-Modified check,
		// modified = 1, meaning the resource has been modified since the last
		// request and should be returned from the DB, instead of a 304 status.
		modified := truncated.Compare(ms)

		return NewETag(*resourceUpdated), modified <= 0
	}

	return NewETag(*resourceUpdated), false
}
