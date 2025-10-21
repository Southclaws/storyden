package rpc

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Southclaws/storyden/app/services/reqinfo"
	"github.com/Southclaws/storyden/app/transports/http/middleware/origin"
)

type requestLogger struct {
	logger *slog.Logger
}

func newRequestLogger(logger *slog.Logger) *requestLogger {
	return &requestLogger{logger: logger}
}

func (m *requestLogger) WithLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			title := "incoming " + r.Method + " /rpc"

			origin := origin.GetOrigin(r.Context())

			clientAddress := reqinfo.GetClientAddress(r.Context())

			logger := m.logger.With(
				slog.String("http.request.header.origin", origin),
				slog.String("client.address", clientAddress),
				slog.String("http.request.method", r.Method),
			)

			defer func() {
				if recovery := recover(); recovery != nil {
					logger.Error(title,
						slog.String("error", fmt.Sprintf("%v", recovery)),
					)
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}
			}()

			logger.Info(title)

			next.ServeHTTP(w, r)
		})
	}
}
