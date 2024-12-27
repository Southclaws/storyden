package reqlog

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/Southclaws/storyden/app/transports/http/middleware/origin"
	"go.uber.org/zap"
)

type Middleware struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

type withStatus struct {
	http.ResponseWriter
	statusCode int
}

func (w *withStatus) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (lrw *withStatus) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (m *Middleware) WithLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			origin := origin.GetOrigin(r.Context())

			// log entries should be in the form "GET /a/b/c".
			title := r.Method + " " + r.URL.Path

			wr := &withStatus{ResponseWriter: w}

			defer func() {
				log := m.logger.With(
					zap.Duration("duration", time.Since(start)),
					zap.String("query", r.URL.Query().Encode()),
					zap.Int64("body", r.ContentLength),
					zap.String("ip", r.RemoteAddr),
					zap.String("origin", origin),
					zap.Int("status", wr.statusCode),
				)

				if recovery := recover(); recovery != nil {
					err := func(v any) error {
						if e, ok := v.(error); ok {
							return e
						} else {
							return fmt.Errorf("%v", v)
						}
					}(recovery)

					trace := debug.Stack()

					errorlog := title + ": " + err.Error()

					log.Error(errorlog, zap.Error(err), zap.String("trace", string(trace)))

					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				log.Info(title)
			}()

			next.ServeHTTP(wr, r)
		})
	}
}
