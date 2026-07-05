package bindings

import (
	"net/http"
	"time"
)

func httpDate(t time.Time) string {
	return t.UTC().Format(http.TimeFormat)
}
