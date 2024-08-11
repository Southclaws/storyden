package reqlog

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Southclaws/storyden/app/transports/http/middleware/origin"
	"go.uber.org/zap"
)

func WithLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			origin := origin.GetOrigin(r.Context())

			// log entries should be in the form "GET /a/b/c".
			title := r.Method + " " + r.URL.Path

			defer func() {
				log := logger.With(
					zap.Duration("duration", time.Since(start)),
					zap.String("query", r.URL.Query().Encode()),
					zap.Int64("body", r.ContentLength),
					zap.String("ip", r.RemoteAddr),
					zap.String("origin", origin),
				)

				if recovery := recover(); recovery != nil {
					err := func(v any) error {
						if e, ok := v.(error); ok {
							return e
						} else {
							return fmt.Errorf("%v", v)
						}
					}(recovery)

					errorlog := title + ": " + err.Error()

					log.Error(errorlog, zap.Error(err))

					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				log.Info(title)
			}()

			next.ServeHTTP(w, r)
		})
	}
}
