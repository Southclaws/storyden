package reqlog

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/Southclaws/storyden/app/transports/http/middleware/origin"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/kv"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/spanner"
)

type Middleware struct {
	ins spanner.Instrumentation
}

func New(ins spanner.Builder) *Middleware {
	return &Middleware{
		ins: ins.Build(),
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

			ctx, span := m.ins.InstrumentNamed(r.Context(), title,
				kv.String("http.request.header.origin", origin),
				kv.String("client.address", r.RemoteAddr),
				kv.String("http.request.method", r.Method),
				kv.String("url.query", r.URL.Query().Encode()),
				kv.Int("http.request.body.size", int(r.ContentLength)),
			)
			defer span.End()

			defer func() {
				span.Annotate(
					kv.Duration("duration", time.Since(start)),
					kv.Int("http.response.status_code", wr.statusCode),
				)

				logger := span.Logger()

				logger.Info(title)

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

					logger.Error(errorlog,
						slog.String("error", span.Wrap(err, errorlog, kv.String("trace", string(trace))).Error()),
					)

					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}()

			next.ServeHTTP(wr, r.WithContext(ctx))
		})
	}
}
