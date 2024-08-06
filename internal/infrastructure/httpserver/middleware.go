package httpserver

import "net/http"

// Apply chains middleware functions onto the http.Handler in reverse order.
func Apply(h http.Handler, m ...func(http.Handler) http.Handler) http.Handler {
	if len(m) < 1 {
		return h
	}

	wrapped := h

	for i := len(m) - 1; i >= 0; i-- {
		fn := m[i]
		wrapped = fn(wrapped)
	}

	return wrapped
}
