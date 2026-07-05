package cachecontrol

import (
	"net/http"
	"time"
)

// httpDate formats a time as a spec-compliant http-date (rfc 7231) for headers
// such as last-modified and retry-after, always rendered in gmt
func HTTPDate(t time.Time) string {
	return t.UTC().Format(http.TimeFormat)
}
