package reqlog

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/Southclaws/storyden/app/services/reqinfo"
	"github.com/Southclaws/storyden/app/transports/http/middleware/origin"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/kv"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/servertiming"
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
	wrote      bool
	beforeSend func(http.Header)
}

func (w *withStatus) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *withStatus) WriteHeader(code int) {
	if w.wrote {
		return
	}
	w.wrote = true
	w.statusCode = code
	if w.beforeSend != nil {
		w.beforeSend(w.Header())
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *withStatus) Write(p []byte) (int, error) {
	if !w.wrote {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}

func (w *withStatus) ReadFrom(r io.Reader) (int64, error) {
	if !w.wrote {
		w.WriteHeader(http.StatusOK)
	}
	if rf, ok := w.ResponseWriter.(io.ReaderFrom); ok {
		return rf.ReadFrom(r)
	}
	return io.Copy(w.ResponseWriter, r)
}

func (m *Middleware) WithLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			origin := origin.GetOrigin(r.Context())

			// log entries should be in the form "GET /a/b/c".
			title := r.Method + " " + r.URL.Path
			clientAddress := reqinfo.GetClientAddress(r.Context())
			if clientAddress == "" {
				clientAddress = r.RemoteAddr
			}

			ctx, span := m.ins.InstrumentNamed(r.Context(), title,
				kv.String("http.request.header.origin", origin),
				kv.String("client.address", clientAddress),
				kv.String("http.request.method", r.Method),
				kv.String("url.query", r.URL.Query().Encode()),
				kv.Int("http.request.body.size", int(r.ContentLength)),
			)
			defer span.End()

			collector := servertiming.NewCollector()
			ctx = servertiming.WithCollector(ctx, collector)

			wr := &withStatus{
				ResponseWriter: w,
				beforeSend: func(h http.Header) {
					if value := collector.HeaderValue(); value != "" {
						h.Set("Server-Timing", value)
					}
				},
			}

			defer func() {
				if !wr.wrote {
					if value := collector.HeaderValue(); value != "" {
						wr.Header().Set("Server-Timing", value)
					}
				}

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

					err = span.Wrap(err, errorlog)

					logger.Error(errorlog,
						slog.String("error", err.Error()),
						slog.Any("trace", trace),
					)

					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}()

			next.ServeHTTP(wr, r.WithContext(ctx))
		})
	}
}
