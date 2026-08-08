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
	ETag opt.Optional[ETag]
}

// NewQuery must be constructed from a HTTP request's conditional headers.
func NewQuery(ifNoneMatch opt.Optional[string]) Query {
	return Query{
		ETag: opt.Map(ifNoneMatch, ParseETag),
	}
}

// Check compares the request ETag with the current resource timestamp. It
// returns true when the representation is unchanged and a 304 can be sent.
func (q Query) Check(fn func() *time.Time) (*ETag, bool) {
	resourceUpdated := fn()
	if resourceUpdated == nil {
		return nil, false
	}

	if etag, ok := q.ETag.Get(); ok {
		return NewETag(*resourceUpdated), resourceUpdated.Equal(etag.Time)
	}

	return NewETag(*resourceUpdated), false
}
